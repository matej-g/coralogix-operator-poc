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
	"fmt"
	"time"

	"coralogix-operator-poc/controllers/clientset"
	rulesgroups "coralogix-operator-poc/controllers/clientset/grpc/rules-groups/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	coralogixv1 "coralogix-operator-poc/api/v1"
)

const (
	defaultRequeuePeriod    = 60 * time.Second
	defaultErrRequeuePeriod = 5 * time.Second
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

	//Get instance
	instance := &coralogixv1.RuleGroup{}
	var result ctrl.Result
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return result, nil
		}
		// Error reading the object - requeue the request
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	//If ID is nil, create the rule-group
	if instance.Status.ID == nil {
		log.Info(fmt.Sprintf("attempt to create rule-group with spec: %v", instance.Spec.ToString()))
		createRuleGroupReq := instance.Spec.ExtractCreateRuleGroupRequest()
		createRuleGroupResp, err := r.CoralogixClientSet.RuleGroups().CreateRuleGroup(ctx, createRuleGroupReq)
		if err != nil {
			log.Error(err, fmt.Sprintf("received an error while creating a rule-group: %s", err.Error()))
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
		id := createRuleGroupResp.GetRuleGroup().GetId().GetValue()
		log.Info(fmt.Sprintf("rule-group %s was created:", id))
		instance.Status.ID = new(string)
		*instance.Status.ID = id
	}

	getRuleGroupReq := &rulesgroups.GetRuleGroupRequest{
		GroupId: *instance.Status.ID,
	}
	actualState, err := r.CoralogixClientSet.RuleGroups().GetRuleGroup(ctx, getRuleGroupReq)
	if err != nil {
		log.Error(err, fmt.Sprintf("received an error while reading a rule-group: %s", err.Error()))
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}
	log.Info(fmt.Sprintf("received a rule-group: %s", getRuleGroupReq.String()))

	if !instance.Spec.DeepEqual(actualState.RuleGroup) {
		log.Info(fmt.Sprintf("find diffs betwen spec and actual state. attempt to update rule-group"))
		updateRuleGroupReq := instance.Spec.ExtractUpdateRuleGroupRequest(*instance.Status.ID)
		updateRuleGroupResp, err := r.CoralogixClientSet.RuleGroups().UpdateRuleGroup(ctx, updateRuleGroupReq)
		if err != nil {
			log.Error(err, fmt.Sprintf("received an error while updating a rule-group: %s", err.Error()))
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
		log.Info(fmt.Sprintf("rule-group was updated: %s", updateRuleGroupResp.String()))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RuleGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&coralogixv1.RuleGroup{}).
		Complete(r)
}
