package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/scopes"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type scopeRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *scopeRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbounds.ScopeRepository = (*scopeRepo)(nil)

func NewScopeRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbounds.ScopeRepository {
	return &scopeRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func scopeTypeFromSQLCScopeType(scopeType string) scopes.ScopeType {
	switch scopes.ScopeType(scopeType) {
	case scopes.ScopeTypeGlobal:
		return scopes.ScopeTypeGlobal
	case scopes.ScopeTypeProjectRoot:
		return scopes.ScopeTypeProjectRoot
	case scopes.ScopeTypeProjectScope:
		return scopes.ScopeTypeProjectScope
	default:
		return scopes.ScopeTypeNone
	}
}

func mapScopeFromDB(dst *scopes.Scope, src *sqlc.Scope) {
	dst.ID = src.ID
	dst.Type = scopeTypeFromSQLCScopeType(src.Type)
	dst.ProjectID = src.ProjectID
	dst.ExternalID = src.ExternalID
	dst.Name = src.Name
	dst.CreatedAt = src.CreatedAt
}

func annotateScope(span trace.Span, toAnnotate scopes.Scope) {
	span.SetAttributes(attribute.String("scope.type", string(toAnnotate.Type)))

	if toAnnotate.Name != nil {
		span.SetAttributes(attribute.String("scope.name", *toAnnotate.Name))
	}

	if toAnnotate.ProjectID != nil {
		span.SetAttributes(attribute.String("scope.project_id", toAnnotate.ProjectID.String()))
	}

	if toAnnotate.ExternalID != nil {
		span.SetAttributes(attribute.String("scope.external_id", *toAnnotate.ExternalID))
	}
}

func (repo *scopeRepo) Create(ctx context.Context, toCreate scopes.Scope) (*scopes.Scope, error) {
	ctx, span := repo.tracer.Start(ctx, "ScopeRepo.Create")
	annotateScope(span, toCreate)
	defer span.End()

	sqlcScope, err := repo.queries(ctx).CreateScope(ctx, sqlc.CreateScopeParams{
		Type:       string(toCreate.Type),
		ProjectID:  toCreate.ProjectID,
		Name:       toCreate.Name,
		ExternalID: toCreate.ExternalID,
	})

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
	}

	var created scopes.Scope
	mapScopeFromDB(&created, &sqlcScope)
	return &created, nil
}

func (repo *scopeRepo) GetByIDInternal(ctx context.Context, id uuid.UUID) (*scopes.Scope, error) {
	ctx, span := repo.tracer.Start(ctx, "ScopeRepo.GetByIDInternal",
		trace.WithAttributes(
			attribute.String("scope.id", id.String()),
		),
	)
	defer span.End()

	sqlcScope, err := repo.queries(ctx).GetScopeByIDInternal(ctx, id)
	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
	}

	var created scopes.Scope
	mapScopeFromDB(&created, &sqlcScope)
	annotateScope(span, created)
	return &created, nil
}

func (repo *scopeRepo) GetByIDExternal(ctx context.Context, id, projectID uuid.UUID) (*scopes.Scope, error) {
	ctx, span := repo.tracer.Start(ctx, "ScopeRepo.GetByIDExternal",
		trace.WithAttributes(
			attribute.String("scope.id", id.String()),
		),
	)
	defer span.End()

	sqlcScope, err := repo.queries(ctx).GetScopeByIDExternal(ctx, sqlc.GetScopeByIDExternalParams{
		ID:        id,
		ProjectID: &projectID,
	})
	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
	}

	var created scopes.Scope
	mapScopeFromDB(&created, &sqlcScope)
	annotateScope(span, created)
	return &created, nil
}

func (repo *scopeRepo) GetRootByProjectID(ctx context.Context, projectID uuid.UUID) (*scopes.Scope, error) {
	ctx, span := repo.tracer.Start(ctx, "ScopeRepo.GetRootByProjectID",
		trace.WithAttributes(
			attribute.String("project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcScope, err := repo.queries(ctx).GetRootByProjectID(ctx, &projectID)
	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
	}

	var created scopes.Scope
	mapScopeFromDB(&created, &sqlcScope)
	annotateScope(span, created)
	return &created, nil
}

func (repo *scopeRepo) GetProjectScopes(ctx context.Context, projectID uuid.UUID) ([]scopes.Scope, error) {
	ctx, span := repo.tracer.Start(ctx, "ScopeRepo.GetProjectScopes")
	defer span.End()

	sqlcScopes, err := repo.queries(ctx).GetProjectScopes(ctx, &projectID)
	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
	}

	span.SetAttributes(attribute.Int("scope.count", len(sqlcScopes)))

	found := make([]scopes.Scope, 0, len(sqlcScopes))
	for _, sqlcScope := range sqlcScopes {
		var created scopes.Scope
		mapScopeFromDB(&created, &sqlcScope)
		found = append(found, created)
	}
	return found, nil
}
