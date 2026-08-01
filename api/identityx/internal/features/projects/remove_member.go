package projects

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/MintzyG/fun"
)

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveProjectMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveMember")
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
		return fun.ErrBadRequest("Cannot remove yourself from the project")
	}

	project, err := o.projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}
	if actor.ID == project.OwnerID {
		return fun.ErrBadRequest("cannot remove the owner of the project")
	}

	err = o.authz.CheckProject(ctx, ident.Sub.ID, payload.ProjectID, nil, models.ProjectRoleAdmin)
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
