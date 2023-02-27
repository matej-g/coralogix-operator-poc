//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Alert) DeepCopyInto(out *Alert) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Alert.
func (in *Alert) DeepCopy() *Alert {
	if in == nil {
		return nil
	}
	out := new(Alert)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Alert) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AlertList) DeepCopyInto(out *AlertList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Alert, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AlertList.
func (in *AlertList) DeepCopy() *AlertList {
	if in == nil {
		return nil
	}
	out := new(AlertList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AlertList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AlertSpec) DeepCopyInto(out *AlertSpec) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ExpirationDate != nil {
		in, out := &in.ExpirationDate, &out.ExpirationDate
		*out = new(ExpirationDate)
		**out = **in
	}
	if in.Notifications != nil {
		in, out := &in.Notifications, &out.Notifications
		*out = new(Notifications)
		(*in).DeepCopyInto(*out)
	}
	if in.Scheduling != nil {
		in, out := &in.Scheduling, &out.Scheduling
		*out = new(Scheduling)
		(*in).DeepCopyInto(*out)
	}
	in.AlertType.DeepCopyInto(&out.AlertType)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AlertSpec.
func (in *AlertSpec) DeepCopy() *AlertSpec {
	if in == nil {
		return nil
	}
	out := new(AlertSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AlertStatus) DeepCopyInto(out *AlertStatus) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AlertStatus.
func (in *AlertStatus) DeepCopy() *AlertStatus {
	if in == nil {
		return nil
	}
	out := new(AlertStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AlertType) DeepCopyInto(out *AlertType) {
	*out = *in
	if in.Standard != nil {
		in, out := &in.Standard, &out.Standard
		*out = new(Standard)
		(*in).DeepCopyInto(*out)
	}
	if in.Ratio != nil {
		in, out := &in.Ratio, &out.Ratio
		*out = new(Ratio)
		(*in).DeepCopyInto(*out)
	}
	if in.NewValue != nil {
		in, out := &in.NewValue, &out.NewValue
		*out = new(NewValue)
		(*in).DeepCopyInto(*out)
	}
	if in.UniqueCount != nil {
		in, out := &in.UniqueCount, &out.UniqueCount
		*out = new(UniqueCount)
		(*in).DeepCopyInto(*out)
	}
	if in.TimeRelative != nil {
		in, out := &in.TimeRelative, &out.TimeRelative
		*out = new(TimeRelative)
		(*in).DeepCopyInto(*out)
	}
	if in.Metric != nil {
		in, out := &in.Metric, &out.Metric
		*out = new(Metric)
		(*in).DeepCopyInto(*out)
	}
	if in.Tracing != nil {
		in, out := &in.Tracing, &out.Tracing
		*out = new(Tracing)
		(*in).DeepCopyInto(*out)
	}
	if in.Flow != nil {
		in, out := &in.Flow, &out.Flow
		*out = new(Flow)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AlertType.
func (in *AlertType) DeepCopy() *AlertType {
	if in == nil {
		return nil
	}
	out := new(AlertType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Block) DeepCopyInto(out *Block) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Block.
func (in *Block) DeepCopy() *Block {
	if in == nil {
		return nil
	}
	out := new(Block)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExpirationDate) DeepCopyInto(out *ExpirationDate) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExpirationDate.
func (in *ExpirationDate) DeepCopy() *ExpirationDate {
	if in == nil {
		return nil
	}
	out := new(ExpirationDate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Extract) DeepCopyInto(out *Extract) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Extract.
func (in *Extract) DeepCopy() *Extract {
	if in == nil {
		return nil
	}
	out := new(Extract)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtractTimestamp) DeepCopyInto(out *ExtractTimestamp) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtractTimestamp.
func (in *ExtractTimestamp) DeepCopy() *ExtractTimestamp {
	if in == nil {
		return nil
	}
	out := new(ExtractTimestamp)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FieldFilter) DeepCopyInto(out *FieldFilter) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]FieldFilterValue, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FieldFilter.
func (in *FieldFilter) DeepCopy() *FieldFilter {
	if in == nil {
		return nil
	}
	out := new(FieldFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Filters) DeepCopyInto(out *Filters) {
	*out = *in
	if in.SearchQuery != nil {
		in, out := &in.SearchQuery, &out.SearchQuery
		*out = new(string)
		**out = **in
	}
	if in.Severities != nil {
		in, out := &in.Severities, &out.Severities
		*out = make([]FiltersLogSeverity, len(*in))
		copy(*out, *in)
	}
	if in.Applications != nil {
		in, out := &in.Applications, &out.Applications
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Subsystems != nil {
		in, out := &in.Subsystems, &out.Subsystems
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Categories != nil {
		in, out := &in.Categories, &out.Categories
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Computers != nil {
		in, out := &in.Computers, &out.Computers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Classes != nil {
		in, out := &in.Classes, &out.Classes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Methods != nil {
		in, out := &in.Methods, &out.Methods
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.IPs != nil {
		in, out := &in.IPs, &out.IPs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Filters.
func (in *Filters) DeepCopy() *Filters {
	if in == nil {
		return nil
	}
	out := new(Filters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Flow) DeepCopyInto(out *Flow) {
	*out = *in
	if in.Stages != nil {
		in, out := &in.Stages, &out.Stages
		*out = make([]FlowStage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Flow.
func (in *Flow) DeepCopy() *Flow {
	if in == nil {
		return nil
	}
	out := new(Flow)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlowStage) DeepCopyInto(out *FlowStage) {
	*out = *in
	out.TimeWindow = in.TimeWindow
	if in.Groups != nil {
		in, out := &in.Groups, &out.Groups
		*out = make([]FlowStageGroup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlowStage.
func (in *FlowStage) DeepCopy() *FlowStage {
	if in == nil {
		return nil
	}
	out := new(FlowStage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlowStageGroup) DeepCopyInto(out *FlowStageGroup) {
	*out = *in
	if in.SubAlerts != nil {
		in, out := &in.SubAlerts, &out.SubAlerts
		*out = make([]SubAlert, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlowStageGroup.
func (in *FlowStageGroup) DeepCopy() *FlowStageGroup {
	if in == nil {
		return nil
	}
	out := new(FlowStageGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlowStageTimeFrame) DeepCopyInto(out *FlowStageTimeFrame) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlowStageTimeFrame.
func (in *FlowStageTimeFrame) DeepCopy() *FlowStageTimeFrame {
	if in == nil {
		return nil
	}
	out := new(FlowStageTimeFrame)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JsonExtract) DeepCopyInto(out *JsonExtract) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JsonExtract.
func (in *JsonExtract) DeepCopy() *JsonExtract {
	if in == nil {
		return nil
	}
	out := new(JsonExtract)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JsonStringify) DeepCopyInto(out *JsonStringify) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JsonStringify.
func (in *JsonStringify) DeepCopy() *JsonStringify {
	if in == nil {
		return nil
	}
	out := new(JsonStringify)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Lucene) DeepCopyInto(out *Lucene) {
	*out = *in
	if in.SearchQuery != nil {
		in, out := &in.SearchQuery, &out.SearchQuery
		*out = new(string)
		**out = **in
	}
	in.LuceneConditions.DeepCopyInto(&out.LuceneConditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Lucene.
func (in *Lucene) DeepCopy() *Lucene {
	if in == nil {
		return nil
	}
	out := new(Lucene)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LuceneConditions) DeepCopyInto(out *LuceneConditions) {
	*out = *in
	if in.ArithmeticOperatorModifier != nil {
		in, out := &in.ArithmeticOperatorModifier, &out.ArithmeticOperatorModifier
		*out = new(int)
		**out = **in
	}
	out.Threshold = in.Threshold.DeepCopy()
	if in.GroupBy != nil {
		in, out := &in.GroupBy, &out.GroupBy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ReplaceMissingValueWithZero != nil {
		in, out := &in.ReplaceMissingValueWithZero, &out.ReplaceMissingValueWithZero
		*out = new(bool)
		**out = **in
	}
	if in.ManageUndetectedValues != nil {
		in, out := &in.ManageUndetectedValues, &out.ManageUndetectedValues
		*out = new(ManageUndetectedValues)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LuceneConditions.
func (in *LuceneConditions) DeepCopy() *LuceneConditions {
	if in == nil {
		return nil
	}
	out := new(LuceneConditions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManageUndetectedValues) DeepCopyInto(out *ManageUndetectedValues) {
	*out = *in
	if in.AutoRetireRatio != nil {
		in, out := &in.AutoRetireRatio, &out.AutoRetireRatio
		*out = new(AutoRetireRatio)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManageUndetectedValues.
func (in *ManageUndetectedValues) DeepCopy() *ManageUndetectedValues {
	if in == nil {
		return nil
	}
	out := new(ManageUndetectedValues)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Metric) DeepCopyInto(out *Metric) {
	*out = *in
	if in.Lucene != nil {
		in, out := &in.Lucene, &out.Lucene
		*out = new(Lucene)
		(*in).DeepCopyInto(*out)
	}
	if in.Promql != nil {
		in, out := &in.Promql, &out.Promql
		*out = new(Promql)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Metric.
func (in *Metric) DeepCopy() *Metric {
	if in == nil {
		return nil
	}
	out := new(Metric)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NewValue) DeepCopyInto(out *NewValue) {
	*out = *in
	if in.Filters != nil {
		in, out := &in.Filters, &out.Filters
		*out = new(Filters)
		(*in).DeepCopyInto(*out)
	}
	out.Conditions = in.Conditions
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NewValue.
func (in *NewValue) DeepCopy() *NewValue {
	if in == nil {
		return nil
	}
	out := new(NewValue)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NewValueConditions) DeepCopyInto(out *NewValueConditions) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NewValueConditions.
func (in *NewValueConditions) DeepCopy() *NewValueConditions {
	if in == nil {
		return nil
	}
	out := new(NewValueConditions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Notifications) DeepCopyInto(out *Notifications) {
	*out = *in
	if in.OnTriggerAndResolved != nil {
		in, out := &in.OnTriggerAndResolved, &out.OnTriggerAndResolved
		*out = new(bool)
		**out = **in
	}
	if in.IgnoreInfinity != nil {
		in, out := &in.IgnoreInfinity, &out.IgnoreInfinity
		*out = new(bool)
		**out = **in
	}
	if in.NotifyOnlyOnTriggeredGroupByValues != nil {
		in, out := &in.NotifyOnlyOnTriggeredGroupByValues, &out.NotifyOnlyOnTriggeredGroupByValues
		*out = new(bool)
		**out = **in
	}
	in.Recipients.DeepCopyInto(&out.Recipients)
	if in.NotifyEveryMin != nil {
		in, out := &in.NotifyEveryMin, &out.NotifyEveryMin
		*out = new(int)
		**out = **in
	}
	if in.PayloadFilters != nil {
		in, out := &in.PayloadFilters, &out.PayloadFilters
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Notifications.
func (in *Notifications) DeepCopy() *Notifications {
	if in == nil {
		return nil
	}
	out := new(Notifications)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Parse) DeepCopyInto(out *Parse) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Parse.
func (in *Parse) DeepCopy() *Parse {
	if in == nil {
		return nil
	}
	out := new(Parse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParseJsonField) DeepCopyInto(out *ParseJsonField) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParseJsonField.
func (in *ParseJsonField) DeepCopy() *ParseJsonField {
	if in == nil {
		return nil
	}
	out := new(ParseJsonField)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Promql) DeepCopyInto(out *Promql) {
	*out = *in
	if in.SearchQuery != nil {
		in, out := &in.SearchQuery, &out.SearchQuery
		*out = new(string)
		**out = **in
	}
	in.PromqlConditions.DeepCopyInto(&out.PromqlConditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Promql.
func (in *Promql) DeepCopy() *Promql {
	if in == nil {
		return nil
	}
	out := new(Promql)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PromqlConditions) DeepCopyInto(out *PromqlConditions) {
	*out = *in
	out.Threshold = in.Threshold.DeepCopy()
	if in.GroupBy != nil {
		in, out := &in.GroupBy, &out.GroupBy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ReplaceMissingValueWithZero != nil {
		in, out := &in.ReplaceMissingValueWithZero, &out.ReplaceMissingValueWithZero
		*out = new(bool)
		**out = **in
	}
	if in.ManageUndetectedValues != nil {
		in, out := &in.ManageUndetectedValues, &out.ManageUndetectedValues
		*out = new(ManageUndetectedValues)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromqlConditions.
func (in *PromqlConditions) DeepCopy() *PromqlConditions {
	if in == nil {
		return nil
	}
	out := new(PromqlConditions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ratio) DeepCopyInto(out *Ratio) {
	*out = *in
	in.Query1Filters.DeepCopyInto(&out.Query1Filters)
	in.Query2Filters.DeepCopyInto(&out.Query2Filters)
	in.Conditions.DeepCopyInto(&out.Conditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ratio.
func (in *Ratio) DeepCopy() *Ratio {
	if in == nil {
		return nil
	}
	out := new(Ratio)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatioConditions) DeepCopyInto(out *RatioConditions) {
	*out = *in
	out.Ratio = in.Ratio.DeepCopy()
	if in.GroupBy != nil {
		in, out := &in.GroupBy, &out.GroupBy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.GroupByFor != nil {
		in, out := &in.GroupByFor, &out.GroupByFor
		*out = new(GroupByFor)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatioConditions.
func (in *RatioConditions) DeepCopy() *RatioConditions {
	if in == nil {
		return nil
	}
	out := new(RatioConditions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Recipients) DeepCopyInto(out *Recipients) {
	*out = *in
	if in.Emails != nil {
		in, out := &in.Emails, &out.Emails
		*out = make([]Email, len(*in))
		copy(*out, *in)
	}
	if in.Webhooks != nil {
		in, out := &in.Webhooks, &out.Webhooks
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Recipients.
func (in *Recipients) DeepCopy() *Recipients {
	if in == nil {
		return nil
	}
	out := new(Recipients)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoveFields) DeepCopyInto(out *RemoveFields) {
	*out = *in
	if in.ExcludedFields != nil {
		in, out := &in.ExcludedFields, &out.ExcludedFields
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoveFields.
func (in *RemoveFields) DeepCopy() *RemoveFields {
	if in == nil {
		return nil
	}
	out := new(RemoveFields)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Replace) DeepCopyInto(out *Replace) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Replace.
func (in *Replace) DeepCopy() *Replace {
	if in == nil {
		return nil
	}
	out := new(Replace)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rule) DeepCopyInto(out *Rule) {
	*out = *in
	if in.Parse != nil {
		in, out := &in.Parse, &out.Parse
		*out = new(Parse)
		**out = **in
	}
	if in.Block != nil {
		in, out := &in.Block, &out.Block
		*out = new(Block)
		**out = **in
	}
	if in.JsonExtract != nil {
		in, out := &in.JsonExtract, &out.JsonExtract
		*out = new(JsonExtract)
		**out = **in
	}
	if in.Replace != nil {
		in, out := &in.Replace, &out.Replace
		*out = new(Replace)
		**out = **in
	}
	if in.ExtractTimestamp != nil {
		in, out := &in.ExtractTimestamp, &out.ExtractTimestamp
		*out = new(ExtractTimestamp)
		**out = **in
	}
	if in.RemoveFields != nil {
		in, out := &in.RemoveFields, &out.RemoveFields
		*out = new(RemoveFields)
		(*in).DeepCopyInto(*out)
	}
	if in.JsonStringify != nil {
		in, out := &in.JsonStringify, &out.JsonStringify
		*out = new(JsonStringify)
		**out = **in
	}
	if in.Extract != nil {
		in, out := &in.Extract, &out.Extract
		*out = new(Extract)
		**out = **in
	}
	if in.ParseJsonField != nil {
		in, out := &in.ParseJsonField, &out.ParseJsonField
		*out = new(ParseJsonField)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rule.
func (in *Rule) DeepCopy() *Rule {
	if in == nil {
		return nil
	}
	out := new(Rule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuleGroup) DeepCopyInto(out *RuleGroup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuleGroup.
func (in *RuleGroup) DeepCopy() *RuleGroup {
	if in == nil {
		return nil
	}
	out := new(RuleGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RuleGroup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuleGroupList) DeepCopyInto(out *RuleGroupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RuleGroup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuleGroupList.
func (in *RuleGroupList) DeepCopy() *RuleGroupList {
	if in == nil {
		return nil
	}
	out := new(RuleGroupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RuleGroupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuleGroupSpec) DeepCopyInto(out *RuleGroupSpec) {
	*out = *in
	if in.Applications != nil {
		in, out := &in.Applications, &out.Applications
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Subsystems != nil {
		in, out := &in.Subsystems, &out.Subsystems
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Severities != nil {
		in, out := &in.Severities, &out.Severities
		*out = make([]RuleSeverity, len(*in))
		copy(*out, *in)
	}
	if in.Order != nil {
		in, out := &in.Order, &out.Order
		*out = new(int32)
		**out = **in
	}
	if in.RuleSubgroups != nil {
		in, out := &in.RuleSubgroups, &out.RuleSubgroups
		*out = make([]RuleSubGroup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuleGroupSpec.
func (in *RuleGroupSpec) DeepCopy() *RuleGroupSpec {
	if in == nil {
		return nil
	}
	out := new(RuleGroupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuleGroupStatus) DeepCopyInto(out *RuleGroupStatus) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuleGroupStatus.
func (in *RuleGroupStatus) DeepCopy() *RuleGroupStatus {
	if in == nil {
		return nil
	}
	out := new(RuleGroupStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuleSubGroup) DeepCopyInto(out *RuleSubGroup) {
	*out = *in
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]Rule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuleSubGroup.
func (in *RuleSubGroup) DeepCopy() *RuleSubGroup {
	if in == nil {
		return nil
	}
	out := new(RuleSubGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Scheduling) DeepCopyInto(out *Scheduling) {
	*out = *in
	if in.DaysEnabled != nil {
		in, out := &in.DaysEnabled, &out.DaysEnabled
		*out = make([]Day, len(*in))
		copy(*out, *in)
	}
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		*out = new(Time)
		**out = **in
	}
	if in.EndTime != nil {
		in, out := &in.EndTime, &out.EndTime
		*out = new(Time)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scheduling.
func (in *Scheduling) DeepCopy() *Scheduling {
	if in == nil {
		return nil
	}
	out := new(Scheduling)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Standard) DeepCopyInto(out *Standard) {
	*out = *in
	if in.Filters != nil {
		in, out := &in.Filters, &out.Filters
		*out = new(Filters)
		(*in).DeepCopyInto(*out)
	}
	in.Conditions.DeepCopyInto(&out.Conditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Standard.
func (in *Standard) DeepCopy() *Standard {
	if in == nil {
		return nil
	}
	out := new(Standard)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StandardConditions) DeepCopyInto(out *StandardConditions) {
	*out = *in
	if in.Threshold != nil {
		in, out := &in.Threshold, &out.Threshold
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.TimeWindow != nil {
		in, out := &in.TimeWindow, &out.TimeWindow
		*out = new(TimeWindow)
		**out = **in
	}
	if in.GroupBy != nil {
		in, out := &in.GroupBy, &out.GroupBy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ManageUndetectedValues != nil {
		in, out := &in.ManageUndetectedValues, &out.ManageUndetectedValues
		*out = new(ManageUndetectedValues)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StandardConditions.
func (in *StandardConditions) DeepCopy() *StandardConditions {
	if in == nil {
		return nil
	}
	out := new(StandardConditions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubAlert) DeepCopyInto(out *SubAlert) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubAlert.
func (in *SubAlert) DeepCopy() *SubAlert {
	if in == nil {
		return nil
	}
	out := new(SubAlert)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TagFilter) DeepCopyInto(out *TagFilter) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TagFilter.
func (in *TagFilter) DeepCopy() *TagFilter {
	if in == nil {
		return nil
	}
	out := new(TagFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimeRelative) DeepCopyInto(out *TimeRelative) {
	*out = *in
	if in.Filters != nil {
		in, out := &in.Filters, &out.Filters
		*out = new(Filters)
		(*in).DeepCopyInto(*out)
	}
	in.Conditions.DeepCopyInto(&out.Conditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimeRelative.
func (in *TimeRelative) DeepCopy() *TimeRelative {
	if in == nil {
		return nil
	}
	out := new(TimeRelative)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimeRelativeConditions) DeepCopyInto(out *TimeRelativeConditions) {
	*out = *in
	out.Threshold = in.Threshold.DeepCopy()
	if in.GroupBy != nil {
		in, out := &in.GroupBy, &out.GroupBy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ManageUndetectedValues != nil {
		in, out := &in.ManageUndetectedValues, &out.ManageUndetectedValues
		*out = new(ManageUndetectedValues)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimeRelativeConditions.
func (in *TimeRelativeConditions) DeepCopy() *TimeRelativeConditions {
	if in == nil {
		return nil
	}
	out := new(TimeRelativeConditions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tracing) DeepCopyInto(out *Tracing) {
	*out = *in
	if in.Filters != nil {
		in, out := &in.Filters, &out.Filters
		*out = new(TracingFilters)
		(*in).DeepCopyInto(*out)
	}
	in.Conditions.DeepCopyInto(&out.Conditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tracing.
func (in *Tracing) DeepCopy() *Tracing {
	if in == nil {
		return nil
	}
	out := new(Tracing)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TracingCondition) DeepCopyInto(out *TracingCondition) {
	*out = *in
	out.Threshold = in.Threshold.DeepCopy()
	if in.GroupBy != nil {
		in, out := &in.GroupBy, &out.GroupBy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TracingCondition.
func (in *TracingCondition) DeepCopy() *TracingCondition {
	if in == nil {
		return nil
	}
	out := new(TracingCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TracingFilters) DeepCopyInto(out *TracingFilters) {
	*out = *in
	if in.TagFilters != nil {
		in, out := &in.TagFilters, &out.TagFilters
		*out = make([]TagFilter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.FieldFilters != nil {
		in, out := &in.FieldFilters, &out.FieldFilters
		*out = make([]FieldFilter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TracingFilters.
func (in *TracingFilters) DeepCopy() *TracingFilters {
	if in == nil {
		return nil
	}
	out := new(TracingFilters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UniqueCount) DeepCopyInto(out *UniqueCount) {
	*out = *in
	if in.Filters != nil {
		in, out := &in.Filters, &out.Filters
		*out = new(Filters)
		(*in).DeepCopyInto(*out)
	}
	in.Conditions.DeepCopyInto(&out.Conditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UniqueCount.
func (in *UniqueCount) DeepCopy() *UniqueCount {
	if in == nil {
		return nil
	}
	out := new(UniqueCount)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UniqueCountConditions) DeepCopyInto(out *UniqueCountConditions) {
	*out = *in
	if in.GroupBy != nil {
		in, out := &in.GroupBy, &out.GroupBy
		*out = new(string)
		**out = **in
	}
	if in.MaxUniqueValuesForGroupBy != nil {
		in, out := &in.MaxUniqueValuesForGroupBy, &out.MaxUniqueValuesForGroupBy
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UniqueCountConditions.
func (in *UniqueCountConditions) DeepCopy() *UniqueCountConditions {
	if in == nil {
		return nil
	}
	out := new(UniqueCountConditions)
	in.DeepCopyInto(out)
	return out
}
