package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEvent, error) {
	ctx, span := q.tracer.Start(ctx, "GetEventByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := q.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = q.checkWalletAccess(ctx, event.WalletID, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return event, nil
}
