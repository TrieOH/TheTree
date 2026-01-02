package outbound

import (
	"GoAuth/internal/domain/schema"
	"context"

	"github.com/google/uuid"
)

type SchemaRepository interface {
	Draft(ctx context.Context, schema schema.Schema) (*schema.Schema, error)
	Publish(ctx context.Context, schemaID uuid.UUID) error
	Archive(ctx context.Context, schemaID uuid.UUID) error
	Delete(ctx context.Context, schemaID uuid.UUID) error
}
