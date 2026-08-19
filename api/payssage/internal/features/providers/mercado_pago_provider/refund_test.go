package mercado_pago_provider

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"payssage/internal/config"
	"payssage/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"resty.dev/v3"
)

func refundTestIntent(walletID uuid.UUID) *models.Intent {
	intent := testIntent(walletID, uuid.MustParse("22222222-2222-2222-2222-222222222222"))
	raw, _ := json.Marshal(&models.MercadoPagoIntentData{TransactionID: "174515095132"})
	intent.ProviderData = raw
	intent.Status = models.IntentStatusSucceeded
	return intent
}

// TestRefund_SendsFullRefund pins the refund op: it POSTs to the MP refunds
// endpoint with the seller token + idempotency key and an empty body (full
// refund — no amount in v1), records the refund object in provider_data, and
// leaves the intent status untouched (the payment.refunded webhook confirms).
func TestRefund_SendsFullRefund(t *testing.T) {
	walletID := uuid.New()
	stub := &stubRoundTripper{
		body: `{"id": 987654321, "status": "approved", "status_detail": "refunded", "amount": 1000.0}`,
		req:  make(chan map[string]any, 1),
		url:  make(chan string, 1),
	}
	p := &Provider{
		cfg:        config.MercadoPagoConfig{MpAccessToken: "TEST_ACCESS_TOKEN"},
		intents:    fakeIntents{},
		collectors: fakeCollectors{},
		sellers:    &fakeSellers{seller: testSeller(walletID)},
		wallets:    &fakeWallets{wallet: &models.Wallet{ID: walletID, FeeBps: 500}},
		httpClient: resty.New().SetTransport(stub),
	}

	intent := refundTestIntent(walletID)
	err := p.Refund(context.Background(), intent)
	if err != nil {
		t.Fatalf("Refund: %v", err)
	}

	if got := <-stub.url; got != "/v1/payments/174515095132/refunds" {
		t.Errorf("refund URL = %q, want /v1/payments/174515095132/refunds", got)
	}
	body := <-stub.req
	if len(body) != 0 {
		t.Errorf("refund body = %v, want empty (full refund)", body)
	}

	if intent.Status != models.IntentStatusSucceeded {
		t.Errorf("intent status = %q, want succeeded (webhook confirms the flip)", intent.Status)
	}
	var got models.MercadoPagoIntentData
	err = json.Unmarshal(intent.ProviderData, &got)
	if err != nil {
		t.Fatalf("unmarshal provider data: %v", err)
	}
	if got.RefundID == nil || *got.RefundID != "987654321" {
		t.Errorf("refund_id = %v, want 987654321", got.RefundID)
	}
	if got.RefundStatus == nil || *got.RefundStatus != "approved" {
		t.Errorf("refund_status = %v, want approved", got.RefundStatus)
	}
	if got.RefundAmountCents == nil || *got.RefundAmountCents != 100000 {
		t.Errorf("refund_amount_cents = %v, want 100000", got.RefundAmountCents)
	}
}

func TestRefund_NoTransactionIDConflicts(t *testing.T) {
	walletID := uuid.New()
	p := &Provider{
		cfg:        config.MercadoPagoConfig{MpAccessToken: "TEST_ACCESS_TOKEN"},
		intents:    fakeIntents{},
		collectors: fakeCollectors{},
		sellers:    &fakeSellers{seller: testSeller(walletID)},
		wallets:    &fakeWallets{wallet: &models.Wallet{ID: walletID, FeeBps: 500}},
		httpClient: resty.New().SetTransport(&stubRoundTripper{body: "{}"}),
	}

	intent := testIntent(walletID, uuid.MustParse("22222222-2222-2222-2222-222222222222"))
	// Valid provider_data JSON but no transaction id — the real pre-checkout
	// shape for an intent that never reached the provider.
	intent.ProviderData = json.RawMessage(`{}`)
	err := p.Refund(context.Background(), intent)
	if !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("expected conflict for missing transaction id, got %v", err)
	}
}

func TestRefund_RevokedSellerConflicts(t *testing.T) {
	walletID := uuid.New()
	seller := testSeller(walletID)
	revokedAt := time.Now()
	seller.RevokedAt = &revokedAt
	p := &Provider{
		cfg:        config.MercadoPagoConfig{MpAccessToken: "TEST_ACCESS_TOKEN"},
		intents:    fakeIntents{},
		collectors: fakeCollectors{},
		sellers:    &fakeSellers{seller: seller},
		wallets:    &fakeWallets{wallet: &models.Wallet{ID: walletID, FeeBps: 500}},
		httpClient: resty.New().SetTransport(&stubRoundTripper{body: "{}"}),
	}

	err := p.Refund(context.Background(), refundTestIntent(walletID))
	if !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("expected conflict for revoked seller, got %v", err)
	}
}
