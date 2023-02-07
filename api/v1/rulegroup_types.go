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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Rule struct {
	// +optional
	Parse *Parse `json:"parse,omitempty"`

	// +optional
	Block *Block `json:"block,omitempty"`

	// +optional
	JsonExtract *JsonExtract `json:"jsonExtract,omitempty"`

	// +optional
	Replace *Replace `json:"replace,omitempty"`

	// +optional
	ExtractTimestamp *ExtractTimestamp `json:"extractTimestamp,omitempty"`

	// +optional
	RemoveFields *RemoveFields `json:"removeFields,omitempty"`

	// +optional
	JsonStringify *JsonStringify `json:"jsonStringify,omitempty"`

	// +optional
	Extract *Extract `json:"Extract,omitempty"`

	// +optional
	ParseJsonField *ParseJsonField `json:"parseJsonField,omitempty"`
}

type RuleParams struct {
	Name string `json:"name,omitempty"`

	// +optional
	Description *string `json:"description,omitempty"`

	// +optional
	Active *bool `json:"active,omitempty"`

	// +optional
	Order *uint `json:"order,omitempty"`
}

type Parse struct {
	RuleParams
	SourceField       string `json:"sourceField,omitempty"`
	DestinationField  string `json:"destinationField,omitempty"`
	RegularExpression string `json:"regularExpression,omitempty"`
}

type Block struct {
	RuleParams
	SourceField               string `json:"sourceField,omitempty"`
	RegularExpression         string `json:"regularExpression,omitempty"`
	KeepBlockedLogs           string `json:"keepBlockedLogs,omitempty"`
	BlockingAllMatchingBlocks bool   `json:"blockingAllMatchingBlocks,omitempty"`
}

type JsonExtract struct {
	RuleParams
	DestinationField string `json:"destinationField,omitempty"`
	JsonKey          string `json:"jsonKey,omitempty"`
}

type Replace struct {
	RuleParams
	SourceField       string `json:"sourceField,omitempty"`
	DestinationField  string `json:"destinationField,omitempty"`
	RegularExpression string `json:"regularExpression,omitempty"`
	ReplacementString string `json:"replacementString,omitempty"`
}

type ExtractTimestamp struct {
	RuleParams
	SourceField         string `json:"sourceField,omitempty"`
	FieldFormatStandard string `json:"fieldFormatStandard,omitempty"`
	TimeFormat          string `json:"timeFormat,omitempty"`
}

type RemoveFields struct {
	RuleParams
	ExcludedFields []string `json:"excludedFields,omitempty"`
}

type JsonStringify struct {
	RuleParams
	SourceField      string `json:"sourceField,omitempty"`
	DestinationField string `json:"destinationField,omitempty"`
	KeepSourceField  bool   `json:"keepSourceField,omitempty"`
}

type Extract struct {
	RuleParams
	SourceField       string `json:"sourceField,omitempty"`
	RegularExpression string `json:"regularExpression,omitempty"`
	KeepSourceField   bool   `json:"keepSourceField,omitempty"`
}

type ParseJsonField struct {
	RuleParams
	SourceField          string `json:"sourceField,omitempty"`
	DestinationField     string `json:"destinationField,omitempty"`
	KeepSourceField      bool   `json:"keepSourceField,omitempty"`
	KeepDestinationField bool   `json:"keepDestinationField,omitempty"`
}

type RuleSubgroup struct {
	// +optional
	Active *bool `json:"active,omitempty"`

	// +optional
	Order *uint `json:"order,omitempty"`

	// +optional
	Rules []Rule `json:"rules,omitempty"`
}

// RuleGroupSpec defines the desired state of RuleGroup
type RuleGroupSpec struct {
	//+kubebuilder:validation:MinLength=0
	Name string `json:"name"`

	// +optional
	Active *bool `json:"active,omitempty"`

	// +optional
	Applications []string `json:"applications,omitempty"`

	// +optional
	Subsystems []string

	// +optional
	Severities []string

	// +optional
	Hidden *bool `json:"hidden,omitempty"`

	// +optional
	Creator *string

	// +optional
	Order *uint `json:"order,omitempty"`

	// +optional
	RuleSubgroups []RuleSubgroup `json:"ruleSubgroups,omitempty"`
}

// RuleGroupStatus defines the observed state of RuleGroup
type RuleGroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RuleGroup is the Schema for the rulegroups API
type RuleGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RuleGroupSpec   `json:"spec,omitempty"`
	Status RuleGroupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RuleGroupList contains a list of RuleGroup
type RuleGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RuleGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RuleGroup{}, &RuleGroupList{})
}
