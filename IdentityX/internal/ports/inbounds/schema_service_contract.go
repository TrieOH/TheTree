package inbounds

import (
	"GoAuth/internal/domain/schema"
	"time"

	"github.com/google/uuid"
)

type SchemaServiceInput struct {
	SchemaType string
	Title      string
	FlowID     string
	ProjectID  uuid.UUID
	SchemaID   uuid.UUID
}

type SchemaVerboseOutput struct {
	SchemaOutput
	Versions []VersionVerboseOutput
}

type SchemaOutput struct {
	ID               uuid.UUID
	ProjectID        uuid.UUID
	Title            string
	FlowID           string
	Type             string
	CurrentVersionID *uuid.UUID
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func SchemaToSchemaOutput(out *schema.Schema) *SchemaOutput {
	if out == nil {
		return nil
	}
	return &SchemaOutput{
		ID:               out.ID,
		ProjectID:        out.ProjectID,
		Title:            out.Title,
		FlowID:           out.FlowID,
		Type:             string(out.Type),
		CurrentVersionID: out.CurrentVersionID,
		Status:           string(out.Status),
		CreatedAt:        out.CreatedAt,
		UpdatedAt:        out.UpdatedAt,
	}
}
