package commands

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"
)

func (c *Commands) Create(ctx context.Context, input models.CreateWebhookEndpointInput) (*models.WebhookEndpoint, error) {
	ctx, span := telemetry.StartSpan(ctx, "CreateWebhookEndpoint")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = c.checkWalletAccess(ctx, input.WalletID, ident.Sub.ID)
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

	return c.endpoints.Create(ctx, endpoint)
}
