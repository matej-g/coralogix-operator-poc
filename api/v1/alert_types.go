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

package v1

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	utils "coralogix-operator-poc/api"
	alerts "coralogix-operator-poc/controllers/clientset/grpc/alerts/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

var (
	alertSchemaSeverityToProtoSeverity = map[AlertSeverity]alerts.AlertSeverity{
		"Info":     alerts.AlertSeverity_ALERT_SEVERITY_INFO_OR_UNSPECIFIED,
		"Warning":  alerts.AlertSeverity_ALERT_SEVERITY_WARNING,
		"Critical": alerts.AlertSeverity_ALERT_SEVERITY_CRITICAL,
		"Error":    alerts.AlertSeverity_ALERT_SEVERITY_ERROR,
	}
	alertSchemaDayToProtoDay = map[Day]alerts.DayOfWeek{
		"Sunday":    alerts.DayOfWeek_DAY_OF_WEEK_SUNDAY,
		"Monday":    alerts.DayOfWeek_DAY_OF_WEEK_MONDAY_OR_UNSPECIFIED,
		"Tuesday":   alerts.DayOfWeek_DAY_OF_WEEK_TUESDAY,
		"Wednesday": alerts.DayOfWeek_DAY_OF_WEEK_WEDNESDAY,
		"Thursday":  alerts.DayOfWeek_DAY_OF_WEEK_THURSDAY,
		"Friday":    alerts.DayOfWeek_DAY_OF_WEEK_FRIDAY,
		"Saturday":  alerts.DayOfWeek_DAY_OF_WEEK_SATURDAY,
	}
	alertSchemaTimeWindowToProtoTimeWindow = map[string]alerts.Timeframe{
		"Minute":          alerts.Timeframe_TIMEFRAME_1_MIN,
		"FiveMinutes":     alerts.Timeframe_TIMEFRAME_5_MIN_OR_UNSPECIFIED,
		"TenMinutes":      alerts.Timeframe_TIMEFRAME_10_MIN,
		"FifteenMinutes":  alerts.Timeframe_TIMEFRAME_15_MIN,
		"TwentyMinutes":   alerts.Timeframe_TIMEFRAME_20_MIN,
		"ThirtyMinutes":   alerts.Timeframe_TIMEFRAME_30_MIN,
		"Hour":            alerts.Timeframe_TIMEFRAME_1_H,
		"TwoHours":        alerts.Timeframe_TIMEFRAME_2_H,
		"FourHours":       alerts.Timeframe_TIMEFRAME_4_H,
		"SixHours":        alerts.Timeframe_TIMEFRAME_6_H,
		"TwelveHours":     alerts.Timeframe_TIMEFRAME_12_H,
		"TwentyFourHours": alerts.Timeframe_TIMEFRAME_24_H,
		"ThirtySixHours":  alerts.Timeframe_TIMEFRAME_36_H,
	}
	alertSchemaAutoRetireRatioToProtoAutoRetireRatio = map[AutoRetireRatio]alerts.CleanupDeadmanDuration{
		"Never":           alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_NEVER_OR_UNSPECIFIED,
		"FiveMinutes":     alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_5MIN,
		"TenMinutes":      alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_10MIN,
		"Hour":            alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_1H,
		"TwoHours":        alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_2H,
		"SixHours":        alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_6H,
		"TwelveHours":     alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_12H,
		"TwentyFourHours": alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_24H,
	}
	alertSchemaFiltersLogSeverityToProtoFiltersLogSeverity = map[FiltersLogSeverity]alerts.AlertFilters_LogSeverity{
		"Debug":    alerts.AlertFilters_LOG_SEVERITY_DEBUG_OR_UNSPECIFIED,
		"Verbose":  alerts.AlertFilters_LOG_SEVERITY_VERBOSE,
		"Info":     alerts.AlertFilters_LOG_SEVERITY_INFO,
		"Warning":  alerts.AlertFilters_LOG_SEVERITY_WARNING,
		"Critical": alerts.AlertFilters_LOG_SEVERITY_CRITICAL,
		"Error":    alerts.AlertFilters_LOG_SEVERITY_ERROR,
	}
	alertSchemaRelativeTimeFrameToProtoTimeFrameAndRelativeTimeFrame = map[RelativeTimeWindow]protoTimeFrameAndRelativeTimeFrame{
		"PreviousHour":      {timeFrame: alerts.Timeframe_TIMEFRAME_1_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_HOUR_OR_UNSPECIFIED},
		"SameHourYesterday": {timeFrame: alerts.Timeframe_TIMEFRAME_1_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_DAY},
		"SameHourLastWeek":  {timeFrame: alerts.Timeframe_TIMEFRAME_1_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_WEEK},
		"Yesterday":         {timeFrame: alerts.Timeframe_TIMEFRAME_24_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_DAY},
		"SameDayLastWeek":   {timeFrame: alerts.Timeframe_TIMEFRAME_24_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_WEEK},
		"SameDayLastMonth":  {timeFrame: alerts.Timeframe_TIMEFRAME_24_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_MONTH},
	}
	alertSchemaArithmeticOperatorToProtoArithmeticOperator = map[ArithmeticOperator]alerts.MetricAlertConditionParameters_ArithmeticOperator{
		"Avg":        alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_AVG_OR_UNSPECIFIED,
		"Min":        alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_MIN,
		"Max":        alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_MAX,
		"Sum":        alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_SUM,
		"Count":      alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_COUNT,
		"Percentile": alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_PERCENTILE,
	}
	alertSchemaTracingFilterFieldToProtoTracingFilterField = map[FieldFilterType]string{
		"Application": "applicationName",
		"Subsystem":   "subsystemName",
		"Service":     "serviceName",
	}
	alertSchemaTracingOperatorToProtoTracingOperator = map[FilterOperator]string{
		"Equals":    "equals",
		"Contains":  "contains",
		"StartWith": "startsWith",
		"EndWith":   "endsWith",
	}
	alertSchemaFlowOperatorToProtoFlowOperator = map[FlowOperator]alerts.FlowOperator{
		"And": alerts.FlowOperator_AND,
		"Or":  alerts.FlowOperator_OR,
	}
	msInHour   = int(time.Hour.Milliseconds())
	msInMinute = int(time.Minute.Milliseconds())
	msInSecond = int(time.Second.Milliseconds())
)

type protoTimeFrameAndRelativeTimeFrame struct {
	timeFrame         alerts.Timeframe
	relativeTimeFrame alerts.RelativeTimeframe
}

// AlertSpec defines the desired state of Alert
type AlertSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:validation:MinLength=0
	Name string `json:"name,omitempty"`

	// +optional
	Description string `json:"description,omitempty"`

	//+kubebuilder:default=true
	Active bool `json:"active,omitempty"`

	Severity AlertSeverity `json:"severity,omitempty"`

	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// +optional
	ExpirationDate *ExpirationDate `json:"expirationDate,omitempty"`

	// +optional
	Notifications *Notifications `json:"notifications,omitempty"`

	Scheduling *Scheduling `json:"scheduling,omitempty"`

	AlertType AlertType `json:"alertType,omitempty"`
}

func (in *AlertSpec) ExtractCreateAlertRequest() *alerts.CreateAlertRequest {
	enabled := wrapperspb.Bool(in.Active)
	name := wrapperspb.String(in.Name)
	description := wrapperspb.String(in.Description)
	severity := alertSchemaSeverityToProtoSeverity[in.Severity]
	metaLabels := expandMetaLabels(in.Labels)
	expirationDate := expandExpirationDate(in.ExpirationDate)
	notifications := expandNotifications(in.Notifications.Recipients)
	notifyEvery := expandNotifyEvery(in.Notifications.NotifyEveryMin)
	payloadFilters := utils.StringSliceToWrappedStringSlice(in.Notifications.PayloadFilters)
	activeWhen := expandActiveWhen(in.Scheduling)
	alertTypeParams := expandAlertType(in.AlertType, in.Notifications.OnTriggerAndResolved,
		in.Notifications.NotifyOnlyOnTriggeredGroupByValues, in.Notifications.IgnoreInfinity)

	return &alerts.CreateAlertRequest{
		Name:                       name,
		Description:                description,
		IsActive:                   enabled,
		Severity:                   severity,
		MetaLabels:                 metaLabels,
		Expiration:                 expirationDate,
		Notifications:              notifications,
		NotifyEvery:                notifyEvery,
		NotificationPayloadFilters: payloadFilters,
		ActiveWhen:                 activeWhen,
		Filters:                    alertTypeParams.filters,
		Condition:                  alertTypeParams.condition,
		TracingAlert:               alertTypeParams.tracingAlert,
	}
}

type alertTypeParams struct {
	filters      *alerts.AlertFilters
	condition    *alerts.AlertCondition
	tracingAlert *alerts.TracingAlert
}

func expandAlertType(alertType AlertType, onTriggerAndResolved, notifyOnlyOnTriggeredGroupByValues, ignoreInfinity bool) alertTypeParams {
	if standard := alertType.Standard; standard != nil {
		return expandStandard(standard, onTriggerAndResolved, notifyOnlyOnTriggeredGroupByValues)
	} else if ratio := alertType.Ratio; ratio != nil {
		return expandRatio(ratio, onTriggerAndResolved, notifyOnlyOnTriggeredGroupByValues, ignoreInfinity)
	} else if newValue := alertType.NewValue; newValue != nil {
		return expandNewValue(newValue)
	} else if uniqueCount := alertType.UniqueCount; uniqueCount != nil {
		return expandUniqueCount(uniqueCount)
	} else if timeRelative := alertType.TimeRelative; newValue != nil {
		return expandTimeRelative(timeRelative, onTriggerAndResolved, notifyOnlyOnTriggeredGroupByValues, ignoreInfinity)
	} else if metric := alertType.Metric; metric != nil {
		return expandMetric(metric, onTriggerAndResolved, notifyOnlyOnTriggeredGroupByValues)
	} else if tracing := alertType.Tracing; tracing != nil {
		return expandTracing(tracing, onTriggerAndResolved)
	} else if flow := alertType.Flow; flow != nil {
		return expandFlow(flow)
	}

	return alertTypeParams{}
}

func expandStandard(standard *Standard, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues bool) alertTypeParams {
	condition := expandStandardCondition(standard.Conditions, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues)
	filters := expandCommonFilters(standard.Filters)
	filters.FilterType = alerts.AlertFilters_FILTER_TYPE_TEXT_OR_UNSPECIFIED
	return alertTypeParams{
		condition: condition,
		filters:   filters,
	}
}

func expandRatio(ratio *Ratio, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues, ignoreInfinity bool) alertTypeParams {
	groupBy := utils.StringSliceToWrappedStringSlice(ratio.Conditions.GroupBy)
	var groupByQ1, groupByQ2 []*wrapperspb.StringValue
	if groupByFor := ratio.Conditions.GroupByFor; groupByFor != nil {
		switch string(*groupByFor) {
		case "Q1":
			groupByQ1 = groupBy
		case "Q2":
			groupByQ2 = groupBy
		case "Both":
			groupByQ1 = groupBy
			groupByQ2 = groupBy
		}
	}

	condition := expandRatioCondition(ratio.Conditions, groupByQ1, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues, ignoreInfinity)
	filters := expandRatioFilters(&ratio.Query1Filters, &ratio.Query2Filters, groupByQ2)

	return alertTypeParams{
		condition: condition,
		filters:   filters,
	}
}

func expandRatioCondition(conditions RatioConditions, q1GroupBy []*wrapperspb.StringValue, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues, ignoreInf bool) *alerts.AlertCondition {
	threshold := wrapperspb.Double(conditions.Ratio.AsApproximateFloat64())
	timeFrame := alertSchemaTimeWindowToProtoTimeWindow[string(conditions.TimeWindow)]
	notifyOnResolved := wrapperspb.Bool(notifyWhenResolved)
	notifyGroupByOnlyAlerts := wrapperspb.Bool(notifyOnlyOnTriggeredGroupByValues)
	ignoreInfinity := wrapperspb.Bool(ignoreInf)
	relatedExtendedData := expandRelatedData(conditions.ManageUndetectedValues)

	parameters := &alerts.ConditionParameters{
		Threshold:               threshold,
		Timeframe:               timeFrame,
		GroupBy:                 q1GroupBy,
		NotifyOnResolved:        notifyOnResolved,
		IgnoreInfinity:          ignoreInfinity,
		NotifyGroupByOnlyAlerts: notifyGroupByOnlyAlerts,
		RelatedExtendedData:     relatedExtendedData,
	}

	switch conditions.AlertWhen {
	case "More":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_MoreThan{
				MoreThan: &alerts.MoreThanCondition{Parameters: parameters},
			},
		}
	case "Less":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_LessThan{
				LessThan: &alerts.LessThanCondition{Parameters: parameters},
			},
		}
	}

	return nil
}

