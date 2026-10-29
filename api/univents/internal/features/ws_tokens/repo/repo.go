package repo

import (
	"context"
	"errors"
	"lib/database"
	"univents/errs"
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.WsTokenRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrWSTokenNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapWsToken(src sqlc.WsToken) models.WsToken {
	return models.WsToken{
		ID:         src.ID,
		PurchaseID: src.PurchaseID,
		UserID:     src.UserID,
		TokenHash:  src.TokenHash,
		ExpiresAt:  src.ExpiresAt,
		UsedAt:     src.UsedAt,
	}
}

// Create inserts a ws_tokens row. Transaction-aware: inside a tx (checkout,
// split 7) the insert joins it via database.Queries — token + purchase
// commit atomically, no separate in-tx variant.
func (repo *Repo) Create(ctx context.Context, toCreate *models.WsToken) (*models.WsToken, error) {
	result, err := repo.db(ctx).CreateWsToken(ctx, sqlc.CreateWsTokenParams{
		PurchaseID: toCreate.PurchaseID,
		UserID:     toCreate.UserID,
		TokenHash:  toCreate.TokenHash,
		ExpiresAt:  toCreate.ExpiresAt,
	})
	return database.One(result, err, repo.dbe, mapWsToken)
}

// Consume is the one-time handshake consume: marks the token used only when
// it exists, is unused, and unexpired. Returns (nil, nil) when the guard
// misses (missing / already used / expired) — the handshake rejects. The
// update is atomic, so two concurrent handshakes with the same token race
// on the same row and exactly one wins.
func (repo *Repo) Consume(ctx context.Context, tokenHash string) (*models.WsToken, error) {
	result, err := repo.db(ctx).ConsumeWsToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//nolint:nilnil // guard missed — missing/already-used/expired; the handshake rejects
			return nil, nil
		}
		return nil, repo.dbe(err)
	}
	return new(mapWsToken(result)), nil
}
