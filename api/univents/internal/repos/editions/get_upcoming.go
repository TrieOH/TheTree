package editions

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetUpcoming(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionsRepo.GetUpcoming")
	defer span.End()
	editions, err := database.Queries(ctx, repo.q).GetUpcomingEditions(ctx, eventID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(editions, mapEdition), nil
}
