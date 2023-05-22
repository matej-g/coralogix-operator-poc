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

package alphacontrollers

import (
	"context"
	"fmt"
	"time"

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

	utils "coralogix-operator-poc/apis"
	coralogixv1alpha1 "coralogix-operator-poc/apis/coralogix/v1alpha1"
	"coralogix-operator-poc/controllers/clientset"
	alerts "coralogix-operator-poc/controllers/clientset/grpc/alerts/v1"
)

var (
	alertProtoSeverityToSchemaSeverity                               = utils.ReverseMap(coralogixv1alpha1.AlertSchemaSeverityToProtoSeverity)
	alertProtoDayToSchemaDay                                         = utils.ReverseMap(coralogixv1alpha1.AlertSchemaDayToProtoDay)
	alertProtoTimeWindowToSchemaTimeWindow                           = utils.ReverseMap(coralogixv1alpha1.AlertSchemaTimeWindowToProtoTimeWindow)
	alertProtoAutoRetireRatioToSchemaAutoRetireRatio                 = utils.ReverseMap(coralogixv1alpha1.AlertSchemaAutoRetireRatioToProtoAutoRetireRatio)
	alertProtoFiltersLogSeverityToSchemaFiltersLogSeverity           = utils.ReverseMap(coralogixv1alpha1.AlertSchemaFiltersLogSeverityToProtoFiltersLogSeverity)
	alertProtoRelativeTimeFrameToSchemaTimeFrameAndRelativeTimeFrame = utils.ReverseMap(coralogixv1alpha1.AlertSchemaRelativeTimeFrameToProtoTimeFrameAndRelativeTimeFrame)
	alertProtoArithmeticOperatorToSchemaArithmeticOperator           = utils.ReverseMap(coralogixv1alpha1.AlertSchemaArithmeticOperatorToProtoArithmeticOperator)
	alertProtoTracingFilterFieldToSchemaTracingFilterField           = utils.ReverseMap(coralogixv1alpha1.AlertSchemaTracingFilterFieldToProtoTracingFilterField)
	alertProtoTracingOperatorToSchemaTracingOperator                 = utils.ReverseMap(coralogixv1alpha1.AlertSchemaTracingOperatorToProtoTracingOperator)
	alertProtoFlowOperatorToProtoFlowOperator                        = utils.ReverseMap(coralogixv1alpha1.AlertSchemaFlowOperatorToProtoFlowOperator)
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
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *AlertReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	jsm := &jsonpb.Marshaler{
		//Indent: "\t",
	}
	alertsClient := r.CoralogixClientSet.Alerts()

	//Get alertCRD
	alertCRD := &coralogixv1alpha1.Alert{}
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
				if status.Code(err) == codes.NotFound {
					controllerutil.RemoveFinalizer(alertCRD, myFinalizerName)
					err := r.Update(ctx, alertCRD)
					return ctrl.Result{}, err
				}

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
	var actualState *coralogixv1alpha1.AlertStatus

	if id := alertCRD.Status.ID; id == nil {
		log.V(1).Info("alert wasn't created")
		notFount = true
	} else if getAlertResp, err := alertsClient.GetAlert(ctx, &alerts.GetAlertByUniqueIdRequest{Id: wrapperspb.String(*id)}); status.Code(err) == codes.NotFound {
		log.V(1).Info("alert doesn't exist in Coralogix backend")
		notFount = true
	} else if err == nil {
		actualState = flattenAlert(getAlertResp.GetAlert(), alertCRD.Spec)
	}

	if notFount {
		createAlertReq := alertCRD.Spec.ExtractCreateAlertRequest()
		jstr, _ := jsm.MarshalToString(createAlertReq)
		log.V(1).Info("Creating Alert", "alert", jstr)
		if createAlertResp, err := alertsClient.CreateAlert(ctx, createAlertReq); err == nil {
			jstr, _ = jsm.MarshalToString(createAlertResp)
			log.V(1).Info("Alert was created", "alert", jstr)
			actualState = flattenAlert(createAlertResp.GetAlert(), alertCRD.Spec)
			alertCRD.Status = *actualState
			if err := r.Status().Update(ctx, alertCRD); err != nil {
				log.V(1).Error(err, "updating crd")
			}
			return ctrl.Result{RequeueAfter: defaultRequeuePeriod}, nil
		} else {
			log.Error(err, "Received an error while creating Alert", "alert", createAlertResp)
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
	} else if err != nil {
		log.Error(err, "Received an error while reading Alert", "alert ID", *alertCRD.Status.ID)
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	if equal, diff := alertCRD.Spec.DeepEqual(actualState); !equal {
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

func flattenAlert(actualAlert *alerts.Alert, spec coralogixv1alpha1.AlertSpec) *coralogixv1alpha1.AlertStatus {
	var status coralogixv1alpha1.AlertStatus

	status.ID = new(string)
	*status.ID = actualAlert.GetUniqueIdentifier().GetValue()

	status.Name = actualAlert.GetName().GetValue()

	status.Description = actualAlert.GetDescription().GetValue()

	status.Active = actualAlert.GetIsActive().GetValue()

	status.Severity = alertProtoSeverityToSchemaSeverity[actualAlert.GetSeverity()]

	status.Labels = flattenMetaLabels(actualAlert.GetMetaLabels())

	status.ExpirationDate = flattenExpirationDate(actualAlert.GetExpiration())

	status.Scheduling = flattenScheduling(actualAlert.GetActiveWhen(), spec)

	alertType, notifyData := flattenAlertType(actualAlert)
	status.AlertType = alertType

	status.Notifications = flattenNotifications(actualAlert.GetNotifications(), actualAlert.GetNotifyEvery(), notifyData)

	return &status
}

type NotificationsAlertTypeData struct {
	OnTriggerAndResolved *wrapperspb.BoolValue

	NotifyOnlyOnTriggeredGroupByValues *wrapperspb.BoolValue
}

func flattenAlertType(actualAlert *alerts.Alert) (coralogixv1alpha1.AlertType, *NotificationsAlertTypeData) {
	actualFilters := actualAlert.GetFilters()
	actualCondition := actualAlert.GetCondition()

	var alertType coralogixv1alpha1.AlertType
	var notifyData = new(NotificationsAlertTypeData)
	switch actualFilters.GetFilterType() {
	case alerts.AlertFilters_FILTER_TYPE_TEXT_OR_UNSPECIFIED:
		if newValueCondition, ok := actualCondition.GetCondition().(*alerts.AlertCondition_NewValue); ok {
			alertType.NewValue, notifyData = flattenNewValueAlert(actualFilters, newValueCondition)
		} else {
			alertType.Standard, notifyData = flattenStandardAlert(actualFilters, actualCondition)
		}
	case alerts.AlertFilters_FILTER_TYPE_RATIO:
		alertType.Ratio, notifyData = flattenRatioAlert(actualFilters, actualCondition)
	case alerts.AlertFilters_FILTER_TYPE_UNIQUE_COUNT:
		alertType.UniqueCount, notifyData = flattenUniqueCountAlert(actualFilters, actualCondition)
	case alerts.AlertFilters_FILTER_TYPE_TIME_RELATIVE:
		alertType.TimeRelative, notifyData = flattenTimeRelativeAlert(actualFilters, actualCondition)
	case alerts.AlertFilters_FILTER_TYPE_METRIC:
		alertType.Metric, notifyData = flattenMetricAlert(actualFilters, actualCondition)
	case alerts.AlertFilters_FILTER_TYPE_TRACING:
		alertType.Tracing = flattenTracingAlert(actualAlert.GetTracingAlert(), actualCondition)
	case alerts.AlertFilters_FILTER_TYPE_FLOW:
		alertType.Flow = flattenFlowAlert(actualCondition.GetFlow())
	}

	return alertType, notifyData
}

func flattenNewValueAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition_NewValue) (*coralogixv1alpha1.NewValue, *NotificationsAlertTypeData) {
	flattenedFilters := flattenFilters(filters)
	newValueCondition, notifyData := flattenNewValueCondition(condition.NewValue.GetParameters())

	newValue := &coralogixv1alpha1.NewValue{
		Filters:    flattenedFilters,
		Conditions: newValueCondition,
	}

	return newValue, notifyData
}

func flattenFilters(filters *alerts.AlertFilters) *coralogixv1alpha1.Filters {
	if filters == nil {
		return nil
	}

	var flattenedFilters = &coralogixv1alpha1.Filters{}

	if actualSearchQuery := filters.GetText(); actualSearchQuery != nil {
		flattenedFilters.SearchQuery = new(string)
		*flattenedFilters.SearchQuery = actualSearchQuery.GetValue()
	}

	if actualAlias := filters.GetAlias(); actualAlias == nil {
		flattenedFilters.Alias = new(string)
		*flattenedFilters.Alias = actualAlias.GetValue()
	}

	flattenedFilters.Severities = flattenSeverities(filters.GetSeverities())

	if metaData := filters.Metadata; metaData != nil {
		flattenedFilters.Subsystems = utils.WrappedStringSliceToStringSlice(metaData.Subsystems)
		flattenedFilters.Categories = utils.WrappedStringSliceToStringSlice(metaData.Categories)
		flattenedFilters.Applications = utils.WrappedStringSliceToStringSlice(metaData.Applications)
		flattenedFilters.Computers = utils.WrappedStringSliceToStringSlice(metaData.Computers)
		flattenedFilters.Classes = utils.WrappedStringSliceToStringSlice(metaData.Classes)
		flattenedFilters.Methods = utils.WrappedStringSliceToStringSlice(metaData.Methods)
		flattenedFilters.IPs = utils.WrappedStringSliceToStringSlice(metaData.IpAddresses)
	}

	return flattenedFilters
}

func flattenSeverities(severities []alerts.AlertFilters_LogSeverity) []coralogixv1alpha1.FiltersLogSeverity {
	flattenedSeverities := make([]coralogixv1alpha1.FiltersLogSeverity, 0, len(severities))
	for _, severity := range severities {
		sev := alertProtoFiltersLogSeverityToSchemaFiltersLogSeverity[severity]
		flattenedSeverities = append(flattenedSeverities, sev)
	}
	return flattenedSeverities
}

func flattenNewValueCondition(conditionParams *alerts.ConditionParameters) (coralogixv1alpha1.NewValueConditions, *NotificationsAlertTypeData) {
	var key string
	if actualKeys := conditionParams.GetGroupBy(); len(actualKeys) != 0 {
		key = actualKeys[0].GetValue()
	}
	timeWindow := coralogixv1alpha1.NewValueTimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])

	newValueCondition := coralogixv1alpha1.NewValueConditions{
		Key:        key,
		TimeWindow: timeWindow,
	}

	notifyData := &NotificationsAlertTypeData{
		NotifyOnlyOnTriggeredGroupByValues: conditionParams.GetNotifyGroupByOnlyAlerts(),
		OnTriggerAndResolved:               conditionParams.GetNotifyOnResolved(),
	}

	return newValueCondition, notifyData
}

func flattenStandardAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition) (*coralogixv1alpha1.Standard, *NotificationsAlertTypeData) {
	flattenedFilters := flattenFilters(filters)
	standardCondition, notifyData := flattenStandardCondition(condition)

	standard := &coralogixv1alpha1.Standard{
		Filters:    flattenedFilters,
		Conditions: standardCondition,
	}

	return standard, notifyData
}

