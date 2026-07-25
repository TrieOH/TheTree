package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetPast(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionsRepo.GetPast")
	defer span.End()
	editions, err := database.Queries(ctx, repo.q).GetPastEditions(ctx, eventID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(editions, mapEdition), nil
}
