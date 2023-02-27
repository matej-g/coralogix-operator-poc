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
	alerts "coralogix-operator-poc/controllers/clientset/grpc/alerts/v1"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	coralogixv1 "coralogix-operator-poc/api/v1"
)

// AlertReconciler reconciles a Alert object
type AlertReconciler struct {
	client.Client
	CoralogixClientSet *clientset.ClientSet
	Scheme             *runtime.Scheme
}

//+kubebuilder:rbac:groups=coralogix.coralogix,resources=alerts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coralogix.coralogix,resources=alerts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=coralogix.coralogix,resources=alerts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Alert object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *AlertReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log := log.FromContext(ctx)
	jsm := &jsonpb.Marshaler{
		//Indent: "\t",
	}
	alertsClient := r.CoralogixClientSet.Alerts()

	//Get alertCRD
	alertCRD := &coralogixv1.Alert{}

	if err := r.Client.Get(ctx, req.NamespacedName, alertCRD); err != nil {
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
	if alertCRD.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(alertCRD, myFinalizerName) {
			controllerutil.AddFinalizer(alertCRD, myFinalizerName)
			if err := r.Update(ctx, alertCRD); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(alertCRD, myFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if alertCRD.Status.ID == nil {
				controllerutil.RemoveFinalizer(alertCRD, myFinalizerName)
				err := r.Update(ctx, alertCRD)
				return ctrl.Result{}, err
			}

			alertId := *alertCRD.Status.ID
			deleteAlertReq := &alerts.DeleteAlertByUniqueIdRequest{Id: wrapperspb.String(alertId)}
			log.V(1).Info("Deleting Alert", "Alert ID", alertId)
			if _, err := alertsClient.DeleteAlert(ctx, deleteAlertReq); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				log.Error(err, "Received an error while Deleting a Alert", "Alert ID", alertId)
				return ctrl.Result{}, err
			}

			log.V(1).Info("Alert was deleted", "Alert ID", alertId)
			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(alertCRD, myFinalizerName)
			if err := r.Update(ctx, alertCRD); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	var notFount bool
	var err error
	var actualState *alerts.GetAlertByUniqueIdResponse
	if alertCRD.Status.ID == nil {
		notFount = true
	} else if actualState, err = alertsClient.GetAlert(ctx,
		&alerts.GetAlertByUniqueIdRequest{Id: wrapperspb.String(*alertCRD.Status.ID)}); status.Code(err) == codes.NotFound {
		notFount = true
	}

	if notFount {
		createAlertReq := alertCRD.Spec.ExtractCreateAlertRequest()
		jstr, _ := jsm.MarshalToString(createAlertReq)
		log.V(1).Info("Creating Alert", "alert", jstr)
		if createAlertResp, err := alertsClient.CreateAlert(ctx, createAlertReq); err == nil {
			jstr, _ := jsm.MarshalToString(createAlertResp)
			log.V(1).Info("Alert was updated", "alert", jstr)
			alertCRD.Status.ID = new(string)
			*alertCRD.Status.ID = createAlertResp.GetAlert().GetId().GetValue()
			r.Status().Update(ctx, alertCRD)
			return ctrl.Result{RequeueAfter: defaultRequeuePeriod}, nil
		} else {
			log.Error(err, "Received an error while creating a Alert", "alert", createAlertResp)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
	} else if err != nil {
		log.Error(err, "Received an error while reading a Alert", "alert ID", *alertCRD.Status.ID)
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	if equal, diff := alertCRD.Spec.DeepEqual(actualState.GetAlert()); !equal {
		log.V(1).Info("Find diffs between spec and the actual state", "Diff", diff)
		updateAlertReq := alertCRD.Spec.ExtractUpdateAlertRequest(*alertCRD.Status.ID)
		updateAlertResp, err := alertsClient.UpdateAlert(ctx, updateAlertReq)
		if err != nil {
			log.Error(err, "Received an error while updating a Alert", "alert", updateAlertReq)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
		jstr, _ := jsm.MarshalToString(updateAlertResp)
		log.V(1).Info("Alert was updated", "alert", jstr)
	}

	return ctrl.Result{RequeueAfter: defaultRequeuePeriod}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlertReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&coralogixv1.Alert{}).
		Complete(r)
}
