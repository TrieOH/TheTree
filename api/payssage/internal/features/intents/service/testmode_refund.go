package service

import (
	"context"
	"encoding/json"
	"fmt"
	"payssage/errs"
	"payssage/models"
	"uuid"
)

// TestmodeRefund simulates a provider refund (only when TestMode is on, and
// only for the wallet's admins): flips a succeeded intent to refunded and
// stamps a fake refund marker in provider_data — no provider is contacted.
// Mirrors the real refund's postconditions so downstream systems observe
// the same state the payment.refunded webhook would produce.
func (o *Operations) TestmodeRefund(ctx context.Context, intentID uuid.UUID) (*models.Intent, error) {
	if !o.TestMode {
		return nil, errs.ErrTestModeOnly
	}
	intent, err := o.Intents.GetByID(ctx, intentID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, intent.WalletID, models.OrganizationRoleAdmin)
	if err != nil {
		return nil, err
	}
	if intent.Status == models.IntentStatusRefunded {
		return intent, nil // idempotent
	}
	if intent.Status != models.IntentStatusSucceeded {
		return nil, &errs.IntentStatusError{Action: "refunded", Status: string(intent.Status)}
	}

	var pd models.MercadoPagoIntentData
	_ = json.Unmarshal(intent.ProviderData, &pd)
	fakeID := "testmode-refund-" + intent.ID.String()[:8]
	pd.RefundID = &fakeID
	status := "approved"
	pd.RefundStatus = &status
	pd.RefundAmountCents = &intent.AmountCents
	raw, err := json.Marshal(&pd)
	if err != nil {
		return nil, fmt.Errorf("marshal provider data: %w", err)
	}
	intent.ProviderData = raw
	intent.Status = models.IntentStatusRefunded

	return o.Intents.Update(ctx, *intent)
}
