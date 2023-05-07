/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"strconv"
	"time"

	coralogixv1 "coralogix-operator-poc/api/v1"
	"coralogix-operator-poc/controllers/clientset"
	rrg "coralogix-operator-poc/controllers/clientset/grpc/recording-rules-groups/v1"
	"github.com/golang/protobuf/jsonpb"
	prometheus "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// RecordingRuleGroupReconciler reconciles a RecordingRuleGroup object
type RecordingRuleGroupReconciler struct {
	client.Client
	CoralogixClientSet *clientset.ClientSet
	Scheme             *runtime.Scheme
}

//+kubebuilder:rbac:groups=coralogix.coralogix,resources=recordingrulegroups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coralogix.coralogix,resources=recordingrulegroups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=coralogix.coralogix,resources=recordingrulegroups/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RecordingRuleGroup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *RecordingRuleGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	jsm := &jsonpb.Marshaler{
		//Indent: "\t",
	}
	rRGClient := r.CoralogixClientSet.RecordingRuleGroups()

	prometheusRuleCRD := &prometheus.PrometheusRule{}
	ruleGroupSetCRD := &coralogixv1.RecordingRuleGroup{}
	if err := r.Get(ctx, req.NamespacedName, prometheusRuleCRD); err != nil && !errors.IsNotFound(err) { //Checking if there's a PrometheusRule to track with that NamespacedName.
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	} else if err == nil && toTrackPrometheusRule(prometheusRuleCRD) { //Meaning there's a PrometheusRule to track with that NamespacedName.
		log.V(1).Info("Fetched PrometheusRule. Trying to create/update a RecordingRuleGroupSet accordingly.", "PrometheusRule name", prometheusRuleCRD.Name)
		//Converting the PrometheusRule to the desired RecordingRuleGroupSet.
		if ruleGroupSetCRD, err = prometheusRuleToRuleGroupSet(prometheusRuleCRD); err != nil {
			log.Error(err, "Received an error while Converting PrometheusRule to RecordingRuleGroupSet", "PrometheusRule Name", prometheusRuleCRD.Name)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}

		ruleGroupSetExist := &coralogixv1.RecordingRuleGroup{}
		if err = r.Client.Get(ctx, req.NamespacedName, ruleGroupSetExist); err != nil {
			if errors.IsNotFound(err) {
				//Meaning there's a PrometheusRule with that NamespacedName but not RecordingRuleGroupSet accordingly (so creating it).
				if err = r.Create(ctx, ruleGroupSetCRD); err != nil {
					log.Error(err, "Received an error while trying to create RecordingRuleGroupSet CRD", "RecordingRuleGroupSet Name", ruleGroupSetCRD.Name)
					return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
				}
				if err = r.Client.Get(ctx, req.NamespacedName, ruleGroupSetExist); err != nil {
					log.Error(err, "Received an error while trying to fetch RecordingRuleGroupSet CRD", "RecordingRuleGroupSet Name", ruleGroupSetCRD.Name)
					return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
				}
			} else {
				return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
			}
		}
		ruleGroupSetCRD.ResourceVersion = ruleGroupSetExist.ResourceVersion
	} else if err = r.Client.Get(ctx, req.NamespacedName, ruleGroupSetCRD); err != nil { //Meaning there isn't PrometheusRule with that NamespacedName (so trying to fetch RecordingRuleGroupSet with that NamespacedName).
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		log.Error(err, "Received an error while trying to create RecordingRuleGroupSet CRD", "RecordingRuleGroupSet Name", ruleGroupSetCRD.Name)
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	// name of our custom finalizer
	myFinalizerName := "batch.tutorial.kubebuilder.io/finalizer"

	// examine DeletionTimestamp to determine if object is under deletion
	if ruleGroupSetCRD.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(ruleGroupSetCRD, myFinalizerName) {
			controllerutil.AddFinalizer(ruleGroupSetCRD, myFinalizerName)
			if err := r.Update(ctx, ruleGroupSetCRD); err != nil {
				log.Error(err, "Received an error while Updating a RecordingRuleGroupSet", "recordingRuleGroup Name", ruleGroupSetCRD.Name)
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(ruleGroupSetCRD, myFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if ruleGroupSetCRD.Status.Name == nil {
				controllerutil.RemoveFinalizer(ruleGroupSetCRD, myFinalizerName)
				err := r.Update(ctx, ruleGroupSetCRD)
				log.Error(err, "Received an error while Updating a RecordingRuleGroupSet", "recordingRuleGroup Name", ruleGroupSetCRD.Name)
				return ctrl.Result{}, err
			}

			rRGName := *ruleGroupSetCRD.Status.Name
			deleteRRGReq := &rrg.DeleteRuleGroup{Name: rRGName}
			log.V(1).Info("Deleting RecordingRuleGroupSet", "recordingRuleGroup Name", rRGName)
			if _, err := rRGClient.DeleteRecordingRuleGroup(ctx, deleteRRGReq); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				log.Error(err, "Received an error while Deleting a RecordingRuleGroupSet", "recordingRuleGroup Name", rRGName)
				return ctrl.Result{}, err
			}

			log.V(1).Info("RecordingRuleGroupSet was deleted", "RecordingRuleGroupSet ID", rRGName)
			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(ruleGroupSetCRD, myFinalizerName)
			if err := r.Update(ctx, ruleGroupSetCRD); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	var notFount bool
	var err error
	var actualState *coralogixv1.RecordingRuleGroupStatus
	if name := ruleGroupSetCRD.Status.Name; name == nil {
		log.V(1).Info("RecordingRuleGroupSet wasn't created in Coralogix backend")
		notFount = true
	} else if getRuleGroupResp, err := rRGClient.GetRecordingRuleGroup(ctx, &rrg.FetchRuleGroup{Name: *name}); status.Code(err) == codes.NotFound {
		log.V(1).Info("RecordingRuleGroupSet doesn't exist in Coralogix backend")
		notFount = true
	} else if err == nil {
		actualState = flattenRecordingRuleGroup(getRuleGroupResp.GetRuleGroup())
	}

	if notFount {
		createRuleGroupReq := ruleGroupSetCRD.Spec.ExtractCreateRecordingRuleGroupRequest()
		jstr, _ := jsm.MarshalToString(createRuleGroupReq)
		log.V(1).Info("Creating RecordingRuleGroupSet", "RecordingRuleGroupSet", jstr)
		if createRRGResp, err := rRGClient.CreateRecordingRuleGroup(ctx, createRuleGroupReq); err == nil {
			jstr, _ := jsm.MarshalToString(createRRGResp)
			log.V(1).Info("RecordingRuleGroupSet was updated", "RecordingRuleGroupSet", jstr)
			getRuleGroupReq := &rrg.FetchRuleGroup{Name: createRuleGroupReq.Name}
			var getRRGResp *rrg.FetchRuleGroupResult
			if getRRGResp, err = rRGClient.GetRecordingRuleGroup(ctx, getRuleGroupReq); err != nil || ruleGroupSetCRD == nil {
				log.Error(err, "Received an error while getting a RecordingRuleGroupSet", "RecordingRuleGroupSet", createRuleGroupReq)
				return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
			}
			ruleGroupSetCRD.Status = *flattenRecordingRuleGroup(getRRGResp.GetRuleGroup())
			if err := r.Status().Update(ctx, ruleGroupSetCRD); err != nil {
				log.V(1).Error(err, "updating crd")
			}
			return ctrl.Result{RequeueAfter: defaultRequeuePeriod}, nil
		} else {
			log.Error(err, "Received an error while creating a RecordingRuleGroupSet", "recordingRuleGroup", createRuleGroupReq)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
	} else if err != nil {
		log.Error(err, "Received an error while reading a RecordingRuleGroupSet", "recordingRuleGroup Name", *ruleGroupSetCRD.Status.Name)
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	if equal, diff := ruleGroupSetCRD.Spec.DeepEqual(*actualState); !equal {
		log.V(1).Info("Find diffs between spec and the actual state", "Diff", diff)
		updateRRGReq := ruleGroupSetCRD.Spec.ExtractCreateRecordingRuleGroupRequest()
		if updateRRGResp, err := rRGClient.UpdateRecordingRuleGroup(ctx, updateRRGReq); err != nil {
			log.Error(err, "Received an error while updating a RecordingRuleGroupSet", "recordingRuleGroup", updateRRGReq)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		} else {
			jstr, _ := jsm.MarshalToString(updateRRGResp)
			log.V(1).Info("RecordingRuleGroupSet was updated on backend", "recordingRuleGroup", jstr)
			var getRuleGroupResp *rrg.FetchRuleGroupResult
			if getRuleGroupResp, err = rRGClient.GetRecordingRuleGroup(ctx, &rrg.FetchRuleGroup{Name: ruleGroupSetCRD.Spec.Name}); err != nil {
				log.Error(err, "Received an error while updating a RecordingRuleGroupSet", "recordingRuleGroup", updateRRGReq)
				return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
			}
			ruleGroupSetCRD.Status = *flattenRecordingRuleGroup(getRuleGroupResp.GetRuleGroup())
			if err := r.Status().Update(ctx, ruleGroupSetCRD); err != nil {
				log.V(1).Error(err, "Error on updating RuleGroupSet crd")
				return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
			}
		}
	}

	return ctrl.Result{RequeueAfter: defaultRequeuePeriod}, nil
}

func toTrackPrometheusRule(rule *prometheus.PrometheusRule) bool {
	if toCreateStr, ok := rule.Labels["cxCreateRecordingRule"]; ok && toCreateStr == "true" {
		return true
	}
	return false
}

func prometheusRuleToRuleGroupSet(prometheusRule *prometheus.PrometheusRule) (*coralogixv1.RecordingRuleGroup, error) {
	if len(prometheusRule.Spec.Groups) > 0 {
		h, _ := time.ParseDuration(string(prometheusRule.Spec.Groups[0].Interval))
		intervalSeconds := int32(h.Seconds())

		limit := int64(60)
		if limitStr, ok := prometheusRule.Annotations["cx_limit"]; ok && limitStr != "" {
			if limitInt, err := strconv.Atoi(limitStr); err != nil {
				return nil, err
			} else {
				limit = int64(limitInt)
			}
		}

		rules := coralogixInnerRulesToPrometheusInnerRules(prometheusRule.Spec.Groups[0].Rules)
		return &coralogixv1.RecordingRuleGroup{
			ObjectMeta: metav1.ObjectMeta{Namespace: prometheusRule.Namespace, Name: prometheusRule.Name},
			Spec: coralogixv1.RecordingRuleGroupSpec{
				Name:            prometheusRule.Name,
				IntervalSeconds: intervalSeconds,
				Limit:           limit,
				Rules:           rules,
			},
		}, nil
	}

	return &coralogixv1.RecordingRuleGroup{
		ObjectMeta: metav1.ObjectMeta{Namespace: prometheusRule.Namespace, Name: prometheusRule.Name},
		Spec: coralogixv1.RecordingRuleGroupSpec{
			Name: prometheusRule.Name,
		},
	}, nil
}

func coralogixInnerRulesToPrometheusInnerRules(rules []prometheus.Rule) []coralogixv1.RecordingRule {
	result := make([]coralogixv1.RecordingRule, 0, len(rules))
	for _, r := range rules {
		rule := coralogixInnerRuleToPrometheusInnerRule(r)
		result = append(result, rule)
	}
	return result
}

func coralogixInnerRuleToPrometheusInnerRule(rule prometheus.Rule) coralogixv1.RecordingRule {
	return coralogixv1.RecordingRule{
		Record: rule.Record,
		Expr:   rule.Expr.StrVal,
		Labels: rule.Labels,
	}
}

func (r *RecordingRuleGroupReconciler) findPrometheusRules(v client.Object) []reconcile.Request {
	prometheusRule := v.(*prometheus.PrometheusRule)

	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: prometheusRule.Namespace,
			Name:      prometheusRule.Name,
		},
	}

	return []reconcile.Request{request}
}

func flattenRecordingRuleGroup(rRG *rrg.RecordingRuleGroup) *coralogixv1.RecordingRuleGroupStatus {
	name := new(string)
	*name = rRG.Name
	interval := int32(*rRG.Interval)
	limit := int64(*rRG.Limit)
	rules := flattenRecordingRules(rRG.Rules)

	return &coralogixv1.RecordingRuleGroupStatus{
		Name:            name,
		IntervalSeconds: interval,
		Limit:           limit,
		Rules:           rules,
	}
}

func flattenRecordingRules(rules []*rrg.RecordingRule) []coralogixv1.RecordingRule {
	result := make([]coralogixv1.RecordingRule, 0, len(rules))
	for _, r := range rules {
		rule := flattenRecordingRule(r)
		result = append(result, rule)
	}
	return result
}

func flattenRecordingRule(rule *rrg.RecordingRule) coralogixv1.RecordingRule {
	return coralogixv1.RecordingRule{
		Record: rule.Record,
		Expr:   rule.Expr,
		Labels: rule.Labels,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *RecordingRuleGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&coralogixv1.RecordingRuleGroup{}).
		Watches(&source.Kind{Type: &prometheus.PrometheusRule{}},
			handler.EnqueueRequestsFromMapFunc(r.findPrometheusRules),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		Complete(r)
}
