package repos

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.GetByID")
	defer span.End()
	event, err := database.Queries(ctx, repo.q).GetEventByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(event)), nil
}
