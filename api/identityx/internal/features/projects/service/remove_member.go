package service

import (
	"IdentityX/errs"
	"IdentityX/models"
	"context"
	"errors"
	"lib/database"
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
		return errs.ErrRemoveSelfFromProject
	}

	project, err := o.Projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}
	if actor.ID == project.OwnerID {
		return errs.ErrRemoveProjectOwner
	}

	err = o.Authz.CheckProject(ctx, payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Projects.GetMember(ctx, actor.ID, project.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err != nil {
		return errs.ErrNotProjectMember
	}

	return o.Projects.RemoveMember(ctx, actor.ID, payload.ProjectID)
}
