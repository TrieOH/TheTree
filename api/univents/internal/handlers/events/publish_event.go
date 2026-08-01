package events

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) PublishEvent(ctx context.Context, req openapi.PublishEventRequestObject) (openapi.PublishEventResponseObject, error) {
	err := h.ops.Publish(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.PublishEvent204Response{}, nil
}
