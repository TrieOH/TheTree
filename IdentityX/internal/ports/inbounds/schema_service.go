package inbounds

import (
	"GoAuth/internal/domain/schema"
	"context"
	"time"

	"github.com/google/uuid"
)

type SchemaService interface {
	Draft(ctx context.Context, in DraftSchemaInput) (*SchemaOutput, error)
	GetByID(ctx context.Context, in GetSchemaByIDInput) (*SchemaOutput, error)
}

type DraftSchemaInput struct {
	SchemaType string
	Title      string
	FlowID     string
	ProjectID  string
}

type GetSchemaByIDInput struct {
	ProjectID string
	SchemaID  string
}

type PublishSchemaInput struct {
	FlowID    string
	ProjectID string
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
