package repo

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetRequestByIdempotencyKey(ctx context.Context, idempotencyKey uuid.UUID) (*models.SignatureRequest, error) {
	err := repo.ExpireStaleRequests(ctx)
	if err != nil {
		return nil, err
	}
	result, err := database.Queries(ctx, repo.q).GetSignatureRequestByIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignatureRequest(result)), nil
}
