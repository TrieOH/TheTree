package commands

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
)

func (c *Commands) RemoveMember(ctx context.Context, payload models.RemoveProjectMemberInput) error {
	ctx, span := c.tracer.Start(ctx, "ProjectService.RemoveMember")
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
		return fun.ErrBadRequest("Cannot remove yourself from the project")
	}

	project, err := c.projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}
	if actor.ID == project.OwnerID {
		return fun.ErrBadRequest("cannot remove the owner of the project")
	}

	if ident.Sub.ID != project.OwnerID {
		member, err := c.projects.GetMember(ctx, ident.Sub.ID, payload.ProjectID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return err
		}
		if err != nil {
			return fun.ErrForbidden("insufficient permissions")
		}
		if member.Role != models.ProjectRoleAdmin {
			return fun.ErrForbidden("insufficient permissions")
		}
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
