package webhooks

import (
	"context"
	"time"

	"github.com/MintzyG/fun"

	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) ReceiveWebhook(ctx context.Context, req openapi.ReceiveWebhookRequestObject) (openapi.ReceiveWebhookResponseObject, error) {
	raw := RawRequestFrom(ctx)
	if raw == nil {
		return nil, fun.ErrBadRequest("invalid request")
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
