package orgs

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
)

func (o *Operations) AddMember(ctx context.Context, payload models.AddOrganizationMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "OrganizationService.AddMember")
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
		return fun.ErrBadRequest("users can't add themselves to organizations")
	}

	org, err := o.orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return fun.ErrBadRequest("owners can't be added to organizations they own")
	}

	err = o.authz.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil {
		return fun.ErrBadRequest("user is already a member of the organization")
	}

	newMember := models.OrganizationMember{
		MemberID:       actor.ID,
		OrganizationID: payload.OrganizationID,
		Role:           payload.Role,
	}

	return o.orgs.AddMember(ctx, newMember)
}
