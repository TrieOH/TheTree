package handlers

import (
	"context"

	"payssage/internal/features/webhook_events/service"
	"payssage/internal/openapi"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) GetWebhookEvent(ctx context.Context, req openapi.GetWebhookEventRequestObject) (openapi.GetWebhookEventResponseObject, error) {
	event, err := h.ops.GetByID(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookEvent200JSONResponse(*event), nil
}

func (h *Handlers) ListWebhookEvents(ctx context.Context, req openapi.ListWebhookEventsRequestObject) (openapi.ListWebhookEventsResponseObject, error) {
	events, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookEvents200JSONResponse{Items: events}, nil
}