func flattenStandardCondition(condition *alerts.AlertCondition) (coralogixv1alpha1.StandardConditions, *NotificationsAlertTypeData) {
	var standardCondition coralogixv1alpha1.StandardConditions
	var conditionParams *alerts.ConditionParameters

	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		conditionParams = condition.LessThan.GetParameters()
		standardCondition.AlertWhen = coralogixv1alpha1.StandardAlertWhenLessThan
		*standardCondition.Threshold = int(conditionParams.GetThreshold().GetValue())
		*standardCondition.TimeWindow = coralogixv1alpha1.TimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])

		if actualManageUndetectedValues := conditionParams.GetRelatedExtendedData(); actualManageUndetectedValues != nil {
			actualShouldTriggerDeadman, actualCleanupDeadmanDuration := actualManageUndetectedValues.GetShouldTriggerDeadman().GetValue(), actualManageUndetectedValues.GetCleanupDeadmanDuration()
			autoRetireRatio := alertProtoAutoRetireRatioToSchemaAutoRetireRatio[actualCleanupDeadmanDuration]
			standardCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: actualShouldTriggerDeadman,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		} else {
			autoRetireRatio := coralogixv1alpha1.AutoRetireRatioNever
			standardCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: true,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		}
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.MoreThan.GetParameters()
		standardCondition.AlertWhen = coralogixv1alpha1.StandardAlertWhenMoreThan
		standardCondition.Threshold = new(int)
		*standardCondition.Threshold = int(conditionParams.GetThreshold().GetValue())
		standardCondition.TimeWindow = new(coralogixv1alpha1.TimeWindow)
		*standardCondition.TimeWindow = coralogixv1alpha1.TimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])
	case *alerts.AlertCondition_MoreThanUsual:
		conditionParams = condition.MoreThanUsual.GetParameters()
		standardCondition.AlertWhen = coralogixv1alpha1.StandardAlertWhenMoreThanUsual
		*standardCondition.Threshold = int(conditionParams.GetThreshold().GetValue())
	case *alerts.AlertCondition_Immediate:
		standardCondition.AlertWhen = coralogixv1alpha1.StandardAlertWhenImmediately
		return standardCondition, new(NotificationsAlertTypeData)
	}

	standardCondition.GroupBy = utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy())
	notifyData := new(NotificationsAlertTypeData)
	notifyData.NotifyOnlyOnTriggeredGroupByValues = conditionParams.NotifyGroupByOnlyAlerts
	notifyData.OnTriggerAndResolved = conditionParams.NotifyOnResolved

	return standardCondition, notifyData
}

func flattenRatioAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition) (*coralogixv1alpha1.Ratio, *NotificationsAlertTypeData) {
	query1Filters := flattenFilters(filters)
	q2Filters := filters.GetRatioAlerts()[0]
	query2Filters := flattenRatioFilters(q2Filters)
	ratioCondition, notifyData := flattenRatioCondition(condition, q2Filters.GetGroupBy())

	ratio := &coralogixv1alpha1.Ratio{
		Query1Filters: *query1Filters,
		Query2Filters: query2Filters,
		Conditions:    ratioCondition,
	}

	return ratio, notifyData
}

func flattenRatioFilters(filters *alerts.AlertFilters_RatioAlert) coralogixv1alpha1.RatioQ2Filters {
	var flattenedFilters coralogixv1alpha1.RatioQ2Filters
	if filters == nil {
		return flattenedFilters
	}

	if actualSearchQuery := filters.GetText(); actualSearchQuery != nil {
		flattenedFilters.SearchQuery = new(string)
		*flattenedFilters.SearchQuery = actualSearchQuery.GetValue()
	}

	if actualAlias := filters.GetAlias(); actualAlias == nil {
		*flattenedFilters.Alias = actualAlias.GetValue()
	}

	flattenedFilters.Severities = flattenSeverities(filters.GetSeverities())
	flattenedFilters.Subsystems = utils.WrappedStringSliceToStringSlice(filters.GetSubsystems())
	flattenedFilters.Applications = utils.WrappedStringSliceToStringSlice(filters.GetApplications())

	return flattenedFilters
}