func expandRatioFilters(q1Filters *Filters, q2Filters *RatioQ2Filters, groupByQ2 []*wrapperspb.StringValue) *alerts.AlertFilters {
	filters := expandCommonFilters(q1Filters)
	if q1Alias := q1Filters.Alias; q1Alias != nil {
		filters.Alias = wrapperspb.String(*q1Alias)
	}
	q2 := expandQ2Filters(q2Filters, groupByQ2)
	filters.RatioAlerts = []*alerts.AlertFilters_RatioAlert{q2}
	filters.FilterType = alerts.AlertFilters_FILTER_TYPE_RATIO
	return filters
}

func expandQ2Filters(q2Filters *RatioQ2Filters, q2GroupBy []*wrapperspb.StringValue) *alerts.AlertFilters_RatioAlert {
	var text *wrapperspb.StringValue
	if searchQuery := q2Filters.SearchQuery; searchQuery != nil {
		text = wrapperspb.String(*searchQuery)
	}
	alias := wrapperspb.String(q2Filters.Alias)
	severities := expandAlertFiltersSeverities(q2Filters.Severities)
	applications := utils.StringSliceToWrappedStringSlice(q2Filters.Applications)
	subsystems := utils.StringSliceToWrappedStringSlice(q2Filters.Subsystems)

	return &alerts.AlertFilters_RatioAlert{
		Alias:        alias,
		Text:         text,
		Severities:   severities,
		Applications: applications,
		Subsystems:   subsystems,
		GroupBy:      q2GroupBy,
	}
}

func expandNewValue(newValue *NewValue) alertTypeParams {
	condition := expandNewValueCondition(&newValue.Conditions)
	filters := expandCommonFilters(newValue.Filters)
	filters.FilterType = alerts.AlertFilters_FILTER_TYPE_TEXT_OR_UNSPECIFIED
	return alertTypeParams{
		condition: condition,
		filters:   filters,
	}
}

func expandNewValueCondition(conditions *NewValueConditions) *alerts.AlertCondition {
	timeFrame := alertSchemaTimeWindowToProtoTimeWindow[string(conditions.TimeWindow)]
	groupBy := []*wrapperspb.StringValue{wrapperspb.String(conditions.Key)}
	parameters := &alerts.ConditionParameters{
		Timeframe: timeFrame,
		GroupBy:   groupBy,
	}

	return &alerts.AlertCondition{
		Condition: &alerts.AlertCondition_NewValue{
			NewValue: &alerts.NewValueCondition{
				Parameters: parameters,
			},
		},
	}
}

func expandUniqueCount(uniqueCount *UniqueCount) alertTypeParams {
	condition := expandUniqueCountCondition(&uniqueCount.Conditions)
	filters := expandCommonFilters(uniqueCount.Filters)
	filters.FilterType = alerts.AlertFilters_FILTER_TYPE_UNIQUE_COUNT
	return alertTypeParams{
		condition: condition,
		filters:   filters,
	}
}

func expandUniqueCountCondition(conditions *UniqueCountConditions) *alerts.AlertCondition {
	uniqueCountKey := []*wrapperspb.StringValue{wrapperspb.String(conditions.Key)}
	threshold := wrapperspb.Double(float64(conditions.MaxUniqueValues))
	timeFrame := alertSchemaTimeWindowToProtoTimeWindow[string(conditions.TimeWindow)]
	var groupBy []*wrapperspb.StringValue
	var groupByThreshold *wrapperspb.UInt32Value
	if groupByKey := conditions.GroupBy; groupByKey != nil {
		wrapperspb.String(*groupByKey)
		groupByThreshold = wrapperspb.UInt32(uint32(*conditions.MaxUniqueValuesForGroupBy))
	}

	parameters := &alerts.ConditionParameters{
		CardinalityFields:                 uniqueCountKey,
		Threshold:                         threshold,
		Timeframe:                         timeFrame,
		GroupBy:                           groupBy,
		MaxUniqueCountValuesForGroupByKey: groupByThreshold,
	}

	return &alerts.AlertCondition{
		Condition: &alerts.AlertCondition_UniqueCount{
			UniqueCount: &alerts.UniqueCountCondition{
				Parameters: parameters,
			},
		},
	}
}

func expandTimeRelative(timeRelative *TimeRelative, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues, ignoreInfinity bool) alertTypeParams {
	condition := expandTimeRelativeCondition(&timeRelative.Conditions, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues, ignoreInfinity)
	filters := expandCommonFilters(timeRelative.Filters)
	filters.FilterType = alerts.AlertFilters_FILTER_TYPE_UNIQUE_COUNT
	return alertTypeParams{
		condition: condition,
		filters:   filters,
	}
}

func expandTimeRelativeCondition(condition *TimeRelativeConditions, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues, ignoreInfinity bool) *alerts.AlertCondition {
	threshold := wrapperspb.Double(condition.Threshold.AsApproximateFloat64())
	timeFrameAndRelativeTimeFrame := alertSchemaRelativeTimeFrameToProtoTimeFrameAndRelativeTimeFrame[condition.TimeWindow]
	groupBy := utils.StringSliceToWrappedStringSlice(condition.GroupBy)
	notifyOnResolved := wrapperspb.Bool(notifyWhenResolved)
	notifyGroupByOnlyAlerts := wrapperspb.Bool(notifyOnlyOnTriggeredGroupByValues)
	ignoreInf := wrapperspb.Bool(ignoreInfinity)
	relatedExtendedData := expandRelatedData(condition.ManageUndetectedValues)

	parameters := &alerts.ConditionParameters{
		Timeframe:               timeFrameAndRelativeTimeFrame.timeFrame,
		RelativeTimeframe:       timeFrameAndRelativeTimeFrame.relativeTimeFrame,
		GroupBy:                 groupBy,
		Threshold:               threshold,
		IgnoreInfinity:          ignoreInf,
		NotifyOnResolved:        notifyOnResolved,
		NotifyGroupByOnlyAlerts: notifyGroupByOnlyAlerts,
		RelatedExtendedData:     relatedExtendedData,
	}

	switch condition.AlertWhen {
	case "More":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_MoreThan{
				MoreThan: &alerts.MoreThanCondition{Parameters: parameters},
			},
		}
	case "Less":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_LessThan{
				LessThan: &alerts.LessThanCondition{Parameters: parameters},
			},
		}
	}

	return nil
}

func expandMetric(metric *Metric, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues bool) alertTypeParams {
	if promql := metric.Promql; promql != nil {
		return expandPromql(promql, notifyWhenResolved)
	} else if lucene := metric.Lucene; lucene != nil {
		return expandLucene(lucene, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues)
	}

	return alertTypeParams{}
}

func expandPromql(promql *Promql, notifyWhenResolved bool) alertTypeParams {
	condition := expandPromqlCondition(&promql.Conditions, promql.SearchQuery, notifyWhenResolved)
	filters := &alerts.AlertFilters{
		FilterType: alerts.AlertFilters_FILTER_TYPE_METRIC,
	}

	return alertTypeParams{
		condition: condition,
		filters:   filters,
	}
}

func expandPromqlCondition(conditions *PromqlConditions, searchQuery *string, notifyWhenResolved bool) *alerts.AlertCondition {
	var text *wrapperspb.StringValue
	if searchQuery != nil {
		text = wrapperspb.String(*searchQuery)
	}
	sampleThresholdPercentage := wrapperspb.UInt32(uint32(conditions.SampleThresholdPercentage))
	var nonNullPercentage *wrapperspb.UInt32Value
	if minNonNullValuesPercentage := conditions.MinNonNullValuesPercentage; minNonNullValuesPercentage != nil {
		nonNullPercentage = wrapperspb.UInt32(uint32(*minNonNullValuesPercentage))
	}
	swapNullValues := wrapperspb.Bool(conditions.ReplaceMissingValueWithZero)
	promqlParams := &alerts.MetricAlertPromqlConditionParameters{
		PromqlText:                text,
		SampleThresholdPercentage: sampleThresholdPercentage,
		NonNullPercentage:         nonNullPercentage,
		SwapNullValues:            swapNullValues,
	}
	threshold := wrapperspb.Double(conditions.Threshold.AsApproximateFloat64())
	timeWindow := alertSchemaTimeWindowToProtoTimeWindow[string(conditions.TimeWindow)]
	relatedExtendedData := expandRelatedData(conditions.ManageUndetectedValues)
	notifyOnResolved := wrapperspb.Bool(notifyWhenResolved)

	parameters := &alerts.ConditionParameters{
		Threshold:                   threshold,
		Timeframe:                   timeWindow,
		RelatedExtendedData:         relatedExtendedData,
		MetricAlertPromqlParameters: promqlParams,
		NotifyOnResolved:            notifyOnResolved,
	}

	switch conditions.AlertWhen {
	case "More":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_MoreThan{
				MoreThan: &alerts.MoreThanCondition{Parameters: parameters},
			},
		}
	case "Less":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_LessThan{
				LessThan: &alerts.LessThanCondition{Parameters: parameters},
			},
		}
	}

	return nil
}

func expandLucene(lucene *Lucene, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues bool) alertTypeParams {
	condition := expandLuceneCondition(&lucene.Conditions, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues)
	var text *wrapperspb.StringValue
	if searchQuery := lucene.SearchQuery; searchQuery != nil {
		text = wrapperspb.String(*searchQuery)
	}

	filters := &alerts.AlertFilters{
		FilterType: alerts.AlertFilters_FILTER_TYPE_METRIC,
		Text:       text,
	}

	return alertTypeParams{
		condition: condition,
		filters:   filters,
	}
}

func expandLuceneCondition(conditions *LuceneConditions, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues bool) *alerts.AlertCondition {
	metricField := wrapperspb.String(conditions.MetricField)
	arithmeticOperator := alertSchemaArithmeticOperatorToProtoArithmeticOperator[conditions.ArithmeticOperator]
	var arithmeticOperatorModifier *wrapperspb.UInt32Value
	if modifier := conditions.ArithmeticOperatorModifier; modifier != nil {
		arithmeticOperatorModifier = wrapperspb.UInt32(uint32(*modifier))
	}
	sampleThresholdPercentage := wrapperspb.UInt32(uint32(conditions.SampleThresholdPercentage))
	swapNullValues := wrapperspb.Bool(conditions.ReplaceMissingValueWithZero)
	var nonNullPercentage *wrapperspb.UInt32Value
	if minNonNullValuesPercentage := conditions.MinNonNullValuesPercentage; minNonNullValuesPercentage != nil {
		nonNullPercentage = wrapperspb.UInt32(uint32(*minNonNullValuesPercentage))
	}
	luceneParams := &alerts.MetricAlertConditionParameters{
		MetricSource:               alerts.MetricAlertConditionParameters_METRIC_SOURCE_LOGS2METRICS_OR_UNSPECIFIED,
		MetricField:                metricField,
		ArithmeticOperator:         arithmeticOperator,
		ArithmeticOperatorModifier: arithmeticOperatorModifier,
		SampleThresholdPercentage:  sampleThresholdPercentage,
		NonNullPercentage:          nonNullPercentage,
		SwapNullValues:             swapNullValues,
	}

	groupBy := utils.StringSliceToWrappedStringSlice(conditions.GroupBy)
	threshold := wrapperspb.Double(conditions.Threshold.AsApproximateFloat64())
	timeWindow := alertSchemaTimeWindowToProtoTimeWindow[string(conditions.TimeWindow)]
	relatedExtendedData := expandRelatedData(conditions.ManageUndetectedValues)
	notifyOnResolved := wrapperspb.Bool(notifyWhenResolved)
	notifyGroupByOnlyAlerts := wrapperspb.Bool(notifyOnlyOnTriggeredGroupByValues)

	parameters := &alerts.ConditionParameters{
		GroupBy:                 groupBy,
		Threshold:               threshold,
		Timeframe:               timeWindow,
		RelatedExtendedData:     relatedExtendedData,
		NotifyOnResolved:        notifyOnResolved,
		NotifyGroupByOnlyAlerts: notifyGroupByOnlyAlerts,
		MetricAlertParameters:   luceneParams,
	}

	switch conditions.AlertWhen {
	case "More":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_MoreThan{
				MoreThan: &alerts.MoreThanCondition{Parameters: parameters},
			},
		}
	case "Less":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_LessThan{
				LessThan: &alerts.LessThanCondition{Parameters: parameters},
			},
		}
	}

	return nil
}

func expandTracing(tracing *Tracing, notifyWhenResolved bool) alertTypeParams {
	filters := &alerts.AlertFilters{
		FilterType: alerts.AlertFilters_FILTER_TYPE_TRACING,
	}
	condition := expandTracingCondition(&tracing.Conditions, notifyWhenResolved)
	tracingAlert := expandTracingAlert(&tracing.Filters)
	return alertTypeParams{
		filters:      filters,
		condition:    condition,
		tracingAlert: tracingAlert,
	}
}

