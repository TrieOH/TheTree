package queries

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEvent, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListEventsByWallet")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	err = q.checkWalletAccess(ctx, walletID, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return q.events.ListByWallet(ctx, walletID)
}
