package events

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.ListJoined")
	defer span.End()
	sqlcEvents, err := database.Queries(ctx, repo.q).ListJoinedEvents(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcEvents, mapEvent), nil
}
