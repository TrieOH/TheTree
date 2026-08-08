package events

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) SetPaymentsConfig(ctx context.Context, id uuid.UUID, sellerID, walletID *uuid.UUID, publicKey *string) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.SetPaymentsConfig")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).SetEventPaymentsConfig(ctx, sqlc.SetEventPaymentsConfigParams{
		PayssageSellerID:  sellerID,
		PayssageWalletID:  walletID,
		PayssagePublicKey: publicKey,
		ID:                id,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(result)), nil
}

func (repo *Repo) ClearPaymentsConfig(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.ClearPaymentsConfig")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).ClearEventPaymentsConfig(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(result)), nil
}
