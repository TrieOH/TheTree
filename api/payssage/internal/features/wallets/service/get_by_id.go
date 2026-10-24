package service

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	wallet, err := o.Wallets.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Wallet(ctx, wallet.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
