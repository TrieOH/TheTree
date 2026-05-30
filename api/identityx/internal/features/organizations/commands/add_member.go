package commands

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
)

func (c *Commands) AddMember(ctx context.Context, payload models.AddOrganizationMemberInput) error {
	ctx, span := c.tracer.Start(ctx, "OrganizationService.AddMember")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := c.actors.GetByEmail(ctx, payload.ActorEmail, nil)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return fun.ErrBadRequest("users can't add themselves to organizations")
	}

	org, err := c.orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return fun.ErrBadRequest("owners can't be added to organizations they own")
	}

	if ident.Sub.ID != org.OwnerID {
		member, err := c.orgs.GetMember(ctx, ident.Sub.ID, payload.OrganizationID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return err
		}
		if err != nil {
			return fun.ErrForbidden("insufficient permissions")
		}
		if member.Role != models.OrganizationRoleAdmin {
			return fun.ErrForbidden("insufficient permissions")
		}
	}

	_, err = c.actors.GetByID(ctx, actor.ID)
	if err != nil {
		return err
	}

	_, err = c.orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil {
		return fun.ErrBadRequest("user is already a member of the organization")
	}

	newMember := models.OrganizationMember{
		ActorID:        actor.ID,
		OrganizationID: payload.OrganizationID,
		Role:           payload.Role,
	}

	return c.orgs.AddMember(ctx, newMember)
}
