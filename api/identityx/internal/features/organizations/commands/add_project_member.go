package commands

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
)

func (c *Commands) AddProjectMember(ctx context.Context, payload models.AddOrgProjectMemberInput) error {
	ctx, span := c.tracer.Start(ctx, "OrganizationService.AddProjectMember")
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
		return fun.ErrBadRequest("users can't add themselves to projects")
	}

	org, err := c.orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return fun.ErrBadRequest("owner of the organization is already a member of the project")
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
		return fun.ErrBadRequest("organization member is already a member of the project")
	}

	project, err := c.projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}

	_, err = c.projects.GetMember(ctx, actor.ID, project.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil {
		return fun.ErrBadRequest("user is already a member of the project")
	}

	newMember, err := models.NewProjectMember(payload.ProjectID, actor.ID, payload.Role)
	if err != nil {
		return err
	}

	return c.projects.AddMember(ctx, *newMember)
}
