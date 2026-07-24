package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) GetPast(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.GetPast")
	defer span.End()
	editions, err := database.Queries(ctx, repo.q).GetPastEditions(ctx, eventID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(editions, mapEdition), nil
}
