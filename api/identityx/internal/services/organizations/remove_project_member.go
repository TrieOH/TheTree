package organizations

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/MintzyG/fun"
)

func (o *Operations) RemoveProjectMember(ctx context.Context, payload models.RemoveOrgProjectMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveProjectMember")
	defer span.End()
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.actors.GetByEmail(ctx, payload.ActorEmail, nil)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return fun.ErrBadRequest("cannot remove yourself from the project")
	}

	org, err := o.orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return fun.ErrBadRequest("cannot remove the owner of the organization from the project")
	}

	err = o.authz.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil {
		return fun.ErrBadRequest("cannot remove organization member from project")
	}

	project, err := o.projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}

	_, err = o.projects.GetMember(ctx, actor.ID, project.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the project")
	}

	return o.projects.RemoveMember(ctx, actor.ID, payload.ProjectID)
}
