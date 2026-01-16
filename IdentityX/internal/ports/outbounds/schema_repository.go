package outbounds

import (
	"GoAuth/internal/domain/schema"
	"context"

	"github.com/google/uuid"
)

type SchemaRepository interface {
	Draft(ctx context.Context, toDraft schema.Schema) (*schema.Schema, error)
	Publish(ctx context.Context, toPublish schema.Schema) error
	Archive(ctx context.Context, toArchive schema.Schema) error
	Delete(ctx context.Context, toDelete schema.Schema) error
	Exists(ctx context.Context, toCheck schema.Schema) (bool, error)
	BelongsToProject(ctx context.Context, toCheck schema.Schema) (bool, error)
	FindByID(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) (*schema.Schema, error)
	FindByFlowIDAndType(ctx context.Context, flowID string, schemaType schema.Type, projectID uuid.UUID) (*schema.Schema, error)
	List(ctx context.Context, projectID uuid.UUID) ([]schema.Schema, error)
	SetVersion(ctx context.Context, toUpdate schema.Schema) error
	GetIDsFromProjectID(ctx context.Context, projectID uuid.UUID) ([]uuid.UUID, error)
}
