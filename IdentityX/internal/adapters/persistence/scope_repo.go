package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/domain/scopes"
	"GoAuth/internal/ports/outbounds"
	"context"
	"encoding/json"

	"github.com/MintzyG/fail/v3"
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
	dst.ParentID = src.ParentID
	dst.ProjectID = src.ProjectID
	dst.ExternalID = src.ExternalID
	dst.Name = src.Name
	dst.Meta = src.Meta
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

	if toAnnotate.ParentID != nil {
		span.SetAttributes(attribute.String("scope.parent_id", toAnnotate.ParentID.String()))
	}
}

func (repo *scopeRepo) Create(ctx context.Context, toCreate scopes.Scope) (*scopes.Scope, error) {
	ctx, span := repo.tracer.Start(ctx, "ScopeRepo.Create")
	annotateScope(span, toCreate)
	defer span.End()

	sqlcScope, err := repo.queries(ctx).CreateScope(ctx, sqlc.CreateScopeParams{
		Type:       string(toCreate.Type),
		ParentID:   toCreate.ParentID,
		ProjectID:  toCreate.ProjectID,
		Name:       toCreate.Name,
		ExternalID: toCreate.ExternalID,
		Meta:       toCreate.Meta,
	})

	if err != nil {
		return nil, fail.From(err).WithArgs(*toCreate.Name, *toCreate.ExternalID).RecordCtx(ctx)
	}

	var created scopes.Scope
	mapScopeFromDB(&created, &sqlcScope)
	return &created, nil
}

func (repo *scopeRepo) UpdateMeta(ctx context.Context, meta *json.RawMessage, id uuid.UUID, projectID *uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ScopeRepo.UpdateMeta",
		trace.WithAttributes(
			attribute.String("scope.id", id.String()),
		),
	)
	defer span.End()

	if projectID != nil {
		span.SetAttributes(attribute.String("scope.project_id", projectID.String()))
	}

	err := repo.queries(ctx).UpdateProjectScopeMeta(ctx, sqlc.UpdateProjectScopeMetaParams{
		ID:        id,
		ProjectID: projectID,
		Meta:      meta,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
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
		return nil, fail.From(err).RecordCtx(ctx)
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
		return nil, fail.From(err).RecordCtx(ctx)
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
		return nil, fail.From(err).RecordCtx(ctx)
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
		return nil, fail.From(err).RecordCtx(ctx)
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

func (repo *scopeRepo) GetAncestors(ctx context.Context, scopeID uuid.UUID) ([]scopes.Scope, error) {
	ctx, span := repo.tracer.Start(ctx, "ScopeRepo.GetAncestors",
		trace.WithAttributes(
			attribute.String("scope.id", scopeID.String()),
		),
	)
	defer span.End()

	sqlcScopes, err := repo.queries(ctx).GetScopeAncestors(ctx, scopeID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Int("ancestor.count", len(sqlcScopes)))

	found := make([]scopes.Scope, 0, len(sqlcScopes))
	for _, sqlcScope := range sqlcScopes {
		scope := scopes.Scope{
			ID:         sqlcScope.ID,
			Type:       scopeTypeFromSQLCScopeType(sqlcScope.Type),
			ParentID:   sqlcScope.ParentID,
			ProjectID:  sqlcScope.ProjectID,
			Name:       sqlcScope.Name,
			ExternalID: sqlcScope.ExternalID,
			CreatedAt:  sqlcScope.CreatedAt,
		}
		found = append(found, scope)
	}
	return found, nil
}

func (repo *scopeRepo) GetGlobalScope(ctx context.Context) (*scopes.Scope, error) {
	ctx, span := repo.tracer.Start(ctx, "ScopeRepo.GetGlobalScope")
	defer span.End()

	sqlcScope, err := repo.queries(ctx).GetGlobalScope(ctx)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var created scopes.Scope
	mapScopeFromDB(&created, &sqlcScope)
	annotateScope(span, created)
	return &created, nil
}

func (repo *scopeRepo) GetChildren(ctx context.Context, scopeID uuid.UUID) ([]scopes.Scope, error) {
	ctx, span := repo.tracer.Start(ctx, "ScopeRepo.GetChildren",
		trace.WithAttributes(
			attribute.String("scope.id", scopeID.String()),
		),
	)
	defer span.End()

	sqlcScopes, err := repo.queries(ctx).GetScopeChildren(ctx, &scopeID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Int("children.count", len(sqlcScopes)))

	found := make([]scopes.Scope, 0, len(sqlcScopes))
	for _, sqlcScope := range sqlcScopes {
		var scope scopes.Scope
		mapScopeFromDB(&scope, &sqlcScope)
		found = append(found, scope)
	}
	return found, nil
}
