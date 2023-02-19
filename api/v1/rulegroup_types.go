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
	"encoding/json"
	"fmt"

	utils "coralogix-operator-poc/api"
	rulesgroups "coralogix-operator-poc/controllers/clientset/grpc/rules-groups/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	//+kubebuilder:validation:MinLength=0
	Name string `json:"name,omitempty"`

	// +optional
	Description string `json:"description,omitempty"`

	//+kubebuilder:default=true
	Active bool `json:"active,omitempty"`

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

func (in *Rule) DeepEqual(rule *rulesgroups.Rule) (bool, utils.Diff) {
	if actualActive := rule.Enabled.GetValue(); in.Active != actualActive {
		return false, utils.Diff{
			Name:    "Active",
			Desired: in.Active,
			Actual:  actualActive,
		}
	}

	if actualDescription := rule.Description.GetValue(); in.Description != actualDescription {
		return false, utils.Diff{
			Name:    "Description",
			Desired: in.Description,
			Actual:  actualDescription,
		}
	}

	if actualName := rule.Name.GetValue(); in.Name != actualName {
		return false, utils.Diff{
			Name:    "Name",
			Desired: in.Name,
			Actual:  actualName,
		}
	}

	return in.DeepEqualRuleType(rule)
}

func (in *Rule) DeepEqualRuleType(parameters *rulesgroups.Rule) (bool, utils.Diff) {
	switch typeParameters := parameters.Parameters.RuleParameters.(type) {
	case *rulesgroups.RuleParameters_ExtractParameters:
		if extractParameters := in.Extract; extractParameters == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Extract",
			}
		} else {
			actualExtractParameters := typeParameters.ExtractParameters
			if actualRegex := actualExtractParameters.Rule.GetValue(); extractParameters.Regex != actualRegex {
				return false, utils.Diff{
					Name:    "Extract.Regex",
					Desired: extractParameters.Regex,
					Actual:  actualRegex,
				}
			}

			if actualSourceField := parameters.SourceField.GetValue(); extractParameters.SourceField != actualSourceField {
				return false, utils.Diff{
					Name:    "Extract.SourceField",
					Desired: extractParameters.SourceField,
					Actual:  actualSourceField,
				}
			}
		}
	case *rulesgroups.RuleParameters_JsonExtractParameters:
		if jsonExtractParameters := in.JsonExtract; jsonExtractParameters == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "JsonExtract",
			}
		} else {
			actualJsonExtractParameters := typeParameters.JsonExtractParameters
			if actualJsonKey := actualJsonExtractParameters.Rule.GetValue(); jsonExtractParameters.JsonKey != actualJsonKey {
				return false, utils.Diff{
					Name:    "JsonExtract.JsonKey",
					Desired: jsonExtractParameters.JsonKey,
					Actual:  actualJsonKey,
				}
			}

			if desiredDestinationField, actualDestinationField := rulesSchemaDestinationFieldToProtoSeverityDestinationField[jsonExtractParameters.DestinationField], actualJsonExtractParameters.DestinationField.String(); desiredDestinationField != actualDestinationField {
				return false, utils.Diff{
					Name:    "JsonExtract.DestinationField",
					Desired: desiredDestinationField,
					Actual:  actualDestinationField,
				}
			}
		}
	case *rulesgroups.RuleParameters_ReplaceParameters:
		if replaceParameters := in.Replace; replaceParameters == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Replace",
			}
		} else {
			actualReplaceParameters := typeParameters.ReplaceParameters
			if actualReplacementString := actualReplaceParameters.ReplaceNewVal.GetValue(); replaceParameters.ReplacementString != actualReplacementString {
				return false, utils.Diff{
					Name:    "Replace.ReplacementString",
					Desired: replaceParameters.ReplacementString,
					Actual:  actualReplacementString,
				}
			}

			if actualDestinationField := actualReplaceParameters.DestinationField.GetValue(); replaceParameters.DestinationField != actualDestinationField {
				return false, utils.Diff{
					Name:    "Replace.DestinationField",
					Desired: replaceParameters.DestinationField,
					Actual:  actualDestinationField,
				}
			}

			if actualRegex := actualReplaceParameters.Rule.GetValue(); replaceParameters.Regex != actualRegex {
				return false, utils.Diff{
					Name:    "Replace.Regex",
					Desired: replaceParameters.Regex,
					Actual:  actualRegex,
				}
			}

			if actualSourceField := parameters.SourceField.GetValue(); replaceParameters.SourceField != actualSourceField {
				return false, utils.Diff{
					Name:    "Replace.SourceField",
					Desired: replaceParameters.SourceField,
					Actual:  actualSourceField,
				}
			}
		}
	case *rulesgroups.RuleParameters_ParseParameters:
		if parseParameters := in.Parse; parseParameters == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Parse",
			}
		} else {
			actualParseParameters := typeParameters.ParseParameters

			if actualRegex := actualParseParameters.Rule.GetValue(); parseParameters.Regex != actualRegex {
				return false, utils.Diff{
					Name:    "Parse.Regex",
					Desired: parseParameters.Regex,
					Actual:  actualRegex,
				}
			}

			if actualDestinationField := actualParseParameters.DestinationField.GetValue(); parseParameters.DestinationField != actualDestinationField {
				return false, utils.Diff{
					Name:    "Parse.DestinationField",
					Desired: parseParameters.DestinationField,
					Actual:  actualDestinationField,
				}
			}

			if actualSourceField := parameters.SourceField.GetValue(); parseParameters.SourceField != actualSourceField {
				return false, utils.Diff{
					Name:    "Parse.SourceField",
					Desired: parseParameters.SourceField,
					Actual:  actualSourceField,
				}
			}
		}
	case *rulesgroups.RuleParameters_AllowParameters:
		if allowParameters := in.Block; allowParameters == nil || allowParameters.BlockingAllMatchingBlocks {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Allow",
			}
		} else {
			actualAllowParameters := typeParameters.AllowParameters
			if actualRegex := actualAllowParameters.Rule.GetValue(); allowParameters.Regex != actualRegex {
				return false, utils.Diff{
					Name:    "Allow.Regex",
					Desired: allowParameters.Regex,
					Actual:  actualRegex,
				}
			}

			if actualKeepBlockedLogs := actualAllowParameters.KeepBlockedLogs.GetValue(); allowParameters.KeepBlockedLogs != actualKeepBlockedLogs {
				return false, utils.Diff{
					Name:    "Allow.KeepBlockedLogs",
					Desired: allowParameters.KeepBlockedLogs,
					Actual:  actualKeepBlockedLogs,
				}
			}

			if actualSourceField := parameters.SourceField.GetValue(); allowParameters.SourceField != actualSourceField {
				return false, utils.Diff{
					Name:    "Allow.SourceField",
					Desired: allowParameters.SourceField,
					Actual:  actualSourceField,
				}
			}
		}
	case *rulesgroups.RuleParameters_BlockParameters:
		if blockParameters := in.Block; blockParameters == nil || !blockParameters.BlockingAllMatchingBlocks {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "Block",
			}
		} else {
			actualBlockParameters := typeParameters.BlockParameters

			if actualRegex := actualBlockParameters.Rule.GetValue(); blockParameters.Regex != actualRegex {
				return false, utils.Diff{
					Name:    "Block.Regex",
					Desired: blockParameters.Regex,
					Actual:  actualRegex,
				}
			}

			if actualKeepBlockedLogs := actualBlockParameters.KeepBlockedLogs.GetValue(); blockParameters.KeepBlockedLogs != actualKeepBlockedLogs {
				return false, utils.Diff{
					Name:    "Block.KeepBlockedLogs",
					Desired: blockParameters.KeepBlockedLogs,
					Actual:  actualKeepBlockedLogs,
				}
			}
			if actualSourceField := parameters.SourceField.GetValue(); blockParameters.SourceField != actualSourceField {
				return false, utils.Diff{
					Name:    "Block.SourceField",
					Desired: blockParameters.SourceField,
					Actual:  actualSourceField,
				}
			}
		}
	case *rulesgroups.RuleParameters_ExtractTimestampParameters:
		if extractTimestampParameters := in.ExtractTimestamp; extractTimestampParameters == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "ExtractTimestamp",
			}
		} else {
			actualExtractTimestampParameters := typeParameters.ExtractTimestampParameters

			if desiredFieldFormatStandard, actualFieldFormatStandard := rulesSchemaFormatStandardToProtoFormatStandard[extractTimestampParameters.FieldFormatStandard], actualExtractTimestampParameters.Standard.String(); desiredFieldFormatStandard != actualFieldFormatStandard {
				return false, utils.Diff{
					Name:    "ExtractTimestamp.FieldFormatStandard",
					Desired: desiredFieldFormatStandard,
					Actual:  actualFieldFormatStandard,
				}
			}

			if actualTimeFormat := actualExtractTimestampParameters.Format.GetValue(); extractTimestampParameters.TimeFormat != actualTimeFormat {
				return false, utils.Diff{
					Name:    "ExtractTimestamp.TimeFormat",
					Desired: extractTimestampParameters.TimeFormat,
					Actual:  actualTimeFormat,
				}
			}

			if actualSourceField := parameters.SourceField.GetValue(); extractTimestampParameters.SourceField != actualSourceField {
				return false, utils.Diff{
					Name:    "ExtractTimestamp.SourceField",
					Desired: extractTimestampParameters.SourceField,
					Actual:  actualSourceField,
				}
			}
		}
	case *rulesgroups.RuleParameters_RemoveFieldsParameters:
		if removeFieldsParameters := in.RemoveFields; removeFieldsParameters == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "RemoveFields",
			}
		} else {
			actualRemoveFieldsParameters := typeParameters.RemoveFieldsParameters
			if len(removeFieldsParameters.ExcludedFields) != len(actualRemoveFieldsParameters.Fields) ||
				!utils.SlicesWithUniqueValuesEqual(removeFieldsParameters.ExcludedFields, actualRemoveFieldsParameters.Fields) {
				return false, utils.Diff{
					Name:    "RemoveFields.ExcludedFields",
					Desired: removeFieldsParameters.ExcludedFields,
					Actual:  actualRemoveFieldsParameters.Fields,
				}
			}
		}

	case *rulesgroups.RuleParameters_JsonStringifyParameters:
		if jsonStringifyParameters := in.JsonStringify; jsonStringifyParameters == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "JsonStringify",
			}
		} else {
			actualJsonStringifyParameters := typeParameters.JsonStringifyParameters

			if actualKeepSourceField := !actualJsonStringifyParameters.DeleteSource.GetValue(); jsonStringifyParameters.KeepSourceField != actualKeepSourceField {
				return false, utils.Diff{
					Name:    "JsonStringify.KeepSourceField",
					Desired: jsonStringifyParameters.KeepSourceField,
					Actual:  actualKeepSourceField,
				}
			}

			if actualDestinationField := actualJsonStringifyParameters.DestinationField.GetValue(); jsonStringifyParameters.DestinationField != actualDestinationField {
				return false, utils.Diff{
					Name:    "JsonStringify.DestinationField",
					Desired: jsonStringifyParameters.DestinationField,
					Actual:  actualDestinationField,
				}
			}

			if actualSourceField := parameters.SourceField.GetValue(); jsonStringifyParameters.SourceField != actualSourceField {
				return false, utils.Diff{
					Name:    "JsonStringify.SourceField",
					Desired: jsonStringifyParameters.SourceField,
					Actual:  actualSourceField,
				}
			}
		}
	case *rulesgroups.RuleParameters_JsonParseParameters:
		if jsonParseParameters := in.ParseJsonField; jsonParseParameters == nil {
			return false, utils.Diff{
				Name:   "Type",
				Actual: "JsonParse",
			}
		} else {
			actualJsonParseParameters := typeParameters.JsonParseParameters

			if actualDestinationField := actualJsonParseParameters.DestinationField.GetValue(); jsonParseParameters.DestinationField != actualDestinationField {
				return false, utils.Diff{
					Name:    "JsonParse.DestinationField",
					Desired: jsonParseParameters.DestinationField,
					Actual:  actualDestinationField,
				}
			}

			if actualSourceField := parameters.SourceField.GetValue(); jsonParseParameters.SourceField != actualSourceField {
				return false, utils.Diff{
					Name:    "JsonParse.SourceField",
					Desired: jsonParseParameters.SourceField,
					Actual:  actualSourceField,
				}
			}

			if actualKeepSourceField := !actualJsonParseParameters.DeleteSource.GetValue(); jsonParseParameters.KeepSourceField == actualKeepSourceField {
				return false, utils.Diff{
					Name:    "JsonParse.KeepSourceField",
					Desired: jsonParseParameters.KeepSourceField,
					Actual:  actualKeepSourceField,
				}
			}

			if actualKeepDestinationField := actualJsonParseParameters.OverrideDest.GetValue(); jsonParseParameters.KeepDestinationField == actualKeepDestinationField {
				return false, utils.Diff{
					Name:    "JsonParse.KeepDestinationField",
					Desired: jsonParseParameters.KeepDestinationField,
					Actual:  actualKeepDestinationField,
				}
			}
		}
	}

	return true, utils.Diff{}
}

