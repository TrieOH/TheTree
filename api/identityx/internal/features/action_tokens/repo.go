package action_tokens

import (
	"IdentityX/errs"
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
		dbe: database.NewErrorHandler(errs.ErrActionTokenNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

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
	sqlcToken, err := repo.db(ctx).InsertActionToken(ctx, sqlc.InsertActionTokenParams{
		Jti:       token.JTI,
		Purpose:   string(token.Purpose),
		ActorID:   token.ActorID,
		ExpiresAt: token.ExpiresAt,
	})
	return database.One(sqlcToken, err, repo.dbe, mapActionToken)
}

func (repo *Repo) GetByJTI(ctx context.Context, jti uuid.UUID) (*models.ActionToken, error) {
	sqlcToken, err := repo.db(ctx).GetActionTokenByJTI(ctx, jti)
	return database.One(sqlcToken, err, repo.dbe, mapActionToken)
}

func (repo *Repo) Consume(ctx context.Context, jti uuid.UUID) (*models.ActionToken, error) {
	sqlcToken, err := repo.db(ctx).ConsumeActionToken(ctx, jti)
	return database.One(sqlcToken, err, repo.dbe, mapActionToken)
}

func (repo *Repo) DeleteExpired(ctx context.Context) error {
	err := repo.db(ctx).DeleteExpiredActionTokens(ctx)
	return repo.dbe(err)
}
