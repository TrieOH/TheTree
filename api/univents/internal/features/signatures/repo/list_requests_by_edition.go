package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListRequestsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.SignatureRequest, error) {
	err := repo.ExpireStaleRequests(ctx)
	if err != nil {
		return nil, err
	}
	results, err := database.Queries(ctx, repo.q).ListSignatureRequestsByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapSignatureRequest), nil
}
