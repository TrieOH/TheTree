package permission

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/permissions"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"

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
	tx           inbounds.TxRunner
}

var _ inbounds.PermissionService = (*UseCase)(nil)

func New(
	permissions outbounds.PermissionRepository,
	projects outbounds.ProjectRepository,
	projectUsers outbounds.ProjectUserRepository,
	sessions outbounds.SessionRepository,
	tx inbounds.TxRunner,
) inbounds.PermissionService {
	return &UseCase{
		permissions:  permissions,
		projects:     projects,
		projectUsers: projectUsers,
		sessions:     sessions,
		tx:           tx,
	}
}

func (uc *UseCase) Create(ctx context.Context, in inbounds.CreatePermissionInput) (*inbounds.PermissionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.Create")
	defer span.End()

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot create permissions for a project you don't own"})
	}

	if err := permissions.ValidatePermission(in.Object, in.Action); err != nil {
		return nil, apierr.FromService(span, err)
	}

	permission, err := uc.permissions.Create(ctx, permissions.Permission{
		ProjectID:  in.ProjectID,
		Object:     in.Object,
		Action:     in.Action,
		Conditions: in.Conditions,
	})
	if err != nil {
		return nil, err
	}

	return inbounds.PermissionToPermissionOutput(*permission), nil
}

func (uc *UseCase) GetByIDExternal(ctx context.Context, in inbounds.GetPermissionInput) (*inbounds.PermissionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.GetByIDExternal")
	defer span.End()

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot get permission for a project you don't own"})
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
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot get permissions for a project you don't own"})
	}

	if in.Object != nil && *in.Object != "" {
		if err = permissions.ValidateObject(*in.Object); err != nil {
			return nil, apierr.FromService(span, err)
		}
	} else {
		in.Object = nil
	}

	if in.Action != nil && *in.Action != "" {
		if err = permissions.ValidateAction(*in.Action); err != nil {
			return nil, apierr.FromService(span, err)
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

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot edit a project you don't own"})
	}

	var permissionBelongs bool
	permissionBelongs, err = uc.permissions.BelongsToProject(ctx, in.PermissionID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !permissionBelongs {
		return apierr.FromService(span, inbounds.ErrRoleNotOwned{Msg: "cannot edit a permission you don't own"})
	}

	var userBelongs bool
	userBelongs, err = uc.projectUsers.BelongsToProject(ctx, in.EntityID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !userBelongs {
		return apierr.FromService(span, inbounds.ErrProjectUserNotFromProject{})
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

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot edit a project you don't own"})
	}

	var permissionBelongs bool
	permissionBelongs, err = uc.permissions.BelongsToProject(ctx, in.PermissionID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !permissionBelongs {
		return apierr.FromService(span, inbounds.ErrRoleNotOwned{Msg: "cannot edit a permission you don't own"})
	}

	var userBelongs bool
	userBelongs, err = uc.projectUsers.BelongsToProject(ctx, in.EntityID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !userBelongs {
		return apierr.FromService(span, inbounds.ErrProjectUserNotFromProject{})
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
