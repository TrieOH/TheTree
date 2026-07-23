package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) Discontinue(ctx context.Context, id uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.Discontinue")
	defer span.End()
	err := database.Queries(ctx, repo.q).DiscontinueEvent(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
