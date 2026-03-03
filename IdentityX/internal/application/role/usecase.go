package role

import (
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/roles"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"strings"

	"github.com/MintzyG/fail/v3"
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
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot create roles for a project you don't own").RecordCtx(ctx)
	}

	role, err := uc.roles.Create(ctx, roles.Role{
		ProjectID:   in.ProjectID,
		Name:        in.Name,
		Description: in.Description,
		Meta:        in.Meta,
	})
	if err != nil {
		return nil, err
	}

	return inbounds.RoleToRoleOutput(*role), nil
}

func (uc *UseCase) UpdateDescription(ctx context.Context, in inbounds.RoleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.UpdateDescription")
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
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot update roles in a project you don't own").RecordCtx(ctx)
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

func (uc *UseCase) UpdateMeta(ctx context.Context, in inbounds.RoleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.UpdateMeta")
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
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot update roles in a project you don't own").RecordCtx(ctx)
	}

	err = uc.roles.UpdateMeta(ctx, in.Meta, in.RoleID, in.ProjectID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) Delete(ctx context.Context, in inbounds.RoleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.Delete")
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
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot delete roles in a project you don't own").RecordCtx(ctx)
	}

	err = uc.roles.Delete(ctx, in.RoleID, in.ProjectID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetByIDExternal(ctx context.Context, in inbounds.GetRoleInput) (*inbounds.RoleOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.GetByIDExternal")
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
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get roles from a project you don't own").RecordCtx(ctx)
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
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get roles from a project you don't own").RecordCtx(ctx)
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
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get roles from a project you don't own").RecordCtx(ctx)
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

	var roleBelongs bool
	roleBelongs, err = uc.roles.BelongsToProject(ctx, in.RoleID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !roleBelongs {
		return fail.New(errx.ROLENotOwnedByPrincipal).RecordCtx(ctx)
	}

	var permissionBelongs bool
	permissionBelongs, err = uc.permissions.BelongsToProject(ctx, in.PermissionID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !permissionBelongs {
		return fail.New(errx.PERMissionNotOwnedByPrincipal).Trace("cannot add permission to a role you don't own").RecordCtx(ctx)
	}

	if err = uc.roles.AddPermission(ctx, in.RoleID, in.PermissionID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) RemovePermission(ctx context.Context, in inbounds.RolePermissionInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.RemovePermission")
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
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot edit a project you don't own").RecordCtx(ctx)
	}

	var roleBelongs bool
	roleBelongs, err = uc.roles.BelongsToProject(ctx, in.RoleID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !roleBelongs {
		return fail.New(errx.ROLENotOwnedByPrincipal).RecordCtx(ctx)
	}

	var permissionBelongs bool
	permissionBelongs, err = uc.permissions.BelongsToProject(ctx, in.PermissionID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !permissionBelongs {
		return fail.New(errx.PERMissionNotOwnedByPrincipal).Trace("cannot remove permission to a role you don't own").RecordCtx(ctx)
	}

	if err = uc.roles.RemovePermission(ctx, in.RoleID, in.PermissionID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetPermissions(ctx context.Context, in inbounds.RolePermissionInput) ([]inbounds.PermissionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.GetPermissions")
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
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot fetch from a project you don't own").RecordCtx(ctx)
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

	var roleBelongs bool
	roleBelongs, err = uc.roles.BelongsToProject(ctx, in.RoleID, *in.ProjectID)
	if err != nil {
		return err
	}

	var userRole *roles.Role
	userRole, err = uc.roles.GetByIDInternal(ctx, in.RoleID)
	if err != nil {
		return err
	}

	if !roleBelongs {
		return fail.New(errx.ROLENotOwnedByPrincipal).RecordCtx(ctx)
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

	if err = uc.roles.GiveRole(ctx, in.RoleID, userIdentity.ID, in.ScopeID, userRole.Name); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) TakeRole(ctx context.Context, in inbounds.ManageRoleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.TakeRole")
	defer span.End()

	isProjectGlobal := in.ScopeID == nil
	span.SetAttributes(attribute.Bool("role.project_global", isProjectGlobal))

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

	var roleBelongs bool
	roleBelongs, err = uc.roles.BelongsToProject(ctx, in.RoleID, *in.ProjectID)
	if err != nil {
		return err
	}

	if !roleBelongs {
		return fail.New(errx.ROLENotOwnedByPrincipal).RecordCtx(ctx)
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

	if err = uc.roles.TakeRole(ctx, in.RoleID, userIdentity.ID, in.ScopeID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GiveRoleByName(ctx context.Context, in inbounds.ManageRoleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.GiveRoleByName")
	defer span.End()

	if in.RoleName == "" {
		return &inbounds.InvalidRoleNameError{Name: in.RoleName}
	}

	span.SetAttributes(attribute.String("role.name", in.RoleName))

	role, err := uc.roles.GetByName(ctx, in.RoleName, in.ProjectID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") || strings.Contains(err.Error(), "not found") {
			return &inbounds.RoleNotFoundByNameError{Name: in.RoleName, ProjectID: in.ProjectID}
		}
		return err
	}

	in.RoleID = role.ID

	return uc.GiveRole(ctx, in)
}

func (uc *UseCase) TakeRoleByName(ctx context.Context, in inbounds.ManageRoleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.TakeRoleByName")
	defer span.End()

	if in.RoleName == "" {
		return &inbounds.InvalidRoleNameError{Name: in.RoleName}
	}

	span.SetAttributes(attribute.String("role.name", in.RoleName))

	role, err := uc.roles.GetByName(ctx, in.RoleName, in.ProjectID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") || strings.Contains(err.Error(), "not found") {
			return &inbounds.RoleNotFoundByNameError{Name: in.RoleName, ProjectID: in.ProjectID}
		}
		return err
	}

	in.RoleID = role.ID

	return uc.TakeRole(ctx, in)
}

func (uc *UseCase) GetUserRoles(ctx context.Context, in inbounds.GetRoleInput) ([]inbounds.RoleOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "RoleService.GetUserRoles")
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
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot fetch from a project you don't own").RecordCtx(ctx)
	}

	var userBelongs bool
	userBelongs, err = uc.projectUsers.BelongsToProject(ctx, in.EntityID, *in.ProjectID)
	if err != nil {
		return nil, err
	}

	if !userBelongs {
		return nil, fail.New(errx.ProjectUserNotFromProject).RecordCtx(ctx)
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
