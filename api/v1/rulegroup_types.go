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
	utils "coralogix-operator-poc/controllers"
	rulesgroups "coralogix-operator-poc/controllers/clientset/grpc/rules-groups/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	rulesSchemaSeverityToProtoSeverity = map[RuleGroupSeverity]string{
		"Debug":    "VALUE_DEBUG_OR_UNSPECIFIED",
		"Verbose":  "VALUE_VERBOSE",
		"Info":     "VALUE_INFO",
		"Warning":  "VALUE_WARNING",
		"Error":    "VALUE_ERROR",
		"Critical": "VALUE_CRITICAL",
	}
	rulesProtoSeverityToSchemaSeverity                         = utils.ReverseMap(rulesSchemaSeverityToProtoSeverity)
	rulesSchemaDestinationFieldToProtoSeverityDestinationField = map[DestinationField]string{
		"Category":  "DESTINATION_FIELD_CATEGORY_OR_UNSPECIFIED",
		"ClassName": "DESTINATION_FIELD_CLASSNAME",
		"Method":    "DESTINATION_FIELD_METHODNAME",
		"ThreadID":  "DESTINATION_FIELD_THREADID",
		"Severity":  "DESTINATION_FIELD_SEVERITY",
	}
	rulesProtoSeverityDestinationFieldToSchemaDestinationField = utils.ReverseMap(rulesSchemaDestinationFieldToProtoSeverityDestinationField)
	rulesSchemaFormatStandardToProtoFormatStandard             = map[FieldFormatStandard]string{
		"Strftime": "FORMAT_STANDARD_STRFTIME_OR_UNSPECIFIED",
		"JavaSDF":  "FORMAT_STANDARD_JAVASDF",
		"Golang":   "FORMAT_STANDARD_GOLANG",
		"SecondTS": "FORMAT_STANDARD_SECONDSTS",
		"MilliTS":  "FORMAT_STANDARD_MILLITS",
		"MicroTS":  "FORMAT_STANDARD_MICROTS",
		"NanoTS":   "FORMAT_STANDARD_NANOTS",
	}
	rulesProtoFormatStandardToSchemaFormatStandard = utils.ReverseMap(rulesSchemaFormatStandardToProtoFormatStandard)
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Rule struct {
	Name string `json:"name,omitempty"`

	// +optional
	Description *string `json:"description,omitempty"`

	//+kubebuilder:default=true
	Active bool `json:"active,omitempty"`

	// +optional
	Order *uint32 `json:"order,omitempty"`

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
	Extract *Extract `json:"extract,omitempty"`

	// +optional
	ParseJsonField *ParseJsonField `json:"parseJsonField,omitempty"`
}

func (in *Rule) DeepEqual(rule *rulesgroups.Rule) bool {
	if in.Active != rule.Enabled.GetValue() {
		return false
	}
	if in.Order == nil || *in.Order != rule.Order.GetValue() {
		return false
	}
	if (in.Description == nil && rule.Description != nil) || (*in.Description != rule.Description.GetValue()) {
		return false
	}
	if in.Name != rule.Name.GetValue() {
		return false
	}
	if !in.DeepEqualRuleType(rule) {
		return false
	}
	return true
}

