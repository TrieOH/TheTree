package service

import (
	"context"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) ListOwned(ctx context.Context) ([]models.Event, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	events, err := o.Events.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
