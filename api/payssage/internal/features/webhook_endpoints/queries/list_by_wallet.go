package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEndpoint, error) {
	ctx, span := q.tracer.Start(ctx, "ListEndpointsByWallet")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	if err := q.checkWalletAccess(ctx, walletID, ident.Sub.ID); err != nil {
		return nil, err
	}

	return q.endpoints.ListByWallet(ctx, walletID)
}
