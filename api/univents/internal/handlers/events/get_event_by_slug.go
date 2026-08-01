package events

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetEventBySlug(ctx context.Context, req openapi.GetEventBySlugRequestObject) (openapi.GetEventBySlugResponseObject, error) {
	event, err := h.ops.GetBySlug(ctx, req.EventSlug)
	if err != nil {
		return nil, err
	}
	return openapi.GetEventBySlug200JSONResponse{
		Code: 200, Data: event, Timestamp: time.Now(), Module: module,
	}, nil
}