func (in *Rule) DeepEqualRuleType(parameters *rulesgroups.Rule) bool {
	switch typeParameters := parameters.Parameters.RuleParameters.(type) {
	case *rulesgroups.RuleParameters_ExtractParameters:
		if extractParameters := in.Extract; extractParameters == nil {
			return false
		} else {
			actualExtractParameters := typeParameters.ExtractParameters
			if extractParameters.RegularExpression != actualExtractParameters.Rule.GetValue() {
				return false
			}
			if extractParameters.SourceField != parameters.SourceField.GetValue() {
				return false
			}
		}
	case *rulesgroups.RuleParameters_JsonExtractParameters:
		if jsonExtractParameters := in.JsonExtract; jsonExtractParameters == nil {
			return false
		} else {
			actualJsonExtractParameters := typeParameters.JsonExtractParameters
			if jsonExtractParameters.JsonKey != actualJsonExtractParameters.Rule.GetValue() {
				return false
			}
			if rulesSchemaDestinationFieldToProtoSeverityDestinationField[jsonExtractParameters.DestinationField] !=
				actualJsonExtractParameters.DestinationField.String() {
				return false
			}
		}
	case *rulesgroups.RuleParameters_ReplaceParameters:
		if replaceParameters := in.Replace; replaceParameters == nil {
			return false
		} else {
			actualReplaceParameters := typeParameters.ReplaceParameters
			if replaceParameters.ReplacementString != actualReplaceParameters.ReplaceNewVal.GetValue() {
				return false
			}
			if rulesSchemaDestinationFieldToProtoSeverityDestinationField[replaceParameters.DestinationField] !=
				actualReplaceParameters.DestinationField.String() {
				return false
			}
			if replaceParameters.RegularExpression != actualReplaceParameters.Rule.GetValue() {
				return false
			}
			if replaceParameters.SourceField != parameters.SourceField.GetValue() {
				return false
			}
		}
	case *rulesgroups.RuleParameters_ParseParameters:
		if parseParameters := in.Parse; parseParameters == nil {
			return false
		} else {
			actualParseParameters := typeParameters.ParseParameters
			if parseParameters.RegularExpression != actualParseParameters.Rule.GetValue() {
				return false
			}
			if parseParameters.DestinationField != actualParseParameters.DestinationField.String() {
				return false
			}
			if parseParameters.SourceField != parameters.SourceField.GetValue() {
				return false
			}
		}
	case *rulesgroups.RuleParameters_AllowParameters:
		if allowParameters := in.Block; allowParameters == nil || allowParameters.BlockingAllMatchingBlocks {
			return false
		} else {
			actualAllowParameters := typeParameters.AllowParameters
			if allowParameters.RegularExpression != actualAllowParameters.Rule.GetValue() {
				return false
			}
			if allowParameters.KeepBlockedLogs != actualAllowParameters.KeepBlockedLogs.GetValue() {
				return false
			}
			if allowParameters.SourceField != parameters.SourceField.GetValue() {
				return false
			}
		}
	case *rulesgroups.RuleParameters_BlockParameters:
		if blockParameters := in.Block; blockParameters == nil || !blockParameters.BlockingAllMatchingBlocks {
			return false
		} else {
			actualBlockParameters := typeParameters.BlockParameters
			if blockParameters.RegularExpression != actualBlockParameters.Rule.GetValue() {
				return false
			}
			if blockParameters.KeepBlockedLogs != actualBlockParameters.KeepBlockedLogs.GetValue() {
				return false
			}
			if blockParameters.SourceField != parameters.SourceField.GetValue() {
				return false
			}
		}
	case *rulesgroups.RuleParameters_ExtractTimestampParameters:
		if extractTimestampParameters := in.ExtractTimestamp; extractTimestampParameters == nil {
			return false
		} else {
			actualExtractTimestampParameters := typeParameters.ExtractTimestampParameters
			if rulesSchemaFormatStandardToProtoFormatStandard[extractTimestampParameters.FieldFormatStandard] != actualExtractTimestampParameters.Standard.String() {
				return false
			}
			if extractTimestampParameters.TimeFormat != actualExtractTimestampParameters.Format.GetValue() {
				return false
			}
			if extractTimestampParameters.SourceField != parameters.SourceField.GetValue() {
				return false
			}
		}
	case *rulesgroups.RuleParameters_RemoveFieldsParameters:
		if removeFieldsParameters := in.RemoveFields; removeFieldsParameters == nil {
			return false
		} else if len(removeFieldsParameters.ExcludedFields) != len(removeFieldsParameters.ExcludedFields) ||
			!utils.SlicesWithUniqueValuesEqual(removeFieldsParameters.ExcludedFields, removeFieldsParameters.ExcludedFields) {
			return false
		}
	case *rulesgroups.RuleParameters_JsonStringifyParameters:
		if jsonStringifyParameters := in.JsonStringify; jsonStringifyParameters == nil {
			return false
		} else {
			actualJsonStringifyParameters := typeParameters.JsonStringifyParameters
			if jsonStringifyParameters.KeepSourceField == actualJsonStringifyParameters.DeleteSource.GetValue() {
				return false
			}
			if rulesSchemaDestinationFieldToProtoSeverityDestinationField[jsonStringifyParameters.DestinationField] !=
				actualJsonStringifyParameters.DestinationField.String() {
				return false
			}
			if jsonStringifyParameters.SourceField != parameters.SourceField.GetValue() {
				return false
			}
		}
	case *rulesgroups.RuleParameters_JsonParseParameters:
		if jsonParseParameters := in.ParseJsonField; jsonParseParameters == nil {
			return false
		} else {
			actualJsonParseParameters := typeParameters.JsonParseParameters
			if jsonParseParameters.DestinationField != actualJsonParseParameters.DestinationField.String() {
				return false
			}
			if jsonParseParameters.SourceField != parameters.SourceField.GetValue() {
				return false
			}
			if jsonParseParameters.KeepSourceField == actualJsonParseParameters.DeleteSource.GetValue() {
				return false
			}
			if jsonParseParameters.KeepDestinationField == actualJsonParseParameters.OverrideDest.GetValue() {
				return false
			}
		}
	}
	return true
}

