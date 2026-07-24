package repos

import (
	"context"
	"lib/database"
	"univents/models"
)

func (repo *repo) GetBySlug(ctx context.Context, slug string) (*models.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.GetBySlug")
	defer span.End()
	event, err := database.Queries(ctx, repo.q).GetEventBySlug(ctx, slug)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(event)), nil
}
