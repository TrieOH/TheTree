package inbounds

import (
	"GoAuth/internal/domain/schema"
	"time"

	"github.com/google/uuid"
)

type SchemaVersionServiceInput struct {
	SchemaID  string
	ProjectID string
}

type VersionVerboseOutput struct {
	SchemaVersionOutput
	Fields []OutputField
}

type SchemaVersionOutput struct {
	ID               uuid.UUID
	SchemaID         uuid.UUID
	BasedOnVersionID *uuid.UUID
	VersionNumber    int
	Status           schema.VersionStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func SchemaVersionToOutput(out *schema.Version) *SchemaVersionOutput {
	if out == nil {
		return nil
	}
	return &SchemaVersionOutput{
		ID:               out.ID,
		SchemaID:         out.SchemaID,
		BasedOnVersionID: out.BasedOnVersionID,
		VersionNumber:    out.VersionNumber,
		Status:           out.Status,
		CreatedAt:        out.CreatedAt,
		UpdatedAt:        out.UpdatedAt,
	}
}
