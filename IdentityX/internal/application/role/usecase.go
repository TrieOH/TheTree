package role

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/domain/roles"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	usecaseTracer = otel.Tracer("role_usecase")
)

type UseCase struct {
	roles        outbounds.RoleRepository
	permissions  outbounds.PermissionRepository
	projects     outbounds.ProjectRepository
	projectUsers outbounds.ProjectUserRepository
	sessions     outbounds.SessionRepository
	tx           inbounds.TxRunner
}

var _ inbounds.RoleService = (*UseCase)(nil)

func New(
	roles outbounds.RoleRepository,
	permissions outbounds.PermissionRepository,
	projects outbounds.ProjectRepository,
	projectUsers outbounds.ProjectUserRepository,
	sessions outbounds.SessionRepository,
	tx inbounds.TxRunner,
) inbounds.RoleService {
	return &UseCase{
		roles:        roles,
		permissions:  permissions,
		projects:     projects,
		projectUsers: projectUsers,
		sessions:     sessions,
		tx:           tx,
	}
}

func (uc *UseCase) Create(ctx context.Context, in inbounds.RoleInput) (*inbounds.RoleOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.Create")
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
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot create roles for a project you don't own"})
	}

	role, err := uc.roles.Create(ctx, roles.Role{
		ProjectID:   in.ProjectID,
		Name:        in.Name,
		Description: in.Description,
	})
	if err != nil {
		return nil, err
	}

	return inbounds.RoleToRoleOutput(*role), nil
}

func (uc *UseCase) UpdateDescription(ctx context.Context, in inbounds.RoleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.UpdateDescription")
	defer span.End()

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
		return apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot update roles in a project you don't own"})
	}

	if in.Description == nil {
		err = uc.roles.UpdateDescription(ctx, "", in.RoleID, in.ProjectID)
	} else {
		err = uc.roles.UpdateDescription(ctx, *in.Description, in.RoleID, in.ProjectID)
	}

	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetByIDExternal(ctx context.Context, in inbounds.GetRoleInput) (*inbounds.RoleOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.GetByIDExternal")
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
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot get roles from a project you don't own"})
	}

	role, err := uc.roles.GetByIDExternal(ctx, in.RoleID, *in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.RoleToRoleOutput(*role), nil
}

func (uc *UseCase) GetByName(ctx context.Context, in inbounds.GetRoleInput) (*inbounds.RoleOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.GetByName")
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
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot get roles from a project you don't own"})
	}

	role, err := uc.roles.GetByName(ctx, in.Name, in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.RoleToRoleOutput(*role), nil
}

func (uc *UseCase) ListByProject(ctx context.Context, in inbounds.GetRoleInput) ([]inbounds.RoleOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.ListByProject")
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
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot get roles from a project you don't own"})
	}

	foundRoles, err := uc.roles.ListByProject(ctx, *in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.RoleSliceToRoleOutputSlice(foundRoles), nil
}

func (uc *UseCase) AddPermission(ctx context.Context, in inbounds.RolePermissionInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.AddPermission")
	defer span.End()

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

	var roleBelongs bool
	roleBelongs, err = uc.roles.BelongsToProject(ctx, in.RoleID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !roleBelongs {
		return fail.New(apierr.ROLENotOwnedByPrincipal)
	}

	var permissionBelongs bool
	permissionBelongs, err = uc.permissions.BelongsToProject(ctx, in.PermissionID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !permissionBelongs {
		return fail.New(apierr.PERMissionNotOwnedByPrincipal).Trace("cannot add permission to a role you don't own")
	}

	if err = uc.roles.AddPermission(ctx, in.RoleID, in.PermissionID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) RemovePermission(ctx context.Context, in inbounds.RolePermissionInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.RemovePermission")
	defer span.End()

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

	var roleBelongs bool
	roleBelongs, err = uc.roles.BelongsToProject(ctx, in.RoleID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !roleBelongs {
		return fail.New(apierr.ROLENotOwnedByPrincipal)
	}

	var permissionBelongs bool
	permissionBelongs, err = uc.permissions.BelongsToProject(ctx, in.PermissionID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !permissionBelongs {
		return fail.New(apierr.PERMissionNotOwnedByPrincipal).Trace("cannot remove permission to a role you don't own")
	}

	if err = uc.roles.RemovePermission(ctx, in.RoleID, in.PermissionID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetPermissions(ctx context.Context, in inbounds.RolePermissionInput) ([]inbounds.PermissionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.GetPermissions")
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
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot fetch from a project you don't own"})
	}

	permissions, err := uc.roles.GetPermissions(ctx, in.RoleID, *in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.PermissionSliceToPermissionOutputSlice(permissions), nil
}

func (uc *UseCase) GiveRole(ctx context.Context, in inbounds.ManageRoleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.GiveRole")
	defer span.End()

	isProjectGlobal := in.ScopeID == nil
	span.SetAttributes(attribute.Bool("role.project_global", isProjectGlobal))

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

	var roleBelongs bool
	roleBelongs, err = uc.roles.BelongsToProject(ctx, in.RoleID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !roleBelongs {
		return fail.New(apierr.ROLENotOwnedByPrincipal)
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

	if err = uc.roles.GiveRole(ctx, in.RoleID, userIdentity.ID, in.ScopeID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) TakeRole(ctx context.Context, in inbounds.ManageRoleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.TakeRole")
	defer span.End()

	isProjectGlobal := in.ScopeID == nil
	span.SetAttributes(attribute.Bool("role.project_global", isProjectGlobal))

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

	var roleBelongs bool
	roleBelongs, err = uc.roles.BelongsToProject(ctx, in.RoleID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !roleBelongs {
		return fail.New(apierr.ROLENotOwnedByPrincipal)
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

	if err = uc.roles.TakeRole(ctx, in.RoleID, userIdentity.ID, in.ScopeID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetUserRoles(ctx context.Context, in inbounds.GetRoleInput) ([]inbounds.RoleOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.GetUserRoles")
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
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot fetch from a project you don't own"})
	}

	var userBelongs bool
	userBelongs, err = uc.projectUsers.BelongsToProject(ctx, in.EntityID, *in.ProjectID)
	if err != nil {
		return nil, err
	}

	if !userBelongs {
		return nil, apierr.FromService(span, inbounds.ErrProjectUserNotFromProject{})
	}

	userIdentity, err := uc.sessions.GetIdentityByEntityIDAndType(ctx, in.EntityID, session.ProjectIdentity)
	if err != nil {
		return nil, err
	}

	foundRoles, err := uc.roles.GetUserRoles(ctx, userIdentity.ID, *in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.RoleSliceToRoleOutputSlice(foundRoles), nil
}
