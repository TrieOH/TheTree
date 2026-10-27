package service

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEvent, error) {
	event, err := o.Events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, event.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return event, nil
}
