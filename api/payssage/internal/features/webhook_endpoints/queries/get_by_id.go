package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEndpoint, error) {
	ctx, span := q.tracer.Start(ctx, "GetEndpointByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	endpoint, err := q.endpoints.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = q.checkWalletAccess(ctx, endpoint.WalletID, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return endpoint, nil
}
