package projects

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/MintzyG/fun"
)

func (o *Operations) AddMember(ctx context.Context, payload models.AddProjectMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "AddMember")
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
		return fun.ErrBadRequest("users can't add themselves to projects")
	}

	project, err := o.projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}

	if actor.ID == project.OwnerID {
		return fun.ErrBadRequest("owners can't be added to projects they own")
	}

	err = o.authz.CheckProject(ctx, ident.Sub.ID, payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.actors.GetByID(ctx, actor.ID)
	if err != nil {
		return err
	}

	_, err = o.projects.GetMember(ctx, actor.ID, project.ID)
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

	return o.projects.AddMember(ctx, *newMember)
}
