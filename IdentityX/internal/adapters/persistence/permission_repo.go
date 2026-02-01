package persistence

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/domain/permissions"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type permissionRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *permissionRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbounds.PermissionRepository = (*permissionRepo)(nil)

func NewPermissionRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbounds.PermissionRepository {
	return &permissionRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func mapPermissionFromDB(dst *permissions.Permission, src *sqlc.Permission) {
	dst.ID = src.ID
	dst.ProjectID = src.ProjectID
	dst.Object = src.Object
	dst.Action = src.Action
	dst.CreatedAt = src.CreatedAt
	// FIXME deal with error
	var err error
	dst.Conditions, err = permissions.DecodeCondition(src.Conditions)
	if err != nil {
		logs.L().Error("error while encoding condition in permission repo", zap.Error(err))
	}
}

func (repo *permissionRepo) Create(ctx context.Context, toCreate outbounds.CreatePermissionInput) (*permissions.Permission, error) {
	ctx, span := repo.tracer.Start(ctx, "PermissionRepo.Create")
	defer span.End()

	if toCreate.ProjectID != nil {
		span.SetAttributes(attribute.String("permission.project_id", toCreate.ProjectID.String()))
	}

	sqlcPermission, err := repo.queries(ctx).CreatePermission(ctx, sqlc.CreatePermissionParams{
		ProjectID:  toCreate.ProjectID,
		Object:     toCreate.Object,
		Action:     toCreate.Action,
		Conditions: toCreate.Conditions,
	})
	if err != nil {
		return nil, fail.From(err)
	}

	span.SetAttributes(attribute.String("permission.id", sqlcPermission.ID.String()))

	var outPermission permissions.Permission
	mapPermissionFromDB(&outPermission, &sqlcPermission)
	return &outPermission, nil
}

func (repo *permissionRepo) GetByIDInternal(ctx context.Context, id uuid.UUID) (*permissions.Permission, error) {
	ctx, span := repo.tracer.Start(ctx, "PermissionRepo.GetByIDInternal",
		trace.WithAttributes(
			attribute.String("permission.id", id.String()),
		),
	)
	defer span.End()

	sqlcPermission, err := repo.queries(ctx).GetPermissionByIDInternal(ctx, id)
	if err != nil {
		return nil, fail.From(err)
	}

	if sqlcPermission.ProjectID != nil {
		span.SetAttributes(attribute.String("permission.project_id", sqlcPermission.ProjectID.String()))
	}

	var outPermission permissions.Permission
	mapPermissionFromDB(&outPermission, &sqlcPermission)
	return &outPermission, nil
}

func (repo *permissionRepo) GetByIDExternal(ctx context.Context, id, projectID uuid.UUID) (*permissions.Permission, error) {
	ctx, span := repo.tracer.Start(ctx, "PermissionRepo.GetByIDExternal",
		trace.WithAttributes(
			attribute.String("permission.id", id.String()),
			attribute.String("permission.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcPermission, err := repo.queries(ctx).GetPermissionByIDExternal(ctx, sqlc.GetPermissionByIDExternalParams{
		ID:        id,
		ProjectID: &projectID,
	})
	if err != nil {
		return nil, fail.From(err)
	}

	var outPermission permissions.Permission
	mapPermissionFromDB(&outPermission, &sqlcPermission)
	return &outPermission, nil
}

func (repo *permissionRepo) ListByProject(ctx context.Context, object, action *string, projectID uuid.UUID) ([]permissions.Permission, error) {
	ctx, span := repo.tracer.Start(ctx, "PermissionRepo.ListByProject",
		trace.WithAttributes(
			attribute.String("permission.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcPermissions, err := repo.queries(ctx).ListPermissionsByProject(ctx, sqlc.ListPermissionsByProjectParams{
		ProjectID: &projectID,
		Object:    object,
		Action:    action,
	})
	if err != nil {
		return nil, fail.From(err)
	}

	span.SetAttributes(attribute.Int("permission.count", len(sqlcPermissions)))

	outPermissions := make([]permissions.Permission, 0, len(sqlcPermissions))
	for _, sqlcPermission := range sqlcPermissions {
		var outPermission permissions.Permission
		mapPermissionFromDB(&outPermission, &sqlcPermission)
		outPermissions = append(outPermissions, outPermission)
	}
	return outPermissions, nil
}

func (repo *permissionRepo) BelongsToProject(ctx context.Context, id, projectID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "PermissionRepo.BelongsToProject",
		trace.WithAttributes(
			attribute.String("permission.id", id.String()),
			attribute.String("permission.project_id", projectID.String()),
		),
	)
	defer span.End()

	belongs, err := repo.queries(ctx).PermissionBelongsToProject(ctx, sqlc.PermissionBelongsToProjectParams{
		ID:        id,
		ProjectID: &projectID,
	})
	if err != nil {
		return false, fail.From(err)
	}

	return belongs, nil
}

func (repo *permissionRepo) GiveDirect(ctx context.Context, id, identityID uuid.UUID, scopeID *uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "PermissionRepo.GiveDirectPermission",
		trace.WithAttributes(
			attribute.String("permission.id", id.String()),
			attribute.String("permission.identity_id", identityID.String()),
		),
	)

	if scopeID != nil {
		span.SetAttributes(attribute.String("permission.scope_id", scopeID.String()))
	}

	defer span.End()

	err := repo.queries(ctx).GiveDirectPermission(ctx, sqlc.GiveDirectPermissionParams{
		PermissionID: id,
		IdentityID:   identityID,
		ScopeID:      scopeID,
	})
	if err != nil {
		return fail.From(err)
	}

	return nil
}

func (repo *permissionRepo) TakeDirect(ctx context.Context, id, identityID uuid.UUID, scopeID *uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "PermissionRepo.TakeDirectPermission",
		trace.WithAttributes(
			attribute.String("permission.id", id.String()),
			attribute.String("permission.identity_id", identityID.String()),
		),
	)

	if scopeID != nil {
		span.SetAttributes(attribute.String("permission.scope_id", scopeID.String()))
	}

	defer span.End()

	err := repo.queries(ctx).TakeDirectPermission(ctx, sqlc.TakeDirectPermissionParams{
		PermissionID: id,
		IdentityID:   identityID,
		ScopeID:      scopeID,
	})
	if err != nil {
		return fail.From(err)
	}

	return nil
}

func (repo *permissionRepo) GetEffective(ctx context.Context, identityID uuid.UUID, projectID, scopeID *uuid.UUID) ([]permissions.Permission, error) {
	ctx, span := repo.tracer.Start(ctx, "PermissionRepo.GetEffective",
		trace.WithAttributes(
			attribute.String("identity_id", identityID.String()),
		),
	)

	if scopeID != nil {
		span.SetAttributes(attribute.String("scope_id", scopeID.String()))
	}

	if projectID != nil {
		span.SetAttributes(attribute.String("project_id", projectID.String()))
	}

	defer span.End()

	sqlcPermissions, err := repo.queries(ctx).GetEffectivePermissions(ctx, sqlc.GetEffectivePermissionsParams{
		ProjectID:  projectID,
		IdentityID: identityID,
		ScopeID:    scopeID,
	})
	if err != nil {
		return nil, fail.From(err)
	}

	span.SetAttributes(attribute.Int("permission.count", len(sqlcPermissions)))

	outPermissions := make([]permissions.Permission, 0, len(sqlcPermissions))
	for _, sqlcPermission := range sqlcPermissions {
		var outPermission permissions.Permission
		mapPermissionFromDB(&outPermission, &sqlcPermission)
		outPermissions = append(outPermissions, outPermission)
	}
	return outPermissions, nil
}
