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
	GetCurrent(ctx context.Context, schemaID uuid.UUID) (*schema.Version, error)
	GetLatest(ctx context.Context, schemaID uuid.UUID) (*schema.Version, error)
	List(ctx context.Context, schemaID uuid.UUID) ([]schema.Version, error)
	CopyOnDraft(ctx context.Context, schemaVersionID uuid.UUID) (*schema.Version, error)
}
