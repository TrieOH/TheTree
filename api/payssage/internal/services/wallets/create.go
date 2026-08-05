package wallets

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateWalletInput) (*models.Wallet, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var org *models.Organization
	if payload.OrganizationID != nil {
		org, err = o.orgs.GetByID(ctx, *payload.OrganizationID)
		if err != nil {
			return nil, err
		}
		err = o.authz.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleAdmin)
		if err != nil {
			return nil, err
		}
	}

	ownerID := ident.Sub.ID
	var orgID *uuid.UUID
	if org != nil {
		ownerID = org.OwnerID
		orgID = &org.ID
	}

	wallet := models.Wallet{
		OwnerID:        ownerID,
		OrganizationID: orgID,
		Name:           payload.Name,
		Sandbox:        false,
		FeeBps:         0,
	}

	return o.wallets.Create(ctx, wallet)
}