func flattenRatioCondition(condition *alerts.AlertCondition, groupByQ2 []*wrapperspb.StringValue) (coralogixv1alpha1.RatioConditions, *NotificationsAlertTypeData) {
	var ratioCondition coralogixv1alpha1.RatioConditions
	var conditionParams *alerts.ConditionParameters

	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		conditionParams = condition.LessThan.GetParameters()
		ratioCondition.AlertWhen = coralogixv1alpha1.AlertWhenLessThan

		if actualManageUndetectedValues := conditionParams.GetRelatedExtendedData(); actualManageUndetectedValues != nil {
			actualShouldTriggerDeadman, actualCleanupDeadmanDuration := actualManageUndetectedValues.GetShouldTriggerDeadman().GetValue(), actualManageUndetectedValues.GetCleanupDeadmanDuration()
			autoRetireRatio := alertProtoAutoRetireRatioToSchemaAutoRetireRatio[actualCleanupDeadmanDuration]
			ratioCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: actualShouldTriggerDeadman,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		} else {
			autoRetireRatio := coralogixv1alpha1.AutoRetireRatioNever
			ratioCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: true,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		}
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.MoreThan.GetParameters()
		ratioCondition.AlertWhen = coralogixv1alpha1.AlertWhenMoreThan
	}

	ratioCondition.Ratio = utils.FloatToQuantity(conditionParams.GetThreshold().GetValue())
	ratioCondition.TimeWindow = coralogixv1alpha1.TimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])

	if groupByQ1 := conditionParams.GetGroupBy(); len(groupByQ1) > 0 && len(groupByQ2) == 0 {
		ratioCondition.GroupBy = utils.WrappedStringSliceToStringSlice(groupByQ1)
		ratioCondition.GroupByFor = new(coralogixv1alpha1.GroupByFor)
		*ratioCondition.GroupByFor = coralogixv1alpha1.GroupByForQ1
	} else if len(groupByQ2) > 0 && len(groupByQ1) == 0 {
		ratioCondition.GroupBy = utils.WrappedStringSliceToStringSlice(groupByQ2)
		ratioCondition.GroupByFor = new(coralogixv1alpha1.GroupByFor)
		*ratioCondition.GroupByFor = coralogixv1alpha1.GroupByForQ2
	} else if len(groupByQ1) > 0 && len(groupByQ2) > 0 {
		ratioCondition.GroupBy = utils.WrappedStringSliceToStringSlice(groupByQ2)
		ratioCondition.GroupByFor = new(coralogixv1alpha1.GroupByFor)
		*ratioCondition.GroupByFor = coralogixv1alpha1.GroupByForBoth
	}

	notifyData := new(NotificationsAlertTypeData)
	notifyData.NotifyOnlyOnTriggeredGroupByValues = conditionParams.NotifyGroupByOnlyAlerts
	notifyData.OnTriggerAndResolved = conditionParams.NotifyOnResolved

	return ratioCondition, notifyData
}

func flattenUniqueCountAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition) (*coralogixv1alpha1.UniqueCount, *NotificationsAlertTypeData) {
	flattenedFilters := flattenFilters(filters)
	uniqueCountCondition, notifyData := flattenUniqueCountCondition(condition)

	ratio := &coralogixv1alpha1.UniqueCount{
		Filters:    flattenedFilters,
		Conditions: uniqueCountCondition,
	}

	return ratio, notifyData
}

func flattenUniqueCountCondition(condition *alerts.AlertCondition) (coralogixv1alpha1.UniqueCountConditions, *NotificationsAlertTypeData) {
	conditionParams := condition.GetCondition().(*alerts.AlertCondition_UniqueCount).UniqueCount.GetParameters()
	var uniqueCountCondition coralogixv1alpha1.UniqueCountConditions

	uniqueCountCondition.Key = conditionParams.GetCardinalityFields()[0].GetValue()
	uniqueCountCondition.MaxUniqueValues = int(conditionParams.GetThreshold().GetValue())
	uniqueCountCondition.TimeWindow = coralogixv1alpha1.UniqueValueTimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])
	if actualGroupBy := conditionParams.GetGroupBy(); len(actualGroupBy) > 0 {
		uniqueCountCondition.GroupBy = new(string)
		*uniqueCountCondition.GroupBy = actualGroupBy[0].GetValue()

		uniqueCountCondition.MaxUniqueValuesForGroupBy = new(int)
		*uniqueCountCondition.MaxUniqueValuesForGroupBy = int(conditionParams.GetMaxUniqueCountValuesForGroupByKey().GetValue())
	}

	notifyData := new(NotificationsAlertTypeData)
	notifyData.NotifyOnlyOnTriggeredGroupByValues = conditionParams.NotifyGroupByOnlyAlerts
	notifyData.OnTriggerAndResolved = conditionParams.NotifyOnResolved

	return uniqueCountCondition, notifyData
}

func flattenTimeRelativeAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition) (*coralogixv1alpha1.TimeRelative, *NotificationsAlertTypeData) {
	flattenedFilters := flattenFilters(filters)
	timeRelativeCondition, notifyData := flattenTimeRelativeCondition(condition)

	timeRelative := &coralogixv1alpha1.TimeRelative{
		Filters:    flattenedFilters,
		Conditions: timeRelativeCondition,
	}

	return timeRelative, notifyData
}

func flattenTimeRelativeCondition(condition *alerts.AlertCondition) (coralogixv1alpha1.TimeRelativeConditions, *NotificationsAlertTypeData) {
	var conditionParams *alerts.ConditionParameters
	var timeRelativeCondition coralogixv1alpha1.TimeRelativeConditions

	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		conditionParams = condition.LessThan.GetParameters()
		timeRelativeCondition.AlertWhen = coralogixv1alpha1.AlertWhenLessThan

		if actualManageUndetectedValues := conditionParams.GetRelatedExtendedData(); actualManageUndetectedValues != nil {
			actualShouldTriggerDeadman, actualCleanupDeadmanDuration := actualManageUndetectedValues.GetShouldTriggerDeadman().GetValue(), actualManageUndetectedValues.GetCleanupDeadmanDuration()
			autoRetireRatio := alertProtoAutoRetireRatioToSchemaAutoRetireRatio[actualCleanupDeadmanDuration]
			timeRelativeCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: actualShouldTriggerDeadman,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		} else {
			autoRetireRatio := coralogixv1alpha1.AutoRetireRatioNever
			timeRelativeCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: true,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		}
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.MoreThan.GetParameters()
		timeRelativeCondition.AlertWhen = coralogixv1alpha1.AlertWhenMoreThan
	}

	timeRelativeCondition.Threshold = utils.FloatToQuantity(conditionParams.GetThreshold().GetValue())
	relativeTimeFrame := coralogixv1alpha1.ProtoTimeFrameAndRelativeTimeFrame{TimeFrame: conditionParams.GetTimeframe(), RelativeTimeFrame: conditionParams.GetRelativeTimeframe()}
	timeRelativeCondition.TimeWindow = alertProtoRelativeTimeFrameToSchemaTimeFrameAndRelativeTimeFrame[relativeTimeFrame]
	timeRelativeCondition.IgnoreInfinity = conditionParams.GetIgnoreInfinity().GetValue()
	timeRelativeCondition.GroupBy = utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy())

	notifyData := new(NotificationsAlertTypeData)
	notifyData.NotifyOnlyOnTriggeredGroupByValues = conditionParams.GetNotifyGroupByOnlyAlerts()
	notifyData.OnTriggerAndResolved = conditionParams.GetNotifyOnResolved()

	return timeRelativeCondition, notifyData
}

func flattenMetricAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition) (*coralogixv1alpha1.Metric, *NotificationsAlertTypeData) {
	metric := new(coralogixv1alpha1.Metric)
	notifyData := new(NotificationsAlertTypeData)

	var conditionParams *alerts.ConditionParameters
	var alertWhen coralogixv1alpha1.AlertWhen
	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		alertWhen = coralogixv1alpha1.AlertWhenLessThan
		conditionParams = condition.LessThan.GetParameters()
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.MoreThan.GetParameters()
		alertWhen = coralogixv1alpha1.AlertWhenMoreThan
	}

	if promqlParams := conditionParams.GetMetricAlertPromqlParameters(); promqlParams != nil {
		metric.Promql, notifyData = flattenPromqlAlert(conditionParams, promqlParams, alertWhen)
	} else {
		metric.Lucene, notifyData = flattenLuceneAlert(conditionParams, filters.GetText(), alertWhen)
	}

	return metric, notifyData
}

func flattenPromqlAlert(conditionParams *alerts.ConditionParameters, promqlParams *alerts.MetricAlertPromqlConditionParameters, alertWhen coralogixv1alpha1.AlertWhen) (*coralogixv1alpha1.Promql, *NotificationsAlertTypeData) {
	promql := new(coralogixv1alpha1.Promql)

	promql.SearchQuery = promqlParams.GetPromqlText().GetValue()
	promql.Conditions = coralogixv1alpha1.PromqlConditions{
		AlertWhen:                   alertWhen,
		Threshold:                   utils.FloatToQuantity(conditionParams.GetThreshold().GetValue()),
		SampleThresholdPercentage:   int(promqlParams.GetSampleThresholdPercentage().GetValue()),
		TimeWindow:                  coralogixv1alpha1.MetricTimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()]),
		GroupBy:                     utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy()),
		ReplaceMissingValueWithZero: promqlParams.GetSwapNullValues().GetValue(),
	}

	if minNonNullValuesPercentage := promqlParams.GetNonNullPercentage(); minNonNullValuesPercentage != nil {
		promql.Conditions.MinNonNullValuesPercentage = new(int)
		*promql.Conditions.MinNonNullValuesPercentage = int(minNonNullValuesPercentage.GetValue())
	}

	if alertWhen == coralogixv1alpha1.AlertWhenLessThan {
		if actualManageUndetectedValues := conditionParams.GetRelatedExtendedData(); actualManageUndetectedValues != nil {
			actualShouldTriggerDeadman, actualCleanupDeadmanDuration := actualManageUndetectedValues.GetShouldTriggerDeadman().GetValue(), actualManageUndetectedValues.GetCleanupDeadmanDuration()
			autoRetireRatio := alertProtoAutoRetireRatioToSchemaAutoRetireRatio[actualCleanupDeadmanDuration]
			promql.Conditions.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: actualShouldTriggerDeadman,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		} else {
			autoRetireRatio := coralogixv1alpha1.AutoRetireRatioNever
			promql.Conditions.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: true,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		}
	}

	notifyData := new(NotificationsAlertTypeData)
	notifyData.NotifyOnlyOnTriggeredGroupByValues = conditionParams.NotifyGroupByOnlyAlerts
	notifyData.OnTriggerAndResolved = conditionParams.NotifyOnResolved

	return promql, notifyData
}

