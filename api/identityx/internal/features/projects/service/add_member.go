package service

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
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
		return fun.ErrBadRequest("users can't add themselves to projects")
	}

	project, err := o.Projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}

	if actor.ID == project.OwnerID {
		return fun.ErrBadRequest("owners can't be added to projects they own")
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

	return o.Projects.AddMember(ctx, *newMember)
}
