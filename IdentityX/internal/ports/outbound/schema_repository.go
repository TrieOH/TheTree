package outbound

import (
	"GoAuth/internal/domain/schema"
	"context"

	"github.com/google/uuid"
)

type SchemaRepository interface {
	Draft(ctx context.Context, newSchema schema.Schema) (*schema.Schema, error)
	Publish(ctx context.Context, publishedSchema schema.Schema) error
	Archive(ctx context.Context, archivedSchema schema.Schema) error
	Delete(ctx context.Context, deletedSchema schema.Schema) error
	Exists(ctx context.Context, existsSchema schema.Schema) (bool, error)
	BelongsToProject(ctx context.Context, belongsToProject schema.Schema) (bool, error)
	FindByID(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) (*schema.Schema, error)
	FindByFlowID(ctx context.Context, flowID string, projectID uuid.UUID) (*schema.Schema, error)
	List(ctx context.Context, projectID uuid.UUID) ([]schema.Schema, error)
	SetVersion(ctx context.Context, setVersion schema.Schema) error
}
