package service

import (
	"context"
	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) ListFromOrg(ctx context.Context, orgID uuid.UUID) ([]models.Wallet, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListFromOrg")
	defer span.End()

	org, err := o.Orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Wallets.ListFromOrg(ctx, org.ID)
}