func expandTracingCondition(conditions *TracingCondition, notifyWhenResolved bool) *alerts.AlertCondition {
	switch conditions.AlertWhen {
	case "More":
		var timeFrame alerts.Timeframe
		if timeWindow := conditions.TimeWindow; timeWindow != nil {
			timeFrame = alertSchemaTimeWindowToProtoTimeWindow[string(*timeWindow)]
		}
		groupBy := utils.StringSliceToWrappedStringSlice(conditions.GroupBy)
		threshold := wrapperspb.Double(conditions.Threshold.AsApproximateFloat64())
		notifyOnResolved := wrapperspb.Bool(notifyWhenResolved)
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_MoreThan{
				MoreThan: &alerts.MoreThanCondition{
					Parameters: &alerts.ConditionParameters{
						Timeframe:        timeFrame,
						Threshold:        threshold,
						GroupBy:          groupBy,
						NotifyOnResolved: notifyOnResolved,
					},
				},
			},
		}
	case "Immediately":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_Immediate{},
		}
	}

	return nil
}

func expandTracingAlert(tracingFilters *TracingFilters) *alerts.TracingAlert {
	conditionLatency := uint32(tracingFilters.LatencyThresholdMS) * uint32(time.Millisecond.Microseconds())
	fieldFilters := expandFieldFilters(tracingFilters.FieldFilters)
	tagFilters := expandTagFilters(tracingFilters.TagFilters)
	return &alerts.TracingAlert{
		ConditionLatency: conditionLatency,
		FieldFilters:     fieldFilters,
		TagFilters:       tagFilters,
	}
}

func expandFieldFilters(fieldFilters []FieldFilter) []*alerts.FilterData {
	result := make([]*alerts.FilterData, 0, len(fieldFilters))
	for _, ff := range fieldFilters {
		field := alertSchemaTracingFilterFieldToProtoTracingFilterField[ff.Field]
		filters := make([]*alerts.Filters, 0, len(ff.Filters))
		for _, f := range ff.Filters {
			values := f.Values
			operator := alertSchemaTracingOperatorToProtoTracingOperator[f.Operator]
			filter := &alerts.Filters{
				Values:   values,
				Operator: operator,
			}
			filters = append(filters, filter)
		}

		filterData := &alerts.FilterData{
			Field:   field,
			Filters: filters,
		}
		result = append(result, filterData)
	}
	return result
}

func expandTagFilters(tagFilters []TagFilter) []*alerts.FilterData {
	result := make([]*alerts.FilterData, 0, len(tagFilters))
	for _, tf := range tagFilters {
		field := tf.Field
		filters := make([]*alerts.Filters, 0, len(tf.Filters))
		for _, f := range tf.Filters {
			values := f.Values
			operator := alertSchemaTracingOperatorToProtoTracingOperator[f.Operator]
			filter := &alerts.Filters{
				Values:   values,
				Operator: operator,
			}
			filters = append(filters, filter)
		}

		filterData := &alerts.FilterData{
			Field:   field,
			Filters: filters,
		}
		result = append(result, filterData)
	}
	return result
}

func expandFlow(flow *Flow) alertTypeParams {
	stages := expandFlowStages(flow.Stages)
	return alertTypeParams{
		condition: &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_Flow{
				Flow: &alerts.FlowCondition{
					Stages: stages,
				},
			},
		},
		filters: &alerts.AlertFilters{
			FilterType: alerts.AlertFilters_FILTER_TYPE_FLOW,
		},
	}
}

func expandFlowStages(stages []FlowStage) []*alerts.FlowStage {
	result := make([]*alerts.FlowStage, 0, len(stages))
	for _, s := range stages {
		stage := expandFlowStage(s)
		result = append(result, stage)
	}
	return result
}

func expandFlowStage(stage FlowStage) *alerts.FlowStage {
	groups := expandFlowStageGroups(stage.Groups)
	var timeFrame *alerts.FlowTimeframe
	if timeWindow := stage.TimeWindow; timeWindow != nil {
		timeFrame = new(alerts.FlowTimeframe)
		timeFrame.Ms = wrapperspb.UInt32(uint32(expandTimeToMS(*timeWindow)))
	}

	return &alerts.FlowStage{
		Groups:    groups,
		Timeframe: timeFrame,
	}
}

func expandTimeToMS(t FlowStageTimeFrame) int {
	timeMS := msInHour * t.Hours
	timeMS += msInMinute * t.Minutes
	timeMS += msInSecond * t.Seconds

	return timeMS
}

func expandFlowStageGroups(groups []FlowStageGroup) []*alerts.FlowGroup {
	result := make([]*alerts.FlowGroup, 0, len(groups))
	for _, g := range groups {
		group := expandFlowStageGroup(g)
		result = append(result, group)
	}
	return result
}

func expandFlowStageGroup(group FlowStageGroup) *alerts.FlowGroup {
	subAlerts := expandFlowSubgroupAlerts(group.SubgroupAlerts)
	nextOp := alertSchemaFlowOperatorToProtoFlowOperator[group.Operator]
	return &alerts.FlowGroup{
		Alerts: subAlerts,
		NextOp: nextOp,
	}
}

func expandFlowSubgroupAlerts(subgroup SubgroupAlerts) *alerts.FlowAlerts {
	return &alerts.FlowAlerts{
		Op:     alertSchemaFlowOperatorToProtoFlowOperator[subgroup.Operator],
		Values: expandFlowInnerAlerts(subgroup.Alerts),
	}
}

func expandFlowInnerAlerts(innerAlerts []InnerFlowAlert) []*alerts.FlowAlert {
	result := make([]*alerts.FlowAlert, 0, len(innerAlerts))
	for _, a := range innerAlerts {
		alert := expandFlowInnerAlert(a)
		result = append(result, alert)
	}
	return result
}

func expandFlowInnerAlert(alert InnerFlowAlert) *alerts.FlowAlert {
	return &alerts.FlowAlert{
		Id:  wrapperspb.String(alert.UserAlertId),
		Not: wrapperspb.Bool(alert.Not),
	}
}

func expandCommonFilters(filters *Filters) *alerts.AlertFilters {
	severities := expandAlertFiltersSeverities(filters.Severities)
	metadata := expandMetadata(filters)
	var text *wrapperspb.StringValue
	if searchQuery := filters.SearchQuery; searchQuery != nil {
		text = wrapperspb.String(*searchQuery)
	}
	return &alerts.AlertFilters{
		Severities: severities,
		Metadata:   metadata,
		Text:       text,
	}
}

func expandAlertFiltersSeverities(severities []FiltersLogSeverity) []alerts.AlertFilters_LogSeverity {
	result := make([]alerts.AlertFilters_LogSeverity, 0, len(severities))
	for _, s := range severities {
		severity := alertSchemaFiltersLogSeverityToProtoFiltersLogSeverity[s]
		result = append(result, severity)
	}
	return result
}

func expandMetadata(filters *Filters) *alerts.AlertFilters_MetadataFilters {
	categories := utils.StringSliceToWrappedStringSlice(filters.Categories)
	applications := utils.StringSliceToWrappedStringSlice(filters.Applications)
	subsystems := utils.StringSliceToWrappedStringSlice(filters.Subsystems)
	ips := utils.StringSliceToWrappedStringSlice(filters.IPs)
	classes := utils.StringSliceToWrappedStringSlice(filters.Classes)
	methods := utils.StringSliceToWrappedStringSlice(filters.Methods)
	computers := utils.StringSliceToWrappedStringSlice(filters.Computers)
	return &alerts.AlertFilters_MetadataFilters{
		Categories:   categories,
		Applications: applications,
		Subsystems:   subsystems,
		IpAddresses:  ips,
		Classes:      classes,
		Methods:      methods,
		Computers:    computers,
	}
}

func expandStandardCondition(condition StandardConditions, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues bool) *alerts.AlertCondition {
	var threshold *wrapperspb.DoubleValue
	if condition.Threshold != nil {
		threshold = wrapperspb.Double(float64(*condition.Threshold))
	}
	var timeFrame alerts.Timeframe
	if condition.TimeWindow != nil {
		timeFrame = alertSchemaTimeWindowToProtoTimeWindow[string(*condition.TimeWindow)]
	}
	groupBy := utils.StringSliceToWrappedStringSlice(condition.GroupBy)
	notifyOnResolved := wrapperspb.Bool(notifyWhenResolved)
	notifyGroupByOnlyAlerts := wrapperspb.Bool(notifyOnlyOnTriggeredGroupByValues)
	relatedExtendedData := expandRelatedData(condition.ManageUndetectedValues)

	parameters := &alerts.ConditionParameters{
		Threshold:               threshold,
		Timeframe:               timeFrame,
		GroupBy:                 groupBy,
		NotifyOnResolved:        notifyOnResolved,
		NotifyGroupByOnlyAlerts: notifyGroupByOnlyAlerts,
		RelatedExtendedData:     relatedExtendedData,
	}

	switch condition.AlertWhen {
	case "More":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_MoreThan{
				MoreThan: &alerts.MoreThanCondition{Parameters: parameters},
			},
		}
	case "Less":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_LessThan{
				LessThan: &alerts.LessThanCondition{Parameters: parameters},
			},
		}
	case "Immediately":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_Immediate{},
		}
	case "MoreThanUsual":
		return &alerts.AlertCondition{
			Condition: &alerts.AlertCondition_MoreThanUsual{
				MoreThanUsual: &alerts.MoreThanUsualCondition{Parameters: parameters},
			},
		}
	}

	return nil
}

func expandRelatedData(manageUndetectedValues *ManageUndetectedValues) *alerts.RelatedExtendedData {
	if manageUndetectedValues != nil {
		shouldTriggerDeadman := wrapperspb.Bool(manageUndetectedValues.EnableTriggeringOnUndetectedValues)
		cleanupDeadmanDuration := alertSchemaAutoRetireRatioToProtoAutoRetireRatio[*manageUndetectedValues.AutoRetireRatio]
		return &alerts.RelatedExtendedData{
			ShouldTriggerDeadman:   shouldTriggerDeadman,
			CleanupDeadmanDuration: &cleanupDeadmanDuration,
		}
	}
	return nil
}

func expandActiveWhen(scheduling *Scheduling) *alerts.AlertActiveWhen {
	if scheduling == nil {
		return nil
	}

	daysOfWeek := expandDaysOfWeek(scheduling.DaysEnabled)
	start := expandTime(scheduling.StartTime)
	end := expandTime(scheduling.EndTime)

	return &alerts.AlertActiveWhen{
		Timeframes: []*alerts.AlertActiveTimeframe{
			{
				DaysOfWeek: daysOfWeek,
				Range: &alerts.TimeRange{
					Start: start,
					End:   end,
				},
			},
		},
	}
}

func expandTime(time *Time) *alerts.Time {
	if time == nil {
		return nil
	}

	timeArr := strings.Split(string(*time), ":")
	hours, _ := strconv.Atoi(timeArr[0])
	minutes, _ := strconv.Atoi(timeArr[1])

	return &alerts.Time{
		Hours:   int32(hours),
		Minutes: int32(minutes),
	}
}

func expandDaysOfWeek(days []Day) []alerts.DayOfWeek {
	daysOfWeek := make([]alerts.DayOfWeek, 0, len(days))
	for _, d := range days {
		daysOfWeek = append(daysOfWeek, alertSchemaDayToProtoDay[d])
	}
	return daysOfWeek
}

func expandMetaLabels(labels map[string]string) []*alerts.MetaLabel {
	result := make([]*alerts.MetaLabel, 0)
	for k, v := range labels {
		result = append(result, &alerts.MetaLabel{
			Key:   wrapperspb.String(k),
			Value: wrapperspb.String(v),
		})
	}
	return result
}

func expandExpirationDate(date *ExpirationDate) *alerts.Date {
	if date == nil {
		return nil
	}

	return &alerts.Date{
		Year:  date.Year,
		Month: date.Month,
		Day:   date.Day,
	}
}

func expandNotifications(recipients Recipients) *alerts.AlertNotifications {
	return &alerts.AlertNotifications{
		Emails:       utils.StringSliceToWrappedStringSlice(recipients.Emails),
		Integrations: utils.StringSliceToWrappedStringSlice(recipients.Webhooks),
	}
}

func expandNotifyEvery(notifyEveryMin *int) *wrapperspb.DoubleValue {
	if notifyEveryMin == nil {
		return nil
	}
	return wrapperspb.Double(float64(60 * *notifyEveryMin))
}

//func emailSliceToWrappedStringSlice(arr []Email) []*wrapperspb.StringValue {
//	result := make([]*wrapperspb.StringValue, 0, len(arr))
//	for _, s := range arr {
//		result = append(result, wrapperspb.String(string(s)))
//	}
//	return result
//}

