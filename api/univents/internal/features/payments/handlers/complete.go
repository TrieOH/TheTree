package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) CompleteEventPayments(ctx context.Context, req openapi.CompleteEventPaymentsRequestObject) (openapi.CompleteEventPaymentsResponseObject, error) {
	event, err := h.ops.Complete(ctx, req.EventId, req.Body.SellerId, req.Body.PublicKey)
	if err != nil {
		return nil, err
	}
	return openapi.CompleteEventPayments200JSONResponse(*event), nil
}
