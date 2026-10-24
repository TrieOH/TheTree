package service

import (
	"context"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) ListJoined(ctx context.Context) ([]models.Event, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	events, err := o.Events.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
