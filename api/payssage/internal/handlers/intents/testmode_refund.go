package intents

import (
	"context"
	"os"
	"time"

	"github.com/MintzyG/fun"

	"payssage/internal/openapi"
)

// TestmodeRefundIntent simulates a provider refund: flips a succeeded intent
// to refunded without contacting the provider. TEST_MODE-gated like
// hardCreateIntent.
func (h *Handlers) TestmodeRefundIntent(ctx context.Context, req openapi.TestmodeRefundIntentRequestObject) (openapi.TestmodeRefundIntentResponseObject, error) {
	if os.Getenv("TEST_MODE") != "true" {
		return nil, fun.Err("test mode only").ServiceUnavailable()
	}
	intent, err := h.ops.TestmodeRefund(ctx, req.IntentId)
	if err != nil {
		return nil, err
	}
	return openapi.TestmodeRefundIntent200JSONResponse{
		Code: 200, Data: intent, Timestamp: time.Now(), Module: module,
	}, nil
}
