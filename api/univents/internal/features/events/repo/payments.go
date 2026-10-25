package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) SetPaymentsConfig(ctx context.Context, id uuid.UUID, sellerID *uuid.UUID, publicKey *string) (*models.Event, error) {
	result, err := database.Queries(ctx, repo.q).SetEventPaymentsConfig(ctx, sqlc.SetEventPaymentsConfigParams{
		PayssageSellerID:  sellerID,
		PayssagePublicKey: publicKey,
		ID:                id,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(result)), nil
}

func (repo *Repo) ClearPaymentsConfig(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	result, err := database.Queries(ctx, repo.q).ClearEventPaymentsConfig(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(result)), nil
}
