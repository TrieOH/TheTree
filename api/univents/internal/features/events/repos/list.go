package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
)

func (repo *repo) List(ctx context.Context) ([]contracts.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.List")
	defer span.End()

	sqlcEvents, err := database.Queries(ctx, repo.q).ListEvents(ctx)
	if err != nil {
		return nil, repo.dbe(err)
	}

	outEvents := make([]contracts.Event, 0, len(sqlcEvents))
	for _, sqlcEvent := range sqlcEvents {
		outEvents = append(outEvents, *mapEventFromDB(&sqlcEvent))
	}
	return outEvents, nil
}
