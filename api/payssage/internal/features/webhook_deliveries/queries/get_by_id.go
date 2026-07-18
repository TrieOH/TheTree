package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	ctx, span := q.tracer.Start(ctx, "GetDeliveryByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	delivery, err := q.deliveries.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	endpoint, err := q.endpoints.GetByID(ctx, delivery.EndpointID)
	if err != nil {
		return nil, err
	}

	if err := q.checkWalletAccess(ctx, endpoint.WalletID, ident.Sub.ID); err != nil {
		return nil, err
	}

	return delivery, nil
}
