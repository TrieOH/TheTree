package inbounds

import (
	"GoAuth/internal/domain/schema"
	"context"
	"time"

	"github.com/google/uuid"
)

type SchemaVersionService interface {
	Draft(ctx context.Context, in DraftSchemaVersionInput) (*DraftSchemaVersionOutput, error)
	// Publish /projects/{id}/schemas/{schemaID}/{version:v[0-9]+}/publish
	Publish(ctx context.Context, in PublishSchemaVersionInput) error
}

type DraftSchemaVersionInput struct {
	SchemaID  string
	ProjectID string
}

type PublishSchemaVersionInput struct {
	VersionID string
	SchemaID  string
	ProjectID string
}

type DraftSchemaVersionOutput struct {
	ID            uuid.UUID
	SchemaID      uuid.UUID
	VersionNumber int
	Status        schema.VersionStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func SchemaVersionToOutput(out *schema.Version) *DraftSchemaVersionOutput {
	if out == nil {
		return nil
	}
	return &DraftSchemaVersionOutput{
		ID:            out.ID,
		SchemaID:      out.SchemaID,
		VersionNumber: out.VersionNumber,
		Status:        out.Status,
		CreatedAt:     out.CreatedAt,
		UpdatedAt:     out.UpdatedAt,
	}
}
