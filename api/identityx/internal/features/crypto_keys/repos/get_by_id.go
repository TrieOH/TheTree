package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.CryptoKey, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByID")
	defer span.End()
	sqlcKeyPair, err := database.Queries(ctx, repo.q).GetCryptoKeyByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapKeys(sqlcKeyPair)), nil
}
