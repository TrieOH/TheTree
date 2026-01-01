package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/ports/outbound"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type schemaVersionRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ outbound.SchemaVersionRepository = (*schemaVersionRepo)(nil)

func NewSchemaVersionRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.SchemaVersionRepository {
	return &schemaVersionRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func (r schemaVersionRepo) Draft(ctx context.Context, version schema.Version) (*schema.Version, error) {
	// TODO Implement me!
	return nil, apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r schemaVersionRepo) Publish(ctx context.Context, versionID uuid.UUID) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r schemaVersionRepo) Archive(ctx context.Context, versionID uuid.UUID) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r schemaVersionRepo) Delete(ctx context.Context, versionID uuid.UUID) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}
