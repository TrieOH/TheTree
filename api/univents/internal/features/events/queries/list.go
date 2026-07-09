package queries

import (
	"context"
	"univents/contracts"
)

func (uc *Queries) List(ctx context.Context) (out []contracts.Event, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EventService.ListEvents")
	defer span.End()

	var outEvents []contracts.Event
	outEvents, err = uc.events.List(ctx)
	if err != nil {
		return nil, err
	}

	return outEvents, nil
}
