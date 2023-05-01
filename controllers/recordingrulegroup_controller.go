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

	coralogixv1 "coralogix-operator-poc/api/v1"
	"coralogix-operator-poc/controllers/clientset"
	rrg "coralogix-operator-poc/controllers/clientset/grpc/recording-rules-groups/v1"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
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

	//Get rRGCRD
	rRGCRD := &coralogixv1.RecordingRuleGroup{}

	if err := r.Client.Get(ctx, req.NamespacedName, rRGCRD); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	// name of our custom finalizer
	myFinalizerName := "batch.tutorial.kubebuilder.io/finalizer"

	// examine DeletionTimestamp to determine if object is under deletion
	if rRGCRD.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(rRGCRD, myFinalizerName) {
			controllerutil.AddFinalizer(rRGCRD, myFinalizerName)
			if err := r.Update(ctx, rRGCRD); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(rRGCRD, myFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if rRGCRD.Status.Name == nil {
				controllerutil.RemoveFinalizer(rRGCRD, myFinalizerName)
				err := r.Update(ctx, rRGCRD)
				return ctrl.Result{}, err
			}

			rRGName := *rRGCRD.Status.Name
			deleteRRGReq := &rrg.DeleteRuleGroup{Name: rRGName}
			log.V(1).Info("Deleting Rule-Group", "recordingRuleGroup Name", rRGName)
			if _, err := rRGClient.DeleteRecordingRuleGroup(ctx, deleteRRGReq); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				log.Error(err, "Received an error while Deleting a Recording-Rule-Group", "recordingRuleGroup Name", rRGName)
				return ctrl.Result{}, err
			}

			log.V(1).Info("Rule-Group was deleted", "Recording-Rule-Group ID", rRGName)
			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(rRGCRD, myFinalizerName)
			if err := r.Update(ctx, rRGCRD); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	var notFount bool
	var err error
	var actualState *coralogixv1.RecordingRuleGroupStatus
	if name := rRGCRD.Status.Name; name == nil {
		log.V(1).Info("recordingRuleGroup wasn't created")
		notFount = true
	} else if getRuleGroupResp, err := rRGClient.GetRecordingRuleGroup(ctx, &rrg.FetchRuleGroup{Name: *name}); status.Code(err) == codes.NotFound {
		log.V(1).Info("Recording-Rule-Group doesn't exist in Coralogix backend")
		notFount = true
	} else if err == nil {
		actualState = flattenRecordingRuleGroup(getRuleGroupResp.GetRuleGroup())
	}

	if notFount {
		createRuleGroupReq := rRGCRD.Spec.ExtractCreateRecordingRuleGroupRequest()
		jstr, _ := jsm.MarshalToString(createRuleGroupReq)
		log.V(1).Info("Creating Rule-Group", "recordingRuleGroup", jstr)
		if createRRGResp, err := rRGClient.CreateRecordingRuleGroup(ctx, createRuleGroupReq); err == nil {
			jstr, _ := jsm.MarshalToString(createRRGResp)
			log.V(1).Info("Rule-Group was updated", "recordingRuleGroup", jstr)
			getRuleGroupReq := &rrg.FetchRuleGroup{Name: createRuleGroupReq.Name}
			var getRRGResp *rrg.FetchRuleGroupResult
			if getRRGResp, err = rRGClient.GetRecordingRuleGroup(ctx, getRuleGroupReq); err != nil || rRGCRD == nil {
				log.Error(err, "Received an error while getting a Recording-Rule-Group", "recordingRuleGroup", createRuleGroupReq)
				return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
			}
			rRGCRD.Status = *flattenRecordingRuleGroup(getRRGResp.GetRuleGroup())
			if err := r.Status().Update(ctx, rRGCRD); err != nil {
				log.V(1).Error(err, "updating crd")
			}
			return ctrl.Result{RequeueAfter: defaultRequeuePeriod}, nil
		} else {
			log.Error(err, "Received an error while creating a Recording-Rule-Group", "recordingRuleGroup", createRuleGroupReq)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
	} else if err != nil {
		log.Error(err, "Received an error while reading a Recording-Rule-Group", "recordingRuleGroup Name", *rRGCRD.Status.Name)
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	if equal, diff := rRGCRD.Spec.DeepEqual(*actualState); !equal {
		log.V(1).Info("Find diffs between spec and the actual state", "Diff", diff)
		updateRRGReq := rRGCRD.Spec.ExtractCreateRecordingRuleGroupRequest()
		updateRRGResp, err := rRGClient.UpdateRecordingRuleGroup(ctx, updateRRGReq)
		if err != nil {
			log.Error(err, "Received an error while updating a Recording-Rule-Group", "recordingRuleGroup", updateRRGReq)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
		jstr, _ := jsm.MarshalToString(updateRRGResp)
		log.V(1).Info("Rule-Group was updated", "recordingRuleGroup", jstr)
	}

	return ctrl.Result{RequeueAfter: defaultRequeuePeriod}, nil
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
		Complete(r)
}
