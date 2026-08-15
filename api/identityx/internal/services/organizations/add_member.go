package organizations

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/MintzyG/fun"
)

func (o *Operations) AddMember(ctx context.Context, payload models.AddOrganizationMemberInput) error {
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
		return fun.ErrBadRequest("users can't add themselves to organizations")
	}

	org, err := o.orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return fun.ErrBadRequest("owners can't be added to organizations they own")
	}

	err = o.authz.CheckOrg(ctx, payload.OrganizationID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.actors.GetByID(ctx, actor.ID)
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
		ActorID:        actor.ID,
		OrganizationID: payload.OrganizationID,
		Role:           payload.Role,
	}

	return o.orgs.AddMember(ctx, newMember)
}
