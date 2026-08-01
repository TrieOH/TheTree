package wallets

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) ListFromOrg(ctx context.Context, orgID uuid.UUID) ([]models.Wallet, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListFromOrg")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := o.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	err = authz.Service.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.wallets.ListFromOrg(ctx, org.ID)
}
