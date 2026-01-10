package inbounds

import (
	"context"
)

type SchemaVersionService interface {
	Draft(ctx context.Context, in SchemaVersionServiceInput) (*SchemaVersionOutput, error)
	Publish(ctx context.Context, in SchemaVersionServiceInput) error
}
