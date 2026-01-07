package inbounds

import (
	"GoAuth/internal/domain/schema"
	"context"
	"time"

	"github.com/google/uuid"
)

type SchemaVersionService interface {
	Draft(ctx context.Context, in DraftSchemaVersionInput) (*SchemaVersionOutput, error)
	Publish(ctx context.Context, in PublishSchemaVersionInput) error
}

type DraftSchemaVersionInput struct {
	SchemaID  string
	ProjectID string
}

type PublishSchemaVersionInput struct {
	SchemaID  string
	ProjectID string
}

type SchemaVersionOutput struct {
	ID            uuid.UUID
	SchemaID      uuid.UUID
	VersionNumber int
	Status        schema.VersionStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func SchemaVersionToOutput(out *schema.Version) *SchemaVersionOutput {
	if out == nil {
		return nil
	}
	return &SchemaVersionOutput{
		ID:            out.ID,
		SchemaID:      out.SchemaID,
		VersionNumber: out.VersionNumber,
		Status:        out.Status,
		CreatedAt:     out.CreatedAt,
		UpdatedAt:     out.UpdatedAt,
	}
}

type VersionVerboseOutput struct {
	SchemaVersionOutput
	Fields []OutputField
}
