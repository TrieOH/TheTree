package events

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) DiscontinueEvent(ctx context.Context, req openapi.DiscontinueEventRequestObject) (openapi.DiscontinueEventResponseObject, error) {
	err := h.ops.Discontinue(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.DiscontinueEvent204Response{}, nil
}
