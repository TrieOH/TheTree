package inbounds

import (
	"context"

	"github.com/google/uuid"
)

type SchemaService interface {
	Draft(ctx context.Context, in SchemaServiceInput) (*SchemaOutput, error)
	Publish(ctx context.Context, in SchemaServiceInput) error
	GetByID(ctx context.Context, in SchemaServiceInput) (*SchemaOutput, error)
	GetVerbose(ctx context.Context, in SchemaServiceInput) (*SchemaVerboseOutput, error)
	GetIDsFromProjectID(ctx context.Context, projectID uuid.UUID) ([]uuid.UUID, error)
	List(ctx context.Context, projectID uuid.UUID) ([]SchemaOutput, error)
}