type Parse struct {
	SourceField string `json:"sourceField,omitempty"`

	DestinationField string `json:"destinationField,omitempty"`

	RegularExpression string `json:"regularExpression,omitempty"`
}

type Block struct {
	SourceField string `json:"sourceField,omitempty"`

	RegularExpression string `json:"regularExpression,omitempty"`

	//+kubebuilder:default=false
	KeepBlockedLogs bool `json:"keepBlockedLogs,omitempty"`

	//+kubebuilder:default=true
	BlockingAllMatchingBlocks bool `json:"blockingAllMatchingBlocks,omitempty"`
}

// +kubebuilder:validation:Enum=Category;CLASSNAME;METHODNAME;THREADID;SEVERITY
type DestinationField string

type JsonExtract struct {
	DestinationField DestinationField `json:"destinationField,omitempty"`

	JsonKey string `json:"jsonKey,omitempty"`
}

type Replace struct {
	SourceField string `json:"sourceField,omitempty"`

	DestinationField DestinationField `json:"destinationField,omitempty"`

	RegularExpression string `json:"regularExpression,omitempty"`

	ReplacementString string `json:"replacementString,omitempty"`
}

// +kubebuilder:validation:Enum=Strftime;JavaSDF;Golang;SecondTS;MilliTS;MicroTS;NanoTS
type FieldFormatStandard string

type ExtractTimestamp struct {
	SourceField string `json:"sourceField,omitempty"`

	FieldFormatStandard FieldFormatStandard `json:"fieldFormatStandard,omitempty"`

	TimeFormat string `json:"timeFormat,omitempty"`
}

type RemoveFields struct {
	ExcludedFields []string `json:"excludedFields,omitempty"`
}

type JsonStringify struct {
	SourceField string `json:"sourceField,omitempty"`

	DestinationField DestinationField `json:"destinationField,omitempty"`

	//+kubebuilder:default=false
	KeepSourceField bool `json:"keepSourceField,omitempty"`
}

type Extract struct {
	SourceField string `json:"sourceField,omitempty"`

	RegularExpression string `json:"regularExpression,omitempty"`
}

type ParseJsonField struct {
	SourceField string `json:"sourceField,omitempty"`

	DestinationField string `json:"destinationField,omitempty"`

	KeepSourceField bool `json:"keepSourceField,omitempty"`

	KeepDestinationField bool `json:"keepDestinationField,omitempty"`
}

