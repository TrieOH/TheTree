package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"
)

func (repo *Repo) GetBySlug(ctx context.Context, slug string) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.GetBySlug")
	defer span.End()
	event, err := database.Queries(ctx, repo.q).GetEventBySlug(ctx, slug)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(event)), nil
}
