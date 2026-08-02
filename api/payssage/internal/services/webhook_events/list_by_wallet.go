package webhook_events

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEvent, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListEventsByWallet")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	err = o.authz.CheckWalletAccess(ctx, ident.Sub.ID, walletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.events.ListByWallet(ctx, walletID)
}
