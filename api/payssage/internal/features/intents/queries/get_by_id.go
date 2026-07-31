package queries

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	intent, err := q.intents.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = q.checkWalletAccess(ctx, intent.WalletID, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return intent, nil
}