func flattenLuceneAlert(conditionParams *alerts.ConditionParameters, searchQuery *wrapperspb.StringValue, alertWhen coralogixv1alpha1.AlertWhen) (*coralogixv1alpha1.Lucene, *NotificationsAlertTypeData) {
	lucene := new(coralogixv1alpha1.Lucene)
	metricParams := conditionParams.GetMetricAlertParameters()

	if searchQuery != nil {
		lucene.SearchQuery = new(string)
		*lucene.SearchQuery = searchQuery.GetValue()
	}

	lucene.Conditions = coralogixv1alpha1.LuceneConditions{
		MetricField:                 metricParams.GetMetricField().GetValue(),
		ArithmeticOperator:          alertProtoArithmeticOperatorToSchemaArithmeticOperator[metricParams.GetArithmeticOperator()],
		AlertWhen:                   alertWhen,
		Threshold:                   utils.FloatToQuantity(conditionParams.GetThreshold().GetValue()),
		SampleThresholdPercentage:   int(metricParams.GetSampleThresholdPercentage().GetValue()),
		TimeWindow:                  coralogixv1alpha1.MetricTimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()]),
		GroupBy:                     utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy()),
		ReplaceMissingValueWithZero: metricParams.GetSwapNullValues().GetValue(),
		MinNonNullValuesPercentage:  int(metricParams.GetNonNullPercentage().GetValue()),
	}

	if arithmeticOperatorModifier := metricParams.GetArithmeticOperatorModifier(); arithmeticOperatorModifier != nil {
		lucene.Conditions.ArithmeticOperatorModifier = new(int)
		*lucene.Conditions.ArithmeticOperatorModifier = int(arithmeticOperatorModifier.GetValue())
	}

	if alertWhen == coralogixv1alpha1.AlertWhenLessThan {
		if actualManageUndetectedValues := conditionParams.GetRelatedExtendedData(); actualManageUndetectedValues != nil {
			actualShouldTriggerDeadman, actualCleanupDeadmanDuration := actualManageUndetectedValues.GetShouldTriggerDeadman().GetValue(), actualManageUndetectedValues.GetCleanupDeadmanDuration()
			autoRetireRatio := alertProtoAutoRetireRatioToSchemaAutoRetireRatio[actualCleanupDeadmanDuration]
			lucene.Conditions.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: actualShouldTriggerDeadman,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		} else {
			autoRetireRatio := coralogixv1alpha1.AutoRetireRatioNever
			lucene.Conditions.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: true,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		}
	}

	notifyData := new(NotificationsAlertTypeData)
	notifyData.NotifyOnlyOnTriggeredGroupByValues = conditionParams.NotifyGroupByOnlyAlerts
	notifyData.OnTriggerAndResolved = conditionParams.NotifyOnResolved

	return lucene, notifyData
}

func flattenTracingAlert(tracingAlert *alerts.TracingAlert, condition *alerts.AlertCondition) *coralogixv1alpha1.Tracing {
	latencyThresholdMS := float64(tracingAlert.GetConditionLatency()) / float64(time.Millisecond.Microseconds())
	tracingFilters := flattenTracingAlertFilters(tracingAlert)
	tracingFilters.LatencyThresholdMilliseconds = utils.FloatToQuantity(latencyThresholdMS)

	tracingCondition := flattenTracingCondition(condition)

	return &coralogixv1alpha1.Tracing{
		Filters:    tracingFilters,
		Conditions: tracingCondition,
	}
}

func flattenTracingCondition(condition *alerts.AlertCondition) coralogixv1alpha1.TracingCondition {
	var tracingCondition coralogixv1alpha1.TracingCondition
	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_Immediate:
		tracingCondition.AlertWhen = coralogixv1alpha1.TracingAlertWhenImmediately
	case *alerts.AlertCondition_MoreThan:
		conditionParams := condition.MoreThan.GetParameters()
		tracingCondition.AlertWhen = coralogixv1alpha1.TracingAlertWhenMore

		tracingCondition.Threshold = new(int)
		*tracingCondition.Threshold = int(conditionParams.GetThreshold().GetValue())

		tracingCondition.TimeWindow = new(coralogixv1alpha1.TimeWindow)
		*tracingCondition.TimeWindow = coralogixv1alpha1.TimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])

		tracingCondition.GroupBy = utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy())
	}

	return tracingCondition
}

