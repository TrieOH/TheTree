package service

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	intent, err := o.Intents.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, intent.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return intent, nil
}