func (in *AlertSpec) DeepEqual(actualAlert *alerts.Alert) (bool, utils.Diff) {
	if actualName := actualAlert.GetName().GetValue(); actualName != in.Name {
		return false, utils.Diff{
			Name:    "Name",
			Desired: in.Name,
			Actual:  actualName,
		}
	}

	if actualDescription := actualAlert.GetDescription().GetValue(); actualDescription != in.Description {
		return false, utils.Diff{
			Name:    "Description",
			Desired: in.Description,
			Actual:  actualDescription,
		}
	}

	if actualActive := actualAlert.GetIsActive().GetValue(); actualActive != in.Active {
		return false, utils.Diff{
			Name:    "Active",
			Desired: in.Active,
			Actual:  actualActive,
		}
	}

	if actualSeverity := actualAlert.GetSeverity(); actualSeverity != alertSchemaSeverityToProtoSeverity[in.Severity] {
		return false, utils.Diff{
			Name:    "Severity",
			Desired: in.Severity,
			Actual:  actualSeverity.String(),
		}
	}

	if !equalLabels(in.Labels, actualAlert.MetaLabels) {
		return false, utils.Diff{
			Name:    "Labels",
			Desired: in.Labels,
			Actual:  actualAlert.MetaLabels,
		}
	}

	if expirationDate, actualExpirationDate := in.ExpirationDate, actualAlert.GetExpiration(); expirationDate == nil && actualExpirationDate != nil {
		return false, utils.Diff{
			Name:    "ExpirationDate",
			Desired: in.ExpirationDate,
			Actual:  *actualExpirationDate,
		}
	} else if expirationDate != nil && !expirationDate.DeepEqual(actualExpirationDate) {
		return false, utils.Diff{
			Name:    "ExpirationDate",
			Desired: *in.ExpirationDate,
			Actual:  *actualExpirationDate,
		}
	}

	var notifyData notificationsAlertTypeData
	if equal, diff := in.AlertType.DeepEqual(actualAlert, &notifyData); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("AlertType.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Notifications.DeepEqual(actualAlert.GetNotifications(),
		actualAlert.GetNotificationPayloadFilters(), actualAlert.GetNotifyEvery(), notifyData); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Notifications.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

type notificationsAlertTypeData struct {
	onTriggerAndResolved *wrapperspb.BoolValue

	notifyOnlyOnTriggeredGroupByValues *wrapperspb.BoolValue

	ignoreInfinity *wrapperspb.BoolValue
}

func equalLabels(labels map[string]string, actualLabels []*alerts.MetaLabel) bool {
	if len(labels) != len(actualLabels) {
		return false
	}

	for _, label := range actualLabels {
		if value, ok := labels[label.GetKey().GetValue()]; !ok || value != label.GetValue().GetValue() {
			return false
		}
	}

	return true
}

func (in *AlertSpec) ExtractUpdateAlertRequest(id string) *alerts.UpdateAlertByUniqueIdRequest {
	uniqueIdentifier := wrapperspb.String(id)
	enabled := wrapperspb.Bool(in.Active)
	name := wrapperspb.String(in.Name)
	description := wrapperspb.String(in.Description)
	severity := alertSchemaSeverityToProtoSeverity[in.Severity]
	metaLabels := expandMetaLabels(in.Labels)
	expirationDate := expandExpirationDate(in.ExpirationDate)
	notifications := expandNotifications(in.Notifications.Recipients)
	notifyEvery := expandNotifyEvery(in.Notifications.NotifyEveryMin)
	payloadFilters := utils.StringSliceToWrappedStringSlice(in.Notifications.PayloadFilters)
	activeWhen := expandActiveWhen(in.Scheduling)
	alertTypeParams := expandAlertType(in.AlertType, in.Notifications.OnTriggerAndResolved,
		in.Notifications.NotifyOnlyOnTriggeredGroupByValues, in.Notifications.IgnoreInfinity)

	return &alerts.UpdateAlertByUniqueIdRequest{
		Alert: &alerts.Alert{
			UniqueIdentifier:           uniqueIdentifier,
			Name:                       name,
			Description:                description,
			IsActive:                   enabled,
			Severity:                   severity,
			MetaLabels:                 metaLabels,
			Expiration:                 expirationDate,
			Notifications:              notifications,
			NotifyEvery:                notifyEvery,
			NotificationPayloadFilters: payloadFilters,
			ActiveWhen:                 activeWhen,
			Filters:                    alertTypeParams.filters,
			Condition:                  alertTypeParams.condition,
			TracingAlert:               alertTypeParams.tracingAlert,
		},
	}
}

// +kubebuilder:validation:Enum=Info;Warning;Critical;Error
type AlertSeverity string

type ExpirationDate struct {
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=31
	Day int32 `json:"day,omitempty"`

	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=12
	Month int32 `json:"month,omitempty"`

	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=9999
	Year int32 `json:"year,omitempty"`
}

func (in *ExpirationDate) DeepEqual(date *alerts.Date) bool {
	return in.Year != date.Year || in.Month != date.Month || in.Day != date.Day
}

type Notifications struct {
	//+kubebuilder:default=false
	OnTriggerAndResolved bool `json:"onTriggerAndResolved,omitempty"`

	//+kubebuilder:default=false
	IgnoreInfinity bool `json:"ignoreInfinity,omitempty"`

	//+kubebuilder:default=false
	NotifyOnlyOnTriggeredGroupByValues bool `json:"notifyOnlyOnTriggeredGroupByValues,omitempty"`

	// +optional
	Recipients Recipients `json:"recipients,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum:=1
	NotifyEveryMin *int `json:"notifyEveryMin,omitempty"`

	// +optional
	PayloadFilters []string `json:"payloadFilters,omitempty"`
}

func (in *Notifications) DeepEqual(actualRecipients *alerts.AlertNotifications,
	actualNotifyPayloadFilters []*wrapperspb.StringValue, actualNotifyEvery *wrapperspb.DoubleValue,
	actualNotifyAlertTypeData notificationsAlertTypeData) (bool, utils.Diff) {

	if emails, actualEmails := in.Recipients.Emails, utils.WrappedStringSliceToStringSlice(actualRecipients.Emails); !utils.SlicesWithUniqueValuesEqual(emails, actualEmails) {
		return false, utils.Diff{
			Name: "Emails",
		}
	}

	if webhooks, actualWebhooks := in.Recipients.Webhooks, utils.WrappedStringSliceToStringSlice(actualRecipients.Integrations); !utils.SlicesWithUniqueValuesEqual(webhooks, actualWebhooks) {
		return false, utils.Diff{
			Name: "Webhooks",
		}
	}

	if actualNotifyPayloadFilters := utils.WrappedStringSliceToStringSlice(actualNotifyPayloadFilters); !utils.SlicesWithUniqueValuesEqual(in.PayloadFilters, actualNotifyPayloadFilters) {
		return false, utils.Diff{
			Name: "PayloadFilters",
		}
	}

	if in.NotifyEveryMin == nil && actualNotifyEvery != nil {
		return false, utils.Diff{
			Name:    "NotifyEveryMin",
			Desired: in.NotifyEveryMin,
			Actual:  actualNotifyEvery.GetValue(),
		}
	} else if desiredNotifyEverySec := float64(*in.NotifyEveryMin) * 60; actualNotifyEvery.GetValue() != desiredNotifyEverySec {
		return false, utils.Diff{
			Name:    "NotifyEveryMin",
			Desired: fmt.Sprintf("%d (minutes)", *in.NotifyEveryMin),
			Actual:  fmt.Sprintf("%f (seconds)", actualNotifyEvery.GetValue()),
		}
	}

	if actualOnTriggerAndResolved := actualNotifyAlertTypeData.onTriggerAndResolved.GetValue(); in.OnTriggerAndResolved != actualOnTriggerAndResolved {
		return false, utils.Diff{
			Name:    "OnTriggerAndResolved",
			Desired: in.OnTriggerAndResolved,
			Actual:  actualOnTriggerAndResolved,
		}
	}

	if actualNotifyOnlyOnTriggeredGroupByValues := actualNotifyAlertTypeData.notifyOnlyOnTriggeredGroupByValues.GetValue(); in.NotifyOnlyOnTriggeredGroupByValues != actualNotifyOnlyOnTriggeredGroupByValues {
		return false, utils.Diff{
			Name:    "NotifyOnlyOnTriggeredGroupByValues",
			Desired: in.NotifyOnlyOnTriggeredGroupByValues,
			Actual:  actualNotifyOnlyOnTriggeredGroupByValues,
		}
	}

	if actualIgnoreInfinity := actualNotifyAlertTypeData.ignoreInfinity.GetValue(); in.IgnoreInfinity != actualIgnoreInfinity {
		return false, utils.Diff{
			Name:    "IgnoreInfinity",
			Desired: in.IgnoreInfinity,
			Actual:  actualIgnoreInfinity,
		}
	}

	return true, utils.Diff{}
}

type Recipients struct {
	// +optional
	Emails []string `json:"emails,omitempty"`

	// +optional
	Webhooks []string `json:"webhooks,omitempty"`
}

type Scheduling struct {
	//+kubebuilder:default=UTC+0
	TimeZone TimeZone `json:"timeZone,omitempty"`

	DaysEnabled []Day `json:"daysEnabled,omitempty"`

	StartTime *Time `json:"startTime,omitempty"`

	EndTime *Time `json:"endTime,omitempty"`
}

/* +kubebuilder:validation:Pattern=/^UTC[+-]\d{2}:\d{2}$/g*/
type TimeZone string

// +kubebuilder:validation:Enum=Sunday;Monday;Tuesday;Wednesday;Thursday;Friday;Saturday;
type Day string

/* +kubebuilder:validation:Pattern=^(0\d|1\d|2[0-3]):[0-5]\d$*/
type Time string

type AlertType struct {
	// +optional
	Standard *Standard `json:"standard,omitempty"`

	// +optional
	Ratio *Ratio `json:"ratio,omitempty"`

	// +optional
	NewValue *NewValue `json:"newValue,omitempty"`

	// +optional
	UniqueCount *UniqueCount `json:"uniqueCount,omitempty"`

	// +optional
	TimeRelative *TimeRelative `json:"timeRelative,omitempty"`

	// +optional
	Metric *Metric `json:"metric,omitempty"`

	// +optional
	Tracing *Tracing `json:"tracing,omitempty"`

	// +optional
	Flow *Flow `json:"flow,omitempty"`
}

func (in *AlertType) DeepEqual(actualAlert *alerts.Alert, notifyData *notificationsAlertTypeData) (bool, utils.Diff) {
	actualFilters := actualAlert.GetFilters()
	actualCondition := actualAlert.GetCondition()

	switch actualFilters.GetFilterType() {
	case alerts.AlertFilters_FILTER_TYPE_TEXT_OR_UNSPECIFIED:
		if newValueCondition, ok := actualCondition.GetCondition().(*alerts.AlertCondition_NewValue); ok {
			if newValue := in.NewValue; newValue == nil {
				return false, utils.Diff{
					Name:   "Type",
					Actual: "NewValue",
				}
			} else if equal, diff := newValue.DeepEqual(actualFilters, newValueCondition); !equal {
				return false, utils.Diff{
					Name:    fmt.Sprintf("NewValue.%s", diff.Name),
					Desired: diff.Desired,
					Actual:  diff.Actual,
				}
			}
		} else {
			if standard := in.Standard; standard == nil {
				return false, utils.Diff{
					Name:   "Type",
					Actual: "Standard",
				}
			} else if equal, diff := standard.DeepEqual(actualFilters, actualAlert.GetCondition(), notifyData); !equal {
				return false, utils.Diff{
					Name:    fmt.Sprintf("Standard.%s", diff.Name),
					Desired: diff.Desired,
					Actual:  diff.Actual,
				}
			}
		}
	case alerts.AlertFilters_FILTER_TYPE_RATIO:
		if ratio := in.Ratio; ratio == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Ratio",
			}
		} else if equal, diff := ratio.DeepEqual(actualFilters, actualAlert.GetCondition(), notifyData); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Ratio.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	case alerts.AlertFilters_FILTER_TYPE_UNIQUE_COUNT:
		if uniqueCount := in.UniqueCount; uniqueCount == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "UniqueCount",
			}
		} else if equal, diff := uniqueCount.DeepEqual(actualFilters, actualAlert.GetCondition()); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("UniqueCount.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	case alerts.AlertFilters_FILTER_TYPE_TIME_RELATIVE:
		if timeRelative := in.TimeRelative; timeRelative == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "TimeRelative",
			}
		} else if equal, diff := timeRelative.DeepEqual(actualFilters, actualAlert.GetCondition(), notifyData); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("TimeRelative.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	case alerts.AlertFilters_FILTER_TYPE_METRIC:
		if metric := in.Metric; metric == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Metric",
			}
		} else if equal, diff := metric.DeepEqual(actualFilters, actualAlert.GetCondition(), notifyData); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Metric.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	case alerts.AlertFilters_FILTER_TYPE_TRACING:
		if tracing := in.Tracing; tracing == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "TimeRelative",
			}
		} else if equal, diff := tracing.DeepEqual(actualAlert.GetTracingAlert(), actualAlert.GetCondition(), notifyData); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("TimeRelative.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	case alerts.AlertFilters_FILTER_TYPE_FLOW:
		if flow := in.Flow; flow == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Flow",
			}
		} else if equal, diff := flow.DeepEqual(actualAlert.GetCondition().GetFlow()); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Flow.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	return true, utils.Diff{}
}

type Standard struct {
	// +optional
	Filters *Filters `json:"filters,omitempty"`

	Conditions StandardConditions `json:"conditions,omitempty"`
}

func (in *Standard) DeepEqual(filters *alerts.AlertFilters, condition *alerts.AlertCondition, data *notificationsAlertTypeData) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(condition, data); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Filters.DeepEqual(filters); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Filters.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

func equalSeverities(severities []FiltersLogSeverity, actualSeverities []alerts.AlertFilters_LogSeverity) bool {
	if len(severities) != len(actualSeverities) {
		return false
	}

	valuesSet := make(map[alerts.AlertFilters_LogSeverity]bool, len(severities))
	for _, _a := range actualSeverities {
		valuesSet[_a] = true
	}

	for _, _b := range severities {
		if !valuesSet[alertSchemaFiltersLogSeverityToProtoFiltersLogSeverity[_b]] {
			return false
		}
	}

	return true
}

type Ratio struct {
	Query1Filters Filters `json:"q1Filters,omitempty"`

	Query2Filters RatioQ2Filters `json:"q2Filters,omitempty"`

	Conditions RatioConditions `json:"conditions,omitempty"`
}

type RatioQ2Filters struct {
	// +optional
	Alias string `json:"alias,omitempty"`

	// +optional
	SearchQuery *string `json:"searchQuery,omitempty"`

	// +optional
	Severities []FiltersLogSeverity `json:"severities,omitempty"`

	// +optional
	Applications []string `json:"applications,omitempty"`

	// +optional
	Subsystems []string `json:"subsystems,omitempty"`
}

func (in *Ratio) DeepEqual(filters *alerts.AlertFilters, condition *alerts.AlertCondition, data *notificationsAlertTypeData) (bool, utils.Diff) {
	if equal, diff := in.Query1Filters.DeepEqual(filters); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Q1Filters.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	actualQ2Filters := filters.GetRatioAlerts()[0]
	if equal, diff := in.Query2Filters.DeepEqual(actualQ2Filters); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Q2Filters.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	actualQ2GroupBy := actualQ2Filters.GetGroupBy()
	if equal, diff := in.Conditions.DeepEqual(condition, actualQ2GroupBy, data); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

type NewValue struct {
	// +optional
	Filters *Filters `json:"filters,omitempty"`

	Conditions NewValueConditions `json:"conditions,omitempty"`
}

func (in *NewValue) DeepEqual(filters *alerts.AlertFilters, condition *alerts.AlertCondition_NewValue) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(condition); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Filters.DeepEqual(filters); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Filters.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

type UniqueCount struct {
	// +optional
	Filters *Filters `json:"filters,omitempty"`

	Conditions UniqueCountConditions `json:"conditions,omitempty"`
}

func (in *UniqueCount) DeepEqual(filters *alerts.AlertFilters, condition *alerts.AlertCondition) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(condition); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Filters.DeepEqual(filters); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Filters.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

type TimeRelative struct {
	// +optional
	Filters *Filters `json:"filters,omitempty"`

	Conditions TimeRelativeConditions `json:"conditions,omitempty"`
}

func (in *TimeRelative) DeepEqual(filters *alerts.AlertFilters, condition *alerts.AlertCondition, data *notificationsAlertTypeData) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(condition, data); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}
	if equal, diff := in.Filters.DeepEqual(filters); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Filters.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

type Metric struct {
	// +optional
	Lucene *Lucene `json:"lucene,omitempty"`

	// +optional
	Promql *Promql `json:"promql,omitempty"`
}

func (in *Metric) DeepEqual(filters *alerts.AlertFilters, condition *alerts.AlertCondition, data *notificationsAlertTypeData) (bool, utils.Diff) {
	if promql := in.Promql; promql != nil {
		return promql.DeepEqual(filters, condition, data)
	} else if lucene := in.Lucene; lucene != nil {
		return lucene.DeepEqual(filters, condition, data)
	}

	return false, utils.Diff{}
}

type Lucene struct {
	// +optional
	SearchQuery *string `json:"searchQuery,omitempty"`

	Conditions LuceneConditions `json:"conditions,omitempty"`
}

func (in *Lucene) DeepEqual(filters *alerts.AlertFilters, condition *alerts.AlertCondition, data *notificationsAlertTypeData) (bool, utils.Diff) {
	var conditionParams *alerts.ConditionParameters
	var actualAlertWhen string
	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		actualAlertWhen = "Less"
		conditionParams = condition.LessThan.GetParameters()
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.MoreThan.GetParameters()
		actualAlertWhen = "More"
	}

	if conditionParams.GetMetricAlertPromqlParameters() != nil {
		return false, utils.Diff{
			Name:    "Type",
			Desired: "Lucene",
			Actual:  "Promql",
		}
	}

	if alertWhen := string(in.Conditions.AlertWhen); actualAlertWhen != alertWhen {
		return false, utils.Diff{
			Name:    "Lucene.Conditions.AlertWhen",
			Desired: alertWhen,
			Actual:  actualAlertWhen,
		}
	}

	if searchQuery, actualSearchQuery := in.SearchQuery, filters.GetText(); searchQuery == nil && actualSearchQuery != nil {
		return false, utils.Diff{
			Name:    "Lucene.SearchQuery",
			Desired: searchQuery,
			Actual:  *actualSearchQuery,
		}
	} else if searchQuery, actualSearchQuery := *searchQuery, actualSearchQuery.GetValue(); searchQuery != actualSearchQuery {
		return false, utils.Diff{
			Name:    "Lucene.SearchQuery",
			Desired: searchQuery,
			Actual:  actualSearchQuery,
		}
	}

	if equal, diff := in.Conditions.DeepEqual(conditionParams, data); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Lucene.Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

type Promql struct {
	// +optional
	SearchQuery *string `json:"searchQuery,omitempty"`

	Conditions PromqlConditions `json:"conditions,omitempty"`
}

func (in *Promql) DeepEqual(filters *alerts.AlertFilters, condition *alerts.AlertCondition, data *notificationsAlertTypeData) (bool, utils.Diff) {
	var conditionParams *alerts.ConditionParameters
	var actualAlertWhen string
	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		actualAlertWhen = "Less"
		conditionParams = condition.LessThan.GetParameters()
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.MoreThan.GetParameters()
		actualAlertWhen = "More"
	}

	promqlParams := conditionParams.GetMetricAlertPromqlParameters()
	if promqlParams == nil {
		return false, utils.Diff{
			Name:    "Type",
			Desired: "Promql",
			Actual:  "Lucene",
		}
	}

	if alertWhen := string(in.Conditions.AlertWhen); actualAlertWhen != alertWhen {
		return false, utils.Diff{
			Name:    "Promql.Conditions.AlertWhen",
			Desired: alertWhen,
			Actual:  actualAlertWhen,
		}
	}

	if searchQuery, actualSearchQuery := in.SearchQuery, filters.GetText(); searchQuery == nil && actualSearchQuery != nil {
		return false, utils.Diff{
			Name:    "Lucene.SearchQuery",
			Desired: searchQuery,
			Actual:  *actualSearchQuery,
		}
	} else if searchQuery, actualSearchQuery := *searchQuery, actualSearchQuery.GetValue(); searchQuery != actualSearchQuery {
		return false, utils.Diff{
			Name:    "Lucene.SearchQuery",
			Desired: searchQuery,
			Actual:  actualSearchQuery,
		}
	}

	if equal, diff := in.Conditions.DeepEqual(conditionParams, data); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Promql.Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

type Tracing struct {
	Filters TracingFilters `json:"filters,omitempty"`

	Conditions TracingCondition `json:"conditions,omitempty"`
}

func (in *Tracing) DeepEqual(filters *alerts.TracingAlert, condition *alerts.AlertCondition, data *notificationsAlertTypeData) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(condition, data); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}
	if equal, diff := in.Filters.DeepEqual(filters); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Filters.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

type Flow struct {
	Stages []FlowStage `json:"stages,omitempty"`
}

func (in *Flow) DeepEqual(flow *alerts.FlowCondition) (bool, utils.Diff) {
	if stages, actualStages := in.Stages, flow.Stages; len(stages) != len(actualStages) {
		return false, utils.Diff{
			Name:    "Stages",
			Desired: stages,
			Actual:  actualStages,
		}
	} else {
		for i, stage := range stages {
			if equal, diff := stage.DeepEqual(actualStages[i]); !equal {
				return false, utils.Diff{
					Name:    fmt.Sprintf("Stages.%d.%s", i, diff.Name),
					Desired: diff.Desired,
					Actual:  diff.Actual,
				}
			}
		}
	}

	return true, utils.Diff{}
}

type StandardConditions struct {
	AlertWhen StandardAlertWhen `json:"alertWhen,omitempty"`

	// +optional
	Threshold *int `json:"threshold,omitempty"`

	// +optional
	TimeWindow *TimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

func (in *StandardConditions) DeepEqual(condition *alerts.AlertCondition, data *notificationsAlertTypeData) (bool, utils.Diff) {
	var conditionParams *alerts.ConditionParameters
	switch condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		conditionParams = condition.GetLessThan().GetParameters()
		if alertWhen := in.AlertWhen; alertWhen != "Less" {
			return false, utils.Diff{
				Name:    "AlertWhen",
				Desired: alertWhen,
				Actual:  "Less",
			}
		}
		if threshold, actualThreshold := float64(*(in.Threshold)), conditionParams.GetThreshold().GetValue(); threshold != actualThreshold {
			return false, utils.Diff{
				Name:    "Threshold",
				Desired: threshold,
				Actual:  actualThreshold,
			}
		}
		if timeWindow, actualTimeWindow := alertSchemaTimeWindowToProtoTimeWindow[string(*in.TimeWindow)], conditionParams.GetTimeframe(); timeWindow != actualTimeWindow {
			return false, utils.Diff{
				Name:    "TimeWindow",
				Desired: timeWindow,
				Actual:  actualTimeWindow,
			}
		}
		if manageUndetectedValues, actualManageUndetectedValues := in.ManageUndetectedValues, conditionParams.GetRelatedExtendedData(); manageUndetectedValues == nil && actualManageUndetectedValues != nil {
			return false, utils.Diff{
				Name:    "ManageUndetectedValues",
				Desired: manageUndetectedValues,
				Actual:  *actualManageUndetectedValues,
			}
		} else if equal, diff := manageUndetectedValues.DeepEqual(actualManageUndetectedValues); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("ManageUndetectedValues,%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.GetMoreThan().GetParameters()
		if alertWhen := in.AlertWhen; alertWhen != "More" {
			return false, utils.Diff{
				Name:    "AlertWhen",
				Desired: alertWhen,
				Actual:  "More",
			}
		}
		if threshold, actualThreshold := float64(*(in.Threshold)), conditionParams.GetThreshold().GetValue(); threshold != actualThreshold {
			return false, utils.Diff{
				Name:    "Threshold",
				Desired: threshold,
				Actual:  actualThreshold,
			}
		}
		if timeWindow, actualTimeWindow := alertSchemaTimeWindowToProtoTimeWindow[string(*in.TimeWindow)], conditionParams.GetTimeframe(); timeWindow != actualTimeWindow {
			return false, utils.Diff{
				Name:    "TimeWindow",
				Desired: timeWindow,
				Actual:  actualTimeWindow,
			}
		}
	case *alerts.AlertCondition_MoreThanUsual:
		conditionParams = condition.GetMoreThanUsual().GetParameters()
		if alertWhen := in.AlertWhen; alertWhen != "MoreThanUsual" {
			return false, utils.Diff{
				Name:    "AlertWhen",
				Desired: alertWhen,
				Actual:  "MoreThanUsual",
			}
		}
		if threshold, actualThreshold := float64(*(in.Threshold)), conditionParams.GetThreshold().GetValue(); threshold != actualThreshold {
			return false, utils.Diff{
				Name:    "Threshold",
				Desired: threshold,
				Actual:  actualThreshold,
			}
		}
	case *alerts.AlertCondition_Immediate:
		conditionParams = condition.GetMoreThanUsual().GetParameters()
		if alertWhen := in.AlertWhen; alertWhen != "Immediately" {
			return false, utils.Diff{
				Name:    "AlertWhen",
				Desired: alertWhen,
				Actual:  "Immediately",
			}
		}
	}

	if groupBy, actualGroupBy := in.GroupBy, utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy()); !utils.SlicesWithUniqueValuesEqual(groupBy, actualGroupBy) {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: groupBy,
			Actual:  actualGroupBy,
		}
	}

	data.notifyOnlyOnTriggeredGroupByValues = conditionParams.NotifyGroupByOnlyAlerts
	data.onTriggerAndResolved = conditionParams.NotifyOnResolved
	data.ignoreInfinity = conditionParams.IgnoreInfinity
	return true, utils.Diff{}
}

type RatioConditions struct {
	AlertWhen AlertWhen `json:"alertWhen,omitempty"`

	Ratio resource.Quantity `json:"ratio,omitempty"`

	TimeWindow TimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	GroupByFor *GroupByFor `json:"groupByFor,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

func (in *RatioConditions) DeepEqual(condition *alerts.AlertCondition, actualQ2GroupBy []*wrapperspb.StringValue, data *notificationsAlertTypeData) (bool, utils.Diff) {
	var conditionParams *alerts.ConditionParameters

	switch condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		conditionParams = condition.GetLessThan().GetParameters()
		if alertWhen := in.AlertWhen; alertWhen != "Less" {
			return false, utils.Diff{
				Name:    "AlertWhen",
				Desired: alertWhen,
				Actual:  "Less",
			}
		}
		if manageUndetectedValues, actualManageUndetectedValues := in.ManageUndetectedValues, conditionParams.GetRelatedExtendedData(); manageUndetectedValues == nil && actualManageUndetectedValues != nil {
			return false, utils.Diff{
				Name:    "ManageUndetectedValues",
				Desired: manageUndetectedValues,
				Actual:  *actualManageUndetectedValues,
			}
		} else if equal, diff := manageUndetectedValues.DeepEqual(actualManageUndetectedValues); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("ManageUndetectedValues,%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.GetMoreThan().GetParameters()
		if alertWhen := in.AlertWhen; alertWhen != "More" {
			return false, utils.Diff{
				Name:    "AlertWhen",
				Desired: alertWhen,
				Actual:  "More",
			}
		}
	}

	if threshold, actualThreshold := in.Ratio.AsApproximateFloat64(), conditionParams.GetThreshold().GetValue(); threshold != actualThreshold {
		return false, utils.Diff{
			Name:    "Threshold",
			Desired: threshold,
			Actual:  actualThreshold,
		}
	}
	if timeWindow, actualTimeWindow := alertSchemaTimeWindowToProtoTimeWindow[string(in.TimeWindow)], conditionParams.GetTimeframe(); timeWindow != actualTimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: timeWindow,
			Actual:  actualTimeWindow,
		}
	}

	if groupByFor := in.GroupByFor; groupByFor != nil && *groupByFor == "Q1" || *groupByFor == "Both" {
		if groupBy, actualGroupBy := in.GroupBy, utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy()); !utils.SlicesWithUniqueValuesEqual(groupBy, actualGroupBy) {
			return false, utils.Diff{
				Name:    "GroupBy (Q1)",
				Desired: groupBy,
				Actual:  actualGroupBy,
			}
		}
	}

	if groupByFor := in.GroupByFor; groupByFor != nil && *groupByFor == "Q2" || *groupByFor == "Both" {
		if groupBy, actualGroupBy := in.GroupBy, utils.WrappedStringSliceToStringSlice(actualQ2GroupBy); !utils.SlicesWithUniqueValuesEqual(groupBy, actualGroupBy) {
			return false, utils.Diff{
				Name:    "GroupBy (Q2)",
				Desired: groupBy,
				Actual:  actualGroupBy,
			}
		}
	}

	data.notifyOnlyOnTriggeredGroupByValues = conditionParams.NotifyGroupByOnlyAlerts
	data.onTriggerAndResolved = conditionParams.NotifyOnResolved
	data.ignoreInfinity = conditionParams.IgnoreInfinity

	return true, utils.Diff{}
}

