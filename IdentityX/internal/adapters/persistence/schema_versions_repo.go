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

type schemaVersionRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	Tracer trace.Tracer
}

var _ outbound.SchemaVersionRepository = (*schemaVersionRepo)(nil)

func NewSchemaVersionRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.SchemaVersionRepository {
	return &schemaVersionRepo{
		q:      q,
		log:    l,
		Tracer: tracer,
	}
}

func (r schemaVersionRepo) Draft(ctx context.Context, version schema.Version) (*schema.Version, error) {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r schemaVersionRepo) Publish(ctx context.Context, versionID uuid.UUID) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r schemaVersionRepo) Archive(ctx context.Context, versionID uuid.UUID) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r schemaVersionRepo) Delete(ctx context.Context, versionID uuid.UUID) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}
