// Package webhook_events is the webhook events feature: its repo, handlers
// and operations. Build wires it.
package webhook_events

import (
	"context"
	"payssage/internal/authz"
	"payssage/models"
	"payssage/ports"
	"uuid"
)

type Deps struct {
	Events  ports.WebhookEventRepo
	Wallets ports.WalletRepo
	Orgs    ports.OrganizationRepo
	Authz   *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

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

func (o *Operations) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEvent, error) {
	_, err := o.Authz.Wallet(ctx, walletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Events.ListByWallet(ctx, walletID)
}
