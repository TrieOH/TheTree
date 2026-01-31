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
	"errors"
	"time"

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
		return nil, err)
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, inbounds.ErrNotProjectOwner{Msg: "cannot create permissions for a project you don't own"}
	}

	if err := permissions.ValidatePermission(in.Object, in.Action); err != nil {
		return nil, apierr.FromService(span, err)
	}

	_, err = permissions.DecodeAndValidateCondition(in.Conditions)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	permission, err := uc.permissions.Create(ctx, outbounds.CreatePermissionInput{
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

func (uc *UseCase) GetEffective(ctx context.Context, in inbounds.ManagePermissionInput) (perms []inbounds.PermissionOutput, err error) {
	ctx, span := usecaseTracer.Start(ctx, "PermissionService.GetEffective")
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

	if err = permissions.ValidatePermission(in.Object, in.Action); err != nil {
		return false, apierr.FromService(span, err)
	}

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return false, apierr.FromService(span, err)
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
			return false, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot check permissions in a project you don't own"})
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
			return false, apierr.FromService(span, inbounds.ErrProjectUserNotFromProject{})
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

			if p.Conditions != nil {
				if in.Resource == nil {
					// FIXME when conditions exist but no resource provided, error message unclear
					// FIXME user doesn't know which permission requires resource data
					return false, errors.New("can't check conditions without resource")
				}
				evalCtx := permissions.ConditionContext{
					Subject: map[string]interface{}{
						"id": userIdentity.ID,
					},
					Resource: *in.Resource,
					Environment: map[string]interface{}{
						"now": time.Now().UTC(),
					},
				}

				var ok bool
				var motive permissions.Motive
				ok, motive, err = p.Conditions.Evaluate(ctx, &evalCtx)
				if err != nil {
					return false, err
				}
				if !ok {
					span.SetAttributes(attribute.String("denial.motive", motive.Code))
					continue
				}
			}

			// PHASE 7: Conditions ignored for now
			span.SetAttributes(attribute.Bool("permission.found", true))
			return true, nil
		}
	}

	span.SetAttributes(attribute.Bool("permission.found", false))
	return false, apierr.FromService(span, permissions.ErrInsufficientPermissions{})
}
