package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/permissions"
	"GoAuth/internal/ports/outbounds"
	"context"
	"database/sql"

	"github.com/google/uuid"
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
	if tx, ok := ctx.Value(transactions.TxKeyValue).(*sql.Tx); ok && tx != nil {
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
	dst.Conditions = src.Conditions
	dst.CreatedAt = src.CreatedAt
}

func (repo *permissionRepo) Create(ctx context.Context, toCreate permissions.Permission) (*permissions.Permission, error) {
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
