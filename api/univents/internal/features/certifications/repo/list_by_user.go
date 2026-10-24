package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Certification, error) {
	results, err := database.Queries(ctx, repo.q).ListCertificationsByUser(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapCertification), nil
}
