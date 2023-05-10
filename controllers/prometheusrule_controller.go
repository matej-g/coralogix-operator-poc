package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	coralogixv1 "coralogix-operator-poc/api/v1"
	"coralogix-operator-poc/controllers/clientset"
	prometheus "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=get;list;watch

//+kubebuilder:rbac:groups=coralogix.coralogix,resources=recordingrulegroups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coralogix.coralogix,resources=recordingrulegroups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=coralogix.coralogix,resources=recordingrulegroups/finalizers,verbs=update

//+kubebuilder:rbac:groups=coralogix.coralogix,resources=alerts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coralogix.coralogix,resources=alerts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=coralogix.coralogix,resources=alerts/finalizers,verbs=update

// PrometheusRuleReconciler reconciles a PrometheusRule object
type PrometheusRuleReconciler struct {
	client.Client
	CoralogixClientSet *clientset.ClientSet
	Scheme             *runtime.Scheme
}

func (r *PrometheusRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	prometheusRuleCRD := &prometheus.PrometheusRule{}
	if err := r.Get(ctx, req.NamespacedName, prometheusRuleCRD); err != nil && !errors.IsNotFound(err) {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	if createAndTrackRecordingRules(prometheusRuleCRD) {
		ruleGroupSetCRD := &coralogixv1.RecordingRuleGroup{}
		if err := r.Client.Get(ctx, req.NamespacedName, ruleGroupSetCRD); err != nil {
			if errors.IsNotFound(err) {
				log.V(1).Info(fmt.Sprintf("Couldn't find RecordingRuleSet Namespace: %s, Name: %s. Trying to create.", req.Namespace, req.Name))
				//Meaning there's a PrometheusRule with that NamespacedName but not RecordingRuleGroupSet accordingly (so creating it).
				if ruleGroupSetCRD.Spec, err = prometheusRuleToRuleGroupSet(prometheusRuleCRD); err != nil {
					log.Error(err, "Received an error while Converting PrometheusRule to RecordingRuleGroupSet", "PrometheusRule Name", prometheusRuleCRD.Name)
					return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
				}
				ruleGroupSetCRD.Namespace = req.Namespace
				ruleGroupSetCRD.Name = req.Name
				if err = r.Create(ctx, ruleGroupSetCRD); err != nil {
					log.Error(err, "Received an error while trying to create RecordingRuleGroupSet CRD", "RecordingRuleGroupSet Name", ruleGroupSetCRD.Name)
					return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
				}

			} else {
				return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
			}
		}

		//Converting the PrometheusRule to the desired RecordingRuleGroupSet.
		var err error
		if ruleGroupSetCRD.Spec, err = prometheusRuleToRuleGroupSet(prometheusRuleCRD); err != nil {
			log.Error(err, "Received an error while Converting PrometheusRule to RecordingRuleGroupSet", "PrometheusRule Name", prometheusRuleCRD.Name)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}

		if err = r.Client.Update(ctx, ruleGroupSetCRD); err != nil {
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
	}

	for _, group := range prometheusRuleCRD.Spec.Groups {
		for _, rule := range group.Rules {
			if alertName := rule.Annotations["cxAlertName"]; rule.Alert != "" && alertName != "" {
				alertCRD := &coralogixv1.Alert{}
				req.Name = alertName
				if err := r.Client.Get(ctx, req.NamespacedName, alertCRD); err != nil {
					if errors.IsNotFound(err) {
						log.V(1).Info(fmt.Sprintf("Couldn't find Alert Namespace: %s, Name: %s. Trying to create.", req.Namespace, req.Name))
						alertCRD.Spec = prometheusInnerRuleToCoralogixAlert(rule)
						alertCRD.Namespace = req.Namespace
						alertCRD.Name = alertName
						if err = r.Create(ctx, alertCRD); err != nil {
							log.Error(err, "Received an error while trying to create Alert CRD", "Alert Name", alertCRD.Name)
							return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
						}
					} else {
						return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
					}
				}

				//Converting the PrometheusRule to the desired Alert.
				alertCRD.Spec = prometheusInnerRuleToCoralogixAlert(rule)
				if err := r.Client.Update(ctx, alertCRD); err != nil {
					return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
				}
			}
		}
	}

	return ctrl.Result{RequeueAfter: defaultRequeuePeriod}, nil

}

func createAndTrackRecordingRules(prometheusRule *prometheus.PrometheusRule) bool {
	if toCreateStr, ok := prometheusRule.Labels["cxCreateAndTrackRecordingRule"]; ok && toCreateStr == "true" {
		return true
	}
	return false
}

func createAndTrackAlerts(prometheusRule *prometheus.PrometheusRule) bool {
	if toCreateStr, ok := prometheusRule.Labels["cxCreateAndTrackAlerts"]; ok && toCreateStr == "true" {
		return true
	}
	return false
}

func prometheusRuleToRuleGroupSet(prometheusRule *prometheus.PrometheusRule) (coralogixv1.RecordingRuleGroupSpec, error) {
	if len(prometheusRule.Spec.Groups) > 0 {
		s, _ := time.ParseDuration(string(prometheusRule.Spec.Groups[0].Interval))
		intervalSeconds := int32(s.Seconds())

		limit := int64(60)
		if limitStr, ok := prometheusRule.Annotations["cx_limit"]; ok && limitStr != "" {
			if limitInt, err := strconv.Atoi(limitStr); err != nil {
				return coralogixv1.RecordingRuleGroupSpec{}, err
			} else {
				limit = int64(limitInt)
			}
		}

		rules := prometheusInnerRulesToCoralogixInnerRules(prometheusRule.Spec.Groups[0].Rules)

		return coralogixv1.RecordingRuleGroupSpec{
			Name:            prometheusRule.Name,
			IntervalSeconds: intervalSeconds,
			Limit:           limit,
			Rules:           rules,
		}, nil
	}

	return coralogixv1.RecordingRuleGroupSpec{}, nil
}

func prometheusInnerRuleToCoralogixAlert(prometheusRule prometheus.Rule) coralogixv1.AlertSpec {
	notifyEveryMin := 0
	if cxNotifyEveryMin, ok := prometheusRule.Annotations["cxNotifyEveryMin"]; ok {
		notifyEveryMin, _ = strconv.Atoi(cxNotifyEveryMin)
	} else {
		duration, _ := time.ParseDuration(string(prometheusRule.For))
		notifyEveryMin = int(duration.Minutes())
	}

	minNonNullValuesPercentage := 0
	if cxMinNonNullValuesPercentage, ok := prometheusRule.Annotations["cxMinNonNullValuesPercentage"]; ok {
		minNonNullValuesPercentage, _ = strconv.Atoi(cxMinNonNullValuesPercentage)
	}

	timeWindow := prometheusAlertForToCoralogixPromqlAlertTimeWindow[prometheusRule.For]
	//expr, _ := promql.ParseExpr(prometheusRule.Expr.StrVal)

	return coralogixv1.AlertSpec{
		Severity: coralogixv1.AlertSeverityInfo,
		Notifications: &coralogixv1.Notifications{
			NotifyEveryMin: &notifyEveryMin,
		},
		Name: prometheusRule.Alert,
		AlertType: coralogixv1.AlertType{
			Metric: &coralogixv1.Metric{
				Promql: &coralogixv1.Promql{
					SearchQuery: prometheusRule.Expr.StrVal,
					Conditions: coralogixv1.PromqlConditions{
						TimeWindow:                 timeWindow,
						AlertWhen:                  coralogixv1.AlertWhenLessThan,
						Threshold:                  resource.MustParse("5"),
						SampleThresholdPercentage:  100,
						MinNonNullValuesPercentage: &minNonNullValuesPercentage,
					},
				},
			},
		},
	}
}

var prometheusAlertForToCoralogixPromqlAlertTimeWindow = map[prometheus.Duration]coralogixv1.MetricTimeWindow{
	"1m":  coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowMinute),
	"5m":  coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowFiveMinutes),
	"10m": coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowTenMinutes),
	"15m": coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowFifteenMinutes),
	"20m": coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowTwentyMinutes),
	"30m": coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowThirtyMinutes),
	"1h":  coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowHour),
	"2h":  coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowTwelveHours),
	"4h":  coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowFourHours),
	"6h":  coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowSixHours),
	"12":  coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowTwelveHours),
	"24h": coralogixv1.MetricTimeWindow(coralogixv1.TimeWindowTwentyFourHours),
}

func prometheusInnerRulesToCoralogixInnerRules(rules []prometheus.Rule) []coralogixv1.RecordingRule {
	result := make([]coralogixv1.RecordingRule, 0)
	for _, r := range rules {
		if r.Record != "" {
			rule := prometheusInnerRuleToCoralogixInnerRule(r)
			result = append(result, rule)
		}
	}
	return result
}

func prometheusInnerRuleToCoralogixInnerRule(rule prometheus.Rule) coralogixv1.RecordingRule {
	return coralogixv1.RecordingRule{
		Record: rule.Record,
		Expr:   rule.Expr.StrVal,
		Labels: rule.Labels,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PrometheusRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&prometheus.PrometheusRule{}).
		Watches(&source.Kind{Type: &prometheus.PrometheusRule{}},
			handler.EnqueueRequestsFromMapFunc(findPrometheusRules),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		Complete(r)
}

func findPrometheusRules(v client.Object) []reconcile.Request {
	prometheusRule := v.(*prometheus.PrometheusRule)
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: prometheusRule.Namespace,
			Name:      prometheusRule.Name,
		},
	}

	return []reconcile.Request{request}
}
