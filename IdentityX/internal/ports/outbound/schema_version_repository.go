package outbound

import (
	"GoAuth/internal/domain/schema"
	"context"

	"github.com/google/uuid"
)

type SchemaVersionRepository interface {
	Draft(ctx context.Context, version schema.Version) (*schema.Version, error)
	Publish(ctx context.Context, versionID uuid.UUID) error
	Archive(ctx context.Context, versionID uuid.UUID) error
	Delete(ctx context.Context, version uuid.UUID) error
}
