package dto

import (
	"GoAuth/internal/ports/inbounds"
	"time"

	"github.com/google/uuid"
)

type SchemaVersionResponse struct {
	ID            uuid.UUID `json:"id"`
	SchemaID      uuid.UUID `json:"schema_id"`
	VersionNumber int       `json:"version_number"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type VersionVerboseResponse struct {
	SchemaVersionResponse
	Fields []FieldResponse `json:"fields"`
}

func SchemaVersionOutputToResponse(out *inbounds.SchemaVersionOutput) *SchemaVersionResponse {
	if out == nil {
		return nil
	}
	return &SchemaVersionResponse{
		ID:            out.ID,
		SchemaID:      out.SchemaID,
		VersionNumber: out.VersionNumber,
		Status:        string(out.Status),
		CreatedAt:     out.CreatedAt,
		UpdatedAt:     out.UpdatedAt,
	}
}
