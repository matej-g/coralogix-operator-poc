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
	"reflect"
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
	AlertSchemaSeverityToProtoSeverity = map[AlertSeverity]alerts.AlertSeverity{
		AlertSeverityInfo:     alerts.AlertSeverity_ALERT_SEVERITY_INFO_OR_UNSPECIFIED,
		AlertSeverityWarning:  alerts.AlertSeverity_ALERT_SEVERITY_WARNING,
		AlertSeverityCritical: alerts.AlertSeverity_ALERT_SEVERITY_CRITICAL,
		AlertSeverityError:    alerts.AlertSeverity_ALERT_SEVERITY_ERROR,
	}
	AlertSchemaDayToProtoDay = map[Day]alerts.DayOfWeek{
		Sunday:    alerts.DayOfWeek_DAY_OF_WEEK_SUNDAY,
		Monday:    alerts.DayOfWeek_DAY_OF_WEEK_MONDAY_OR_UNSPECIFIED,
		Tuesday:   alerts.DayOfWeek_DAY_OF_WEEK_TUESDAY,
		Wednesday: alerts.DayOfWeek_DAY_OF_WEEK_WEDNESDAY,
		Thursday:  alerts.DayOfWeek_DAY_OF_WEEK_THURSDAY,
		Friday:    alerts.DayOfWeek_DAY_OF_WEEK_FRIDAY,
		Saturday:  alerts.DayOfWeek_DAY_OF_WEEK_SATURDAY,
	}
	AlertSchemaTimeWindowToProtoTimeWindow = map[string]alerts.Timeframe{
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
	AlertSchemaAutoRetireRatioToProtoAutoRetireRatio = map[AutoRetireRatio]alerts.CleanupDeadmanDuration{
		AutoRetireRatioNever:           alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_NEVER_OR_UNSPECIFIED,
		AutoRetireRatioFiveMinutes:     alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_5MIN,
		AutoRetireRatioTenMinutes:      alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_10MIN,
		AutoRetireRatioHour:            alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_1H,
		AutoRetireRatioTwoHours:        alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_2H,
		AutoRetireRatioSixHours:        alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_6H,
		AutoRetireRatioTwelveHours:     alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_12H,
		AutoRetireRatioTwentyFourHours: alerts.CleanupDeadmanDuration_CLEANUP_DEADMAN_DURATION_24H,
	}
	AlertSchemaFiltersLogSeverityToProtoFiltersLogSeverity = map[FiltersLogSeverity]alerts.AlertFilters_LogSeverity{
		FiltersLogSeverityDebug:    alerts.AlertFilters_LOG_SEVERITY_DEBUG_OR_UNSPECIFIED,
		FiltersLogSeverityVerbose:  alerts.AlertFilters_LOG_SEVERITY_VERBOSE,
		FiltersLogSeverityInfo:     alerts.AlertFilters_LOG_SEVERITY_INFO,
		FiltersLogSeverityWarning:  alerts.AlertFilters_LOG_SEVERITY_WARNING,
		FiltersLogSeverityCritical: alerts.AlertFilters_LOG_SEVERITY_CRITICAL,
		FiltersLogSeverityError:    alerts.AlertFilters_LOG_SEVERITY_ERROR,
	}
	AlertSchemaRelativeTimeFrameToProtoTimeFrameAndRelativeTimeFrame = map[RelativeTimeWindow]protoTimeFrameAndRelativeTimeFrame{
		RelativeTimeWindowPreviousHour:      {timeFrame: alerts.Timeframe_TIMEFRAME_1_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_HOUR_OR_UNSPECIFIED},
		RelativeTimeWindowSameHourYesterday: {timeFrame: alerts.Timeframe_TIMEFRAME_1_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_DAY},
		RelativeTimeWindowSameHourLastWeek:  {timeFrame: alerts.Timeframe_TIMEFRAME_1_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_WEEK},
		RelativeTimeWindowYesterday:         {timeFrame: alerts.Timeframe_TIMEFRAME_24_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_DAY},
		RelativeTimeWindowSameDayLastWeek:   {timeFrame: alerts.Timeframe_TIMEFRAME_24_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_WEEK},
		RelativeTimeWindowSameDayLastMonth:  {timeFrame: alerts.Timeframe_TIMEFRAME_24_H, relativeTimeFrame: alerts.RelativeTimeframe_RELATIVE_TIMEFRAME_MONTH},
	}
	AlertSchemaArithmeticOperatorToProtoArithmeticOperator = map[ArithmeticOperator]alerts.MetricAlertConditionParameters_ArithmeticOperator{
		ArithmeticOperatorAvg:        alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_AVG_OR_UNSPECIFIED,
		ArithmeticOperatorMin:        alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_MIN,
		ArithmeticOperatorMax:        alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_MAX,
		ArithmeticOperatorSum:        alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_SUM,
		ArithmeticOperatorCount:      alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_COUNT,
		ArithmeticOperatorPercentile: alerts.MetricAlertConditionParameters_ARITHMETIC_OPERATOR_PERCENTILE,
	}
	AlertSchemaTracingFilterFieldToProtoTracingFilterField = map[FieldFilterType]string{
		"Application": "applicationName",
		"Subsystem":   "subsystemName",
		"Service":     "serviceName",
	}
	AlertSchemaTracingOperatorToProtoTracingOperator = map[FilterOperator]string{
		"Equals":    "equals",
		"Contains":  "contains",
		"StartWith": "startsWith",
		"EndWith":   "endsWith",
	}
	AlertSchemaFlowOperatorToProtoFlowOperator = map[FlowOperator]alerts.FlowOperator{
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

	// +optional
	PayloadFilters []string `json:"payloadFilters,omitempty"`

	Scheduling *Scheduling `json:"scheduling,omitempty"`

	AlertType AlertType `json:"alertType,omitempty"`
}

func (in *AlertSpec) ExtractCreateAlertRequest() *alerts.CreateAlertRequest {
	enabled := wrapperspb.Bool(in.Active)
	name := wrapperspb.String(in.Name)
	description := wrapperspb.String(in.Description)
	severity := AlertSchemaSeverityToProtoSeverity[in.Severity]
	metaLabels := expandMetaLabels(in.Labels)
	expirationDate := expandExpirationDate(in.ExpirationDate)
	notifications := expandNotifications(in.Notifications.Recipients)
	notifyEvery := expandNotifyEvery(in.Notifications.NotifyEveryMin)
	payloadFilters := utils.StringSliceToWrappedStringSlice(in.PayloadFilters)
	activeWhen := expandActiveWhen(in.Scheduling)
	alertTypeParams := expandAlertType(in.AlertType, in.Notifications.OnTriggerAndResolved,
		in.Notifications.NotifyOnlyOnTriggeredGroupByValues)

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

func expandAlertType(alertType AlertType, onTriggerAndResolved, notifyOnlyOnTriggeredGroupByValues bool) alertTypeParams {
	if standard := alertType.Standard; standard != nil {
		return expandStandard(standard, onTriggerAndResolved, notifyOnlyOnTriggeredGroupByValues)
	} else if ratio := alertType.Ratio; ratio != nil {
		return expandRatio(ratio, onTriggerAndResolved, notifyOnlyOnTriggeredGroupByValues)
	} else if newValue := alertType.NewValue; newValue != nil {
		return expandNewValue(newValue)
	} else if uniqueCount := alertType.UniqueCount; uniqueCount != nil {
		return expandUniqueCount(uniqueCount)
	} else if timeRelative := alertType.TimeRelative; timeRelative != nil {
		return expandTimeRelative(timeRelative, onTriggerAndResolved, notifyOnlyOnTriggeredGroupByValues)
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

func expandRatio(ratio *Ratio, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues bool) alertTypeParams {
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

	condition := expandRatioCondition(ratio.Conditions, groupByQ1, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues)
	filters := expandRatioFilters(&ratio.Query1Filters, &ratio.Query2Filters, groupByQ2)

	return alertTypeParams{
		condition: condition,
		filters:   filters,
	}
}

func expandRatioCondition(conditions RatioConditions, q1GroupBy []*wrapperspb.StringValue, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues bool) *alerts.AlertCondition {
	threshold := wrapperspb.Double(conditions.Ratio.AsApproximateFloat64())
	timeFrame := AlertSchemaTimeWindowToProtoTimeWindow[string(conditions.TimeWindow)]
	notifyOnResolved := wrapperspb.Bool(notifyWhenResolved)
	notifyGroupByOnlyAlerts := wrapperspb.Bool(notifyOnlyOnTriggeredGroupByValues)
	ignoreInfinity := wrapperspb.Bool(conditions.IgnoreInfinity)
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

	var alias *wrapperspb.StringValue
	if desiredAlias := q2Filters.Alias; desiredAlias != nil {
		alias = wrapperspb.String(*desiredAlias)
	}
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
	timeFrame := AlertSchemaTimeWindowToProtoTimeWindow[string(conditions.TimeWindow)]
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
	timeFrame := AlertSchemaTimeWindowToProtoTimeWindow[string(conditions.TimeWindow)]
	var groupBy []*wrapperspb.StringValue
	var maxUniqueValuesForGroupBy *wrapperspb.UInt32Value
	if groupByKey := conditions.GroupBy; groupByKey != nil {
		groupBy = []*wrapperspb.StringValue{wrapperspb.String(*groupByKey)}
		maxUniqueValuesForGroupBy = wrapperspb.UInt32(uint32(*conditions.MaxUniqueValuesForGroupBy))
	}

	parameters := &alerts.ConditionParameters{
		CardinalityFields:                 uniqueCountKey,
		Threshold:                         threshold,
		Timeframe:                         timeFrame,
		GroupBy:                           groupBy,
		MaxUniqueCountValuesForGroupByKey: maxUniqueValuesForGroupBy,
	}

	return &alerts.AlertCondition{
		Condition: &alerts.AlertCondition_UniqueCount{
			UniqueCount: &alerts.UniqueCountCondition{
				Parameters: parameters,
			},
		},
	}
}

func expandTimeRelative(timeRelative *TimeRelative, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues bool) alertTypeParams {
	condition := expandTimeRelativeCondition(&timeRelative.Conditions, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues)
	filters := expandCommonFilters(timeRelative.Filters)
	filters.FilterType = alerts.AlertFilters_FILTER_TYPE_TIME_RELATIVE
	return alertTypeParams{
		condition: condition,
		filters:   filters,
	}
}

func expandTimeRelativeCondition(condition *TimeRelativeConditions, notifyWhenResolved, notifyOnlyOnTriggeredGroupByValues bool) *alerts.AlertCondition {
	threshold := wrapperspb.Double(condition.Threshold.AsApproximateFloat64())
	timeFrameAndRelativeTimeFrame := AlertSchemaRelativeTimeFrameToProtoTimeFrameAndRelativeTimeFrame[condition.TimeWindow]
	groupBy := utils.StringSliceToWrappedStringSlice(condition.GroupBy)
	notifyOnResolved := wrapperspb.Bool(notifyWhenResolved)
	notifyGroupByOnlyAlerts := wrapperspb.Bool(notifyOnlyOnTriggeredGroupByValues)
	ignoreInf := wrapperspb.Bool(condition.IgnoreInfinity)
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

func expandPromqlCondition(conditions *PromqlConditions, searchQuery string, notifyWhenResolved bool) *alerts.AlertCondition {
	text := wrapperspb.String(searchQuery)
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
	timeWindow := AlertSchemaTimeWindowToProtoTimeWindow[string(conditions.TimeWindow)]
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
	arithmeticOperator := AlertSchemaArithmeticOperatorToProtoArithmeticOperator[conditions.ArithmeticOperator]
	var arithmeticOperatorModifier *wrapperspb.UInt32Value
	if modifier := conditions.ArithmeticOperatorModifier; modifier != nil {
		arithmeticOperatorModifier = wrapperspb.UInt32(uint32(*modifier))
	}
	sampleThresholdPercentage := wrapperspb.UInt32(uint32(conditions.SampleThresholdPercentage))
	swapNullValues := wrapperspb.Bool(conditions.ReplaceMissingValueWithZero)
	nonNullPercentage := wrapperspb.UInt32(uint32(*conditions.MinNonNullValuesPercentage))

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
	timeWindow := AlertSchemaTimeWindowToProtoTimeWindow[string(conditions.TimeWindow)]
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
			timeFrame = AlertSchemaTimeWindowToProtoTimeWindow[string(*timeWindow)]
		}
		groupBy := utils.StringSliceToWrappedStringSlice(conditions.GroupBy)
		threshold := wrapperspb.Double(float64(*conditions.Threshold))
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
	conditionLatency := uint32(tracingFilters.LatencyThresholdMilliseconds.AsApproximateFloat64() * float64(time.Millisecond.Microseconds()))
	fieldFilters := expandFiltersData(tracingFilters.Applications, tracingFilters.Subsystems, tracingFilters.Services)
	tagFilters := expandTagFilters(tracingFilters.TagFilters)
	return &alerts.TracingAlert{
		ConditionLatency: conditionLatency,
		FieldFilters:     fieldFilters,
		TagFilters:       tagFilters,
	}
}

func expandFiltersData(applications, subsystems, services []string) []*alerts.FilterData {
	result := make([]*alerts.FilterData, 0, len(applications)+len(subsystems)+len(services))
	result = append(result, expandSpecificFilter("applicationName", applications))
	result = append(result, expandSpecificFilter("subsystemName", subsystems))
	result = append(result, expandSpecificFilter("serviceName", services))
	return result
}

func expandTagFilters(tagFilters []TagFilter) []*alerts.FilterData {
	result := make([]*alerts.FilterData, 0, len(tagFilters))
	for _, tagFilter := range tagFilters {
		result = append(result, expandSpecificFilter(tagFilter.Field, tagFilter.Values))
	}
	return result
}

func expandSpecificFilter(filterName string, values []string) *alerts.FilterData {
	operatorToFilterValues := make(map[string]*alerts.Filters)
	for _, val := range values {
		operator, filterValue := expandFilter(val)
		if _, ok := operatorToFilterValues[operator]; !ok {
			operatorToFilterValues[operator] = new(alerts.Filters)
			operatorToFilterValues[operator].Operator = operator
			operatorToFilterValues[operator].Values = make([]string, 0)
		}
		operatorToFilterValues[operator].Values = append(operatorToFilterValues[operator].Values, filterValue)
	}

	filterResult := make([]*alerts.Filters, 0, len(operatorToFilterValues))
	for _, filters := range operatorToFilterValues {
		filterResult = append(filterResult, filters)
	}

	return &alerts.FilterData{
		Field:   filterName,
		Filters: filterResult,
	}
}

func expandFilter(filterString string) (operator, filterValue string) {
	operator, filterValue = "Equals", filterString
	if strings.HasPrefix(filterValue, "filter:") {
		arr := strings.SplitN(filterValue, ":", 3)
		operator, filterValue = arr[1], arr[2]
	}

	return
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
	subAlerts := expandFlowSubgroupAlerts(group.InnerFlowAlerts)
	nextOp := AlertSchemaFlowOperatorToProtoFlowOperator[group.NextOperator]
	return &alerts.FlowGroup{
		Alerts: subAlerts,
		NextOp: nextOp,
	}
}

func expandFlowSubgroupAlerts(subgroup InnerFlowAlerts) *alerts.FlowAlerts {
	return &alerts.FlowAlerts{
		Op:     AlertSchemaFlowOperatorToProtoFlowOperator[subgroup.Operator],
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
		severity := AlertSchemaFiltersLogSeverityToProtoFiltersLogSeverity[s]
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
		timeFrame = AlertSchemaTimeWindowToProtoTimeWindow[string(*condition.TimeWindow)]
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
		cleanupDeadmanDuration := AlertSchemaAutoRetireRatioToProtoAutoRetireRatio[*manageUndetectedValues.AutoRetireRatio]
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

	timeFrames := expandTimeFrames(scheduling)

	return &alerts.AlertActiveWhen{
		Timeframes: timeFrames,
	}
}

func expandTimeFrames(scheduling *Scheduling) []*alerts.AlertActiveTimeframe {
	utc := ExtractUTC(scheduling.TimeZone)
	daysOfWeek := expandDaysOfWeek(scheduling.DaysEnabled)
	start := expandTime(scheduling.StartTime)
	end := expandTime(scheduling.EndTime)
	timeRange := &alerts.TimeRange{
		Start: start,
		End:   end,
	}
	timeRange, daysOfWeek = convertTimeFramesToGMT(timeRange, daysOfWeek, utc)

	alertActiveTimeframe := &alerts.AlertActiveTimeframe{
		DaysOfWeek: daysOfWeek,
		Range:      timeRange,
	}

	return []*alerts.AlertActiveTimeframe{
		alertActiveTimeframe,
	}
}

func ExtractUTC(timeZone TimeZone) int32 {
	utcStr := strings.Split(string(timeZone), "UTC")[1]
	utc, _ := strconv.Atoi(utcStr)
	return int32(utc)
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
		daysOfWeek = append(daysOfWeek, AlertSchemaDayToProtoDay[d])
	}
	return daysOfWeek
}

func convertTimeFramesToGMT(frameRange *alerts.TimeRange, daysOfWeek []alerts.DayOfWeek, utc int32) (*alerts.TimeRange, []alerts.DayOfWeek) {
	daysOfWeekOffset := daysOfWeekOffsetToGMT(frameRange, utc)
	frameRange.Start.Hours = convertUtcToGmt(frameRange.GetStart().GetHours(), utc)
	frameRange.End.Hours = convertUtcToGmt(frameRange.GetEnd().GetHours(), utc)
	if daysOfWeekOffset != 0 {
		for i, d := range daysOfWeek {
			daysOfWeek[i] = alerts.DayOfWeek((int32(d) + daysOfWeekOffset) % 7)
		}
	}

	return frameRange, daysOfWeek
}

func daysOfWeekOffsetToGMT(frameRange *alerts.TimeRange, utc int32) int32 {
	daysOfWeekOffset := int32(frameRange.Start.Hours-utc) / 24
	if daysOfWeekOffset < 0 {
		daysOfWeekOffset += 7
	}
	return daysOfWeekOffset
}

func convertUtcToGmt(hours, utc int32) int32 {
	hours -= utc
	if hours < 0 {
		hours += 24
	} else if hours >= 24 {
		hours -= 24
	}

	return hours
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

func (in *AlertSpec) DeepEqual(actualAlert *AlertStatus) (bool, utils.Diff) {
	if actualName := actualAlert.Name; actualName != in.Name {
		return false, utils.Diff{
			Name:    "Name",
			Desired: in.Name,
			Actual:  actualName,
		}
	}

	if actualDescription := actualAlert.Description; actualDescription != in.Description {
		return false, utils.Diff{
			Name:    "Description",
			Desired: in.Description,
			Actual:  actualDescription,
		}
	}

	if actualActive := actualAlert.Active; actualActive != in.Active {
		return false, utils.Diff{
			Name:    "Active",
			Desired: in.Active,
			Actual:  actualActive,
		}
	}

	if actualSeverity := actualAlert.Severity; actualSeverity != in.Severity {
		return false, utils.Diff{
			Name:    "Severity",
			Desired: in.Severity,
			Actual:  actualSeverity,
		}
	}

	if !reflect.DeepEqual(in.Labels, actualAlert.Labels) {
		return false, utils.Diff{
			Name:    "Labels",
			Desired: in.Labels,
			Actual:  actualAlert.Labels,
		}
	}

	if !reflect.DeepEqual(in.ExpirationDate, actualAlert.ExpirationDate) {
		return false, utils.Diff{
			Name:    "ExpirationDate",
			Desired: utils.PointerToString(in.ExpirationDate),
			Actual:  utils.PointerToString(actualAlert.ExpirationDate),
		}
	}

	if equal, diff := in.AlertType.DeepEqual(actualAlert.AlertType); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("AlertType.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Notifications.DeepEqual(actualAlert.Notifications); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Notifications.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.PayloadFilters, actualAlert.PayloadFilters) {
		return false, utils.Diff{
			Name:    "PayloadFilters",
			Desired: in.PayloadFilters,
			Actual:  actualAlert.PayloadFilters,
		}
	}

	if scheduling, actualScheduling := in.Scheduling, actualAlert.Scheduling; scheduling == nil && actualScheduling != nil {
		return false, utils.Diff{
			Name:    "Scheduling",
			Desired: scheduling,
			Actual:  actualScheduling,
		}
	} else if equal, diff := scheduling.DeepEqual(*actualScheduling); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Scheduling.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

func (in *AlertSpec) ExtractUpdateAlertRequest(id string) *alerts.UpdateAlertByUniqueIdRequest {
	uniqueIdentifier := wrapperspb.String(id)
	enabled := wrapperspb.Bool(in.Active)
	name := wrapperspb.String(in.Name)
	description := wrapperspb.String(in.Description)
	severity := AlertSchemaSeverityToProtoSeverity[in.Severity]
	metaLabels := expandMetaLabels(in.Labels)
	expirationDate := expandExpirationDate(in.ExpirationDate)
	notifications := expandNotifications(in.Notifications.Recipients)
	notifyEvery := expandNotifyEvery(in.Notifications.NotifyEveryMin)
	payloadFilters := utils.StringSliceToWrappedStringSlice(in.PayloadFilters)
	activeWhen := expandActiveWhen(in.Scheduling)
	alertTypeParams := expandAlertType(in.AlertType, in.Notifications.OnTriggerAndResolved,
		in.Notifications.NotifyOnlyOnTriggeredGroupByValues)

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

const (
	AlertSeverityInfo     AlertSeverity = "Info"
	AlertSeverityWarning  AlertSeverity = "Warning"
	AlertSeverityCritical AlertSeverity = "Critical"
	AlertSeverityError    AlertSeverity = "Error"
)

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
	OnTriggerAndResolved bool `json:"OnTriggerAndResolved,omitempty"`

	//+kubebuilder:default=false
	NotifyOnlyOnTriggeredGroupByValues bool `json:"NotifyOnlyOnTriggeredGroupByValues,omitempty"`

	// +optional
	Recipients Recipients `json:"recipients,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum:=1
	NotifyEveryMin *int `json:"notifyEveryMin,omitempty"`
}

func (in *Notifications) DeepEqual(actualNotifications *Notifications) (bool, utils.Diff) {
	if equal, diff := in.Recipients.DeepEqual(actualNotifications.Recipients); !equal {
		return false, utils.Diff{
			Name:    "Recipients",
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if !reflect.DeepEqual(in.NotifyEveryMin, actualNotifications.NotifyEveryMin) {
		return false, utils.Diff{
			Name:    "NotifyEveryMin",
			Desired: in.NotifyEveryMin,
			Actual:  utils.PointerToString(actualNotifications.NotifyEveryMin),
		}
	}

	if in.OnTriggerAndResolved != actualNotifications.OnTriggerAndResolved {
		return false, utils.Diff{
			Name:    "OnTriggerAndResolved",
			Desired: in.OnTriggerAndResolved,
			Actual:  actualNotifications.OnTriggerAndResolved,
		}
	}

	if in.NotifyOnlyOnTriggeredGroupByValues != actualNotifications.NotifyOnlyOnTriggeredGroupByValues {
		return false, utils.Diff{
			Name:    "NotifyOnlyOnTriggeredGroupByValues",
			Desired: in.NotifyOnlyOnTriggeredGroupByValues,
			Actual:  actualNotifications.NotifyOnlyOnTriggeredGroupByValues,
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

func (in *Recipients) DeepEqual(actualRecipients Recipients) (bool, utils.Diff) {
	if !utils.SlicesWithUniqueValuesEqual(in.Emails, actualRecipients.Emails) {
		return false, utils.Diff{
			Name:    "Emails",
			Desired: in.Emails,
			Actual:  actualRecipients.Emails,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Webhooks, actualRecipients.Webhooks) {
		return false, utils.Diff{
			Name:    "Webhooks",
			Desired: in.Webhooks,
			Actual:  actualRecipients.Webhooks,
		}
	}

	return false, utils.Diff{}
}

type Scheduling struct {
	//+kubebuilder:default=UTC+00
	TimeZone TimeZone `json:"timeZone,omitempty"`

	DaysEnabled []Day `json:"daysEnabled,omitempty"`

	StartTime *Time `json:"startTime,omitempty"`

	EndTime *Time `json:"endTime,omitempty"`
}

func (in *Scheduling) DeepEqual(scheduling Scheduling) (bool, utils.Diff) {
	if in.TimeZone != scheduling.TimeZone {
		return false, utils.Diff{
			Name:    "TimeZone",
			Desired: in.TimeZone,
			Actual:  scheduling.TimeZone,
		}
	}

	if utils.SlicesWithUniqueValuesEqual(in.DaysEnabled, scheduling.DaysEnabled) {
		return false, utils.Diff{
			Name:    "DaysEnabled",
			Desired: in.DaysEnabled,
			Actual:  scheduling.DaysEnabled,
		}
	}

	if !reflect.DeepEqual(in.StartTime, scheduling.StartTime) {
		return false, utils.Diff{
			Name:    "StartTime",
			Desired: utils.PointerToString(in.TimeZone),
			Actual:  utils.PointerToString(scheduling.TimeZone),
		}
	}

	if !reflect.DeepEqual(in.EndTime, scheduling.EndTime) {
		return false, utils.Diff{
			Name:    "StartTime",
			Desired: utils.PointerToString(in.TimeZone),
			Actual:  utils.PointerToString(scheduling.TimeZone),
		}
	}

	return true, utils.Diff{}
}

func DeepEqualTimeFrames(timeframe, actualTimeframe *alerts.AlertActiveTimeframe) (bool, utils.Diff) {
	if equal, diff := DeepEqualTimeRanges(timeframe.GetRange(), actualTimeframe.GetRange()); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Range.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if daysOfWeek, actualDaysOfWeek := timeframe.GetDaysOfWeek(), actualTimeframe.GetDaysOfWeek(); !utils.SlicesWithUniqueValuesEqual(daysOfWeek, actualDaysOfWeek) {
		return false, utils.Diff{
			Name:    "DaysEnabled",
			Desired: daysOfWeek,
			Actual:  actualDaysOfWeek,
		}
	}

	return true, utils.Diff{}
}

func DeepEqualTimeRanges(timeRange, actualTimeRange *alerts.TimeRange) (bool, utils.Diff) {
	if equal, diff := DeepEqualTimes(timeRange.GetStart(), actualTimeRange.GetStart()); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Start.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}
	if equal, diff := DeepEqualTimes(timeRange.GetEnd(), actualTimeRange.GetEnd()); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("End.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

func DeepEqualTimes(time, actualTime *alerts.Time) (bool, utils.Diff) {
	if time.GetHours() != actualTime.GetHours() || time.GetMinutes() != actualTime.GetMinutes() || time.GetSeconds() != actualTime.GetSeconds() {
		return false, utils.Diff{
			Name:    "Hour",
			Desired: time.String(),
			Actual:  actualTime.String(),
		}
	}

	return true, utils.Diff{}
}

// +kubebuilder:validation:Pattern=`^UTC[+-]\d{2}$`
type TimeZone string

// +kubebuilder:validation:Enum=Sunday;Monday;Tuesday;Wednesday;Thursday;Friday;Saturday;
type Day string

const (
	Sunday    Day = "Sunday"
	Monday    Day = "Monday"
	Tuesday   Day = "Tuesday"
	Wednesday Day = "Wednesday"
	Thursday  Day = "Thursday"
	Friday    Day = "Friday"
	Saturday  Day = "Saturday"
)

// +kubebuilder:validation:Pattern=`^(0\d|1\d|2[0-3]):[0-5]\d$`
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

func (in *AlertType) DeepEqual(actualAlert AlertType) (bool, utils.Diff) {
	if newValue := in.NewValue; newValue != nil {
		if actualNewValue := actualAlert.NewValue; actualNewValue == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "NewValue",
			}
		} else if equal, diff := newValue.DeepEqual(*actualNewValue); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("NewValue.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	if standard := in.Standard; standard != nil {
		if actualStandard := actualAlert.Standard; actualStandard == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Standard",
			}
		} else if equal, diff := standard.DeepEqual(*actualStandard); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Standard.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	if ratio := in.Ratio; ratio != nil {
		if actualRatio := actualAlert.Ratio; actualRatio == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Ratio",
			}
		} else if equal, diff := ratio.DeepEqual(*actualRatio); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Ratio.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	if uniqueCount := in.UniqueCount; uniqueCount != nil {
		if actualUniqueCount := actualAlert.UniqueCount; actualUniqueCount == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "UniqueCount",
			}
		} else if equal, diff := uniqueCount.DeepEqual(*actualUniqueCount); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("UniqueCount.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	if timeRelative := in.TimeRelative; timeRelative != nil {
		if actualTimeRelative := actualAlert.TimeRelative; actualTimeRelative == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "TimeRelative",
			}
		} else if equal, diff := timeRelative.DeepEqual(*actualTimeRelative); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("TimeRelative.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	if metric := in.Metric; metric != nil {
		if actualMetric := actualAlert.Metric; actualMetric == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Metric",
			}
		} else if equal, diff := metric.DeepEqual(*actualMetric); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Metric.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	if tracing := in.Tracing; tracing != nil {
		if actulTracing := actualAlert.Tracing; actulTracing == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Tracing",
			}
		} else if equal, diff := tracing.DeepEqual(*actulTracing); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Tracing.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	if flow := in.Flow; flow != nil {
		if actualFlow := actualAlert.Flow; actualFlow == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Flow",
			}
		} else if equal, diff := flow.DeepEqual(*actualFlow); !equal {
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

func (in *Standard) DeepEqual(actualStandard Standard) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(actualStandard.Conditions); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Filters.DeepEqual(actualStandard.Filters); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Filters.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

type Ratio struct {
	Query1Filters Filters `json:"q1Filters,omitempty"`

	Query2Filters RatioQ2Filters `json:"q2Filters,omitempty"`

	Conditions RatioConditions `json:"conditions,omitempty"`
}

type RatioQ2Filters struct {
	// +optional
	Alias *string `json:"alias,omitempty"`

	// +optional
	SearchQuery *string `json:"searchQuery,omitempty"`

	// +optional
	Severities []FiltersLogSeverity `json:"severities,omitempty"`

	// +optional
	Applications []string `json:"applications,omitempty"`

	// +optional
	Subsystems []string `json:"subsystems,omitempty"`
}

func (in *Ratio) DeepEqual(actualRatio Ratio) (bool, utils.Diff) {
	if equal, diff := in.Query1Filters.DeepEqual(&actualRatio.Query1Filters); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Q1Filters.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Query2Filters.DeepEqual(actualRatio.Query2Filters); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Q2Filters.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Conditions.DeepEqual(actualRatio.Conditions); !equal {
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

func (in *NewValue) DeepEqual(newValue NewValue) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(newValue.Conditions); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Filters.DeepEqual(newValue.Filters); !equal {
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

func (in *UniqueCount) DeepEqual(actualUniqueCount UniqueCount) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(actualUniqueCount.Conditions); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Filters.DeepEqual(actualUniqueCount.Filters); !equal {
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

func (in *TimeRelative) DeepEqual(actualTimeRelative TimeRelative) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(actualTimeRelative.Conditions); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}
	if equal, diff := in.Filters.DeepEqual(actualTimeRelative.Filters); !equal {
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

func (in *Metric) DeepEqual(actualMetric Metric) (bool, utils.Diff) {
	if promql := in.Promql; promql != nil {
		if actualMetric.Promql == nil {
			return false, utils.Diff{
				Name:    "Promql",
				Desired: promql,
				Actual:  "nil",
			}
		} else if equal, diff := promql.DeepEqual(*actualMetric.Promql); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Promql.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	} else if lucene := in.Lucene; lucene != nil {
		if actualMetric.Lucene == nil {
			return false, utils.Diff{
				Name:    "Lucene",
				Desired: promql,
				Actual:  "nil",
			}
		} else if equal, diff := lucene.DeepEqual(*actualMetric.Lucene); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Lucene.%s", diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	return false, utils.Diff{}
}

type Lucene struct {
	// +optional
	SearchQuery *string `json:"searchQuery,omitempty"`

	Conditions LuceneConditions `json:"conditions,omitempty"`
}

func (in *Lucene) DeepEqual(actualLucene Lucene) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(actualLucene.Conditions); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Condition.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if !reflect.DeepEqual(in.SearchQuery, actualLucene.SearchQuery) {
		return false, utils.Diff{
			Name:    "SearchQuery",
			Desired: utils.PointerToString(in.SearchQuery),
			Actual:  utils.PointerToString(actualLucene.SearchQuery),
		}
	}

	return true, utils.Diff{}
}

type Promql struct {
	SearchQuery string `json:"searchQuery,omitempty"`

	Conditions PromqlConditions `json:"conditions,omitempty"`
}

func (in *Promql) DeepEqual(actualPromql Promql) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(actualPromql.Conditions); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Condition.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if in.SearchQuery != actualPromql.SearchQuery {
		return false, utils.Diff{
			Name:    "SearchQuery",
			Desired: in.SearchQuery,
			Actual:  actualPromql.SearchQuery,
		}
	}

	return true, utils.Diff{}
}

type Tracing struct {
	Filters TracingFilters `json:"filters,omitempty"`

	Conditions TracingCondition `json:"conditions,omitempty"`
}

func (in *Tracing) DeepEqual(actualTracing Tracing) (bool, utils.Diff) {
	if equal, diff := in.Conditions.DeepEqual(actualTracing.Conditions); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("Conditions.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if equal, diff := in.Filters.DeepEqual(actualTracing.Filters); !equal {
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

func (in *Flow) DeepEqual(actualFlow Flow) (bool, utils.Diff) {
	if stages, actualStages := in.Stages, actualFlow.Stages; len(stages) != len(actualStages) {
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

func (in *StandardConditions) DeepEqual(actualCondition StandardConditions) (bool, utils.Diff) {
	if in.AlertWhen != actualCondition.AlertWhen {
		return false, utils.Diff{
			Name:    "AlertWhen",
			Desired: in.AlertWhen,
			Actual:  actualCondition.AlertWhen,
		}
	}

	if !reflect.DeepEqual(in.Threshold, actualCondition.Threshold) {
		return false, utils.Diff{
			Name:    "Threshold",
			Desired: utils.PointerToString(in.Threshold),
			Actual:  utils.PointerToString(actualCondition.Threshold),
		}
	}

	if !reflect.DeepEqual(in.TimeWindow, actualCondition.TimeWindow) {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: utils.PointerToString(in.TimeWindow),
			Actual:  utils.PointerToString(actualCondition.TimeWindow),
		}
	}

	if manageUndetectedValues, actualManageUndetectedValues := in.ManageUndetectedValues, actualCondition.ManageUndetectedValues; manageUndetectedValues == nil && actualManageUndetectedValues != nil {
		return false, utils.Diff{
			Name:    "ManageUndetectedValues",
			Desired: utils.PointerToString(manageUndetectedValues),
			Actual:  utils.PointerToString(actualManageUndetectedValues),
		}
	} else if equal, diff := manageUndetectedValues.DeepEqual(actualManageUndetectedValues); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("ManageUndetectedValues.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.GroupBy, actualCondition.GroupBy) {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: in.GroupBy,
			Actual:  actualCondition.GroupBy,
		}
	}

	return true, utils.Diff{}
}

type RatioConditions struct {
	AlertWhen AlertWhen `json:"alertWhen,omitempty"`

	Ratio resource.Quantity `json:"ratio,omitempty"`

	//+kubebuilder:default=false
	IgnoreInfinity bool `json:"ignoreInfinity,omitempty"`

	TimeWindow TimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	GroupByFor *GroupByFor `json:"groupByFor,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

func (in *RatioConditions) DeepEqual(actualCondition RatioConditions) (bool, utils.Diff) {
	if in.AlertWhen != actualCondition.AlertWhen {
		return false, utils.Diff{
			Name:    "AlertWhen",
			Desired: in.AlertWhen,
			Actual:  actualCondition.AlertWhen,
		}
	}

	if manageUndetectedValues, actualManageUndetectedValues := in.ManageUndetectedValues, actualCondition.ManageUndetectedValues; manageUndetectedValues == nil && actualManageUndetectedValues != nil {
		return false, utils.Diff{
			Name:    "ManageUndetectedValues",
			Desired: utils.PointerToString(manageUndetectedValues),
			Actual:  utils.PointerToString(actualManageUndetectedValues),
		}
	} else if equal, diff := manageUndetectedValues.DeepEqual(actualManageUndetectedValues); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("ManageUndetectedValues.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if in.Ratio != actualCondition.Ratio {
		return false, utils.Diff{
			Name:    "Ratio",
			Desired: in.Ratio,
			Actual:  actualCondition.Ratio,
		}
	}

	if in.TimeWindow != actualCondition.TimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: in.TimeWindow,
			Actual:  actualCondition.TimeWindow,
		}
	}

	if !reflect.DeepEqual(in.GroupByFor, actualCondition.GroupByFor) {
		return false, utils.Diff{
			Name:    "GroupByFor",
			Desired: utils.PointerToString(in.GroupByFor),
			Actual:  utils.PointerToString(actualCondition.GroupByFor),
		}
	}

	return true, utils.Diff{}
}

type NewValueConditions struct {
	Key string `json:"key,omitempty"`

	TimeWindow NewValueTimeWindow `json:"timeWindow,omitempty"`
}

func (in *NewValueConditions) DeepEqual(actualConditions NewValueConditions) (bool, utils.Diff) {
	if in.Key != actualConditions.Key {
		return false, utils.Diff{
			Name:    "Key",
			Desired: in.Key,
			Actual:  actualConditions.Key,
		}
	}

	if in.TimeWindow != actualConditions.TimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: in.TimeWindow,
			Actual:  actualConditions.TimeWindow,
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

func (in *UniqueCountConditions) DeepEqual(actualCondition UniqueCountConditions) (bool, utils.Diff) {
	if in.Key != actualCondition.Key {
		return false, utils.Diff{
			Name:    "Key",
			Desired: in.Key,
			Actual:  actualCondition.Key,
		}
	}

	if in.TimeWindow != actualCondition.TimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: in.TimeWindow,
			Actual:  actualCondition.TimeWindow,
		}
	}

	if in.MaxUniqueValues != actualCondition.MaxUniqueValues {
		return false, utils.Diff{
			Name:    "MaxUniqueValues",
			Desired: in.MaxUniqueValues,
			Actual:  actualCondition.MaxUniqueValues,
		}
	}

	if !reflect.DeepEqual(in.GroupBy, actualCondition.GroupBy) {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: utils.PointerToString(in.GroupBy),
			Actual:  utils.PointerToString(actualCondition.GroupBy),
		}
	}

	if !reflect.DeepEqual(in.MaxUniqueValuesForGroupBy, actualCondition.MaxUniqueValuesForGroupBy) {
		return false, utils.Diff{
			Name:    "MaxUniqueValuesForGroupBy",
			Desired: utils.PointerToString(in.MaxUniqueValuesForGroupBy),
			Actual:  utils.PointerToString(actualCondition.MaxUniqueValuesForGroupBy),
		}
	}

	return true, utils.Diff{}
}

type TimeRelativeConditions struct {
	AlertWhen AlertWhen `json:"alertWhen,omitempty"`

	Threshold resource.Quantity `json:"threshold,omitempty"`

	//+kubebuilder:default=false
	IgnoreInfinity bool `json:"ignoreInfinity,omitempty"`

	TimeWindow RelativeTimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

func (in *TimeRelativeConditions) DeepEqual(actualCondition TimeRelativeConditions) (bool, utils.Diff) {
	if in.AlertWhen != actualCondition.AlertWhen {
		return false, utils.Diff{
			Name:    "AlertWhen",
			Desired: in.AlertWhen,
			Actual:  actualCondition.AlertWhen,
		}
	}

	if manageUndetectedValues, actualManageUndetectedValues := in.ManageUndetectedValues, actualCondition.ManageUndetectedValues; manageUndetectedValues == nil && actualManageUndetectedValues != nil {
		return false, utils.Diff{
			Name:    "ManageUndetectedValues",
			Desired: utils.PointerToString(manageUndetectedValues),
			Actual:  utils.PointerToString(actualManageUndetectedValues),
		}
	} else if equal, diff := manageUndetectedValues.DeepEqual(actualManageUndetectedValues); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("ManageUndetectedValues.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if in.Threshold != actualCondition.Threshold {
		return false, utils.Diff{
			Name:    "Threshold",
			Desired: in.Threshold,
			Actual:  actualCondition.Threshold,
		}
	}

	if in.TimeWindow != actualCondition.TimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: in.TimeWindow,
			Actual:  actualCondition.TimeWindow,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.GroupBy, actualCondition.GroupBy) {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: in.GroupBy,
			Actual:  actualCondition.GroupBy,
		}
	}

	return true, utils.Diff{}
}

// +kubebuilder:validation:Enum=Avg;Min;Max;Sum;Count;Percentile;
type ArithmeticOperator string

const (
	ArithmeticOperatorAvg        ArithmeticOperator = "Avg"
	ArithmeticOperatorMin        ArithmeticOperator = "Min"
	ArithmeticOperatorMax        ArithmeticOperator = "Max"
	ArithmeticOperatorSum        ArithmeticOperator = "Sum"
	ArithmeticOperatorCount      ArithmeticOperator = "Count"
	ArithmeticOperatorPercentile ArithmeticOperator = "Percentile"
)

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

func (in *LuceneConditions) DeepEqual(actualCondition LuceneConditions) (bool, utils.Diff) {
	if in.Threshold != actualCondition.Threshold {
		return false, utils.Diff{
			Name:    "Threshold",
			Desired: in.Threshold,
			Actual:  actualCondition.Threshold,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.GroupBy, actualCondition.GroupBy) {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: in.GroupBy,
			Actual:  actualCondition.GroupBy,
		}
	}

	if in.TimeWindow != actualCondition.TimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: in.TimeWindow,
			Actual:  actualCondition.TimeWindow,
		}
	}

	if in.MetricField != actualCondition.MetricField {
		return false, utils.Diff{
			Name:    "MetricField",
			Desired: in.MetricField,
			Actual:  actualCondition.MetricField,
		}
	}

	if in.ArithmeticOperator != actualCondition.ArithmeticOperator {
		return false, utils.Diff{
			Name:    "ArithmeticOperator",
			Desired: in.ArithmeticOperator,
			Actual:  actualCondition.ArithmeticOperator,
		}
	}

	if reflect.DeepEqual(in.ArithmeticOperatorModifier, actualCondition.ArithmeticOperatorModifier) {
		return false, utils.Diff{
			Name:    "ArithmeticOperatorModifier",
			Desired: utils.PointerToString(in.ArithmeticOperatorModifier),
			Actual:  utils.PointerToString(actualCondition.ArithmeticOperatorModifier),
		}
	}

	if in.SampleThresholdPercentage != actualCondition.SampleThresholdPercentage {
		return false, utils.Diff{
			Name:    "SampleThresholdPercentage",
			Desired: in.SampleThresholdPercentage,
			Actual:  actualCondition.SampleThresholdPercentage,
		}
	}

	if in.ReplaceMissingValueWithZero != actualCondition.ReplaceMissingValueWithZero {
		return false, utils.Diff{
			Name:    "MissingValueWithZero",
			Desired: in.ReplaceMissingValueWithZero,
			Actual:  actualCondition.ReplaceMissingValueWithZero,
		}
	}

	if !reflect.DeepEqual(in.MinNonNullValuesPercentage, actualCondition.MinNonNullValuesPercentage) {
		return false, utils.Diff{
			Name:    "MinNonNullValuesPercentage",
			Desired: utils.PointerToString(in.MinNonNullValuesPercentage),
			Actual:  utils.PointerToString(actualCondition.MinNonNullValuesPercentage),
		}
	}

	if manageUndetectedValues, actualManageUndetectedValues := in.ManageUndetectedValues, actualCondition.ManageUndetectedValues; manageUndetectedValues == nil && actualManageUndetectedValues != nil {
		return false, utils.Diff{
			Name:    "ManageUndetectedValues",
			Desired: utils.PointerToString(manageUndetectedValues),
			Actual:  utils.PointerToString(actualManageUndetectedValues),
		}
	} else if equal, diff := manageUndetectedValues.DeepEqual(actualManageUndetectedValues); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("ManageUndetectedValues.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

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

func (in *PromqlConditions) DeepEqual(actualConditions PromqlConditions) (bool, utils.Diff) {
	if in.Threshold != actualConditions.Threshold {
		return false, utils.Diff{
			Name:    "Threshold",
			Desired: in.Threshold,
			Actual:  actualConditions.Threshold,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.GroupBy, actualConditions.GroupBy) {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: in.GroupBy,
			Actual:  actualConditions.GroupBy,
		}
	}

	if in.TimeWindow != actualConditions.TimeWindow {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: in.TimeWindow,
			Actual:  actualConditions.TimeWindow,
		}
	}

	if in.SampleThresholdPercentage != actualConditions.SampleThresholdPercentage {
		return false, utils.Diff{
			Name:    "SampleThresholdPercentage",
			Desired: in.SampleThresholdPercentage,
			Actual:  actualConditions.SampleThresholdPercentage,
		}
	}

	if !reflect.DeepEqual(in.MinNonNullValuesPercentage, actualConditions.MinNonNullValuesPercentage) {
		return false, utils.Diff{
			Name:    "MinNonNullValuesPercentage",
			Desired: utils.PointerToString(in.MinNonNullValuesPercentage),
			Actual:  utils.PointerToString(actualConditions.MinNonNullValuesPercentage),
		}
	}

	if manageUndetectedValues, actualManageUndetectedValues := in.ManageUndetectedValues, actualConditions.ManageUndetectedValues; manageUndetectedValues == nil && actualManageUndetectedValues != nil {
		return false, utils.Diff{
			Name:    "ManageUndetectedValues",
			Desired: utils.PointerToString(manageUndetectedValues),
			Actual:  utils.PointerToString(actualManageUndetectedValues),
		}
	} else if equal, diff := manageUndetectedValues.DeepEqual(actualManageUndetectedValues); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("ManageUndetectedValues.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	return true, utils.Diff{}
}

type TracingCondition struct {
	AlertWhen TracingAlertWhen `json:"alertWhen,omitempty"`

	// +optional
	Threshold *int `json:"threshold,omitempty"`

	// +optional
	TimeWindow *TimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`
}

func (in *TracingCondition) DeepEqual(actualCondition TracingCondition) (bool, utils.Diff) {
	if in.AlertWhen != actualCondition.AlertWhen {
		return false, utils.Diff{
			Name:    "AlertWhen",
			Desired: in.AlertWhen,
			Actual:  actualCondition.AlertWhen,
		}
	}
	if !reflect.DeepEqual(in.Threshold, actualCondition.Threshold) {
		return false, utils.Diff{
			Name:    "Threshold",
			Desired: utils.PointerToString(in.Threshold),
			Actual:  utils.PointerToString(actualCondition.Threshold),
		}
	}
	if !reflect.DeepEqual(in.TimeWindow, actualCondition.TimeWindow) {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: utils.PointerToString(in.TimeWindow),
			Actual:  utils.PointerToString(actualCondition.TimeWindow),
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.GroupBy, actualCondition.GroupBy) {
		return false, utils.Diff{
			Name:    "GroupBy",
			Desired: in.GroupBy,
			Actual:  actualCondition.GroupBy,
		}
	}

	return true, utils.Diff{}
}

// +kubebuilder:validation:Enum=Never;FiveMinutes;TenMinutes;Hour;TwoHours;SixHours;TwelveHours;TwentyFourHours
type AutoRetireRatio string

const (
	AutoRetireRatioNever           AutoRetireRatio = "Never"
	AutoRetireRatioFiveMinutes     AutoRetireRatio = "FiveMinutes"
	AutoRetireRatioTenMinutes      AutoRetireRatio = "TenMinutes"
	AutoRetireRatioHour            AutoRetireRatio = "Hour"
	AutoRetireRatioTwoHours        AutoRetireRatio = "TwoHours"
	AutoRetireRatioSixHours        AutoRetireRatio = "SixHours"
	AutoRetireRatioTwelveHours     AutoRetireRatio = "TwelveHours"
	AutoRetireRatioTwentyFourHours AutoRetireRatio = "TwentyFourHours"
)

// +kubebuilder:validation:Enum=More;Less
type AlertWhen string

const (
	AlertWhenLessThan AlertWhen = "Less"
	AlertWhenMoreThan AlertWhen = "More"
)

// +kubebuilder:validation:Enum=More;Less;Immediately;MoreThanUsual
type StandardAlertWhen string

const (
	StandardAlertWhenLessThan      StandardAlertWhen = "Less"
	StandardAlertWhenMoreThan      StandardAlertWhen = "More"
	StandardAlertWhenMoreThanUsual StandardAlertWhen = "MoreThanUsual"
	StandardAlertWhenImmediately   StandardAlertWhen = "Immediately"
)

// +kubebuilder:validation:Enum=More;Immediately
type TracingAlertWhen string

const (
	TracingAlertWhenMore        TracingAlertWhen = "More"
	TracingAlertWhenImmediately TracingAlertWhen = "Immediately"
)

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

const (
	RelativeTimeWindowPreviousHour      RelativeTimeWindow = "PreviousHour"
	RelativeTimeWindowSameHourYesterday RelativeTimeWindow = "SameHourYesterday"
	RelativeTimeWindowSameHourLastWeek  RelativeTimeWindow = "SameHourLastWeek"
	RelativeTimeWindowYesterday         RelativeTimeWindow = "Yesterday"
	RelativeTimeWindowSameDayLastWeek   RelativeTimeWindow = "SameDayLastWeek"
	RelativeTimeWindowSameDayLastMonth  RelativeTimeWindow = "SameDayLastMonth"
)

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

func (in *Filters) DeepEqual(actualFilters *Filters) (bool, utils.Diff) {
	if !reflect.DeepEqual(in.SearchQuery, actualFilters.SearchQuery) {
		return false, utils.Diff{
			Name:    "SearchQuery",
			Desired: utils.PointerToString(in.SearchQuery),
			Actual:  utils.PointerToString(actualFilters.SearchQuery),
		}
	}

	if !reflect.DeepEqual(in.Alias, actualFilters.Alias) {
		return false, utils.Diff{
			Name:    "Alias",
			Desired: utils.PointerToString(in.Alias),
			Actual:  utils.PointerToString(actualFilters.Alias),
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Severities, actualFilters.Severities) {
		return false, utils.Diff{
			Name:    "Severities",
			Desired: in.Severities,
			Actual:  actualFilters.Severities,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Applications, actualFilters.Applications) {
		return false, utils.Diff{
			Name:    "Application",
			Desired: in.Applications,
			Actual:  actualFilters.Applications,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Subsystems, actualFilters.Subsystems) {
		return false, utils.Diff{
			Name:    "Subsystems",
			Desired: in.Subsystems,
			Actual:  actualFilters.Subsystems,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Categories, actualFilters.Categories) {
		return false, utils.Diff{
			Name:    "Categories",
			Desired: in.Categories,
			Actual:  actualFilters.Categories,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Computers, actualFilters.Computers) {
		return false, utils.Diff{
			Name:    "Computers",
			Desired: in.Computers,
			Actual:  actualFilters.Computers,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Classes, actualFilters.Classes) {
		return false, utils.Diff{
			Name:    "Classes",
			Desired: in.Classes,
			Actual:  actualFilters.Classes,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Methods, actualFilters.Methods) {
		return false, utils.Diff{
			Name:    "Methods",
			Desired: in.Methods,
			Actual:  actualFilters.Methods,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.IPs, actualFilters.IPs) {
		return false, utils.Diff{
			Name:    "IPs",
			Desired: in.IPs,
			Actual:  actualFilters.IPs,
		}
	}

	return true, utils.Diff{}
}

func (in *RatioQ2Filters) DeepEqual(actualRatioQ2Filters RatioQ2Filters) (bool, utils.Diff) {
	if !reflect.DeepEqual(in.SearchQuery, actualRatioQ2Filters.SearchQuery) {
		return false, utils.Diff{
			Name:    "SearchQuery",
			Desired: utils.PointerToString(in.SearchQuery),
			Actual:  utils.PointerToString(actualRatioQ2Filters.SearchQuery),
		}
	}

	if !reflect.DeepEqual(in.Alias, actualRatioQ2Filters.Alias) {
		return false, utils.Diff{
			Name:    "Alias",
			Desired: utils.PointerToString(in.Alias),
			Actual:  utils.PointerToString(actualRatioQ2Filters.Alias),
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Severities, actualRatioQ2Filters.Severities) {
		return false, utils.Diff{
			Name:    "Severities",
			Desired: in.Severities,
			Actual:  actualRatioQ2Filters.Severities,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Applications, actualRatioQ2Filters.Applications) {
		return false, utils.Diff{
			Name:    "Application",
			Desired: in.Applications,
			Actual:  actualRatioQ2Filters.Applications,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Subsystems, actualRatioQ2Filters.Subsystems) {
		return false, utils.Diff{
			Name:    "Subsystems",
			Desired: in.Subsystems,
			Actual:  actualRatioQ2Filters.Subsystems,
		}
	}

	return true, utils.Diff{}
}

// +kubebuilder:validation:Enum=Debug;Verbose;Info;Warning;Critical;Error;
type FiltersLogSeverity string

const (
	FiltersLogSeverityDebug    FiltersLogSeverity = "Debug"
	FiltersLogSeverityVerbose  FiltersLogSeverity = "Verbose"
	FiltersLogSeverityInfo     FiltersLogSeverity = "Info"
	FiltersLogSeverityWarning  FiltersLogSeverity = "Warning"
	FiltersLogSeverityCritical FiltersLogSeverity = "Critical"
	FiltersLogSeverityError    FiltersLogSeverity = "Error"
)

type TracingFilters struct {
	LatencyThresholdMilliseconds resource.Quantity `json:"latencyThresholdMilliseconds,omitempty"`

	// +optional
	TagFilters []TagFilter `json:"tagFilters,omitempty"`

	// +optional
	Applications []string `json:"applications,omitempty"`

	// +optional
	Subsystems []string `json:"subsystems,omitempty"`

	// +optional
	Services []string `json:"services,omitempty"`
}

func (in *TracingFilters) DeepEqual(actualFilters TracingFilters) (bool, utils.Diff) {
	if in.LatencyThresholdMilliseconds != actualFilters.LatencyThresholdMilliseconds {
		return false, utils.Diff{
			Name:    "LatencyThresholdMilliseconds",
			Desired: in.LatencyThresholdMilliseconds,
			Actual:  actualFilters.LatencyThresholdMilliseconds,
		}
	}

	if !deepEqualTagFilters(in.TagFilters, actualFilters.TagFilters) {
		return false, utils.Diff{
			Name:    "TagFilters",
			Desired: in.TagFilters,
			Actual:  actualFilters.TagFilters,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Applications, actualFilters.Applications) {
		return false, utils.Diff{
			Name:    "Applications",
			Desired: in.Applications,
			Actual:  actualFilters.Applications,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Subsystems, actualFilters.Subsystems) {
		return false, utils.Diff{
			Name:    "Subsystems",
			Desired: in.Subsystems,
			Actual:  actualFilters.Subsystems,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Services, actualFilters.Services) {
		return false, utils.Diff{
			Name:    "Services",
			Desired: in.Services,
			Actual:  actualFilters.Services,
		}
	}

	return true, utils.Diff{}
}

func deepEqualTagFilters(tagFilters1, tagFilters2 []TagFilter) bool {
	if len(tagFilters1) != len(tagFilters2) {
		return false
	}

	fieldsToValues := make(map[string][]string, len(tagFilters1))
	for _, tagFilter := range tagFilters1 {
		fieldsToValues[tagFilter.Field] = tagFilter.Values
	}

	for _, tagFilter := range tagFilters2 {
		if tagFilterValues, ok := fieldsToValues[tagFilter.Field]; !ok || !utils.SlicesWithUniqueValuesEqual(tagFilterValues, tagFilter.Values) {
			return false
		}
	}

	return true

}

func flattenTracingFilters(fieldFilters []*alerts.FilterData) (applications, subsystems, services []string) {
	filtersData := flattenFiltersData(fieldFilters)
	applications = filtersData["applicationName"]
	subsystems = filtersData["subsystemName"]
	services = filtersData["serviceName"]
	return
}

func flattenTagFilters(tagFilters []*alerts.FilterData) []TagFilter {
	fieldToFilters := flattenFiltersData(tagFilters)
	result := make([]TagFilter, 0, len(fieldToFilters))
	for field, filters := range fieldToFilters {
		filterSchema := TagFilter{
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
		result[field] = flattenFilters(filter.GetFilters())
	}
	return result
}

func flattenFilters(filters []*alerts.Filters) []string {
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

type TagFilter struct {
	Field  string   `json:"field,omitempty"`
	Values []string `json:"values,omitempty"`
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

func (in *ManageUndetectedValues) DeepEqual(actualManageUndetectedValues *ManageUndetectedValues) (bool, utils.Diff) {
	if actualManageUndetectedValues == nil {
		return false, utils.Diff{}
	}

	if in.EnableTriggeringOnUndetectedValues != actualManageUndetectedValues.EnableTriggeringOnUndetectedValues {
		return false, utils.Diff{
			Name:    "EnableTriggeringOnUndetectedValues",
			Desired: in.EnableTriggeringOnUndetectedValues,
			Actual:  actualManageUndetectedValues.EnableTriggeringOnUndetectedValues,
		}
	}

	if !reflect.DeepEqual(in.AutoRetireRatio, actualManageUndetectedValues.AutoRetireRatio) {
		return false, utils.Diff{
			Name:    "AutoRetireRatio",
			Desired: utils.PointerToString(in.AutoRetireRatio),
			Actual:  utils.PointerToString(actualManageUndetectedValues.AutoRetireRatio),
		}
	}

	return true, utils.Diff{}
}

type FlowStage struct {
	// +optional
	TimeWindow *FlowStageTimeFrame `json:"timeWindow,omitempty"`

	Groups []FlowStageGroup `json:"groups,omitempty"`
}

func (in *FlowStage) DeepEqual(actualStage FlowStage) (bool, utils.Diff) {
	if groups, actualGroups := in.Groups, actualStage.Groups; len(groups) != len(actualGroups) {
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

	if !reflect.DeepEqual(in.TimeWindow, actualStage.TimeWindow) {
		return false, utils.Diff{
			Name:    "TimeWindow",
			Desired: utils.PointerToString(in.TimeWindow),
			Actual:  utils.PointerToString(actualStage.TimeWindow),
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
	InnerFlowAlerts InnerFlowAlerts `json:"innerFlowAlerts,omitempty"`

	NextOperator FlowOperator `json:"nextOperator,omitempty"`
}

func (in *FlowStageGroup) DeepEqual(actualGroup FlowStageGroup) (bool, utils.Diff) {
	if equal, diff := in.InnerFlowAlerts.DeepEqual(actualGroup.InnerFlowAlerts); !equal {
		return false, utils.Diff{
			Name:    fmt.Sprintf("InnerFlowAlerts.%s", diff.Name),
			Desired: diff.Desired,
			Actual:  diff.Actual,
		}
	}

	if in.NextOperator != actualGroup.NextOperator {
		return false, utils.Diff{
			Name:    "NextOperator",
			Desired: in.NextOperator,
			Actual:  actualGroup.NextOperator,
		}
	}

	return true, utils.Diff{}
}

func (in *InnerFlowAlerts) DeepEqual(actualInnerFlowAlerts InnerFlowAlerts) (bool, utils.Diff) {
	if alerts, actualAlerts := in.Alerts, actualInnerFlowAlerts.Alerts; len(alerts) != len(actualAlerts) {
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
	return false, utils.Diff{}
}

type InnerFlowAlerts struct {
	Operator FlowOperator `json:"operator,omitempty"`

	Alerts []InnerFlowAlert `json:"alerts,omitempty"`
}

type InnerFlowAlert struct {
	// +kubebuilder:default=false
	Not bool `json:"not,omitempty"`

	// +optional
	UserAlertId string `json:"userAlertId,omitempty"`
}

func (in *InnerFlowAlert) DeepEqual(actualInnerFlowAlert InnerFlowAlert) (bool, utils.Diff) {
	if in.Not == actualInnerFlowAlert.Not {
		return false, utils.Diff{
			Name:    "Not",
			Desired: in.Not,
			Actual:  actualInnerFlowAlert.Not,
		}
	}

	if in.UserAlertId != actualInnerFlowAlert.UserAlertId {
		return false, utils.Diff{
			Name:    "UserAlertId",
			Desired: in.UserAlertId,
			Actual:  actualInnerFlowAlert.UserAlertId,
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

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	Active bool `json:"active,omitempty"`

	Severity AlertSeverity `json:"severity,omitempty"`

	Labels map[string]string `json:"labels,omitempty"`

	ExpirationDate *ExpirationDate `json:"expirationDate,omitempty"`

	Notifications *Notifications `json:"notifications,omitempty"`

	PayloadFilters []string `json:"payloadFilters,omitempty"`

	Scheduling *Scheduling `json:"scheduling,omitempty"`

	AlertType AlertType `json:"alertType,omitempty"`
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
