package webhook_events

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) GetWebhookEvent(ctx context.Context, req openapi.GetWebhookEventRequestObject) (openapi.GetWebhookEventResponseObject, error) {
	event, err := h.ops.GetByID(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookEvent200JSONResponse{
		Code: 200, Data: event, Timestamp: time.Now(), Module: module,
	}, nil
}
