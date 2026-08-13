package jobs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"payssage/models"

	"github.com/google/uuid"
)

func TestBuildEnvelope(t *testing.T) {
	intentID := uuid.New()
	walletID := uuid.New()
	event := &models.WebhookEvent{
		ID:         uuid.New(),
		WalletID:   walletID,
		IntentID:   intentID,
		Provider:   "mercado_pago",
		ExternalID: "1234567890",
		EventType:  "payment.succeeded",
		Payload:    json.RawMessage(`{"id":1234567890,"status":"approved"}`),
	}

	envelope, err := buildEnvelope(event)
	if err != nil {
		t.Fatalf("buildEnvelope: %v", err)
	}

	var decoded map[string]any
	err = json.Unmarshal(envelope, &decoded)
	if err != nil {
		t.Fatalf("envelope not valid json: %v", err)
	}

	if decoded["intent_id"] != intentID.String() {
		t.Errorf("intent_id = %v, want %s", decoded["intent_id"], intentID)
	}
	if decoded["wallet_id"] != walletID.String() {
		t.Errorf("wallet_id = %v, want %s", decoded["wallet_id"], walletID)
	}
	if decoded["provider"] != "mercado_pago" {
		t.Errorf("provider = %v, want mercado_pago", decoded["provider"])
	}
	if decoded["external_id"] != "1234567890" {
		t.Errorf("external_id = %v, want 1234567890", decoded["external_id"])
	}
	if decoded["event_type"] != "payment.succeeded" {
		t.Errorf("event_type = %v, want payment.succeeded", decoded["event_type"])
	}

	// status_detail is omitted when the event has no outcome detail.
	if _, ok := decoded["status_detail"]; ok {
		t.Errorf("status_detail present for success event, want omitted")
	}

	// payload must be embedded as an object, not a string (D2 envelope).
	payload, ok := decoded["payload"].(map[string]any)
	if !ok {
		t.Fatalf("payload = %T, want embedded object", decoded["payload"])
	}
	if payload["status"] != "approved" {
		t.Errorf("payload.status = %v, want approved", payload["status"])
	}
}

// TestEnvelopeSignatureCoversEnvelopeBytes pins the contract consumers rely
// on (split 4 verifies `X-Payssage-Signature` over the raw body): the
// signature is HMAC-SHA256(secret, envelope bytes) — the exact bytes that
// get POSTed.
func TestEnvelopeSignatureCoversEnvelopeBytes(t *testing.T) {
	secret := "test-secret"
	statusDetail := "high_risk"
	event := &models.WebhookEvent{
		IntentID:     uuid.New(),
		WalletID:     uuid.New(),
		Provider:     "mercado_pago",
		ExternalID:   "42",
		EventType:    "payment.rejected",
		StatusDetail: &statusDetail,
		Payload:      json.RawMessage(`{"id":42}`),
	}

	envelope, err := buildEnvelope(event)
	if err != nil {
		t.Fatalf("buildEnvelope: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(envelope, &decoded); err != nil {
		t.Fatalf("envelope not valid json: %v", err)
	}
	if decoded["status_detail"] != "high_risk" {
		t.Errorf("status_detail = %v, want high_risk", decoded["status_detail"])
	}

	got := signPayload(secret, envelope)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(envelope)
	want := hex.EncodeToString(mac.Sum(nil))

	if got != want {
		t.Errorf("signature = %s, want %s (must cover the envelope bytes)", got, want)
	}
}
