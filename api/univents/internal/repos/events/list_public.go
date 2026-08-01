package events

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"
)

func (repo *Repo) ListPublic(ctx context.Context) ([]models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.List")
	defer span.End()

	sqlcEvents, err := database.Queries(ctx, repo.q).ListPublicEvents(ctx)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcEvents, mapEvent), nil
}
