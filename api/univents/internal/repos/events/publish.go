package events

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) Publish(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.Publish")
	defer span.End()
	err := database.Queries(ctx, repo.q).PublishEvent(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
