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
	utils "coralogix-operator-poc/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

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

	Severity utils.Severity `json:"severity,omitempty"`

	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// +optional
	ExpirationDate ExpirationDate `json:"expirationDate,omitempty"`

	// +optional
	Notification Notification `json:"notification,omitempty"`

	Scheduling Scheduling `json:"scheduling,omitempty"`

	AlertType AlertType `json:"notification,omitempty"`
}

type ExpirationDate struct {
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=31
	Day int `json:"day,omitempty"`

	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=12
	Month int `json:"month,omitempty"`

	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=9999
	Year int `json:"day,omitempty"`
}

type Notification struct {
	// +optional
	OnTriggerAndResolved *bool `json:"onTriggerAndResolved,omitempty"`

	// +optional
	IgnoreInfinity *bool `json:"ignoreInfinity,omitempty"`

	// +optional
	NotifyOnlyOnTriggeredGroupByValues *bool `json:"notifyOnlyOnTriggeredGroupByValues,omitempty"`

	// +optional
	Recipients Recipients `json:"recipients,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum:=1
	NotifyEveryMin *int `json:"NotifyEveryMin,omitempty"`

	// +optional
	PayloadFields []string `json:"payloadFields,omitempty"`
}

type Recipients struct {
	// +optional
	Emails []Email `json:"emails,omitempty"`

	// +optional
	Webhooks []string `json:"webhooks,omitempty"`
}

// +kubebuilder:validation:Pattern:=^[a-z/d._%+\-]+@[a-z/d.\-]+\.[a-z]{2,4}$
type Email string

type Scheduling struct {
	//+kubebuilder:default=UTC+0
	TimeZone TimeZone `json:"timeZone,omitempty"`

	DaysEnabled []Day `json:"daysEnabled,omitempty"`

	StartTime Time `json:"startTime,omitempty"`

	EndTime Time `json:"endTime,omitempty"`
}

// +kubebuilder:validation:Enum=UTC-11;UTC-10";UTC-9;UTC-8;UTC-7;UTC-6;UTC-5;UTC-4;UTC-3;UTC-2;UTC-1;UTC+0;UTC+1;UTC+2;UTC+3;UTC+4;UTC+5;UTC+6;UTC+7;UTC+8;UTC+9;UTC+10;UTC+11;UTC+12;UTC+13;UTC+14
type TimeZone string

// +kubebuilder:validation:Enum=Sunday;Monday;Tuesday;Wednesday;Thursday;Friday;Saturday;
type Day string

// +kubebuilder:validation:Pattern:=^(0\d|1\d|2[0-3]):[0-5]\d$
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

type Standard struct {
	// +optional
	Filters *Filters `json:"filters,omitempty"`

	Conditions StandardConditions `json:"conditions,omitempty"`
}

type Ratio struct {
	Query1Filters Filters `json:"q1Filters,omitempty"`

	Query2Filters Filters `json:"q2Filters,omitempty"`

	Conditions RatioConditions `json:"conditions,omitempty"`
}

type NewValue struct {
	// +optional
	Filters *Filters `json:"filters,omitempty"`

	Conditions NewValueConditions `json:"conditions,omitempty"`
}

type UniqueCount struct {
	// +optional
	Filters *Filters `json:"filters,omitempty"`

	Conditions UniqueCountConditions `json:"conditions,omitempty"`
}

type TimeRelative struct {
	// +optional
	Filters *Filters `json:"filters,omitempty"`

	Conditions TimeRelativeConditions `json:"conditions,omitempty"`
}

type Metric struct {
	// +optional
	Lucene *Lucene `json:"lucene,omitempty"`

	// +optional
	Promql *Promql `json:"promql,omitempty"`
}

type Lucene struct {
	// +optional
	SearchQuery *string `json:"searchQuery,omitempty"`

	LuceneConditions LuceneConditions `json:"conditions,omitempty"`
}

type Promql struct {
	// +optional
	SearchQuery *string `json:"searchQuery,omitempty"`

	PromqlConditions PromqlConditions `json:"conditions,omitempty"`
}

type Tracing struct {
	// +optional
	Filters *TracingFilters `json:"filters,omitempty"`

	Conditions TracingCondition `json:"conditions,omitempty"`
}

type Flow struct {
	Stages []FlowStage `json:"stages,omitempty"`
}

type StandardConditions struct {
	AlertWhen StandardAlertWhen `json:"alertWhen,omitempty"`

	// +optional
	Threshold *float64 `json:"threshold,omitempty"`

	// +optional
	TimeWindow *TimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

type RatioConditions struct {
	AlertWhen AlertWhen `json:"alertWhen,omitempty"`

	Ratio float64 `json:"ratio,omitempty"`

	TimeWindow TimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	GroupByFor *GroupByFor `json:"groupByFor,omitempty"`
}

type NewValueConditions struct {
	Key string `json:"key,omitempty"`

	TimeWindow NewValueTimeWindow `json:"timeWindow,omitempty"`
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

type TimeRelativeConditions struct {
	AlertWhen AlertWhen `json:"alertWhen,omitempty"`

	// +kubebuilder:validation:Minimum:=0
	Threshold float64 `json:"threshold,omitempty"`

	TimeWindow RelativeTimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

type ArithmeticOperator string

type LuceneConditions struct {
	MetricField string `json:"metricField,omitempty"`

	ArithmeticOperator ArithmeticOperator `json:"arithmeticOperator,omitempty"`

	// +optional
	ArithmeticOperatorModifier *int `json:"arithmeticOperatorModifier,omitempty"`

	AlertWhen AlertWhen `json:"alertWhen,omitempty"`

	// +kubebuilder:validation:Minimum:=0
	Threshold float64 `json:"threshold,omitempty"`

	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:MultipleOf:=10
	SampleThresholdPercentage int `json:"sampleThresholdPercentage,omitempty"`

	TimeWindow MetricTimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	ReplaceMissingValueWithZero *bool `json:"replaceMissingValueWithZero,omitempty"`

	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:MultipleOf:=10
	MinNonNullValuesPercentage int `json:"minNonNullValuesPercentage,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

type PromqlConditions struct {
	AlertWhen AlertWhen `json:"alertWhen,omitempty"`

	// +kubebuilder:validation:Minimum:=0
	Threshold float64 `json:"threshold,omitempty"`

	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:MultipleOf:=10
	SampleThresholdPercentage int `json:"sampleThresholdPercentage,omitempty"`

	TimeWindow MetricTimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`

	// +optional
	ReplaceMissingValueWithZero *bool `json:"replaceMissingValueWithZero,omitempty"`

	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:MultipleOf:=10
	MinNonNullValuesPercentage int `json:"minNonNullValuesPercentage,omitempty"`

	// +optional
	ManageUndetectedValues *ManageUndetectedValues `json:"manageUndetectedValues,omitempty"`
}

type TracingCondition struct {
	AlertWhen TracingAlertWhen `json:"alertWhen,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum:=0
	Threshold float64 `json:"threshold,omitempty"`

	// +optional
	TimeWindow TimeWindow `json:"timeWindow,omitempty"`

	// +optional
	GroupBy []string `json:"groupBy,omitempty"`
}

// +kubebuilder:validation:Enum=Never;5Min;10Min;1H;2H;6H;12H;24H;
type AutoRetireRatio string

// +kubebuilder:validation:Enum=More;Less;
type AlertWhen string

// +kubebuilder:validation:Enum=More;Less;Immediately;MoreThanUsual;
type StandardAlertWhen string

// +kubebuilder:validation:Enum=More;Immediately;
type TracingAlertWhen string

// +kubebuilder:validation:Enum=Q1;Q2;Both;
type GroupByFor string

// +kubebuilder:validation:Enum=5Min;10Min;15Min;20Min;30Min;1H;2H;4H;6H;12H;24H;36H;
type TimeWindow string

// +kubebuilder:validation:Enum=12H;24H;48H;72H;1W;1Month;2Month;3Month;
type NewValueTimeWindow string

// +kubebuilder:validation:Enum=1Min:5Min;10Min;15Min;20Min;30Min;1H;2H;4H;6H;12H;24H;36H;
type UniqueValueTimeWindow string

// +kubebuilder:validation:Enum=1Min:5Min;10Min;15Min;20Min;30Min;1H;2H;4H;6H;12H;24H;
type MetricTimeWindow string

// +kubebuilder:validation:Enum=Previous_hour;SameHourYesterday;SameHourLastWeek;Yesterday;SameDayLastWeek;SameDayLastMonth;
type RelativeTimeWindow string

type Filters struct {
	// +optional
	SearchQuery *string `json:"searchQuery,omitempty"`

	// +optional
	Severities []utils.Severity `json:"severities,omitempty"`

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
}

type TracingFilters struct {
	LatencyThresholdMS int `json:"latencyThresholdMS,omitempty"`

	// +optional
	TagFilters []TagFilter `json:"tagFilters,omitempty"`

	// +optional
	FieldFilters []FieldFilter `json:"fieldFilters,omitempty"`
}

type TagFilter struct {
	Field    string         `json:"field,omitempty"`
	Values   []string       `json:"values,omitempty"`
	Operator FilterOperator `json:"operator,omitempty"`
}

type FieldFilter struct {
	Field    string             `json:"field,omitempty"`
	Values   []FieldFilterValue `json:"values,omitempty"`
	Operator FilterOperator     `json:"operator,omitempty"`
}

// +kubebuilder:validation:Enum=Equals:Contains;StartWith;EndWith;
type FilterOperator string

// +kubebuilder:validation:Enum=Application:Subsystem;Service;
type FieldFilterValue string

type ManageUndetectedValues struct {
	EnableTriggeringOnUndetectedValues bool             `json:"enableTriggeringOnUndetectedValues,omitempty"`
	AutoRetireRatio                    *AutoRetireRatio `json:"autoRetireRatio,omitempty"`
}

type FlowStage struct {
	TimeWindow FlowStageTimeFrame `json:"timeWindow,omitempty"`

	Groups []FlowStageGroup `json:"groups,omitempty"`
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
	SubAlerts []SubAlert `json:"subAlerts,omitempty"`

	Operator FlowStageGroupOperator `json:"operator,omitempty"`
}

type SubAlert struct {
	// +kubebuilder:default=false
	Not bool `json:"not,omitempty"`

	// +optional
	UserAlertId string `json:"userAlertId,omitempty"`
}

// +kubebuilder:validation:Enum=And:Or;
type FlowStageGroupOperator string

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
