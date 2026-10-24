package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Certification, error) {
	results, err := database.Queries(ctx, repo.q).ListCertificationsByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapCertification), nil
}
