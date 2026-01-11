package dto

import (
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
)

type CreateFieldRequest struct {
	Fields []FieldParam `json:"fields" validate:"required"`
}

type FieldParam struct {
	Key          string           `json:"key"`
	Type         string           `json:"type"`
	Owner        string           `json:"owner"`
	Title        string           `json:"title"`
	Description  *string          `json:"description"`
	Placeholder  *string          `json:"placeholder"`
	Required     bool             `json:"required"`
	Mutable      bool             `json:"mutable"`
	DefaultValue *json.RawMessage `json:"default_value"`
	Position     int              `json:"position"`
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
	return &inbounds.InputField{
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
}
