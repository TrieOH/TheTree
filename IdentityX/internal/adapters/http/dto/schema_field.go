package dto

import (
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateFieldRequest struct {
	Fields []FieldParam `json:"fields"`
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

func FieldParamSliceToInputFieldSlice(fps []FieldParam) []inbounds.InputField {
	out := make([]inbounds.InputField, 0, len(fps))
	for _, f := range fps {
		out = append(out, *FieldParamToInputField(&f))
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

func OutputFieldSliceToFieldResponseSlice(fps []inbounds.OutputField) []FieldResponse {
	out := make([]FieldResponse, 0, len(fps))
	for _, f := range fps {
		out = append(out, *OutputFieldToFieldResponse(&f))
	}
	return out
}
