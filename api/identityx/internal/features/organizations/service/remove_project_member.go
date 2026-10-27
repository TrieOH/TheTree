package service

import (
	"IdentityX/errs"
	"IdentityX/models"
	"context"
	"errors"
	"lib/database"
)

func (o *Operations) RemoveProjectMember(ctx context.Context, payload models.RemoveOrgProjectMemberInput) error {
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

	org, err := o.Orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return errs.ErrRemoveOrgOwnerFromProject
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
		return errs.ErrRemoveOrgMemberFromProject
	}

	project, err := o.Projects.GetByID(ctx, payload.ProjectID)
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
