package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/ports/outbound"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type schemaRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	Tracer trace.Tracer
}

var _ outbound.SchemaRepository = (*schemaRepo)(nil)

func NewSchemaRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.SchemaRepository {
	return &schemaRepo{
		q:      q,
		log:    l,
		Tracer: tracer,
	}
}

func (r schemaRepo) Draft(ctx context.Context, schema schema.Schema) (*schema.Schema, error) {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r schemaRepo) Publish(ctx context.Context, schemaID uuid.UUID) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r schemaRepo) Archive(ctx context.Context, schemaID uuid.UUID) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r schemaRepo) Delete(ctx context.Context, schemaID uuid.UUID) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}
