package commands

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (c *Commands) Create(ctx context.Context, payload models.CreateWalletInput) (*models.Wallet, error) {
	ctx, span := c.tracer.Start(ctx, "Create")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var org *models.Organization
	if payload.OrganizationID != nil {
		org, err = c.orgs.GetByID(ctx, *payload.OrganizationID)
		if err != nil {
			return nil, err
		}

		err = c.checkRole(ctx, org, ident.Sub.ID, models.OrganizationRoleAdmin)
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

	return c.wallets.Create(ctx, wallet)
}
