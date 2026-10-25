package service

import (
	"IdentityX/errs"
	"IdentityX/models"
	"context"
	"errors"
	"lib/database"
)

func (o *Operations) AddMember(ctx context.Context, payload models.AddProjectMemberInput) error {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.Actors.GetByEmail(ctx, payload.ActorEmail, nil)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return errs.ErrAddSelfToProject
	}

	project, err := o.Projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}

	if actor.ID == project.OwnerID {
		return errs.ErrAddOwnerToProject
	}

	err = o.Authz.CheckProject(ctx, payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Actors.GetByID(ctx, actor.ID)
	if err != nil {
		return err
	}

	_, err = o.Projects.GetMember(ctx, actor.ID, project.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrAlreadyProjectMember
	}

	newMember, err := models.NewProjectMember(payload.ProjectID, actor.ID, payload.Role)
	if err != nil {
		return err
	}

	return o.Projects.AddMember(ctx, *newMember)
}
