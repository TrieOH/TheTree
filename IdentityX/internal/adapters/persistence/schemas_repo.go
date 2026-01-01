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

type schemaRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ outbound.SchemaRepository = (*schemaRepo)(nil)

func NewSchemaRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.SchemaRepository {
	return &schemaRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func (r schemaRepo) Draft(ctx context.Context, schema schema.Schema) (*schema.Schema, error) {
	// TODO Implement me!
	return nil, apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r schemaRepo) Publish(ctx context.Context, schemaID uuid.UUID) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r schemaRepo) Archive(ctx context.Context, schemaID uuid.UUID) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r schemaRepo) Delete(ctx context.Context, schemaID uuid.UUID) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}
