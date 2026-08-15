package organizations

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/MintzyG/fun"
)

func (o *Operations) AddProjectMember(ctx context.Context, payload models.AddOrgProjectMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "AddProjectMember")
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

	org, err := o.orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return fun.ErrBadRequest("owner of the organization is already a member of the project")
	}

	err = o.authz.CheckOrg(ctx, org.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil {
		return fun.ErrBadRequest("organization member is already a member of the project")
	}

	project, err := o.projects.GetByID(ctx, payload.ProjectID)
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
