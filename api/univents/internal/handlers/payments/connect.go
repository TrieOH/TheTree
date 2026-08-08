package payments

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ConnectEventPayments(ctx context.Context, req openapi.ConnectEventPaymentsRequestObject) (openapi.ConnectEventPaymentsResponseObject, error) {
	res, err := h.ops.Connect(ctx, req.EventId, req.Body.Provider)
	if err != nil {
		return nil, err
	}
	return openapi.ConnectEventPayments200JSONResponse{
		Code: 200,
		Data: &openapi.ConnectEventPaymentsResult{
			AuthUrl:  res.AuthURL,
			WalletId: res.WalletID,
		},
		Timestamp: time.Now(),
		Module:    module,
	}, nil
}
