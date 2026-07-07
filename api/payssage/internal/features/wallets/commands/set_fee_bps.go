package commands

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
)

func (c *Commands) SetFeeBPS(ctx context.Context, payload models.SetFeeBPSInput) error {
	ctx, span := c.tracer.Start(ctx, "SetFeeBPS")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	var org *models.Organization
	if payload.OrganizationID != nil {
		org, err = c.orgs.GetByID(ctx, *payload.OrganizationID)
		if err != nil {
			return err
		}
	}

	if org != nil && org.OwnerID != ident.Sub.ID {
		member, err := c.orgs.GetMember(ctx, ident.Sub.ID, org.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return err
		}
		if err != nil {
			return fun.ErrForbidden("insufficient permissions")
		}
		if member.Role != models.OrganizationRoleAdmin {
			return fun.ErrForbidden("insufficient permissions")
		}
	}

	wallet, err := c.wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return err
	}
	if org != nil && wallet.OrganizationID != &org.ID {
		return fun.ErrForbidden("insufficient permissions")
	}
	if org == nil && wallet.OwnerID != ident.Sub.ID {
		return fun.ErrForbidden("insufficient permissions")
	}

	return c.wallets.SetFeeBPS(ctx, wallet.ID, payload.FeeBps)
}
