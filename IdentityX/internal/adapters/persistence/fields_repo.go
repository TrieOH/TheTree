package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/ports/outbound"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type fieldsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	Tracer trace.Tracer
}

var _ outbound.FieldsRepository = (*fieldsRepo)(nil)

func NewFieldsRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.FieldsRepository {
	return &fieldsRepo{
		q:      q,
		log:    l,
		Tracer: tracer,
	}
}

func (r fieldsRepo) Create(ctx context.Context, field field.Field) (*field.Field, error) {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r fieldsRepo) Update(ctx context.Context, field field.Field) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r fieldsRepo) SetOptions(ctx context.Context, options []field.Option) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r fieldsRepo) SetRequiredRules(ctx context.Context, required []field.RequiredRule) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r fieldsRepo) SetVisibilityRules(ctx context.Context, visibilityRules []field.VisibilityRule) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}

func (r fieldsRepo) Delete(ctx context.Context, FieldID uuid.UUID) error {
	// TODO Implement me!
	panic(errors.New("not implemented"))
}
