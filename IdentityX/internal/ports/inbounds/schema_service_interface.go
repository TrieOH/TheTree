package inbounds

import (
	"context"
)

type SchemaService interface {
	Draft(ctx context.Context, in SchemaServiceInput) (*SchemaOutput, error)
	Publish(ctx context.Context, in SchemaServiceInput) error
	GetByID(ctx context.Context, in SchemaServiceInput) (*SchemaOutput, error)
	GetVerbose(ctx context.Context, in SchemaServiceInput) (*SchemaVerboseOutput, error)
}
