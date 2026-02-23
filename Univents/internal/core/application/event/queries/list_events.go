package queries

import (
	"context"
	"univents/internal/core/domain"
)

func (uc *QueryService) ListEvents(ctx context.Context) (out []domain.Event, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EventService.ListEvents")
	defer span.End()

	var outEvents []domain.Event
	outEvents, err = uc.events.ListEvents(ctx)
	if err != nil {
		return nil, err
	}

	return outEvents, nil
}
