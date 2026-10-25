package service

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEndpoint, error) {
	endpoint, err := o.Endpoints.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return endpoint, nil
}
