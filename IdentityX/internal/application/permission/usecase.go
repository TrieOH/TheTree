package permission

import (
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/permissions"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	usecaseTracer = otel.Tracer("permission_usecase")
)

type UseCase struct {
	permissions  outbounds.PermissionRepository
	projects     outbounds.ProjectRepository
	projectUsers outbounds.ProjectUserRepository
	sessions     outbounds.SessionRepository
	schema       inbounds.SchemaService
	tx           inbounds.TxRunner
}

var _ inbounds.PermissionService = (*UseCase)(nil)

func New(
	permissions outbounds.PermissionRepository,
	projects outbounds.ProjectRepository,
	projectUsers outbounds.ProjectUserRepository,
	sessions outbounds.SessionRepository,
	schema inbounds.SchemaService,
	tx inbounds.TxRunner,
) inbounds.PermissionService {
	return &UseCase{
		permissions:  permissions,
		projects:     projects,
		projectUsers: projectUsers,
		sessions:     sessions,
		schema:       schema,
		tx:           tx,
	}
}

func (uc *UseCase) Create(ctx context.Context, in inbounds.CreatePermissionInput) (*inbounds.PermissionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.Create")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot create permissions for a project you don't own").RecordCtx(ctx)
	}

	if err := permissions.ValidatePermission(ctx, in.Object, in.Action); err != nil {
		return nil, err
	}

	permission, err := uc.permissions.Create(ctx, permissions.Permission{
		ProjectID: in.ProjectID,
		Object:    in.Object,
		Action:    in.Action,
		Meta:      in.Meta,
	})
	if err != nil {
		return nil, err
	}

	return inbounds.PermissionToPermissionOutput(*permission), nil
}

