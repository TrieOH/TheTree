package outbound

import (
	"GoAuth/internal/domain/schema"
	"context"

	"github.com/google/uuid"
)

type SchemaRepository interface {
	Draft(ctx context.Context, schema schema.Schema) (*schema.Schema, error)
	Publish(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) error
	Archive(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) error
	Delete(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) error
	FindByID(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) (*schema.Schema, error)
	List(ctx context.Context, projectID uuid.UUID) ([]schema.Schema, error)
}
