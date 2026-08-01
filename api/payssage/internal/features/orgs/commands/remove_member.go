package commands

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
)

func (c *Commands) RemoveMember(ctx context.Context, payload models.RemoveOrganizationMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "OrganizationService.RemoveMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := c.idx.Actors.GetByEmail(ctx, payload.ActorEmail)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return fun.ErrBadRequest("Cannot remove yourself from the organization")
	}

	org, err := c.orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return fun.ErrBadRequest("cannot remove the owner of the organization")
	}

	err = authz.Service.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	_, err = c.orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the organization")
	}

	return c.orgs.RemoveMember(ctx, actor.ID, payload.OrganizationID)
}
