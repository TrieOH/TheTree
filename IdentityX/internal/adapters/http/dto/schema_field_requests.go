package dto

import (
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
)

type CreateFieldRequest struct {
	Fields []FieldParam `json:"fields" validate:"required"`
}

type FieldParam struct {
	Key             string                `json:"key"`
	Type            string                `json:"type"`
	Owner           string                `json:"owner"`
	Title           string                `json:"title"`
	Description     *string               `json:"description"`
	Placeholder     *string               `json:"placeholder"`
	Required        bool                  `json:"required"`
	Mutable         bool                  `json:"mutable"`
	DefaultValue    *json.RawMessage      `json:"default_value"`
	Position        int                   `json:"position"`
	Options         []OptionParam         `json:"options"`
	VisibilityRules []VisibilityRuleParam `json:"visibility_rules"`
	RequiredRules   []RequiredRuleParam   `json:"required_rules"`
}

type OptionParam struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Position int    `json:"position"`
}

type VisibilityRuleParam struct {
	DependsOnFieldKey string           `json:"depends_on_field_key"`
	Operator          string           `json:"operator"`
	Value             *json.RawMessage `json:"value"`
}

type RequiredRuleParam struct {
	DependsOnFieldKey string           `json:"depends_on_field_key"`
	Operator          string           `json:"operator"`
	Value             *json.RawMessage `json:"value"`
}

func FieldParamSliceToInputFieldSlice(fps []FieldParam) []inbounds.InputField {
	out := make([]inbounds.InputField, 0, len(fps))
	for _, f := range fps {
		out = append(out, *FieldParamToInputField(&f))
	}
	return out
}

func FieldParamToInputField(f *FieldParam) *inbounds.InputField {
	if f == nil {
		return nil
	}

	input := &inbounds.InputField{
		Key:          f.Key,
		Type:         f.Type,
		Owner:        f.Owner,
		Title:        f.Title,
		Description:  f.Description,
		Placeholder:  f.Placeholder,
		Required:     f.Required,
		Mutable:      f.Mutable,
		DefaultValue: f.DefaultValue,
		Position:     f.Position,
	}

	for _, opt := range f.Options {
		input.Options = append(input.Options, inbounds.InputOption{
			Value:    opt.Value,
			Label:    opt.Label,
			Position: opt.Position,
		})
	}

	for _, rule := range f.VisibilityRules {
		input.VisibilityRules = append(input.VisibilityRules, inbounds.InputVisibilityRule{
			DependsOnFieldKey: rule.DependsOnFieldKey,
			Operator:          rule.Operator,
			Value:             rule.Value,
		})
	}

	for _, rule := range f.RequiredRules {
		input.RequiredRules = append(input.RequiredRules, inbounds.InputRequiredRule{
			DependsOnFieldKey: rule.DependsOnFieldKey,
			Operator:          rule.Operator,
			Value:             rule.Value,
		})
	}

	return input
}

type EditFieldRequest struct {
	Key          *string          `json:"key,omitempty"`
	Type         *string          `json:"type,omitempty"`
	Title        *string          `json:"title,omitempty"`
	Description  *string          `json:"description,omitempty"`
	Placeholder  *string          `json:"placeholder,omitempty"`
	Required     *bool            `json:"required,omitempty"`
	Mutable      *bool            `json:"mutable,omitempty"`
	DefaultValue *json.RawMessage `json:"default_value,omitempty"`
	Position     *int             `json:"position,omitempty"`
}

type SetFieldOptionsRequest struct {
	Options []OptionParam `json:"options" validate:"required"`
}

type SetVisibilityRulesRequest struct {
	VisibilityRules []VisibilityRuleParam `json:"visibility_rules" validate:"required"`
}

type EditVisibilityRuleRequest struct {
	DependsOnFieldID *string          `json:"depends_on_field_id,omitempty"`
	Operator         *string          `json:"operator,omitempty"`
	Value            *json.RawMessage `json:"value,omitempty"`
}
