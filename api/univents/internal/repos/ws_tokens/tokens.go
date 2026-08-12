package ws_tokens

import (
	"context"
	"errors"

	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/jackc/pgx/v5"
)

// Create inserts a ws_tokens row. Transaction-aware: inside a tx (checkout,
// split 7) the insert joins it via database.Queries — token + purchase
// commit atomically, no separate in-tx variant.
func (repo *Repo) Create(ctx context.Context, toCreate *models.WsToken) (*models.WsToken, error) {
	ctx, span := telemetry.StartSpan(ctx, "WsTokensRepo.Create")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).CreateWsToken(ctx, sqlc.CreateWsTokenParams{
		PurchaseID: toCreate.PurchaseID,
		UserID:     toCreate.UserID,
		TokenHash:  toCreate.TokenHash,
		ExpiresAt:  toCreate.ExpiresAt,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapWsToken(result)), nil
}

// Consume is the one-time handshake consume: marks the token used only when
// it exists, is unused, and unexpired. Returns (nil, nil) when the guard
// misses (missing / already used / expired) — the handshake rejects. The
// update is atomic, so two concurrent handshakes with the same token race
// on the same row and exactly one wins.
func (repo *Repo) Consume(ctx context.Context, tokenHash string) (*models.WsToken, error) {
	ctx, span := telemetry.StartSpan(ctx, "WsTokensRepo.Consume")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).ConsumeWsToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//nolint:nilnil // guard missed — missing/already-used/expired; the handshake rejects
			return nil, nil
		}
		return nil, repo.dbe(err)
	}
	return new(mapWsToken(result)), nil
}