type RuleSubgroup struct {
	//+kubebuilder:default=true
	Active bool `json:"active,omitempty"`

	// +optional
	Order *uint32 `json:"order,omitempty"`

	// +optional
	Rules []Rule `json:"rules,omitempty"`
}

func (in *RuleSubgroup) DeepEqual(subgroup *rulesgroups.RuleSubgroup) bool {
	if in.Active != subgroup.Enabled.GetValue() {
		return false
	}
	if in.Order == nil || *in.Order != subgroup.Order.GetValue() {
		return false
	}
	if len(in.Rules) != len(subgroup.Rules) {
		return false
	}
	for i := range in.Rules {
		if !in.Rules[i].DeepEqual(subgroup.Rules[i]) {
			return false
		}
	}
	return true
}

// +kubebuilder:validation:Enum=Debug;Verbose;Info;Warning;Eror;Critical
type RuleGroupSeverity string

// RuleGroupSpec defines the desired state of RuleGroup
type RuleGroupSpec struct {
	//+kubebuilder:validation:MinLength=0
	Name string `json:"name"`

	// +optional
	Description string `json:"description,omitempty"`

	//+kubebuilder:default=true
	Active bool `json:"active,omitempty"`

	// +optional
	// +kubebuilder:validation:UniqueItems
	Applications []string `json:"applications,omitempty"`

	// +optional
	// +kubebuilder:validation:UniqueItems
	Subsystems []string `json:"subsystems,omitempty"`

	// +optional
	// +kubebuilder:validation:UniqueItems
	Severities []RuleGroupSeverity `json:"severities,omitempty"`

	//+kubebuilder:default=false
	Hidden bool `json:"hidden,omitempty"`

	// +optional
	Creator string `json:"creator,omitempty"`

	// +optional
	Order *uint32 `json:"order,omitempty"`

	// +optional
	RuleSubgroups []RuleSubgroup `json:"ruleSubgroups,omitempty"`
}

func (in *RuleGroupSpec) DeepEqual(actualState *rulesgroups.RuleGroup) bool {
	if in.Name != actualState.GetName().GetValue() {
		return false
	}

	if in.Description != actualState.GetDescription().GetValue() {
		return false
	}

	if in.Active != actualState.GetEnabled().GetValue() {
		return false
	}

	if in.Hidden != actualState.GetHidden().GetValue() {
		return false
	}

	if in.Creator != actualState.GetCreator().GetValue() {
		return false
	}

	if in.Order == nil || *in.Order != actualState.GetOrder().GetValue() {
		return false
	}

	applications, subsystems, severities := flattenRuleMatchers(actualState.RuleMatchers)
	if !utils.SlicesWithUniqueValuesEqual(in.Applications, applications) {
		return false
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Subsystems, subsystems) {
		return false
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Severities, severities) {
		return false
	}

	if len(in.RuleSubgroups) != len(actualState.RuleSubgroups) {
		return false
	}

	for i := range in.RuleSubgroups {
		if !in.RuleSubgroups[i].DeepEqual(actualState.RuleSubgroups[i]) {
			return false
		}
	}

	return true
}

func flattenRuleMatchers(matchers []*rulesgroups.RuleMatcher) (applications []string, subsystems []string, severities []RuleGroupSeverity) {
	applications = make([]string, 0)
	subsystems = make([]string, 0)
	severities = make([]RuleGroupSeverity, 0)

	for _, m := range matchers {
		switch m.Constraint.(type) {
		case *rulesgroups.RuleMatcher_ApplicationName:
			applications = append(applications, m.GetApplicationName().GetValue().GetValue())
		case *rulesgroups.RuleMatcher_SubsystemName:
			subsystems = append(subsystems, m.GetSubsystemName().GetValue().GetValue())
		case *rulesgroups.RuleMatcher_Severity:
			severities = append(severities, rulesProtoSeverityToSchemaSeverity[m.GetSeverity().GetValue().String()])
		}
	}

	return applications, subsystems, severities
}

// RuleGroupStatus defines the observed state of RuleGroup
type RuleGroupStatus struct {
	ID string `json:"id"`
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
