package webhooks

import (
	"context"
	"time"

	"github.com/MintzyG/fun"

	wh "univents/internal/services/webhooks"

	"univents/internal/openapi"
)

// ReceivePayssageWebhook is the payment-confirmation entry point (D3): the
// only component that approves a purchase. The raw body (for signature
// verification) comes from RawRequestMiddleware; the decoded envelope from
// the strict server. Any error returned becomes a non-2xx response —
// Payssage retries; nil stops retries.
func (h *Handlers) ReceivePayssageWebhook(ctx context.Context, req openapi.ReceivePayssageWebhookRequestObject) (openapi.ReceivePayssageWebhookResponseObject, error) {
	raw := RawRequestFrom(ctx)
	if raw == nil || req.Body == nil {
		return nil, fun.ErrBadRequest("invalid request")
	}
	err := h.ops.Receive(ctx, wh.ReceiveInput{
		IntentID:  req.Body.IntentId,
		EventType: req.Body.EventType,
		RawBody:   raw.Body,
		Signature: raw.Req.Header.Get("X-Payssage-Signature"),
	})
	if err != nil {
		return nil, err
	}
	return openapi.ReceivePayssageWebhook200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
