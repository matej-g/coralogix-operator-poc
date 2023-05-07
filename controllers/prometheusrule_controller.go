package controllers

import (
	"context"
	"strconv"
	"time"

	coralogixv1 "coralogix-operator-poc/api/v1"
	"coralogix-operator-poc/controllers/clientset"
	prometheus "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
	if err := r.Get(ctx, req.NamespacedName, prometheusRuleCRD); err != nil && !errors.IsNotFound(err) { //Checking if there's a PrometheusRule to track with that NamespacedName.
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
			log.V(1).Info("Fetched PrometheusRule. Trying to create/update a RecordingRuleGroupSet accordingly.", "PrometheusRule name", prometheusRuleCRD.Name)
			if errors.IsNotFound(err) {
				//Meaning there's a PrometheusRule with that NamespacedName but not RecordingRuleGroupSet accordingly (so creating it).
				if ruleGroupSetCRD.Spec, err = prometheusRuleToRuleGroupSet(prometheusRuleCRD); err != nil {
					log.Error(err, "Received an error while Converting PrometheusRule to RecordingRuleGroupSet", "PrometheusRule Name", prometheusRuleCRD.Name)
					return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
				}
				if err = r.Create(ctx, ruleGroupSetCRD); err != nil {
					log.Error(err, "Received an error while trying to create RecordingRuleGroupSet CRD", "RecordingRuleGroupSet Name", ruleGroupSetCRD.Name)
					return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
				}
				if err = r.Client.Get(ctx, req.NamespacedName, ruleGroupSetCRD); err != nil {
					log.Error(err, "Received an error while trying to fetch RecordingRuleGroupSet CRD", "RecordingRuleGroupSet Name", ruleGroupSetCRD.Name)
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

	//if createAndTrackAlerts(prometheusRuleCRD) {
	//	//Converting the PrometheusRule to the desired RecordingRuleGroupSet.
	//	var alertCRDs []*coralogixv1.Alert
	//	var err error
	//	if alertCRDs, err = prometheusRuleToAlertCRDs(prometheusRuleCRD); err != nil {
	//		log.Error(err, "Received an error while Converting PrometheusRule to RecordingRuleGroupSet", "PrometheusRule Name", prometheusRuleCRD.Name)
	//		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	//	}
	//
	//	ruleGroupSetExist := &coralogixv1.RecordingRuleGroup{}
	//	if err = r.Client.Get(ctx, req.NamespacedName, ruleGroupSetExist); err != nil {
	//		log.V(1).Info("Fetched PrometheusRule. Trying to create/update a RecordingRuleGroupSet accordingly.", "PrometheusRule name", prometheusRuleCRD.Name)
	//		if errors.IsNotFound(err) {
	//			//Meaning there's a PrometheusRule with that NamespacedName but not RecordingRuleGroupSet accordingly (so creating it).
	//			if err = r.Create(ctx, ruleGroupSetCRD); err != nil {
	//				log.Error(err, "Received an error while trying to create RecordingRuleGroupSet CRD", "RecordingRuleGroupSet Name", ruleGroupSetCRD.Name)
	//				return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	//			}
	//			if err = r.Client.Get(ctx, req.NamespacedName, ruleGroupSetExist); err != nil {
	//				log.Error(err, "Received an error while trying to fetch RecordingRuleGroupSet CRD", "RecordingRuleGroupSet Name", ruleGroupSetCRD.Name)
	//				return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	//			}
	//		} else {
	//			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	//		}
	//	}
	//
	//	ruleGroupSetCRD.ResourceVersion = ruleGroupSetExist.ResourceVersion
	//	if err = r.Client.Update(ctx, ruleGroupSetCRD); err != nil {
	//		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	//	}
	//}

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
		h, _ := time.ParseDuration(string(prometheusRule.Spec.Groups[0].Interval))
		intervalSeconds := int32(h.Seconds())

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

func prometheusRuleToAlertCRDs(prometheusRule *prometheus.PrometheusRule) ([]*coralogixv1.Alert, error) {
	if len(prometheusRule.Spec.Groups) > 0 {
		alerts := prometheusInnerRulesToCoralogixAlerts(prometheusRule.Spec.Groups[0].Rules)
		return alerts, nil
	}

	return nil, nil
}

func prometheusInnerRulesToCoralogixAlerts(rules []prometheus.Rule) []*coralogixv1.Alert {
	result := make([]*coralogixv1.Alert, 0)
	for _, rule := range rules {
		if rule.Alert != "" {
			alert := prometheusInnerRuleToCoralogixAlert(rule)
			result = append(result, alert)
		}
	}

	return result
}

func prometheusInnerRuleToCoralogixAlert(prometheusRule prometheus.Rule) *coralogixv1.Alert {
	return &coralogixv1.Alert{
		Spec: coralogixv1.AlertSpec{},
	}
}

func prometheusInnerRulesToCoralogixInnerRules(rules []prometheus.Rule) []coralogixv1.RecordingRule {
	result := make([]coralogixv1.RecordingRule, 0, len(rules))
	for _, r := range rules {
		rule := prometheusInnerRuleToCoralogixInnerRule(r)
		result = append(result, rule)
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
