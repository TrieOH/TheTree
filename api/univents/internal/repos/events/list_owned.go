package events

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListOwned(ctx context.Context, ownerID uuid.UUID) ([]models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.ListOwned")
	defer span.End()
	sqlcEvents, err := database.Queries(ctx, repo.q).ListOwnedEvents(ctx, ownerID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcEvents, mapEvent), nil
}
