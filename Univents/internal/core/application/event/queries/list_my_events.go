package queries

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
)

func (uc *QueryService) ListOwnEvents(ctx context.Context) (out []domain.Event, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EventService.ListOwnEvents")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var outEvents []domain.Event
	outEvents, err = uc.events.ListOwnEvents(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return outEvents, nil
}
