package payments

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) DisconnectEventPayments(ctx context.Context, req openapi.DisconnectEventPaymentsRequestObject) (openapi.DisconnectEventPaymentsResponseObject, error) {
	_, err := h.ops.Disconnect(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.DisconnectEventPayments204Response{}, nil
}
