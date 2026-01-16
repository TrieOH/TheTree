package outbounds

import (
	"GoAuth/internal/domain/version"
	"context"

	"github.com/google/uuid"
)

type SchemaVersionRepository interface {
	Draft(ctx context.Context, toDraft version.Version) (*version.Version, error)
	Publish(ctx context.Context, toPublish version.Version) error
	Archive(ctx context.Context, toArchive version.Version) error
	GetCurrent(ctx context.Context, schemaID uuid.UUID) (*version.Version, error)
	GetLatest(ctx context.Context, schemaID uuid.UUID) (*version.Version, error)
	GetLatestForUpdate(ctx context.Context, schemaID uuid.UUID) (*version.Version, error)
	List(ctx context.Context, schemaID uuid.UUID) ([]version.Version, error)
	CopyOnDraft(ctx context.Context, schemaVersionID uuid.UUID) (*version.Version, error)
}
