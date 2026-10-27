package handlers

import (
	"context"
	"time"

	"payssage/errs"
	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) ReceiveWebhook(ctx context.Context, req openapi.ReceiveWebhookRequestObject) (openapi.ReceiveWebhookResponseObject, error) {
	raw := RawRequestFrom(ctx)
	if raw == nil {
		return nil, errs.ErrInvalidRequest
	}
	err := h.ops.Receive(ctx, models.ReceiveWebhookInput{
		Provider: req.Provider,
		Request:  raw.Req,
		RawBody:  raw.Body,
	})
	if err != nil {
		return nil, err
	}
	return openapi.ReceiveWebhook200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
