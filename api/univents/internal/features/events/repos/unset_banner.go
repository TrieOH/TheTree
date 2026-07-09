package repos

import (
	"context"
	"lib/database"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *repo) UnsetBanner(ctx context.Context, id uuid.UUID) (*contracts.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.UnsetBanner")
	defer span.End()

	sqlcEvent, err := database.Queries(ctx, repo.q).UnsetEventBanner(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return mapEventFromDB(&sqlcEvent), nil
}
