package service

import (
	"context"
	"univents/models"
)

func (o *Operations) ListPublic(ctx context.Context) ([]models.Event, error) {
	events, err := o.Events.ListPublic(ctx)
	if err != nil {
		return nil, err
	}

	return events, nil
}
