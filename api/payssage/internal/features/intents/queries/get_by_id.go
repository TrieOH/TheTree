package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	ctx, span := q.tracer.Start(ctx, "GetByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	intent, err := q.intents.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := q.checkWalletAccess(ctx, intent.WalletID, ident.Sub.ID); err != nil {
		return nil, err
	}

	return intent, nil
}
