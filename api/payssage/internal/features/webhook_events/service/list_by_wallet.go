package service

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEvent, error) {
	_, err := o.Authz.Wallet(ctx, walletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Events.ListByWallet(ctx, walletID)
}
