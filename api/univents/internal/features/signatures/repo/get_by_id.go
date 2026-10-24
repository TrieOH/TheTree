package repo

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Signature, error) {
	result, err := database.Queries(ctx, repo.q).GetSignatureByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignature(result)), nil
}
