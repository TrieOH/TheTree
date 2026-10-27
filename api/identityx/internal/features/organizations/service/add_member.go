package service

import (
	"IdentityX/errs"
	"IdentityX/models"
	"context"
	"errors"
	"lib/database"
)

func (o *Operations) AddMember(ctx context.Context, payload models.AddOrganizationMemberInput) error {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.Actors.GetByEmail(ctx, payload.ActorEmail, nil)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return errs.ErrAddSelfToOrg
	}

	org, err := o.Orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return errs.ErrAddOwnerToOrg
	}

	err = o.Authz.CheckOrg(ctx, payload.OrganizationID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Actors.GetByID(ctx, actor.ID)
	if err != nil {
		return err
	}

	_, err = o.Orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrAlreadyOrgMember
	}

	newMember := models.OrganizationMember{
		ActorID:        actor.ID,
		OrganizationID: payload.OrganizationID,
		Role:           payload.Role,
	}

	return o.Orgs.AddMember(ctx, newMember)
}
