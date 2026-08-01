package events

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListPublicEvents(ctx context.Context, _ openapi.ListPublicEventsRequestObject) (openapi.ListPublicEventsResponseObject, error) {
	events, err := h.ops.ListPublic(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListPublicEvents200JSONResponse{
		Code: 200, Data: &events, Timestamp: time.Now(), Module: module,
	}, nil
}
