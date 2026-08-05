package events

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListJoinedEvents(ctx context.Context, _ openapi.ListJoinedEventsRequestObject) (openapi.ListJoinedEventsResponseObject, error) {
	events, err := h.ops.ListJoined(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListJoinedEvents200JSONResponse{
		Code: 200, Data: &events, Timestamp: time.Now(), Module: module,
	}, nil
}
