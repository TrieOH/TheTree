// Package webhook_endpoints is the webhook endpoints feature: its repo,
// handlers and operations. Build wires it.
package webhook_endpoints

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"payssage/internal/authz"
	"payssage/models"
	"payssage/ports"
	"uuid"
)

type Deps struct {
	Endpoints ports.WebhookEndpointRepo
	Wallets   ports.WalletRepo
	Orgs      ports.OrganizationRepo
	Authz     *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

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

func (o *Operations) Delete(ctx context.Context, id uuid.UUID) error {
	endpoint, err := o.Endpoints.GetByID(ctx, id)
	if err != nil {
		return err
	}
	_, err = o.Authz.Wallet(ctx, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return err
	}

	return o.Endpoints.Delete(ctx, id)
}

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

func (o *Operations) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEndpoint, error) {
	_, err := o.Authz.Wallet(ctx, walletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Endpoints.ListByWallet(ctx, walletID)
}