type Parse struct {
	SourceField string `json:"sourceField,omitempty"`

	DestinationField string `json:"destinationField,omitempty"`

	Regex string `json:"regex,omitempty"`
}

type Block struct {
	SourceField string `json:"sourceField,omitempty"`

	Regex string `json:"regex,omitempty"`

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

	DestinationField string `json:"destinationField,omitempty"`

	Regex string `json:"regex,omitempty"`

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

	DestinationField string `json:"destinationField,omitempty"`

	//+kubebuilder:default=false
	KeepSourceField bool `json:"keepSourceField,omitempty"`
}

type Extract struct {
	SourceField string `json:"sourceField,omitempty"`

	Regex string `json:"regex,omitempty"`
}

type ParseJsonField struct {
	SourceField string `json:"sourceField,omitempty"`

	DestinationField string `json:"destinationField,omitempty"`

	KeepSourceField bool `json:"keepSourceField,omitempty"`

	KeepDestinationField bool `json:"keepDestinationField,omitempty"`
}

type RuleSubGroup struct {
	//+kubebuilder:default=true
	Active bool `json:"active,omitempty"`

	// +optional
	Rules []Rule `json:"rules,omitempty"`
}

func (in *RuleSubGroup) DeepEqual(subgroup *rulesgroups.RuleSubgroup) (bool, utils.Diff) {
	if actualActive := subgroup.Enabled.GetValue(); in.Active != actualActive {
		return false, utils.Diff{
			Name:    "Active",
			Desired: in.Active,
			Actual:  actualActive,
		}
	}

	if len(in.Rules) != len(subgroup.Rules) {
		return false, utils.Diff{
			Name:    "Rules.length",
			Desired: len(in.Rules),
			Actual:  len(subgroup.Rules),
		}
	}

	for i := range in.Rules {
		if equal, diff := in.Rules[i].DeepEqual(subgroup.Rules[i]); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("Rules[%d].%s", i, diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	return true, utils.Diff{}
}

// +kubebuilder:validation:Enum=Debug;Verbose;Info;Warning;Error;Critical
type RuleGroupSeverity string

// RuleGroupSpec defines the Desired state of RuleGroup
type RuleGroupSpec struct {
	//+kubebuilder:validation:MinLength=0
	Name string `json:"name,omitempty"`

	// +optional
	Description string `json:"description,omitempty"`

	//+kubebuilder:default=true
	Active bool `json:"active,omitempty"`

	// +optional
	Applications []string `json:"applications,omitempty"`

	// +optional
	Subsystems []string `json:"subsystems,omitempty"`

	// +optional
	Severities []RuleGroupSeverity `json:"severities,omitempty"`

	//+kubebuilder:default=false
	Hidden bool `json:"hidden,omitempty"`

	// +optional
	Creator string `json:"creator,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum:=1
	Order *int32 `json:"order,omitempty"`

	// +optional
	RuleSubgroups []RuleSubGroup `json:"subgroups,omitempty"`
}

func (in *RuleGroupSpec) ToString() string {
	str, _ := json.Marshal(*in)
	return string(str)
}

func (in *RuleGroupSpec) DeepEqual(actualState *rulesgroups.RuleGroup) (bool, utils.Diff) {
	if actualName := actualState.GetName().GetValue(); in.Name != actualName {
		return false, utils.Diff{
			Name:    "Name",
			Desired: in.Name,
			Actual:  actualName,
		}
	}

	if actualDescription := actualState.GetDescription().GetValue(); in.Description != actualDescription {
		return false, utils.Diff{
			Name:    "Description",
			Desired: in.Description,
			Actual:  actualDescription,
		}
	}

	if actualActive := actualState.GetEnabled().GetValue(); in.Active != actualActive {
		return false, utils.Diff{
			Name:    "Active",
			Desired: in.Active,
			Actual:  actualActive,
		}
	}

	if actualHidden := actualState.GetHidden().GetValue(); in.Hidden != actualHidden {
		return false, utils.Diff{
			Name:    "Hidden",
			Desired: in.Hidden,
			Actual:  actualHidden,
		}
	}

	if actualCreator := actualState.GetCreator().GetValue(); in.Creator != actualCreator {
		return false, utils.Diff{
			Name:    "Creator",
			Desired: in.Creator,
			Actual:  actualCreator,
		}
	}

	if in.Order == nil {
		in.Order = new(int32)
		*in.Order = int32(actualState.GetOrder().GetValue())
	} else if actualOrder := actualState.GetOrder().GetValue(); uint32(*in.Order) != actualOrder {
		return false, utils.Diff{
			Name:    "Order",
			Desired: *in.Order,
			Actual:  actualOrder,
		}
	}

	applications, subsystems, severities := flattenRuleMatchers(actualState.RuleMatchers)

	if !utils.SlicesWithUniqueValuesEqual(in.Applications, applications) {
		return false, utils.Diff{
			Name:    "Applications",
			Desired: in.Applications,
			Actual:  applications,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Subsystems, subsystems) {
		return false, utils.Diff{
			Name:    "Subsystems",
			Desired: in.Subsystems,
			Actual:  subsystems,
		}
	}

	if !utils.SlicesWithUniqueValuesEqual(in.Severities, severities) {
		return false, utils.Diff{
			Name:    "Severities",
			Desired: in.Severities,
			Actual:  severities,
		}
	}

	if len(in.RuleSubgroups) != len(actualState.RuleSubgroups) {
		return false, utils.Diff{
			Name:    "RuleSubgroups length",
			Desired: len(in.RuleSubgroups),
			Actual:  len(actualState.RuleSubgroups),
		}
	}

	for i := range in.RuleSubgroups {
		if equal, diff := in.RuleSubgroups[i].DeepEqual(actualState.RuleSubgroups[i]); !equal {
			return false, utils.Diff{
				Name:    fmt.Sprintf("RuleSubgroups[%d].%s", i, diff.Name),
				Desired: diff.Desired,
				Actual:  diff.Actual,
			}
		}
	}

	return true, utils.Diff{}
}

func (in *RuleGroupSpec) ExtractUpdateRuleGroupRequest(id string) *rulesgroups.UpdateRuleGroupRequest {
	ruleGroup := in.ExtractCreateRuleGroupRequest()
	return &rulesgroups.UpdateRuleGroupRequest{
		GroupId:   wrapperspb.String(id),
		RuleGroup: ruleGroup,
	}
}

func (in *RuleGroupSpec) ExtractCreateRuleGroupRequest() *rulesgroups.CreateRuleGroupRequest {
	name := wrapperspb.String(in.Name)
	description := wrapperspb.String(in.Description)
	enabled := wrapperspb.Bool(in.Active)
	hidden := wrapperspb.Bool(in.Hidden)
	creator := wrapperspb.String(in.Creator)
	ruleMatchers := expandRuleMatchers(in.Applications, in.Subsystems, in.Severities)
	ruleSubGroups := expandRuleSubGroups(in.RuleSubgroups)
	order := expandOrder(in.Order)

	return &rulesgroups.CreateRuleGroupRequest{
		Name:          name,
		Description:   description,
		Enabled:       enabled,
		Hidden:        hidden,
		Creator:       creator,
		RuleMatchers:  ruleMatchers,
		RuleSubgroups: ruleSubGroups,
		Order:         order,
	}
}

func expandOrder(order *int32) *wrapperspb.UInt32Value {
	if order != nil {
		return wrapperspb.UInt32(uint32(*order))
	}
	return nil
}

func expandRuleSubGroups(subGroups []RuleSubGroup) []*rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup {
	ruleSubGroups := make([]*rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup, 0, len(subGroups))
	for i, subGroup := range subGroups {
		rsg := expandRuleSubGroup(subGroup)
		rsg.Order = wrapperspb.UInt32(uint32(i + 1))
		ruleSubGroups = append(ruleSubGroups, rsg)
	}
	return ruleSubGroups
}

func expandRuleSubGroup(subGroup RuleSubGroup) *rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup {
	enabled := wrapperspb.Bool(subGroup.Active)
	rules := expandRules(subGroup.Rules)
	return &rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup{
		Enabled: enabled,
		Rules:   rules,
	}
}

func expandRules(rules []Rule) []*rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup_CreateRule {
	expandedRules := make([]*rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup_CreateRule, 0, len(rules))
	for i, rule := range rules {
		r := expandRule(rule)
		r.Order = wrapperspb.UInt32(uint32(i + 1))
		expandedRules = append(expandedRules, r)
	}
	return expandedRules
}

func expandRule(rule Rule) *rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup_CreateRule {
	name := wrapperspb.String(rule.Name)
	description := wrapperspb.String(rule.Description)
	enabled := wrapperspb.Bool(rule.Active)
	sourceFiled, parameters := expandSourceFiledAndParameters(rule)

	return &rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup_CreateRule{
		Name:        name,
		Description: description,
		SourceField: sourceFiled,
		Parameters:  parameters,
		Enabled:     enabled,
	}
}

func expandSourceFiledAndParameters(rule Rule) (sourceField *wrapperspb.StringValue, parameters *rulesgroups.RuleParameters) {
	if parse := rule.Parse; parse != nil {
		sourceField = wrapperspb.String(parse.SourceField)
		parameters = &rulesgroups.RuleParameters{
			RuleParameters: &rulesgroups.RuleParameters_ParseParameters{
				ParseParameters: &rulesgroups.ParseParameters{
					DestinationField: wrapperspb.String(parse.DestinationField),
					Rule:             wrapperspb.String(parse.Regex),
				},
			},
		}
	} else if parseJsonField := rule.ParseJsonField; parseJsonField != nil {
		sourceField = wrapperspb.String(parseJsonField.SourceField)
		parameters = &rulesgroups.RuleParameters{
			RuleParameters: &rulesgroups.RuleParameters_JsonParseParameters{
				JsonParseParameters: &rulesgroups.JsonParseParameters{
					DestinationField: wrapperspb.String(parseJsonField.DestinationField),
					DeleteSource:     wrapperspb.Bool(!parseJsonField.KeepSourceField),
					OverrideDest:     wrapperspb.Bool(!parseJsonField.KeepDestinationField),
					EscapedValue:     wrapperspb.Bool(true),
				},
			},
		}
	} else if jsonStringify := rule.JsonStringify; jsonStringify != nil {
		sourceField = wrapperspb.String(jsonStringify.SourceField)
		parameters = &rulesgroups.RuleParameters{
			RuleParameters: &rulesgroups.RuleParameters_JsonStringifyParameters{
				JsonStringifyParameters: &rulesgroups.JsonStringifyParameters{
					DestinationField: wrapperspb.String(jsonStringify.DestinationField),
					DeleteSource:     wrapperspb.Bool(!jsonStringify.KeepSourceField),
				},
			},
		}
	} else if jsonExtract := rule.JsonExtract; jsonExtract != nil {
		sourceField = wrapperspb.String("text")
		destinationField := expandJsonExtractDestinationField(jsonExtract.DestinationField)
		jsonKey := wrapperspb.String(jsonExtract.JsonKey)
		parameters = &rulesgroups.RuleParameters{
			RuleParameters: &rulesgroups.RuleParameters_JsonExtractParameters{
				JsonExtractParameters: &rulesgroups.JsonExtractParameters{
					DestinationField: destinationField,
					Rule:             jsonKey,
				},
			},
		}
	} else if removeFields := rule.RemoveFields; removeFields != nil {
		sourceField = wrapperspb.String("text")
		parameters = &rulesgroups.RuleParameters{
			RuleParameters: &rulesgroups.RuleParameters_RemoveFieldsParameters{
				RemoveFieldsParameters: &rulesgroups.RemoveFieldsParameters{
					Fields: removeFields.ExcludedFields,
				},
			},
		}
	} else if extractTimestamp := rule.ExtractTimestamp; extractTimestamp != nil {
		sourceField = wrapperspb.String(extractTimestamp.SourceField)
		standard := expandTimeStampStandard(extractTimestamp)
		format := wrapperspb.String(extractTimestamp.TimeFormat)
		parameters = &rulesgroups.RuleParameters{
			RuleParameters: &rulesgroups.RuleParameters_ExtractTimestampParameters{
				ExtractTimestampParameters: &rulesgroups.ExtractTimestampParameters{
					Standard: standard,
					Format:   format,
				},
			},
		}
	} else if block := rule.Block; block != nil {
		sourceField = wrapperspb.String(block.SourceField)
		if block.BlockingAllMatchingBlocks {
			parameters = &rulesgroups.RuleParameters{
				RuleParameters: &rulesgroups.RuleParameters_BlockParameters{
					BlockParameters: &rulesgroups.BlockParameters{
						KeepBlockedLogs: wrapperspb.Bool(block.KeepBlockedLogs),
						Rule:            wrapperspb.String(block.Regex),
					},
				},
			}
		} else {
			parameters = &rulesgroups.RuleParameters{
				RuleParameters: &rulesgroups.RuleParameters_AllowParameters{
					AllowParameters: &rulesgroups.AllowParameters{
						KeepBlockedLogs: wrapperspb.Bool(block.KeepBlockedLogs),
						Rule:            wrapperspb.String(block.Regex),
					},
				},
			}
		}
	} else if replace := rule.Replace; replace != nil {
		sourceField = wrapperspb.String(replace.SourceField)
		parameters = &rulesgroups.RuleParameters{
			RuleParameters: &rulesgroups.RuleParameters_ReplaceParameters{
				ReplaceParameters: &rulesgroups.ReplaceParameters{
					DestinationField: wrapperspb.String(replace.DestinationField),
					ReplaceNewVal:    wrapperspb.String(replace.ReplacementString),
					Rule:             wrapperspb.String(replace.Regex),
				},
			},
		}
	}

	return
}

func expandTimeStampStandard(extractTimestamp *ExtractTimestamp) rulesgroups.ExtractTimestampParameters_FormatStandard {
	return rulesgroups.ExtractTimestampParameters_FormatStandard(
		rulesgroups.ExtractTimestampParameters_FormatStandard_value[rulesSchemaFormatStandardToProtoFormatStandard[extractTimestamp.FieldFormatStandard]],
	)
}

func expandJsonExtractDestinationField(destinationField DestinationField) rulesgroups.JsonExtractParameters_DestinationField {
	return rulesgroups.JsonExtractParameters_DestinationField(
		rulesgroups.JsonExtractParameters_DestinationField_value[rulesSchemaDestinationFieldToProtoSeverityDestinationField[destinationField]],
	)
}

func expandRuleMatchers(applications, subsystems []string, severities []RuleGroupSeverity) []*rulesgroups.RuleMatcher {
	ruleMatchers := make([]*rulesgroups.RuleMatcher, 0, len(applications)+len(subsystems)+len(severities))

	for _, app := range applications {
		constraintStr := wrapperspb.String(app)
		applicationNameConstraint := rulesgroups.ApplicationNameConstraint{Value: constraintStr}
		ruleMatcherApplicationName := rulesgroups.RuleMatcher_ApplicationName{ApplicationName: &applicationNameConstraint}
		ruleMatchers = append(ruleMatchers, &rulesgroups.RuleMatcher{Constraint: &ruleMatcherApplicationName})
	}

	for _, subSys := range subsystems {
		constraintStr := wrapperspb.String(subSys)
		subsystemNameConstraint := rulesgroups.SubsystemNameConstraint{Value: constraintStr}
		ruleMatcherApplicationName := rulesgroups.RuleMatcher_SubsystemName{SubsystemName: &subsystemNameConstraint}
		ruleMatchers = append(ruleMatchers, &rulesgroups.RuleMatcher{Constraint: &ruleMatcherApplicationName})
	}

	for _, sev := range severities {
		constraintEnum := expandRuledSeverity(sev)
		severityConstraint := rulesgroups.SeverityConstraint{Value: constraintEnum}
		ruleMatcherSeverity := rulesgroups.RuleMatcher_Severity{Severity: &severityConstraint}
		ruleMatchers = append(ruleMatchers, &rulesgroups.RuleMatcher{Constraint: &ruleMatcherSeverity})
	}

	return ruleMatchers
}

func expandRuledSeverity(severity RuleGroupSeverity) rulesgroups.SeverityConstraint_Value {
	sevStr := rulesSchemaSeverityToProtoSeverity[severity]
	return rulesgroups.SeverityConstraint_Value(rulesgroups.SeverityConstraint_Value_value[sevStr])
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
	ID *string `json:"id"`
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
