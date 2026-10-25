package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"payssage/models"
)

func (o *Operations) Create(ctx context.Context, input models.CreateWebhookEndpointInput) (*models.WebhookEndpoint, error) {
	_, err := o.Authz.Wallet(ctx, input.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	secretBytes := make([]byte, 32)
	_, err = rand.Read(secretBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}
	secret := hex.EncodeToString(secretBytes)

	endpoint := models.WebhookEndpoint{
		WalletID: input.WalletID,
		Name:     input.Name,
		URL:      input.URL,
		Secret:   secret,
	}

	return o.Endpoints.Create(ctx, endpoint)
}
