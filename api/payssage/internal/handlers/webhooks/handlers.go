// Package webhooks implements the StrictServerInterface methods for the
// webhooks feature. The provider receive endpoint is public; the raw body
// and request are captured by the chi-level middleware (rawRequestCapture)
// because provider signature verification needs the exact bytes.
package webhooks

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/internal/services"
	"payssage/models"

	"github.com/MintzyG/fun"
)

const module = "Payssage"

type Handlers struct {
	ops *services.Webhooks
}

func New(ops *services.Webhooks) *Handlers { return &Handlers{ops: ops} }

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
