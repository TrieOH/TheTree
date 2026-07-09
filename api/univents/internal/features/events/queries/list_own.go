package queries

import (
	"context"
	idx "sdk/identityx"
	"univents/contracts"
)

func (uc *Queries) ListOwn(ctx context.Context) (out []contracts.Event, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EventService.ListOwnEvents")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var outEvents []contracts.Event
	outEvents, err = uc.events.ListOwn(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return outEvents, nil
}
