package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]models.WebhookDelivery, error) {
	ctx, span := q.tracer.Start(ctx, "ListDeliveriesByEndpoint")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	endpoint, err := q.endpoints.GetByID(ctx, endpointID)
	if err != nil {
		return nil, err
	}

	err = q.checkWalletAccess(ctx, endpoint.WalletID, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return q.deliveries.ListByEndpoint(ctx, endpointID)
}