type NewValueConditions struct {
	Key string `json:"key,omitempty"`

	TimeWindow NewValueTimeWindow `json:"timeWindow,omitempty"`
}

func (in *NewValueConditions) DeepEqual(condition *alerts.AlertCondition_NewValue) (bool, utils.Diff) {
	conditionParams := condition.NewValue.GetParameters()

	if key, actualKey := in.Key, conditionParams.GetCardinalityFields()[0].GetValue(); key != actualKey {
		return false, utils.Diff{
			Name:    "Key",
			Desired: key,
			Actual:  actualKey,
		}
	}
	if timeWindow, actualTimeWindow := alertSchemaTimeWindowToProtoTimeWindow[string(in.TimeWindow)], conditionParams.GetTimeframe(); timeWindow != actualTimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: timeWindow,
			Actual:  actualTimeWindow,
		}
	}
	return true, utils.Diff{}
}

type UniqueCountConditions struct {
	Key string `json:"key,omitempty"`

	// +kubebuilder:validation:Minimum:=1
	MaxUniqueValues int `json:"maxUniqueValues,omitempty"`

	TimeWindow UniqueValueTimeWindow `json:"timeWindow,omitempty"`

	GroupBy *string `json:"groupBy,omitempty"`

	// +kubebuilder:validation:Minimum:=1
	MaxUniqueValuesForGroupBy *int `json:"maxUniqueValuesForGroupBy,omitempty"`
}

