package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/ports/outbound"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type fieldsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

var _ outbound.FieldsRepository = (*fieldsRepo)(nil)

func NewFieldsRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.FieldsRepository {
	return &fieldsRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func (r fieldsRepo) Create(ctx context.Context, field field.Field) (*field.Field, error) {
	// TODO Implement me!
	return nil, apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r fieldsRepo) Update(ctx context.Context, field field.Field) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r fieldsRepo) SetOptions(ctx context.Context, options []field.Option) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r fieldsRepo) SetRequiredRules(ctx context.Context, required []field.RequiredRule) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r fieldsRepo) SetVisibilityRules(ctx context.Context, visibilityRules []field.VisibilityRule) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (r fieldsRepo) Delete(ctx context.Context, fieldID uuid.UUID) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}
