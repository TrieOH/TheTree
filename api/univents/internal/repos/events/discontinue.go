package events

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) Discontinue(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.Discontinue")
	defer span.End()
	err := database.Queries(ctx, repo.q).DiscontinueEvent(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