func (in *UniqueCountConditions) DeepEqual(condition *alerts.AlertCondition) (bool, utils.Diff) {
	conditionParams := condition.GetUniqueCount().GetParameters()

	if key, actualKey := in.Key, conditionParams.GetCardinalityFields()[0].GetValue(); key != actualKey {
		return false, utils.Diff{
			Name:    "Key",
			Desired: key,
			Actual:  actualKey,
		}
	}
	if timeWindow, actualTimeWindow := alertSchemaTimeWindowToProtoTimeWindow[string(in.TimeWindow)], conditionParams.GetTimeframe(); timeWindow != actualTimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: timeWindow,
			Actual:  actualTimeWindow,
		}
	}
	if maxUniqueValues, actualMaxUniqueValues := in.MaxUniqueValues, int(conditionParams.GetMaxUniqueCountValuesForGroupByKey().GetValue()); maxUniqueValues != actualMaxUniqueValues {
		return false, utils.Diff{
			Name:    "UniqueValues",
			Desired: maxUniqueValues,
			Actual:  actualMaxUniqueValues,
		}
	}
	if groupBy, actualGroupBys := in.GroupBy, conditionParams.GetGroupBy(); groupBy != nil && len(actualGroupBys) < 1 {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: *groupBy,
			Actual:  actualGroupBys,
		}
	} else if actualGroupBy := actualGroupBys[0].GetValue(); groupBy == nil {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: groupBy,
			Actual:  actualGroupBy,
		}
	} else if groupBy := *groupBy; actualGroupBy != groupBy {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: groupBy,
			Actual:  actualGroupBys,
		}
	}
	if maxUniqueValuesForGroupBy, actualMaxUniqueValuesForGroupBy := in.MaxUniqueValuesForGroupBy, conditionParams.GetMaxUniqueCountValuesForGroupByKey(); maxUniqueValuesForGroupBy == nil && actualMaxUniqueValuesForGroupBy != nil {
		return false, utils.Diff{
			Name:    "MaxUniqueValuesForGroupBy",
			Desired: maxUniqueValuesForGroupBy,
			Actual:  actualMaxUniqueValuesForGroupBy.GetValue(),
		}
	} else if maxUniqueValuesForGroupBy, actualMaxUniqueValuesForGroupBy := *maxUniqueValuesForGroupBy, int(actualMaxUniqueValuesForGroupBy.GetValue()); maxUniqueValuesForGroupBy != actualMaxUniqueValuesForGroupBy {
		return false, utils.Diff{
			Name:    "MaxUniqueValuesForGroupBy",
			Desired: maxUniqueValuesForGroupBy,
			Actual:  actualMaxUniqueValuesForGroupBy,
		}
	}
	return true, utils.Diff{}
}

type TimeRelativeConditions struct {
	AlertWhen AlertWhen `json:"alertWhen,omitempty"`

	Threshold resource.Quantity `json:"threshold,omitempty"`

	TimeWindow RelativeTimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

func (in *TimeRelativeConditions) DeepEqual(condition *alerts.AlertCondition, data *notificationsAlertTypeData) (bool, utils.Diff) {
	var conditionParams *alerts.ConditionParameters

	switch condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		if alertWhen := in.AlertWhen; alertWhen != "Less" {
			return false, utils.Diff{
				Name:    "AlertWhen",
				Desired: alertWhen,
				Actual:  "Less",
			}
		}
		conditionParams = condition.GetLessThan().GetParameters()
		if manageUndetectedValues, actualManageUndetectedValues := in.ManageUndetectedValues, conditionParams.GetRelatedExtendedData(); manageUndetectedValues == nil && actualManageUndetectedValues != nil {
			return false, utils.Diff{
				Name:    "ManageUndetectedValues",
				Desired: manageUndetectedValues,
				Actual:  *actualManageUndetectedValues,
			}
		} else if equal, diff := manageUndetectedValues.DeepEqual(actualManageUndetectedValues); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("ManageUndetectedValues,%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	case *alerts.AlertCondition_MoreThan:
		if alertWhen := in.AlertWhen; alertWhen != "More" {
			return false, utils.Diff{
				Name:    "AlertWhen",
				Desired: alertWhen,
				Actual:  "More",
			}
		}
		conditionParams = condition.GetMoreThan().GetParameters()
	}

	if threshold, actualThreshold := in.Threshold.AsApproximateFloat64(), conditionParams.GetThreshold().GetValue(); threshold != actualThreshold {
		return false, utils.Diff{
			Name:    "Threshold",
			Desired: threshold,
			Actual:  actualThreshold,
		}
	}

	relativeTimeWindow := alertSchemaRelativeTimeFrameToProtoTimeFrameAndRelativeTimeFrame[in.TimeWindow]
	actualRelativeTimeWindow := protoTimeFrameAndRelativeTimeFrame{timeFrame: conditionParams.GetTimeframe(), relativeTimeFrame: conditionParams.GetRelativeTimeframe()}
	if relativeTimeWindow != actualRelativeTimeWindow {
		return false, utils.Diff{
			Name:    "RelativeTimeWindow",
			Desired: relativeTimeWindow,
			Actual:  actualRelativeTimeWindow,
		}
	}

	if groupBy, actualGroupBy := in.GroupBy, utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy()); !utils.SlicesWithUniqueValuesEqual(groupBy, actualGroupBy) {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: groupBy,
			Actual:  actualGroupBy,
		}
	}

	data.notifyOnlyOnTriggeredGroupByValues = conditionParams.NotifyGroupByOnlyAlerts
	data.onTriggerAndResolved = conditionParams.NotifyOnResolved
	data.ignoreInfinity = conditionParams.IgnoreInfinity
	return true, utils.Diff{}
}

// +kubebuilder:validation:Enum=Avg;Min;Max;Sum;Count;Percentile;
type ArithmeticOperator string

