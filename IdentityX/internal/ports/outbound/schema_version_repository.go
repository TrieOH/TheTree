package outbound

import (
	"GoAuth/internal/domain/schema"
	"context"

	"github.com/google/uuid"
)

type SchemaVersionRepository interface {
	Draft(ctx context.Context, version schema.Version) (*schema.Version, error)
	Publish(ctx context.Context, version schema.Version) error
	Archive(ctx context.Context, version schema.Version) error
	GetLatest(ctx context.Context, schemaID uuid.UUID) (*schema.Version, error)
}
