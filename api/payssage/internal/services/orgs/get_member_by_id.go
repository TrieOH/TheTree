package orgs

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) GetMemberByID(ctx context.Context, id, orgID uuid.UUID) (*models.OrganizationMember, error) {
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

	member, err := o.orgs.GetMember(ctx, actor.ID, orgID)
	if err != nil {
		return nil, err
	}

	return member, nil
}
