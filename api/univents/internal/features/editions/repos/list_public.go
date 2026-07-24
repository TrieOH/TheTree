package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) ListPublic(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.ListPublic")
	defer span.End()
	editions, err := database.Queries(ctx, repo.q).ListPublicEditions(ctx, eventID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(editions, mapEdition), nil
}
