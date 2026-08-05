package orgs

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) GetMemberByEmail(ctx context.Context, email string, orgID uuid.UUID) (*models.OrganizationMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "OrganizationService.GetMemberByEmail")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckOrg(ctx, ident.Sub.ID, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	actor, err := o.idx.Actors.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	member, err := o.orgs.GetMember(ctx, actor.ID, orgID)
	if err != nil {
		return nil, err
	}

	return member, nil
}
