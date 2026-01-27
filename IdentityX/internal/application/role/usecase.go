package role

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/domain/roles"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"

	"go.opentelemetry.io/otel"
)

var (
	usecaseTracer = otel.Tracer("role_usecase")
)

type UseCase struct {
	roles    outbounds.RoleRepository
	projects outbounds.ProjectRepository
	tx       inbounds.TxRunner
}

var _ inbounds.RoleService = (*UseCase)(nil)

func New(
	roles outbounds.RoleRepository,
	projects outbounds.ProjectRepository,
	tx inbounds.TxRunner,
) inbounds.RoleService {
	return &UseCase{
		roles:    roles,
		projects: projects,
		tx:       tx,
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

	foundRoles, err := uc.roles.ListByProject(ctx, *in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.RoleSliceToRoleOutputSlice(foundRoles), nil
}
