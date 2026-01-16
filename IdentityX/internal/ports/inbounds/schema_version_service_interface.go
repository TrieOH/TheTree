package inbounds

import (
	"context"
)

type SchemaVersionService interface {
	Draft(ctx context.Context, in SchemaVersionServiceInput) (*SchemaVersionOutput, error)
	Publish(ctx context.Context, in SchemaVersionServiceInput) error
	GetCurrent(ctx context.Context, in SchemaVersionServiceInput) (*SchemaVersionOutput, error)
	GetLatest(ctx context.Context, in SchemaVersionServiceInput) (*SchemaVersionOutput, error)
	GetVerbose(ctx context.Context, in SchemaVersionServiceInput) (*VersionVerboseOutput, error)
}
