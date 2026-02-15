package dto

import (
	"GoAuth/internal/domain/field"
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type FieldResponse struct {
	ObjectID        uuid.UUID        `json:"object_id"`
	ID              uuid.UUID        `json:"id"`
	Key             string           `json:"key"`
	SchemaID        uuid.UUID        `json:"schema_id"`
	SchemaVersionID uuid.UUID        `json:"schema_version_id"`
	Type            string           `json:"type"`
	Owner           string           `json:"owner"`
	Title           string           `json:"title"`
	Description     *string          `json:"description"`
	Placeholder     *string          `json:"placeholder"`
	Required        bool             `json:"required"`
	Mutable         bool             `json:"mutable"`
	DefaultValue    *json.RawMessage `json:"default_value"`
	Position        int              `json:"position"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

func OutputFieldSliceToFieldResponseSlice(fps []inbounds.OutputField) []FieldResponse {
	out := make([]FieldResponse, 0, len(fps))
	for _, f := range fps {
		out = append(out, *OutputFieldToFieldResponse(&f))
	}
	return out
}

func OutputFieldToFieldResponse(field *inbounds.OutputField) *FieldResponse {
	if field == nil {
		return nil
	}
	return &FieldResponse{
		ObjectID:        field.ObjectID,
		ID:              field.ID,
		Key:             field.Key,
		SchemaID:        field.SchemaID,
		SchemaVersionID: field.SchemaVersionID,
		Type:            field.Type,
		Owner:           field.Owner,
		Title:           field.Title,
		Description:     field.Description,
		Placeholder:     field.Placeholder,
		Required:        field.Required,
		Mutable:         field.Mutable,
		DefaultValue:    field.DefaultValue,
		Position:        field.Position,
		CreatedAt:       field.CreatedAt,
		UpdatedAt:       field.UpdatedAt,
	}
}

type OptionResponse struct {
	ID       uuid.UUID `json:"id"`
	Value    string    `json:"value"`
	Label    string    `json:"label"`
	Position int       `json:"position"`
}

func OptionSliceToResponse(opts []field.Option) []OptionResponse {
	out := make([]OptionResponse, len(opts))
	for i, opt := range opts {
		out[i] = OptionResponse{
			ID:       opt.ID,
			Value:    opt.Value,
			Label:    opt.Label,
			Position: opt.Position,
		}
	}
	return out
}

type VisibilityRuleResponse struct {
	ID               uuid.UUID        `json:"id"`
	DependsOnFieldID uuid.UUID        `json:"depends_on_field_id"`
	Operator         string           `json:"operator"`
	Value            *json.RawMessage `json:"value,omitempty"`
}

func VisibilityRuleSliceToResponse(rules []field.VisibilityRule) []VisibilityRuleResponse {
	out := make([]VisibilityRuleResponse, len(rules))
	for i, rule := range rules {
		out[i] = VisibilityRuleResponse{
			ID:               rule.ID,
			DependsOnFieldID: rule.DependsOnFieldID,
			Operator:         string(rule.Operator),
			Value:            rule.Value,
		}
	}
	return out
}

type RequiredRuleResponse struct {
	ID               uuid.UUID        `json:"id"`
	DependsOnFieldID uuid.UUID        `json:"depends_on_field_id"`
	Operator         string           `json:"operator"`
	Value            *json.RawMessage `json:"value,omitempty"`
}

func RequiredRuleSliceToResponse(rules []field.RequiredRule) []RequiredRuleResponse {
	out := make([]RequiredRuleResponse, len(rules))
	for i, rule := range rules {
		out[i] = RequiredRuleResponse{
			ID:               rule.ID,
			DependsOnFieldID: rule.DependsOnFieldID,
			Operator:         string(rule.Operator),
			Value:            rule.Value,
		}
	}
	return out
}
