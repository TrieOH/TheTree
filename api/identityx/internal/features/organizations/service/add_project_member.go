package service

import (
	"IdentityX/errs"
	"IdentityX/models"
	"context"
	"errors"
	"lib/database"
)

func (o *Operations) AddProjectMember(ctx context.Context, payload models.AddOrgProjectMemberInput) error {
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

	org, err := o.Orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return errs.ErrOrgOwnerAlreadyInProject
	}

	err = o.Authz.CheckOrg(ctx, org.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrOrgMemberAlreadyInProject
	}

	project, err := o.Projects.GetByID(ctx, payload.ProjectID)
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
