package service

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
)

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveProjectMemberInput) error {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.Actors.GetByEmail(ctx, payload.ActorEmail, nil)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return fun.ErrBadRequest("Cannot remove yourself from the project")
	}

	project, err := o.Projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}
	if actor.ID == project.OwnerID {
		return fun.ErrBadRequest("cannot remove the owner of the project")
	}

	err = o.Authz.CheckProject(ctx, payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Projects.GetMember(ctx, actor.ID, project.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the project")
	}

	return o.Projects.RemoveMember(ctx, actor.ID, payload.ProjectID)
}
