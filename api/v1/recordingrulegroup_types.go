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

	utils "coralogix-operator-poc/api"
	recordingrules "coralogix-operator-poc/controllers/clientset/grpc/recording-rules-groups/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RecordingRuleGroupSpec defines the desired state of RecordingRuleGroup
type RecordingRuleGroupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Name string `json:"name,omitempty"`

	// +optional
	IntervalSeconds int32 `json:"intervalSeconds,omitempty"`

	// +optional
	Limit int64 `json:"limit,omitempty"`

	Rules []RecordingRule `json:"rules,omitempty"`
}

func (in *RecordingRuleGroupSpec) ExtractCreateRecordingRuleGroupRequest() *recordingrules.RecordingRuleGroup {
	interval := new(uint32)
	*interval = uint32(in.IntervalSeconds)

	limit := new(uint64)
	*limit = uint64(in.Limit)

	rules := expandRecordingRules(in.Rules)

	return &recordingrules.RecordingRuleGroup{
		Name:     in.Name,
		Interval: interval,
		Limit:    limit,
		Rules:    rules,
	}
}

func expandRecordingRules(rules []RecordingRule) []*recordingrules.RecordingRule {
	result := make([]*recordingrules.RecordingRule, 0, len(rules))
	for _, r := range rules {
		rule := extractRecordingRule(r)
		result = append(result, rule)
	}
	return result
}

func extractRecordingRule(rule RecordingRule) *recordingrules.RecordingRule {
	return &recordingrules.RecordingRule{
		Record: rule.Record,
		Expr:   rule.Expr,
		Labels: rule.Labels,
	}
}

func (in *RecordingRuleGroupSpec) DeepEqual(status RecordingRuleGroupStatus) (bool, utils.Diff) {
	if limit, actualLimit := in.Limit, status.Limit; limit != actualLimit {
		return false, utils.Diff{
			Name:    "Limit",
			Desired: limit,
			Actual:  actualLimit,
		}
	}

	if name, actualName := in.Name, *status.Name; name != actualName {
		return false, utils.Diff{
			Name:    "Name",
			Desired: name,
			Actual:  actualName,
		}
	}

	if interval, actualInterval := in.IntervalSeconds, status.IntervalSeconds; interval != actualInterval {
		return false, utils.Diff{
			Name:    "Interval",
			Desired: interval,
			Actual:  actualInterval,
		}
	}

	if equal, diff := DeepEqual(in.Rules, status.Rules); !equal {
		return false, diff
	}

	return true, utils.Diff{}
}

func DeepEqual(desiredRules, actualRule []RecordingRule) (bool, utils.Diff) {
	if length, actualLength := len(desiredRules), len(actualRule); length != actualLength {
		return false, utils.Diff{
			Name:    "Rules.length",
			Desired: length,
			Actual:  actualLength,
		}
	}

	for i := range desiredRules {
		if equal, diff := desiredRules[i].DeepEqual(actualRule[i]); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Rules.%d.%s", i, diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	return true, utils.Diff{}
}

type RecordingRule struct {
	Record string `json:"record,omitempty"`

	Expr string `json:"expr,omitempty"`

	Labels map[string]string `json:"labels,omitempty"`
}

func (in *RecordingRule) DeepEqual(rule RecordingRule) (bool, utils.Diff) {
	if expr, actualExpr := in.Expr, rule.Expr; expr != actualExpr {
		return false, utils.Diff{
			Name:    "Expr",
			Desired: expr,
			Actual:  actualExpr,
		}
	}

	if record, actualRecord := in.Record, rule.Record; record != actualRecord {
		return false, utils.Diff{
			Name:    "Record",
			Desired: record,
			Actual:  actualRecord,
		}
	}

	if labels, actualLabels := in.Labels, rule.Labels; !reflect.DeepEqual(labels, actualLabels) {
		return false, utils.Diff{
			Name:    "Labels",
			Desired: labels,
			Actual:  actualLabels,
		}
	}

	return true, utils.Diff{}
}

// RecordingRuleGroupStatus defines the observed state of RecordingRuleGroup
type RecordingRuleGroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Name *string `json:"name,omitempty"`

	IntervalSeconds int32 `json:"intervalSeconds,omitempty"`

	Limit int64 `json:"limit,omitempty"`

	Rules []RecordingRule `json:"rules,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RecordingRuleGroup is the Schema for the recordingrulegroups API
type RecordingRuleGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordingRuleGroupSpec   `json:"spec,omitempty"`
	Status RecordingRuleGroupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RecordingRuleGroupList contains a list of RecordingRuleGroup
type RecordingRuleGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecordingRuleGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RecordingRuleGroup{}, &RecordingRuleGroupList{})
}
