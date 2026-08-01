package events

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.GetByID")
	defer span.End()
	event, err := database.Queries(ctx, repo.q).GetEventByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(event)), nil
}
