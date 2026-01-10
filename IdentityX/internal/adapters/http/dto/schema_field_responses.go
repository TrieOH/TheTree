package dto

import (
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
