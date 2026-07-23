package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"
)

func (repo *repo) ListPublic(ctx context.Context) ([]models.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.List")
	defer span.End()

	sqlcEvents, err := database.Queries(ctx, repo.q).ListPublicEvents(ctx)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcEvents, mapEvent), nil
}