func flattenTracingAlertFilters(tracingAlert *alerts.TracingAlert) coralogixv1alpha1.TracingFilters {
	applications, subsystems, services := flattenTracingFilters(tracingAlert.GetFieldFilters())
	tagFilters := flattenTagFiltersData(tracingAlert.GetTagFilters())

	return coralogixv1alpha1.TracingFilters{
		TagFilters:   tagFilters,
		Applications: applications,
		Subsystems:   subsystems,
		Services:     services,
	}
}

func flattenTracingFilters(tracingFilters []*alerts.FilterData) (applications, subsystems, services []string) {
	filtersData := flattenFiltersData(tracingFilters)
	applications = filtersData["applicationName"]
	subsystems = filtersData["subsystemName"]
	services = filtersData["serviceName"]
	return
}

func flattenTagFiltersData(filtersData []*alerts.FilterData) []coralogixv1alpha1.TagFilter {
	fieldToFilters := flattenFiltersData(filtersData)
	result := make([]coralogixv1alpha1.TagFilter, 0, len(fieldToFilters))
	for field, filters := range fieldToFilters {
		filterSchema := coralogixv1alpha1.TagFilter{
			Field:  field,
			Values: filters,
		}
		result = append(result, filterSchema)
	}
	return result
}

func flattenFiltersData(filtersData []*alerts.FilterData) map[string][]string {
	result := make(map[string][]string, len(filtersData))
	for _, filter := range filtersData {
		field := filter.GetField()
		result[field] = flattenTracingFilter(filter.GetFilters())
	}
	return result
}

func flattenTracingFilter(filters []*alerts.Filters) []string {
	result := make([]string, 0)
	for _, f := range filters {
		values := f.GetValues()
		switch operator := f.GetOperator(); operator {
		case "contains", "startsWith", "endsWith":
			for i, val := range values {
				values[i] = fmt.Sprintf("filter:%s:%s", operator, val)
			}
		}
		result = append(result, values...)
	}
	return result
}

func flattenFlowAlert(flow *alerts.FlowCondition) *coralogixv1alpha1.Flow {
	stages := flattenFlowStages(flow.Stages)
	return &coralogixv1alpha1.Flow{
		Stages: stages,
	}
}

func flattenFlowStages(stages []*alerts.FlowStage) []coralogixv1alpha1.FlowStage {
	result := make([]coralogixv1alpha1.FlowStage, 0, len(stages))
	for _, s := range stages {
		stage := flattenFlowStage(s)
		result = append(result, stage)
	}
	return result
}

func flattenFlowStage(stage *alerts.FlowStage) coralogixv1alpha1.FlowStage {
	groups := flattenFlowStageGroups(stage.Groups)

	var timeFrame *coralogixv1alpha1.FlowStageTimeFrame
	if timeWindow := stage.GetTimeframe(); timeWindow != nil {
		timeFrame = convertMillisecondToTime(int(timeWindow.GetMs().GetValue()))
	}

	return coralogixv1alpha1.FlowStage{
		Groups:     groups,
		TimeWindow: timeFrame,
	}
}

