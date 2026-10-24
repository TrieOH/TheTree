package service

import (
	"context"
	"payssage/models"
)

func (o *Operations) SetFeeBPS(ctx context.Context, payload models.SetFeeBPSInput) error {
	wallet, err := o.Wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Wallet(ctx, wallet.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	return o.Wallets.SetFeeBPS(ctx, wallet.ID, payload.FeeBps)
}
