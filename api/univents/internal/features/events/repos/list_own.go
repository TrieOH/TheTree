package repos

import (
	"context"
	"lib/database"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *repo) ListOwn(ctx context.Context, ownerID uuid.UUID) ([]contracts.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.ListOwn")
	defer span.End()

	sqlcEvents, err := database.Queries(ctx, repo.q).ListOwnEvents(ctx, &ownerID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	outEvents := make([]contracts.Event, 0, len(sqlcEvents))
	for _, sqlcEvent := range sqlcEvents {
		outEvents = append(outEvents, *mapEventFromDB(&sqlcEvent))
	}
	return outEvents, nil
}
