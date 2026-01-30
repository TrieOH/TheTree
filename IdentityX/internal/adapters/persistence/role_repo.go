package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/permissions"
	"GoAuth/internal/domain/roles"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type roleRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *roleRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbounds.RoleRepository = (*roleRepo)(nil)

func NewRoleRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbounds.RoleRepository {
	return &roleRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func mapRoleFromDB(dst *roles.Role, src *sqlc.Role) {
	dst.ID = src.ID
	dst.ProjectID = src.ProjectID
	dst.Name = src.Name
	dst.Description = src.Description
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func mapGetUserRolesRowFromDB(dst *roles.Role, src *sqlc.GetUserRolesRow) {
	dst.ID = src.ID
	dst.ProjectID = src.ProjectID
	dst.Name = src.Name
	dst.Description = src.Description
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.ScopeName = src.ScopeName
	dst.ScopeID = src.ScopeID
	dst.ExternalID = src.ExternalID
}

func (repo *roleRepo) Create(ctx context.Context, toCreate roles.Role) (*roles.Role, error) {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.Create")
	defer span.End()

	if toCreate.ProjectID != nil {
		span.SetAttributes(attribute.String("role.project_id", toCreate.ProjectID.String()))
	}

	sqlcRole, err := repo.queries(ctx).CreateRole(ctx, sqlc.CreateRoleParams{
		Name:        toCreate.Name,
		Description: toCreate.Description,
		ProjectID:   toCreate.ProjectID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("role.id", sqlcRole.ID.String()))

	var outRole roles.Role
	mapRoleFromDB(&outRole, &sqlcRole)
	return &outRole, nil
}

func (repo *roleRepo) UpdateDescription(ctx context.Context, description string, id uuid.UUID, projectID *uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.UpdateDescription",
		trace.WithAttributes(
			attribute.String("role.id", id.String()),
		),
	)
	defer span.End()

	if projectID != nil {
		span.SetAttributes(attribute.String("role.project_id", projectID.String()))
	}

	err := repo.queries(ctx).UpdateRoleDescription(ctx, sqlc.UpdateRoleDescriptionParams{
		ID:          id,
		ProjectID:   projectID,
		Description: &description,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (repo *roleRepo) GetByIDInternal(ctx context.Context, id uuid.UUID) (*roles.Role, error) {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.GetByIDInternal",
		trace.WithAttributes(
			attribute.String("role.id", id.String()),
		),
	)
	defer span.End()

	sqlcRole, err := repo.queries(ctx).GetRoleByIDInternal(ctx, id)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	var outRole roles.Role
	mapRoleFromDB(&outRole, &sqlcRole)
	return &outRole, nil
}

func (repo *roleRepo) GetByIDExternal(ctx context.Context, id, projectID uuid.UUID) (*roles.Role, error) {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.GetByIDExternal",
		trace.WithAttributes(
			attribute.String("role.project_id", projectID.String()),
			attribute.String("role.id", id.String()),
		),
	)
	defer span.End()

	sqlcRole, err := repo.queries(ctx).GetRoleByIDExternal(ctx, sqlc.GetRoleByIDExternalParams{
		ID:        id,
		ProjectID: &projectID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	var outRole roles.Role
	mapRoleFromDB(&outRole, &sqlcRole)
	return &outRole, nil
}

func (repo *roleRepo) GetByName(ctx context.Context, name string, projectID *uuid.UUID) (*roles.Role, error) {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.GetByName")
	defer span.End()

	if projectID != nil {
		span.SetAttributes(attribute.String("role.project_id", projectID.String()))
	}

	sqlcRole, err := repo.queries(ctx).GetRoleByName(ctx, sqlc.GetRoleByNameParams{
		Name:      name,
		ProjectID: projectID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("role.id", sqlcRole.ID.String()))

	var outRole roles.Role
	mapRoleFromDB(&outRole, &sqlcRole)
	return &outRole, nil
}

func (repo *roleRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]roles.Role, error) {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.ListByProject",
		trace.WithAttributes(
			attribute.String("project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcRoles, err := repo.queries(ctx).ListRolesByProject(ctx, &projectID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.Int("role.count", len(sqlcRoles)))

	outRoles := make([]roles.Role, 0, len(sqlcRoles))
	for _, sqlcRole := range sqlcRoles {
		var outRole roles.Role
		mapRoleFromDB(&outRole, &sqlcRole)
		outRoles = append(outRoles, outRole)
	}
	return outRoles, nil
}

func (repo *roleRepo) BelongsToProject(ctx context.Context, id, projectID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.BelongsToProject",
		trace.WithAttributes(
			attribute.String("role.id", id.String()),
			attribute.String("role.project_id", projectID.String()),
		),
	)
	defer span.End()

	belongs, err := repo.queries(ctx).RoleBelongsToProject(ctx, sqlc.RoleBelongsToProjectParams{
		ID:        id,
		ProjectID: &projectID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return false, sqlcErr
	}

	return belongs, nil
}

func (repo *roleRepo) AddPermission(ctx context.Context, id uuid.UUID, permissionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.AddPermission",
		trace.WithAttributes(
			attribute.String("role.id", id.String()),
			attribute.String("role.permission_id", permissionID.String()),
		),
	)
	defer span.End()

	err := repo.queries(ctx).AddPermissionToRole(ctx, sqlc.AddPermissionToRoleParams{
		RoleID:       id,
		PermissionID: permissionID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (repo *roleRepo) RemovePermission(ctx context.Context, id uuid.UUID, permissionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.RemovePermission",
		trace.WithAttributes(
			attribute.String("role.id", id.String()),
			attribute.String("role.permission_id", permissionID.String()),
		),
	)
	defer span.End()

	err := repo.queries(ctx).RemovePermissionFromRole(ctx, sqlc.RemovePermissionFromRoleParams{
		RoleID:       id,
		PermissionID: permissionID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (repo *roleRepo) GetPermissions(ctx context.Context, id, projectID uuid.UUID) ([]permissions.Permission, error) {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.GetPermissions",
		trace.WithAttributes(
			attribute.String("role.id", id.String()),
			attribute.String("role.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcPermissions, err := repo.queries(ctx).GetRolePermissions(ctx, sqlc.GetRolePermissionsParams{
		RoleID:    id,
		ProjectID: &projectID,
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

func (repo *roleRepo) GiveRole(ctx context.Context, id, identityID uuid.UUID, scopeID *uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.GiveRole",
		trace.WithAttributes(
			attribute.String("role.id", id.String()),
			attribute.String("role.identity_id", identityID.String()),
		),
	)

	if scopeID != nil {
		span.SetAttributes(attribute.String("role.scope_id", scopeID.String()))
	}

	defer span.End()

	err := repo.queries(ctx).GiveRole(ctx, sqlc.GiveRoleParams{
		RoleID:     id,
		IdentityID: identityID,
		ScopeID:    scopeID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (repo *roleRepo) TakeRole(ctx context.Context, id, identityID uuid.UUID, scopeID *uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.TakeRole",
		trace.WithAttributes(
			attribute.String("role.id", id.String()),
			attribute.String("role.identity_id", identityID.String()),
		),
	)

	if scopeID != nil {
		span.SetAttributes(attribute.String("role.scope_id", scopeID.String()))
	}

	defer span.End()

	err := repo.queries(ctx).TakeRole(ctx, sqlc.TakeRoleParams{
		RoleID:     id,
		IdentityID: identityID,
		ScopeID:    scopeID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (repo *roleRepo) GetUserRoles(ctx context.Context, identityID, projectID uuid.UUID) ([]roles.Role, error) {
	ctx, span := repo.tracer.Start(ctx, "RoleRepo.GetUserRoles",
		trace.WithAttributes(
			attribute.String("user_id", identityID.String()),
			attribute.String("project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcRoles, err := repo.queries(ctx).GetUserRoles(ctx, sqlc.GetUserRolesParams{
		IdentityID: identityID,
		ProjectID:  &projectID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.Int("roles.count", len(sqlcRoles)))

	outRoles := make([]roles.Role, 0, len(sqlcRoles))
	for _, sqlcRole := range sqlcRoles {
		var outRole roles.Role
		mapGetUserRolesRowFromDB(&outRole, &sqlcRole)
		outRoles = append(outRoles, outRole)
	}
	return outRoles, nil
}
