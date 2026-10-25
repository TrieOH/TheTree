package repo

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetRequestByID(ctx context.Context, id uuid.UUID) (*models.SignatureRequest, error) {
	err := repo.ExpireStaleRequests(ctx)
	if err != nil {
		return nil, err
	}
	result, err := database.Queries(ctx, repo.q).GetSignatureRequestByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignatureRequest(result)), nil
}
