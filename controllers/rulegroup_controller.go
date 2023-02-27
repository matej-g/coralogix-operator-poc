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

	"coralogix-operator-poc/controllers/clientset"
	rulesgroups "coralogix-operator-poc/controllers/clientset/grpc/rules-groups/v1"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	coralogixv1 "coralogix-operator-poc/api/v1"
	"google.golang.org/grpc/status"
)

// RuleGroupReconciler reconciles a RuleGroup object
type RuleGroupReconciler struct {
	client.Client
	CoralogixClientSet *clientset.ClientSet
	Scheme             *runtime.Scheme
}

//+kubebuilder:rbac:groups=coralogix.coralogix,resources=rulegroups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coralogix.coralogix,resources=rulegroups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=coralogix.coralogix,resources=rulegroups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RuleGroup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *RuleGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	jsm := &jsonpb.Marshaler{
		//Indent: "\t",
	}
	rulesGroupsClient := r.CoralogixClientSet.RuleGroups()

	//Get ruleGroupCRD
	ruleGroupCRD := &coralogixv1.RuleGroup{}

	if err := r.Client.Get(ctx, req.NamespacedName, ruleGroupCRD); err != nil {
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
	if ruleGroupCRD.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(ruleGroupCRD, myFinalizerName) {
			controllerutil.AddFinalizer(ruleGroupCRD, myFinalizerName)
			if err := r.Update(ctx, ruleGroupCRD); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(ruleGroupCRD, myFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			ruleGroupId := *ruleGroupCRD.Status.ID
			deleteRuleGroupReq := &rulesgroups.DeleteRuleGroupRequest{GroupId: ruleGroupId}
			log.V(1).Info("Deleting Rule-Group", "Rule-Group ID", ruleGroupId)
			if _, err := rulesGroupsClient.DeleteRuleGroup(ctx, deleteRuleGroupReq); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				log.Error(err, "Received an error while Deleting a Rule-Group", "Rule-Group ID", ruleGroupId)
				return ctrl.Result{}, err
			}

			log.V(1).Info("Rule-Group was deleted", "Rule-Group ID", ruleGroupId)
			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(ruleGroupCRD, myFinalizerName)
			if err := r.Update(ctx, ruleGroupCRD); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	var notFount bool
	var err error
	var actualState *rulesgroups.GetRuleGroupResponse
	if ruleGroupCRD.Status.ID == nil {
		notFount = true
	} else if actualState, err = rulesGroupsClient.GetRuleGroup(ctx,
		&rulesgroups.GetRuleGroupRequest{GroupId: *ruleGroupCRD.Status.ID}); status.Code(err) == codes.NotFound {
		notFount = true
	}

	if notFount {
		createRuleGroupReq := ruleGroupCRD.Spec.ExtractCreateRuleGroupRequest()
		jstr, _ := jsm.MarshalToString(createRuleGroupReq)
		log.V(1).Info("Creating Rule-Group", "ruleGroup", jstr)
		if createRuleGroupResp, err := rulesGroupsClient.CreateRuleGroup(ctx, createRuleGroupReq); err == nil {
			jstr, _ := jsm.MarshalToString(createRuleGroupResp)
			log.V(1).Info("Rule-Group was updated", "ruleGroup", jstr)
			ruleGroupCRD.Status.ID = new(string)
			*ruleGroupCRD.Status.ID = createRuleGroupResp.GetRuleGroup().GetId().GetValue()
			r.Status().Update(ctx, ruleGroupCRD)
			return ctrl.Result{RequeueAfter: defaultRequeuePeriod}, nil
		} else {
			log.Error(err, "Received an error while creating a Rule-Group", "ruleGroup", createRuleGroupReq)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
	} else if err != nil {
		log.Error(err, "Received an error while reading a Rule-Group", "ruleGroup ID", *ruleGroupCRD.Status.ID)
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	if equal, diff := ruleGroupCRD.Spec.DeepEqual(actualState.RuleGroup); !equal {
		log.V(1).Info("Find diffs between spec and the actual state", "Diff", diff)
		updateRuleGroupReq := ruleGroupCRD.Spec.ExtractUpdateRuleGroupRequest(*ruleGroupCRD.Status.ID)
		updateRuleGroupResp, err := rulesGroupsClient.UpdateRuleGroup(ctx, updateRuleGroupReq)
		if err != nil {
			log.Error(err, "Received an error while updating a Rule-Group", "ruleGroup", updateRuleGroupReq)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
		jstr, _ := jsm.MarshalToString(updateRuleGroupResp)
		log.V(1).Info("Rule-Group was updated", "ruleGroup", jstr)
	}

	return ctrl.Result{RequeueAfter: defaultRequeuePeriod}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RuleGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&coralogixv1.RuleGroup{}).
		Complete(r)
}
