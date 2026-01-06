package dto

import (
	"GoAuth/internal/ports/inbounds"
	"time"

	"github.com/google/uuid"
)

type DraftSchemaRequest struct {
	SchemaType string `json:"schema_type" validate:"required,oneof=core context sub-context"`
	Title      string `json:"title" validate:"required,max=255"`
	FlowID     string `json:"flow_id" validate:"required,max=63"`
}

type DraftSchemaResponse struct {
	ID               uuid.UUID  `json:"id"`
	ProjectID        uuid.UUID  `json:"project_id"`
	Title            string     `json:"title"`
	FlowID           string     `json:"flow_id"`
	Type             string     `json:"type"`
	CurrentVersionID *uuid.UUID `json:"current_version_id"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func SchemaOutputToResponse(out *inbounds.SchemaOutput) *DraftSchemaResponse {
	if out == nil {
		return nil
	}
	return &DraftSchemaResponse{
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
