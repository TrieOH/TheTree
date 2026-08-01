package orgs

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) GetMemberByID(ctx context.Context, id, orgID uuid.UUID) (*idx.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "OrganizationService.GetMemberByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckOrg(ctx, ident.Sub.ID, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	actor, err := o.idx.Actors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = o.orgs.GetMember(ctx, actor.ID, orgID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}

	return actor, nil
}
