package action_tokens

import (
	"context"

	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ActionTokenRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("action_token"),
	}
}

func mapActionToken(src sqlc.ActionToken) models.ActionToken {
	return models.ActionToken{
		JTI:       src.Jti,
		Purpose:   models.ActionTokenPurpose(src.Purpose),
		ActorID:   src.ActorID,
		ExpiresAt: src.ExpiresAt,
		UsedAt:    src.UsedAt,
		CreatedAt: src.CreatedAt,
	}
}

func (repo *Repo) Insert(ctx context.Context, token models.ActionToken) (*models.ActionToken, error) {
	sqlcToken, err := database.Queries(ctx, repo.q).InsertActionToken(ctx, sqlc.InsertActionTokenParams{
		Jti:       token.JTI,
		Purpose:   string(token.Purpose),
		ActorID:   token.ActorID,
		ExpiresAt: token.ExpiresAt,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActionToken(sqlcToken)), nil
}

func (repo *Repo) GetByJTI(ctx context.Context, jti uuid.UUID) (*models.ActionToken, error) {
	sqlcToken, err := database.Queries(ctx, repo.q).GetActionTokenByJTI(ctx, jti)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActionToken(sqlcToken)), nil
}

func (repo *Repo) Consume(ctx context.Context, jti uuid.UUID) (*models.ActionToken, error) {
	sqlcToken, err := database.Queries(ctx, repo.q).ConsumeActionToken(ctx, jti)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActionToken(sqlcToken)), nil
}

func (repo *Repo) DeleteExpired(ctx context.Context) error {
	return database.Queries(ctx, repo.q).DeleteExpiredActionTokens(ctx)
}
