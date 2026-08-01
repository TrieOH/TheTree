// Package webhook_events implements the StrictServerInterface methods for
// the webhook_events feature.
package webhook_events

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/internal/services"
)

const module = "Payssage"

type Handlers struct {
	ops *services.WebhookEvents
}

func New(ops *services.WebhookEvents) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListWebhookEvents(ctx context.Context, req openapi.ListWebhookEventsRequestObject) (openapi.ListWebhookEventsResponseObject, error) {
	events, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookEvents200JSONResponse{
		Code: 200, Data: &events, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetWebhookEvent(ctx context.Context, req openapi.GetWebhookEventRequestObject) (openapi.GetWebhookEventResponseObject, error) {
	event, err := h.ops.GetByID(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookEvent200JSONResponse{
		Code: 200, Data: event, Timestamp: time.Now(), Module: module,
	}, nil
}
