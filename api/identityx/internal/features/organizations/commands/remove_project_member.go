package commands

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
)

func (c *Commands) RemoveProjectMember(ctx context.Context, payload models.RemoveOrgProjectMemberInput) error {
	ctx, span := c.tracer.Start(ctx, "OrganizationService.RemoveProjectMember")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := c.actors.GetByEmail(ctx, payload.ActorEmail, nil)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return fun.ErrBadRequest("cannot remove yourself from the project")
	}

	org, err := c.orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return fun.ErrBadRequest("cannot remove the owner of the organization from the project")
	}

	if ident.Sub.ID != org.OwnerID {
		member, err := c.orgs.GetMember(ctx, ident.Sub.ID, org.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return err
		}
		if err != nil {
			return fun.ErrForbidden("insufficient permissions")
		}
		if member.Role != models.OrganizationRoleAdmin {
			return fun.ErrForbidden("insufficient permissions")
		}
	}

	_, err = c.orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil {
		return fun.ErrBadRequest("cannot remove organization member from project")
	}

	project, err := c.projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}

	_, err = c.projects.GetMember(ctx, actor.ID, project.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the project")
	}

	return c.projects.RemoveMember(ctx, actor.ID, payload.ProjectID)
}