func convertMillisecondToTime(timeMS int) *coralogixv1alpha1.FlowStageTimeFrame {
	if timeMS == 0 {
		return nil
	}

	msInHour := int(time.Hour.Milliseconds())
	msInMinute := int(time.Minute.Milliseconds())
	msInSecond := int(time.Second.Milliseconds())

	hours := timeMS / msInHour
	timeMS -= hours * msInHour

	minutes := timeMS / msInMinute
	timeMS -= minutes * msInMinute

	seconds := timeMS / msInSecond

	return &coralogixv1alpha1.FlowStageTimeFrame{
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
}

func flattenFlowStageGroups(groups []*alerts.FlowGroup) []coralogixv1alpha1.FlowStageGroup {
	result := make([]coralogixv1alpha1.FlowStageGroup, 0, len(groups))
	for _, g := range groups {
		group := flattenFlowStageGroup(g)
		result = append(result, group)
	}
	return result
}

func flattenFlowStageGroup(group *alerts.FlowGroup) coralogixv1alpha1.FlowStageGroup {
	subAlerts := expandFlowSubgroupAlerts(group.GetAlerts())
	nextOp := alertProtoFlowOperatorToProtoFlowOperator[group.GetNextOp()]
	return coralogixv1alpha1.FlowStageGroup{
		InnerFlowAlerts: subAlerts,
		NextOperator:    nextOp,
	}
}

func expandFlowSubgroupAlerts(subgroup *alerts.FlowAlerts) coralogixv1alpha1.InnerFlowAlerts {
	return coralogixv1alpha1.InnerFlowAlerts{
		Operator: alertProtoFlowOperatorToProtoFlowOperator[subgroup.GetOp()],
		Alerts:   expandFlowInnerAlerts(subgroup.GetValues()),
	}
}

func expandFlowInnerAlerts(innerAlerts []*alerts.FlowAlert) []coralogixv1alpha1.InnerFlowAlert {
	result := make([]coralogixv1alpha1.InnerFlowAlert, 0, len(innerAlerts))
	for _, a := range innerAlerts {
		alert := expandFlowInnerAlert(a)
		result = append(result, alert)
	}
	return result
}

func expandFlowInnerAlert(alert *alerts.FlowAlert) coralogixv1alpha1.InnerFlowAlert {
	return coralogixv1alpha1.InnerFlowAlert{
		UserAlertId: alert.GetId().GetValue(),
		Not:         alert.GetNot().GetValue(),
	}
}

func flattenMetaLabels(labels []*alerts.MetaLabel) map[string]string {
	if len(labels) == 0 {
		return nil
	}

	result := make(map[string]string)
	for _, label := range labels {
		result[label.GetKey().GetValue()] = label.GetValue().GetValue()
	}
	return result
}

func flattenExpirationDate(expirationDate *alerts.Date) *coralogixv1alpha1.ExpirationDate {
	if expirationDate == nil {
		return nil
	}

	return &coralogixv1alpha1.ExpirationDate{
		Day:   expirationDate.Day,
		Month: expirationDate.Month,
		Year:  expirationDate.Year,
	}
}

func flattenScheduling(scheduling *alerts.AlertActiveWhen, spec coralogixv1alpha1.AlertSpec) *coralogixv1alpha1.Scheduling {
	if scheduling == nil || len(scheduling.GetTimeframes()) == 0 {
		return nil
	}

	timeZone := coralogixv1alpha1.TimeZone("UTC+00")
	var utc int32
	if spec.Scheduling != nil {
		timeZone = spec.Scheduling.TimeZone
		utc = coralogixv1alpha1.ExtractUTC(timeZone)
	}

	timeframe := scheduling.GetTimeframes()[0]
	timeRange := timeframe.GetRange()
	activityStartGMT, activityEndGMT := timeRange.GetStart(), timeRange.GetEnd()
	daysOffset := getDaysOffsetFromGMT(activityStartGMT, utc)
	daysEnabled := flattenDaysOfWeek(timeframe.GetDaysOfWeek(), daysOffset)
	activityStartUTC := flattenTimeInDay(activityStartGMT, utc)
	activityEndUTC := flattenTimeInDay(activityEndGMT, utc)

	return &coralogixv1alpha1.Scheduling{
		TimeZone:    timeZone,
		DaysEnabled: daysEnabled,
		StartTime:   activityStartUTC,
		EndTime:     activityEndUTC,
	}
}

func getDaysOffsetFromGMT(activityStartGMT *alerts.Time, utc int32) int32 {
	daysOffset := int32(activityStartGMT.GetHours()+utc) / 24
	if daysOffset < 0 {
		daysOffset += 7
	}

	return daysOffset
}

func flattenTimeInDay(time *alerts.Time, utc int32) *coralogixv1alpha1.Time {
	hours := convertGmtToUtc(time.GetHours(), utc)
	hoursStr := toTwoDigitsFormat(hours)
	minStr := toTwoDigitsFormat(time.GetMinutes())
	result := coralogixv1alpha1.Time(fmt.Sprintf("%s:%s", hoursStr, minStr))
	return &result
}

func convertGmtToUtc(hours, utc int32) int32 {
	hours += utc
	if hours < 0 {
		hours += 24
	} else if hours >= 24 {
		hours -= 24
	}

	return hours
}

func toTwoDigitsFormat(digit int32) string {
	digitStr := fmt.Sprintf("%d", digit)
	if len(digitStr) == 1 {
		digitStr = "0" + digitStr
	}
	return digitStr
}

func flattenDaysOfWeek(daysOfWeek []alerts.DayOfWeek, daysOffset int32) []coralogixv1alpha1.Day {
	result := make([]coralogixv1alpha1.Day, 0, len(daysOfWeek))
	for _, d := range daysOfWeek {
		dayConvertedFromGmtToUtc := alerts.DayOfWeek((int32(d) + daysOffset) % 7)
		day := alertProtoDayToSchemaDay[dayConvertedFromGmtToUtc]
		result = append(result, day)
	}
	return result
}

func flattenNotifications(recipients *alerts.AlertNotifications, notifyEverySec *wrapperspb.DoubleValue, notifyData *NotificationsAlertTypeData) *coralogixv1alpha1.Notifications {
	flattenedRecipients := flattenRecipients(recipients)
	notifyEveryMin := int(notifyEverySec.GetValue() / 60)
	onTriggerAndResolved := notifyData.OnTriggerAndResolved.GetValue()
	notifyOnlyOnTriggeredGroupByValues := notifyData.NotifyOnlyOnTriggeredGroupByValues.GetValue()

	return &coralogixv1alpha1.Notifications{
		Recipients:                         flattenedRecipients,
		NotifyEveryMin:                     &notifyEveryMin,
		OnTriggerAndResolved:               onTriggerAndResolved,
		NotifyOnlyOnTriggeredGroupByValues: notifyOnlyOnTriggeredGroupByValues,
	}
}

func flattenRecipients(recipients *alerts.AlertNotifications) coralogixv1alpha1.Recipients {
	emails := utils.WrappedStringSliceToStringSlice(recipients.GetEmails())
	webhooks := utils.WrappedStringSliceToStringSlice(recipients.GetIntegrations())
	return coralogixv1alpha1.Recipients{
		Emails:   emails,
		Webhooks: webhooks,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlertReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&coralogixv1alpha1.Alert{}).
		Complete(r)
}