type LuceneConditions struct {
	MetricField string `json:"metricField,omitempty"`

	ArithmeticOperator ArithmeticOperator `json:"arithmeticOperator,omitempty"`

	// +optional
	ArithmeticOperatorModifier *int `json:"arithmeticOperatorModifier,omitempty"`

	AlertWhen AlertWhen `json:"alertWhen,omitempty"`

	Threshold resource.Quantity `json:"threshold,omitempty"`

	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:MultipleOf:=10
	SampleThresholdPercentage int `json:"sampleThresholdPercentage,omitempty"`

	TimeWindow MetricTimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	//+kubebuilder:default=false
	ReplaceMissingValueWithZero bool `json:"replaceMissingValueWithZero,omitempty"`

	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:MultipleOf:=10
	MinNonNullValuesPercentage *int `json:"minNonNullValuesPercentage,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

func (in *LuceneConditions) DeepEqual(conditionParams *alerts.ConditionParameters, data *notificationsAlertTypeData) (bool, utils.Diff) {
	if threshold, actualThreshold := in.Threshold.AsApproximateFloat64(), conditionParams.Threshold.GetValue(); threshold != actualThreshold {
		return false, utils.Diff{
			Name:    "Threshold",
			Desired: threshold,
			Actual:  actualThreshold,
		}
	}

	if groupBy, actualGroupBy := in.GroupBy, utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy()); !utils.SlicesWithUniqueValuesEqual(groupBy, actualGroupBy) {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: groupBy,
			Actual:  actualGroupBy,
		}
	}

	if timeWindow, actualTimeWindow := alertSchemaTimeWindowToProtoTimeWindow[string(in.TimeWindow)], conditionParams.GetTimeframe(); timeWindow != actualTimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: timeWindow,
			Actual:  actualTimeWindow,
		}
	}

	metricParams := conditionParams.GetMetricAlertParameters()

	if metricField, actualMetricField := in.MetricField, metricParams.MetricField.GetValue(); metricField != actualMetricField {
		return false, utils.Diff{
			Name:    "MetricField",
			Desired: metricField,
			Actual:  actualMetricField,
		}
	}

	if arithmeticOperator, actualArithmeticOperator := alertSchemaArithmeticOperatorToProtoArithmeticOperator[in.ArithmeticOperator], metricParams.GetArithmeticOperator(); arithmeticOperator != actualArithmeticOperator {
		return false, utils.Diff{
			Name:    "ArithmeticOperator",
			Desired: arithmeticOperator,
			Actual:  actualArithmeticOperator,
		}
	}

	if arithmeticOperatorModifier, actualArithmeticOperatorModifier := in.ArithmeticOperatorModifier, metricParams.ArithmeticOperatorModifier; arithmeticOperatorModifier == nil && actualArithmeticOperatorModifier != nil {
		return false, utils.Diff{
			Name:    "ArithmeticOperatorModifier",
			Desired: arithmeticOperatorModifier,
			Actual:  *actualArithmeticOperatorModifier,
		}
	} else if arithmeticOperatorModifier, actualArithmeticOperatorModifier := *arithmeticOperatorModifier, int(actualArithmeticOperatorModifier.GetValue()); arithmeticOperatorModifier != actualArithmeticOperatorModifier {
		return false, utils.Diff{
			Name:    "ArithmeticOperatorModifier",
			Desired: arithmeticOperatorModifier,
			Actual:  actualArithmeticOperatorModifier,
		}
	}

	if sampleThresholdPercentage, actualSampleThresholdPercentage := in.SampleThresholdPercentage, int(metricParams.SampleThresholdPercentage.GetValue()); sampleThresholdPercentage != actualSampleThresholdPercentage {
		return false, utils.Diff{
			Name:    "SampleThresholdPercentage",
			Desired: sampleThresholdPercentage,
			Actual:  actualSampleThresholdPercentage,
		}
	}

	if replaceMissingValueWithZero, actualReplaceMissingValueWithZero := in.ReplaceMissingValueWithZero, metricParams.GetSwapNullValues().GetValue(); replaceMissingValueWithZero != actualReplaceMissingValueWithZero {
		return false, utils.Diff{
			Name:    "MissingValueWithZero",
			Desired: replaceMissingValueWithZero,
			Actual:  actualReplaceMissingValueWithZero,
		}
	}

	if minNonNullValuesPercentage, actualMinNonNullValuesPercentage := in.MinNonNullValuesPercentage, metricParams.GetNonNullPercentage(); minNonNullValuesPercentage == nil && actualMinNonNullValuesPercentage != nil {
		return false, utils.Diff{
			Name:    "MinNonNullValuesPercentage",
			Desired: minNonNullValuesPercentage,
			Actual:  actualMinNonNullValuesPercentage.GetValue(),
		}
	} else if minNonNullValuesPercentage, actualMinNonNullValuesPercentage := *minNonNullValuesPercentage, int(actualMinNonNullValuesPercentage.GetValue()); minNonNullValuesPercentage != actualMinNonNullValuesPercentage {
		return false, utils.Diff{
			Name:    "MinNonNullValuesPercentage",
			Desired: minNonNullValuesPercentage,
			Actual:  actualMinNonNullValuesPercentage,
		}
	}

	if manageUndetectedValues, actualManageUndetectedValues := in.ManageUndetectedValues, conditionParams.GetRelatedExtendedData(); manageUndetectedValues == nil && actualManageUndetectedValues != nil {
		return false, utils.Diff{
			Name:    "ManageUndetectedValues",
			Desired: manageUndetectedValues,
			Actual:  *actualManageUndetectedValues,
		}
	} else if equal, diff := manageUndetectedValues.DeepEqual(actualManageUndetectedValues); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("ManageUndetectedValues,%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	data.notifyOnlyOnTriggeredGroupByValues = conditionParams.NotifyGroupByOnlyAlerts
	data.onTriggerAndResolved = conditionParams.NotifyOnResolved

	return true, utils.Diff{}
}

type PromqlConditions struct {
	AlertWhen AlertWhen `json:"alertWhen,omitempty"`

	Threshold resource.Quantity `json:"threshold,omitempty"`

	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:MultipleOf:=10
	SampleThresholdPercentage int `json:"sampleThresholdPercentage,omitempty"`

	TimeWindow MetricTimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	ReplaceMissingValueWithZero bool `json:"replaceMissingValueWithZero,omitempty"`

	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:MultipleOf:=10
	MinNonNullValuesPercentage *int `json:"minNonNullValuesPercentage,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

func (in *PromqlConditions) DeepEqual(conditionParams *alerts.ConditionParameters, data *notificationsAlertTypeData) (bool, utils.Diff) {
	if threshold, actualThreshold := in.Threshold.AsApproximateFloat64(), conditionParams.Threshold.GetValue(); threshold != actualThreshold {
		return false, utils.Diff{
			Name:    "Threshold",
			Desired: threshold,
			Actual:  actualThreshold,
		}
	}

	if groupBy, actualGroupBy := in.GroupBy, utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy()); !utils.SlicesWithUniqueValuesEqual(groupBy, actualGroupBy) {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: groupBy,
			Actual:  actualGroupBy,
		}
	}

	if timeWindow, actualTimeWindow := alertSchemaTimeWindowToProtoTimeWindow[string(in.TimeWindow)], conditionParams.GetTimeframe(); timeWindow != actualTimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: timeWindow,
			Actual:  actualTimeWindow,
		}
	}

	metricParams := conditionParams.GetMetricAlertParameters()

	if sampleThresholdPercentage, actualSampleThresholdPercentage := in.SampleThresholdPercentage, int(metricParams.SampleThresholdPercentage.GetValue()); sampleThresholdPercentage != actualSampleThresholdPercentage {
		return false, utils.Diff{
			Name:    "SampleThresholdPercentage",
			Desired: sampleThresholdPercentage,
			Actual:  actualSampleThresholdPercentage,
		}
	}

	if minNonNullValuesPercentage, actualMinNonNullValuesPercentage := in.MinNonNullValuesPercentage, metricParams.GetNonNullPercentage(); minNonNullValuesPercentage == nil && actualMinNonNullValuesPercentage != nil {
		return false, utils.Diff{
			Name:    "MinNonNullValuesPercentage",
			Desired: minNonNullValuesPercentage,
			Actual:  actualMinNonNullValuesPercentage.GetValue(),
		}
	} else if minNonNullValuesPercentage, actualMinNonNullValuesPercentage := *minNonNullValuesPercentage, int(actualMinNonNullValuesPercentage.GetValue()); minNonNullValuesPercentage != actualMinNonNullValuesPercentage {
		return false, utils.Diff{
			Name:    "MinNonNullValuesPercentage",
			Desired: minNonNullValuesPercentage,
			Actual:  actualMinNonNullValuesPercentage,
		}
	}

	if manageUndetectedValues, actualManageUndetectedValues := in.ManageUndetectedValues, conditionParams.GetRelatedExtendedData(); manageUndetectedValues == nil && actualManageUndetectedValues != nil {
		return false, utils.Diff{
			Name:    "ManageUndetectedValues",
			Desired: manageUndetectedValues,
			Actual:  *actualManageUndetectedValues,
		}
	} else if equal, diff := manageUndetectedValues.DeepEqual(actualManageUndetectedValues); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("ManageUndetectedValues,%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	data.notifyOnlyOnTriggeredGroupByValues = conditionParams.NotifyGroupByOnlyAlerts
	data.onTriggerAndResolved = conditionParams.NotifyOnResolved

	return true, utils.Diff{}
}

