// Package webhook_deliveries is the webhook deliveries feature: its repo,
// handlers and operations. Build wires it.
package webhook_deliveries

import (
	"context"
	"payssage/internal/authz"
	"payssage/models"
	"payssage/ports"
	"uuid"
)

type Deps struct {
	Deliveries ports.WebhookDeliveryRepo
	Endpoints  ports.WebhookEndpointRepo
	Wallets    ports.WalletRepo
	Orgs       ports.OrganizationRepo
	Authz      *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	delivery, err := o.Deliveries.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	endpoint, err := o.Endpoints.GetByID(ctx, delivery.EndpointID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (o *Operations) ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]models.WebhookDelivery, error) {
	endpoint, err := o.Endpoints.GetByID(ctx, endpointID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Deliveries.ListByEndpoint(ctx, endpointID)
}
