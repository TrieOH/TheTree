package payments

import (
	"context"

	"univents/internal/features/payments/service"
	"univents/internal/openapi"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) CompleteEventPayments(ctx context.Context, req openapi.CompleteEventPaymentsRequestObject) (openapi.CompleteEventPaymentsResponseObject, error) {
	event, err := h.ops.Complete(ctx, req.EventId, req.Body.SellerId, req.Body.PublicKey)
	if err != nil {
		return nil, err
	}
	return openapi.CompleteEventPayments200JSONResponse(*event), nil
}

func (h *Handlers) ConnectEventPayments(ctx context.Context, req openapi.ConnectEventPaymentsRequestObject) (openapi.ConnectEventPaymentsResponseObject, error) {
	res, err := h.ops.Connect(ctx, req.EventId, req.Body.Provider)
	if err != nil {
		return nil, err
	}
	return openapi.ConnectEventPayments200JSONResponse(openapi.ConnectEventPaymentsResult{
		AuthUrl:  res.AuthURL,
		WalletId: res.WalletID,
	}), nil
}

func (h *Handlers) DisconnectEventPayments(ctx context.Context, req openapi.DisconnectEventPaymentsRequestObject) (openapi.DisconnectEventPaymentsResponseObject, error) {
	_, err := h.ops.Disconnect(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.DisconnectEventPayments204Response{}, nil
}