type TracingCondition struct {
	AlertWhen TracingAlertWhen `json:"alertWhen,omitempty"`

	// +optional
	Threshold *resource.Quantity `json:"threshold,omitempty"`

	// +optional
	TimeWindow *TimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`
}

func (in *TracingCondition) DeepEqual(condition *alerts.AlertCondition, data *notificationsAlertTypeData) (bool, utils.Diff) {
	var conditionParams *alerts.ConditionParameters
	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_MoreThan:
		if alertWhen, actualAlertWhen := in.AlertWhen, "More"; string(alertWhen) != actualAlertWhen {
			return false, utils.Diff{
				Name:    "AlertWhen",
				Desired: alertWhen,
				Actual:  actualAlertWhen,
			}
		}
		conditionParams = condition.MoreThan.GetParameters()
		if threshold, actualThreshold := in.Threshold.AsApproximateFloat64(), conditionParams.GetThreshold().GetValue(); threshold != actualThreshold {
			return false, utils.Diff{
				Name:    "Threshold",
				Desired: threshold,
				Actual:  actualThreshold,
			}
		}
		if timeWindow, actualTimeWindow := in.TimeWindow, conditionParams.GetTimeframe(); timeWindow == nil {
			return false, utils.Diff{
				Name:    "TimeWindow",
				Desired: timeWindow,
				Actual:  actualTimeWindow,
			}
		} else if timeWindow := alertSchemaTimeWindowToProtoTimeWindow[string(*timeWindow)]; timeWindow != actualTimeWindow {
			return false, utils.Diff{
				Name:    "TimeWindow",
				Desired: timeWindow,
				Actual:  actualTimeWindow,
			}
		}
		if groupBy, actualGroupBy := in.GroupBy, utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy()); !utils.SlicesWithUniqueValuesEqual(groupBy, actualGroupBy) {
			return false, utils.Diff{
				Name:    "GroupBy",
				Desired: groupBy,
				Actual:  actualGroupBy,
			}
		}
		data.onTriggerAndResolved = conditionParams.GetNotifyOnResolved()
	case *alerts.AlertCondition_Immediate:
		if alertWhen, actualAlertWhen := in.AlertWhen, "Immediately"; string(alertWhen) != actualAlertWhen {
			return false, utils.Diff{
				Name:    "AlertWhen",
				Desired: alertWhen,
				Actual:  actualAlertWhen,
			}
		}
	}

	return true, utils.Diff{}
}

// +kubebuilder:validation:Enum=Never;FiveMinutes;TenMinutes;Hour;TwoHours;SixHours;TwelveHours;TwentyFourHours
type AutoRetireRatio string

// +kubebuilder:validation:Enum=More;Less
type AlertWhen string

// +kubebuilder:validation:Enum=More;Less;Immediately;MoreThanUsual
type StandardAlertWhen string

// +kubebuilder:validation:Enum=More;Immediately
type TracingAlertWhen string

// +kubebuilder:validation:Enum=Q1;Q2;Both
type GroupByFor string

// +kubebuilder:validation:Enum=FiveMinutes;TenMinutes;FifteenMinutes;TwentyMinutes;ThirtyMinutes;Hour;TwoHours;FourHours;SixHours;TwelveHours;TwentyFourHours;ThirtySixHours
type TimeWindow string

// +kubebuilder:validation:Enum=TwelveHours;TwentyFourHours;FortyEightHours;SeventTwoHours;Week;Month;TwoMonths;ThreeMonths;
type NewValueTimeWindow string

// +kubebuilder:validation:Enum=Minute;FiveMinutes;TenMinutes;FifteenMinutes;TwentyMinutes;ThirtyMinutes;Hour;TwoHours;FourHours;SixHours;TwelveHours;TwentyFourHours;ThirtySixHours
type UniqueValueTimeWindow string

// +kubebuilder:validation:Enum=Minute;FiveMinutes;TenMinutes;FifteenMinutes;TwentyMinutes;ThirtyMinutes;Hour;TwoHours;FourHours;SixHours;TwelveHours;TwentyFourHours
type MetricTimeWindow string

// +kubebuilder:validation:Enum=PreviousHour;SameHourYesterday;SameHourLastWeek;Yesterday;SameDayLastWeek;SameDayLastMonth;
type RelativeTimeWindow string

type Filters struct {
	// +optional
	SearchQuery *string `json:"searchQuery,omitempty"`

	// +optional
	Severities []FiltersLogSeverity `json:"severities,omitempty"`

	// +optional
	Applications []string `json:"applications,omitempty"`

	// +optional
	Subsystems []string `json:"subsystems,omitempty"`

	// +optional
	Categories []string `json:"categories,omitempty"`

	// +optional
	Computers []string `json:"computers,omitempty"`

	// +optional
	Classes []string `json:"classes,omitempty"`

	// +optional
	Methods []string `json:"methods,omitempty"`

	// +optional
	IPs []string `json:"ips,omitempty"`

	// +optional
	Alias *string `json:"alias,omitempty"`
}

func (in *Filters) DeepEqual(filters *alerts.AlertFilters) (bool, utils.Diff) {
	if searchQuery, actualSearchQuery := in.SearchQuery, filters.GetText(); searchQuery == nil && actualSearchQuery != nil {
		return false, utils.Diff{
			Name:    "SearchQuery",
			Desired: searchQuery,
			Actual:  *actualSearchQuery,
		}
	} else if searchQuery, actualSearchQuery := *searchQuery, actualSearchQuery.GetValue(); searchQuery != actualSearchQuery {
		return false, utils.Diff{
			Name:    "SearchQuery",
			Desired: searchQuery,
			Actual:  actualSearchQuery,
		}
	}

	if alias, actualAlias := in.Alias, filters.GetAlias(); alias == nil && actualAlias != nil {
		return false, utils.Diff{
			Name:    "Alias",
			Desired: alias,
			Actual:  *actualAlias,
		}
	} else if alias, actualAlias := *alias, actualAlias.GetValue(); alias != actualAlias {
		return false, utils.Diff{
			Name:    "Alias",
			Desired: alias,
			Actual:  actualAlias,
		}
	}

	if severities, actualSeverities := in.Severities, filters.Severities; !equalSeverities(severities, actualSeverities) {
		return false, utils.Diff{
			Name:    "Severities",
			Desired: severities,
			Actual:  actualSeverities,
		}
	}

	metadata := filters.Metadata

	if applications, actualApplications := in.Applications, utils.WrappedStringSliceToStringSlice(metadata.Applications); !utils.SlicesWithUniqueValuesEqual(applications, actualApplications) {
		return false, utils.Diff{
			Name:    "Application",
			Desired: applications,
			Actual:  actualApplications,
		}
	}

	if subsystems, actualSubsystems := in.Subsystems, utils.WrappedStringSliceToStringSlice(metadata.Subsystems); !utils.SlicesWithUniqueValuesEqual(subsystems, actualSubsystems) {
		return false, utils.Diff{
			Name:    "Subsystems",
			Desired: subsystems,
			Actual:  actualSubsystems,
		}
	}

	if categories, actualCategories := in.Categories, utils.WrappedStringSliceToStringSlice(metadata.Categories); !utils.SlicesWithUniqueValuesEqual(categories, actualCategories) {
		return false, utils.Diff{
			Name:    "Categories",
			Desired: categories,
			Actual:  actualCategories,
		}
	}

	if computers, actualComputers := in.Computers, utils.WrappedStringSliceToStringSlice(metadata.Computers); !utils.SlicesWithUniqueValuesEqual(computers, actualComputers) {
		return false, utils.Diff{
			Name:    "Computers",
			Desired: computers,
			Actual:  actualComputers,
		}
	}

	if classes, actualClasses := in.Classes, utils.WrappedStringSliceToStringSlice(metadata.Classes); !utils.SlicesWithUniqueValuesEqual(classes, actualClasses) {
		return false, utils.Diff{
			Name:    "Classes",
			Desired: classes,
			Actual:  actualClasses,
		}
	}

	if methods, actualMethods := in.Methods, utils.WrappedStringSliceToStringSlice(metadata.Methods); !utils.SlicesWithUniqueValuesEqual(methods, actualMethods) {
		return false, utils.Diff{
			Name:    "Methods",
			Desired: methods,
			Actual:  actualMethods,
		}
	}

	if IPs, actualIPs := in.IPs, utils.WrappedStringSliceToStringSlice(metadata.IpAddresses); !utils.SlicesWithUniqueValuesEqual(IPs, actualIPs) {
		return false, utils.Diff{
			Name:    "IPs",
			Desired: IPs,
			Actual:  actualIPs,
		}
	}

	return true, utils.Diff{}
}

func (in *RatioQ2Filters) DeepEqual(filters *alerts.AlertFilters_RatioAlert) (bool, utils.Diff) {
	if alias, actualAlias := in.Alias, filters.GetAlias().GetValue(); actualAlias != alias {
		return false, utils.Diff{
			Name:    "Alias",
			Desired: alias,
			Actual:  actualAlias,
		}
	}

	if searchQuery, actualSearchQuery := in.SearchQuery, filters.GetText(); searchQuery == nil && actualSearchQuery != nil {
		return false, utils.Diff{
			Name:    "SearchQuery",
			Desired: searchQuery,
			Actual:  *actualSearchQuery,
		}
	} else if searchQuery, actualSearchQuery := *searchQuery, actualSearchQuery.GetValue(); searchQuery != actualSearchQuery {
		return false, utils.Diff{
			Name:    "SearchQuery",
			Desired: searchQuery,
			Actual:  actualSearchQuery,
		}
	}

	if severities, actualSeverities := in.Severities, filters.Severities; !equalSeverities(severities, actualSeverities) {
		return false, utils.Diff{
			Name:    "Severities",
			Desired: severities,
			Actual:  actualSeverities,
		}
	}

	if applications, actualApplications := in.Applications, utils.WrappedStringSliceToStringSlice(filters.Applications); !utils.SlicesWithUniqueValuesEqual(applications, actualApplications) {
		return false, utils.Diff{
			Name:    "Application",
			Desired: applications,
			Actual:  actualApplications,
		}
	}

	if subsystems, actualSubsystems := in.Subsystems, utils.WrappedStringSliceToStringSlice(filters.Subsystems); !utils.SlicesWithUniqueValuesEqual(subsystems, actualSubsystems) {
		return false, utils.Diff{
			Name:    "Subsystems",
			Desired: subsystems,
			Actual:  actualSubsystems,
		}
	}

	return true, utils.Diff{}
}

// +kubebuilder:validation:Enum=Debug;Verbose;Info;Warning;Critical;Error;
type FiltersLogSeverity string

type TracingFilters struct {
	LatencyThresholdMS int `json:"latencyThresholdMS,omitempty"`

	// +optional
	TagFilters []TagFilter `json:"tagFilters,omitempty"`

	// +optional
	FieldFilters []FieldFilter `json:"fieldFilters,omitempty"`
}

func (in *TracingFilters) DeepEqual(filters *alerts.TracingAlert) (bool, utils.Diff) {
	if latencyThresholdMS, actualLatencyThresholdMS := in.LatencyThresholdMS, int(filters.ConditionLatency); latencyThresholdMS != actualLatencyThresholdMS {
		return false, utils.Diff{
			Name:    "LatencyThresholdMS",
			Desired: latencyThresholdMS,
			Actual:  actualLatencyThresholdMS,
		}
	}

	//if tagFilters, actualTagFilters := in.TagFilters, filters.TagFilters; latencyThresholdMS != actualLatencyThresholdMS{
	//	return false, utils.Diff{
	//		Name: "LatencyThresholdMS",
	//		Desired: latencyThresholdMS,
	//		Actual: actualLatencyThresholdMS,
	//	}
	//}
	//
	//if latencyThresholdMS, actualLatencyThresholdMS := in.LatencyThresholdMS, int(filters.ConditionLatency); latencyThresholdMS != actualLatencyThresholdMS{
	//	return false, utils.Diff{
	//		Name: "LatencyThresholdMS",
	//		Desired: latencyThresholdMS,
	//		Actual: actualLatencyThresholdMS,
	//	}
	//}

	return true, utils.Diff{}
}

type TagFilter struct {
	Field   string          `json:"field,omitempty"`
	Filters []TracingFilter `json:"filters,omitempty"`
}

type FieldFilter struct {
	Field   FieldFilterType `json:"field,omitempty"`
	Filters []TracingFilter `json:"filters,omitempty"`
}

type TracingFilter struct {
	Values   []string       `json:"values,omitempty"`
	Operator FilterOperator `json:"operator,omitempty"`
}

// +kubebuilder:validation:Enum=Equals;Contains;StartWith;EndWith;
type FilterOperator string

// +kubebuilder:validation:Enum=Application;Subsystem;Service;
type FieldFilterType string

type ManageUndetectedValues struct {
	//+kubebuilder:default=true
	EnableTriggeringOnUndetectedValues bool `json:"enableTriggeringOnUndetectedValues,omitempty"`

	//+kubebuilder:default=Never
	AutoRetireRatio *AutoRetireRatio `json:"autoRetireRatio,omitempty"`
}

func (in *ManageUndetectedValues) DeepEqual(manageUndetectedValues *alerts.RelatedExtendedData) (bool, utils.Diff) {
	if enableTriggeringOnUndetectedValues, actualEnableTriggeringOnUndetectedValues := in.EnableTriggeringOnUndetectedValues, manageUndetectedValues.GetShouldTriggerDeadman().GetValue(); enableTriggeringOnUndetectedValues != actualEnableTriggeringOnUndetectedValues {
		return false, utils.Diff{
			Name:    "EnableTriggeringOnUndetectedValues",
			Desired: enableTriggeringOnUndetectedValues,
			Actual:  actualEnableTriggeringOnUndetectedValues,
		}
	}
	if autoRetireRatio, actualAutoRetireRatio := alertSchemaAutoRetireRatioToProtoAutoRetireRatio[*in.AutoRetireRatio], manageUndetectedValues.GetCleanupDeadmanDuration(); autoRetireRatio != actualAutoRetireRatio {
		return false, utils.Diff{
			Name:    "AutoRetireRatio",
			Desired: autoRetireRatio,
			Actual:  actualAutoRetireRatio,
		}
	}

	return true, utils.Diff{}
}

type FlowStage struct {
	// +optional
	TimeWindow *FlowStageTimeFrame `json:"timeWindow,omitempty"`

	Groups []FlowStageGroup `json:"groups,omitempty"`
}

func (in *FlowStage) DeepEqual(stage *alerts.FlowStage) (bool, utils.Diff) {
	if groups, actualGroups := in.Groups, stage.Groups; len(groups) != len(actualGroups) {
		return false, utils.Diff{
			Name:    "Groups",
			Desired: groups,
			Actual:  actualGroups,
		}
	} else {
		for i, group := range groups {
			if equal, diff := group.DeepEqual(actualGroups[i]); !equal {
				return false, utils.Diff{
					Name:    fmt.Sprintf("Groups.%d.%s", i, diff.Name),
					Desired: diff.Desired,
					Actual:  diff.Actual,
				}
			}
		}
	}

	if timeFrame, actualTimeFrame := in.TimeWindow, stage.GetTimeframe(); timeFrame == nil && actualTimeFrame != nil {
		return false, utils.Diff{
			Name:    "TimeFrame",
			Desired: timeFrame,
			Actual:  *actualTimeFrame,
		}
	} else if timeFrame != nil && actualTimeFrame == nil {
		return false, utils.Diff{
			Name:    "TimeFrame",
			Desired: *timeFrame,
			Actual:  actualTimeFrame,
		}
	} else if timeFrame != nil && actualTimeFrame != nil {
		if timeMS, actualTimeMS := expandTimeToMS(*timeFrame), int(actualTimeFrame.GetMs().GetValue()); timeMS != actualTimeMS {
			return false, utils.Diff{
				Name:    "TimeFrameMS",
				Desired: timeMS,
				Actual:  actualTimeMS,
			}
		}
	}

	return true, utils.Diff{}
}

type FlowStageTimeFrame struct {
	// +optional
	Hours int `json:"hours,omitempty"`

	// +optional
	Minutes int `json:"minutes,omitempty"`

	// +optional
	Seconds int `json:"seconds,omitempty"`
}

type FlowStageGroup struct {
	SubgroupAlerts SubgroupAlerts `json:"subgroupAlerts,omitempty"`

	Operator FlowOperator `json:"operator,omitempty"`
}

func (in *FlowStageGroup) DeepEqual(group *alerts.FlowGroup) (bool, utils.Diff) {
	if alerts, actualAlerts := in.SubgroupAlerts.Alerts, group.Alerts.Values; len(alerts) != len(actualAlerts) {
		return false, utils.Diff{
			Name:    "Alerts",
			Desired: alerts,
			Actual:  actualAlerts,
		}
	} else {
		for i, alert := range alerts {
			if equal, diff := alert.DeepEqual(actualAlerts[i]); !equal {
				return false, utils.Diff{
					Name:    fmt.Sprintf("Alerts.%d.%s", i, diff.Name),
					Desired: diff.Desired,
					Actual:  diff.Actual,
				}
			}
		}
	}

	if operator, actualOperator := alertSchemaFlowOperatorToProtoFlowOperator[in.Operator], group.NextOp; operator != actualOperator {
		return false, utils.Diff{
			Name:    "Operator",
			Desired: operator,
			Actual:  actualOperator,
		}
	}

	return true, utils.Diff{}
}

type SubgroupAlerts struct {
	Operator FlowOperator `json:"operator,omitempty"`

	Alerts []InnerFlowAlert `json:"alerts,omitempty"`
}

type InnerFlowAlert struct {
	// +kubebuilder:default=false
	Not bool `json:"not,omitempty"`

	// +optional
	UserAlertId string `json:"userAlertId,omitempty"`
}

func (in *InnerFlowAlert) DeepEqual(alert *alerts.FlowAlert) (bool, utils.Diff) {
	if not, actualNot := in.Not, alert.GetNot().GetValue(); not != actualNot {
		return false, utils.Diff{
			Name:    "Not",
			Desired: not,
			Actual:  actualNot,
		}
	}
	if id, actualID := in.UserAlertId, alert.GetId().GetValue(); id != actualID {
		return false, utils.Diff{
			Name:    "ID",
			Desired: id,
			Actual:  actualID,
		}
	}

	return true, utils.Diff{}
}

// +kubebuilder:validation:Enum=And;Or
type FlowOperator string

// AlertStatus defines the observed state of Alert
type AlertStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ID *string `json:"id"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Alert is the Schema for the alerts API
type Alert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertSpec   `json:"spec,omitempty"`
	Status AlertStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AlertList contains a list of Alert
type AlertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Alert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Alert{}, &AlertList{})
}