func (uc *UseCase) UpdateMeta(ctx context.Context, in inbounds.UpdatePermissionInput) error {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.UpdateMeta")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot update permissions in a project you don't own").RecordCtx(ctx)
	}

	var permissionBelongs bool
	permissionBelongs, err = uc.permissions.BelongsToProject(ctx, in.ID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !permissionBelongs {
		return fail.New(errx.PERMissionNotOwnedByPrincipal).RecordCtx(ctx)
	}

	err = uc.permissions.UpdateMeta(ctx, in.Meta, in.ID, in.ProjectID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) Delete(ctx context.Context, in inbounds.DeletePermissionInput) error {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.Delete")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot delete permissions in a project you don't own").RecordCtx(ctx)
	}

	var permissionBelongs bool
	permissionBelongs, err = uc.permissions.BelongsToProject(ctx, in.ID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !permissionBelongs {
		return fail.New(errx.PERMissionNotOwnedByPrincipal).RecordCtx(ctx)
	}

	err = uc.permissions.Delete(ctx, in.ID, in.ProjectID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetByIDExternal(ctx context.Context, in inbounds.GetPermissionInput) (*inbounds.PermissionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.GetByIDExternal")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get permission for a project you don't own").RecordCtx(ctx)
	}

	permission, err := uc.permissions.GetByIDExternal(ctx, in.PermissionID, *in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.PermissionToPermissionOutput(*permission), nil
}

func (uc *UseCase) ListByProject(ctx context.Context, in inbounds.GetPermissionInput) (out []inbounds.PermissionOutput, err error) {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.ListByProject")
	defer span.End()

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get permissions for a project you don't own").RecordCtx(ctx)
	}

	if in.Object != nil && *in.Object != "" {
		if err = permissions.ValidateObject(ctx, *in.Object); err != nil {
			return nil, err
		}
	} else {
		in.Object = nil
	}

	if in.Action != nil && *in.Action != "" {
		if err = permissions.ValidateAction(ctx, *in.Action); err != nil {
			return nil, err
		}
	} else {
		in.Action = nil
	}

	var foundPermissions []permissions.Permission
	foundPermissions, err = uc.permissions.ListByProject(ctx, in.Object, in.Action, *in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.PermissionSliceToPermissionOutputSlice(foundPermissions), nil
}

func (uc *UseCase) GiveDirect(ctx context.Context, in inbounds.ManagePermissionInput) error {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.GiveDirect")
	defer span.End()

	isProjectGlobal := in.ScopeID == nil
	span.SetAttributes(attribute.Bool("permission.project_global", isProjectGlobal))

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot edit a project you don't own").RecordCtx(ctx)
	}

	var permissionBelongs bool
	permissionBelongs, err = uc.permissions.BelongsToProject(ctx, in.PermissionID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !permissionBelongs {
		return fail.New(errx.PERMissionNotOwnedByPrincipal).RecordCtx(ctx)
	}

	var userBelongs bool
	userBelongs, err = uc.projectUsers.BelongsToProject(ctx, in.EntityID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !userBelongs {
		return fail.New(errx.ProjectUserNotFromProject).RecordCtx(ctx)
	}

	userIdentity, err := uc.sessions.GetIdentityByEntityIDAndType(ctx, in.EntityID, session.ProjectIdentity)
	if err != nil {
		return err
	}

	if err = uc.permissions.GiveDirect(ctx, in.PermissionID, userIdentity.ID, in.ScopeID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) TakeDirect(ctx context.Context, in inbounds.ManagePermissionInput) error {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.TakeDirect")
	defer span.End()

	isProjectGlobal := in.ScopeID == nil
	span.SetAttributes(attribute.Bool("permission.project_global", isProjectGlobal))

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot edit a project you don't own").RecordCtx(ctx)
	}

	var permissionBelongs bool
	permissionBelongs, err = uc.permissions.BelongsToProject(ctx, in.PermissionID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !permissionBelongs {
		return fail.New(errx.PERMissionNotOwnedByPrincipal).RecordCtx(ctx)
	}

	var userBelongs bool
	userBelongs, err = uc.projectUsers.BelongsToProject(ctx, in.EntityID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !userBelongs {
		return fail.New(errx.ProjectUserNotFromProject).RecordCtx(ctx)
	}

	userIdentity, err := uc.sessions.GetIdentityByEntityIDAndType(ctx, in.EntityID, session.ProjectIdentity)
	if err != nil {
		return err
	}

	if err = uc.permissions.TakeDirect(ctx, in.PermissionID, userIdentity.ID, in.ScopeID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetEffective(ctx context.Context, in inbounds.ManagePermissionInput) (perms []inbounds.PermissionOutput, err error) {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.GetEffective")
	defer span.End()

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get permissions for a project you don't own").RecordCtx(ctx)
	}

	var userBelongs bool
	userBelongs, err = uc.projectUsers.BelongsToProject(ctx, in.EntityID, *in.ProjectID)
	if err != nil {
		return nil, err
	}

	if !userBelongs {
		return nil, fail.New(errx.ProjectUserNotFromProject).RecordCtx(ctx)
	}

	// CHECK COMPATIBILITY
	isUpToDate, err := uc.schema.CheckSchemaCompatibility(ctx, in.EntityID, *in.ProjectID)
	if err != nil {
		return nil, err
	}
	if !isUpToDate {
		return nil, fail.New(errx.AuthUserSchemaOutdated).RecordCtx(ctx)
	}

	userIdentity, err := uc.sessions.GetIdentityByEntityIDAndType(ctx, in.EntityID, session.ProjectIdentity)
	if err != nil {
		return nil, err
	}

	var foundPermissions []permissions.Permission
	foundPermissions, err = uc.permissions.GetEffective(ctx, userIdentity.ID, in.ProjectID, in.ScopeID)
	if err != nil {
		return nil, err
	}

	return inbounds.PermissionSliceToPermissionOutputSlice(foundPermissions), nil
}

func (uc *UseCase) Check(ctx context.Context, in inbounds.CheckPermissionInput) (hasPermission bool, err error) {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.Check")
	defer span.End()

	if err = permissions.ValidatePermission(ctx, in.Object, in.Action); err != nil {
		return false, err
	}

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return false, err
	}

	// Optional: Check if caller can query permissions (meta-auth)
	// Instead of hardcoded owner check, you might want:
	// if !uc.canCheckPermissions(ctx, principal, in.ProjectID) { ... }

	if in.ProjectID != nil {
		var isOwner bool
		isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
		if err != nil {
			return false, err
		}
		if !isOwner {
			return false, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot check permissions in a project you don't own").RecordCtx(ctx)
		}
	}

	// Validate target user belongs to project (if scoped)
	if in.ProjectID != nil {
		var belongs bool
		belongs, err = uc.projectUsers.BelongsToProject(ctx, in.EntityID, *in.ProjectID)
		if err != nil {
			return false, err
		}
		if !belongs {
			return false, fail.New(errx.ProjectUserNotFromProject).RecordCtx(ctx)
		}
	}

	var userIdentity *session.Identity
	if in.ProjectID != nil {
		userIdentity, err = uc.sessions.GetIdentityByEntityIDAndType(ctx, in.EntityID, session.ProjectIdentity)
		if err != nil {
			return false, err
		}
	} else {
		userIdentity, err = uc.sessions.GetIdentityByEntityIDAndType(ctx, in.EntityID, session.ClientIdentity)
		if err != nil {
			return false, err
		}
	}

	var perms []permissions.Permission
	perms, err = uc.permissions.GetEffective(ctx, userIdentity.ID, in.ProjectID, in.ScopeID)
	if err != nil {
		return false, err
	}

	span.SetAttributes(
		attribute.Int("permissions.checked", len(perms)),
		attribute.String("check.object", in.Object),
		attribute.String("check.action", in.Action),
	)

	// Check each permission - object AND action must match
	for _, p := range perms {
		if permissions.ObjectMatch(p.Object, in.Object) &&
			permissions.ActionMatch(p.Action, in.Action) {
			span.SetAttributes(attribute.Bool("permission.found", true))
			return true, nil
		}
	}

	span.SetAttributes(attribute.Bool("permission.found", false))
	return false, fail.New(errx.PERMissionInsufficient).RecordCtx(ctx)
}
