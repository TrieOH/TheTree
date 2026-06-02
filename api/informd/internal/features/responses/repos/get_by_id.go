package repos

import (
	"Informd/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Response, error) {
	ctx, span := repo.tracer.Start(ctx, "ResponseRepo.GetByID")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetResponseByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapResponse(row)), nil
}
