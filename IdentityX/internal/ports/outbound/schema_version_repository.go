package outbound

import (
	"GoAuth/internal/domain/schema"
	"context"

	"github.com/google/uuid"
)

type SchemaVersionRepository interface {
	Draft(ctx context.Context, toDraft schema.Version) (*schema.Version, error)
	Publish(ctx context.Context, toPublish schema.Version) error
	Archive(ctx context.Context, toArchive schema.Version) error
	GetLatest(ctx context.Context, schemaID uuid.UUID) (*schema.Version, error)
}
