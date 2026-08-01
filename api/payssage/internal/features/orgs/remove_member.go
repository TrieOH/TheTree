package orgs

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
)

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveOrganizationMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "OrganizationService.RemoveMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.idx.Actors.GetByEmail(ctx, payload.ActorEmail)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return fun.ErrBadRequest("Cannot remove yourself from the organization")
	}

	org, err := o.orgs.GetByID(ctx, payload.OrganizationID)
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

	_, err = o.orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the organization")
	}

	return o.orgs.RemoveMember(ctx, actor.ID, payload.OrganizationID)
}
