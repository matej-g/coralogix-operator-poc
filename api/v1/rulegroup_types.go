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
	Name string `json:"name,omitempty"`

	// +optional
	Description string `json:"description,omitempty"`

	//+kubebuilder:default=true
	Active bool `json:"active,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum:=1
	Order *int32 `json:"order,omitempty"`

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
	if in.Order == nil {
		in.Order = new(int32)
		*in.Order = int32(rule.Order.GetValue())
	} else if uint32(*in.Order) != rule.Order.GetValue() {
		return false
	}
	if in.Description != rule.Description.GetValue() {
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
			if replaceParameters.DestinationField != actualReplaceParameters.DestinationField.String() {
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
			if jsonStringifyParameters.DestinationField != actualJsonStringifyParameters.DestinationField.String() {
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

	DestinationField string `json:"destinationField,omitempty"`

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

	DestinationField string `json:"destinationField,omitempty"`

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

type RuleSubGroup struct {
	//+kubebuilder:default=true
	Active bool `json:"active,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum:=1
	Order *int32 `json:"order,omitempty"`

	// +optional
	Rules []Rule `json:"rules,omitempty"`
}

func (in *RuleSubGroup) DeepEqual(subgroup *rulesgroups.RuleSubgroup) bool {
	if in.Active != subgroup.Enabled.GetValue() {
		return false
	}
	if in.Order == nil {
		in.Order = new(int32)
		*in.Order = int32(subgroup.Order.GetValue())
	} else if uint32(*in.Order) != subgroup.Order.GetValue() {
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

	if in.Order == nil {
		in.Order = new(int32)
		*in.Order = int32(actualState.GetOrder().GetValue())
	} else if uint32(*in.Order) != actualState.GetOrder().GetValue() {
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
		if rsg.Order == nil {
			rsg.Order = wrapperspb.UInt32(uint32(i))
		}
		ruleSubGroups = append(ruleSubGroups, rsg)
	}
	return ruleSubGroups
}

func expandRuleSubGroup(subGroup RuleSubGroup) *rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup {
	enabled := wrapperspb.Bool(subGroup.Active)
	order := expandOrder(subGroup.Order)
	rules := expandRules(subGroup.Rules)
	return &rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup{
		Enabled: enabled,
		Order:   order,
		Rules:   rules,
	}
}

func expandRules(rules []Rule) []*rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup_CreateRule {
	expandedRules := make([]*rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup_CreateRule, 0, len(rules))
	for _, rule := range rules {
		r := expandRule(rule)
		expandedRules = append(expandedRules, r)
	}
	return expandedRules
}

func expandRule(rule Rule) *rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup_CreateRule {
	name := wrapperspb.String(rule.Name)
	description := wrapperspb.String(rule.Description)
	enabled := wrapperspb.Bool(rule.Active)
	sourceFiled, paramerters := expandSourceFiledAndParameters(rule)
	order := expandOrder(rule.Order)

	return &rulesgroups.CreateRuleGroupRequest_CreateRuleSubgroup_CreateRule{
		Name:        name,
		Description: description,
		SourceField: sourceFiled,
		Parameters:  paramerters,
		Enabled:     enabled,
		Order:       order,
	}
}

func expandSourceFiledAndParameters(rule Rule) (sourceField *wrapperspb.StringValue, parameters *rulesgroups.RuleParameters) {
	if parse := rule.Parse; parse != nil {
		sourceField = wrapperspb.String(parse.SourceField)
		parameters = &rulesgroups.RuleParameters{
			RuleParameters: &rulesgroups.RuleParameters_ParseParameters{
				ParseParameters: &rulesgroups.ParseParameters{
					DestinationField: wrapperspb.String(parse.DestinationField),
					Rule:             wrapperspb.String(parse.RegularExpression),
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
		sourceField = wrapperspb.String(jsonExtract.JsonKey)
		destinationField := expandJsonExtractDestinationField(jsonExtract.DestinationField)
		parameters = &rulesgroups.RuleParameters{
			RuleParameters: &rulesgroups.RuleParameters_JsonExtractParameters{
				JsonExtractParameters: &rulesgroups.JsonExtractParameters{
					DestinationField: destinationField,
					Rule:             wrapperspb.String(parse.RegularExpression),
				},
			},
		}
	} else if removeFields := rule.RemoveFields; removeFields != nil {
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
						Rule:            wrapperspb.String(block.RegularExpression),
					},
				},
			}
		} else {
			parameters = &rulesgroups.RuleParameters{
				RuleParameters: &rulesgroups.RuleParameters_AllowParameters{
					AllowParameters: &rulesgroups.AllowParameters{
						KeepBlockedLogs: wrapperspb.Bool(block.KeepBlockedLogs),
						Rule:            wrapperspb.String(block.RegularExpression),
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
					Rule:             wrapperspb.String(replace.RegularExpression),
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
