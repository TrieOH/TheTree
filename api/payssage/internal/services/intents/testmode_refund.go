package intents

import (
	"context"
	"encoding/json"
	"lib/utils"
	"payssage/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// TestmodeRefund simulates a provider refund (TEST_MODE-gated in the
// handler): flips a succeeded intent to refunded and stamps a fake refund
// marker in provider_data — no provider is contacted. Mirrors the real
// refund's postconditions so downstream systems observe the same state the
// payment.refunded webhook would produce.
func (o *Operations) TestmodeRefund(ctx context.Context, intentID uuid.UUID) (*models.Intent, error) {
	intent, err := o.intents.GetByID(ctx, intentID)
	if err != nil {
		return nil, err
	}
	if intent.Status == models.IntentStatusRefunded {
		return intent, nil // idempotent
	}
	if intent.Status != models.IntentStatusSucceeded {
		return nil, fun.Errf("intent cannot be refunded from status %q", intent.Status).BadRequest()
	}

	var pd models.MercadoPagoIntentData
	_ = utils.MapTo(&pd, intent.ProviderData)
	fakeID := "testmode-refund-" + intent.ID.String()[:8]
	pd.RefundID = &fakeID
	status := "approved"
	pd.RefundStatus = &status
	pd.RefundAmountCents = &intent.AmountCents
	raw, err := json.Marshal(&pd)
	if err != nil {
		return nil, fun.Errf("marshal provider data: %v", err).Internal()
	}
	intent.ProviderData = raw
	intent.Status = models.IntentStatusRefunded

	return o.intents.Update(ctx, *intent)
}
