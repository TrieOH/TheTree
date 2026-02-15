package queries

import (
	"context"
	"univents/internal/eventcore/domain"
)

func (uc *QueryService) List(ctx context.Context) (out []domain.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.List")
	defer span.End()

	var outEvents []domain.Event
	outEvents, err = uc.events.List(ctx)
	if err != nil {
		return nil, err
	}

	return outEvents, nil
}
