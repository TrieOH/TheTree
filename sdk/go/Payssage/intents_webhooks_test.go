package payssage

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGetWallet(t *testing.T) {
	walletID := uuid.New()
	ownerID := uuid.New()
	client, getReq, _ := newTestClient(t, http.StatusOK, `{"data": {
		"id": "`+walletID.String()+`", "owner_id": "`+ownerID.String()+`",
		"organization_id": null, "name": "platform", "sandbox": true,
		"fee_bps": 500, "collector_id": null, "created_at": "2026-08-01T00:00:00Z"
	}}`)

	wallet, err := client.GetWallet(context.Background(), walletID)
	if err != nil {
		t.Fatalf("GetWallet: %v", err)
	}
	if wallet.ID != walletID {
		t.Fatalf("wallet id = %s, want %s", wallet.ID, walletID)
	}
	if wallet.FeeBps != 500 {
		t.Fatalf("fee_bps = %d, want 500", wallet.FeeBps)
	}
	if r := getReq(); r.Method != http.MethodGet || r.URL.Path != "/wallets/"+walletID.String() {
		t.Fatalf("got %s %s, want GET /wallets/%s", r.Method, r.URL.Path, walletID)
	}
}

func TestGetWallet_NotFound(t *testing.T) {
	client, _, _ := newTestClient(t, http.StatusNotFound, `{"code":404,"message":"wallet not found","timestamp":"2026-08-01T00:00:00Z"}`)
	_, err := client.GetWallet(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("want error")
	}
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false", err)
	}
}

func TestCheckout(t *testing.T) {
	walletID := uuid.New()
	intentID := uuid.New()
	client, getReq, getBody := newTestClient(t, http.StatusCreated, `{"data": {
		"id": "`+intentID.String()+`", "wallet_id": "`+walletID.String()+`",
		"seller_id": "`+uuid.New().String()+`", "collector_id": null,
		"amount_cents": 1234, "currency": "BRL", "sandbox": false,
		"provider": "mercado_pago", "status": "pending", "status_detail": null,
		"provider_data": {"pix_qr_code": "QR"}, "metadata": null,
		"created_at": "2026-08-01T00:00:00Z", "updated_at": "2026-08-01T00:00:00Z"
	}}`)

	intent, err := client.Checkout(context.Background(), walletID, CreateIntentRequest{
		SellerID:    uuid.New(),
		Currency:    "BRL",
		AmountCents: 1234,
	})
	if err != nil {
		t.Fatalf("Checkout: %v", err)
	}
	if intent.ID != intentID {
		t.Fatalf("intent id = %s, want %s", intent.ID, intentID)
	}
	if intent.Status != IntentStatusPending {
		t.Fatalf("status = %s, want pending", intent.Status)
	}
	if r := getReq(); r.Method != http.MethodPost || r.URL.Path != "/wallets/"+walletID.String()+"/intents" {
		t.Fatalf("got %s %s, want POST /wallets/%s/intents", r.Method, r.URL.Path, walletID)
	}
	body := getBody()
	if !strings.Contains(body, `"amount_cents":1234`) || !strings.Contains(body, `"currency":"BRL"`) {
		t.Fatalf("body missing checkout fields: %s", body)
	}
}

func TestGetIntent(t *testing.T) {
	intentID := uuid.New()
	client, getReq, _ := newTestClient(t, http.StatusOK, `{"data": {
		"id": "`+intentID.String()+`", "wallet_id": "`+uuid.New().String()+`",
		"seller_id": "`+uuid.New().String()+`", "collector_id": null,
		"amount_cents": 100, "currency": "BRL", "sandbox": false,
		"provider": "mercado_pago", "status": "succeeded", "status_detail": null,
		"provider_data": {}, "metadata": null,
		"created_at": "2026-08-01T00:00:00Z", "updated_at": "2026-08-01T00:00:00Z"
	}}`)

	intent, err := client.GetIntent(context.Background(), intentID)
	if err != nil {
		t.Fatalf("GetIntent: %v", err)
	}
	if intent.ID != intentID {
		t.Fatalf("intent id = %s, want %s", intent.ID, intentID)
	}
	if r := getReq(); r.Method != http.MethodGet || r.URL.Path != "/intents/"+intentID.String() {
		t.Fatalf("got %s %s, want GET /intents/%s", r.Method, r.URL.Path, intentID)
	}
}

func TestCancelIntent(t *testing.T) {
	intentID := uuid.New()
	client, getReq, _ := newTestClient(t, http.StatusOK, `{"data": {
		"id": "`+intentID.String()+`", "wallet_id": "`+uuid.New().String()+`",
		"seller_id": "`+uuid.New().String()+`", "collector_id": null,
		"amount_cents": 100, "currency": "BRL", "sandbox": false,
		"provider": "mercado_pago", "status": "cancelled", "status_detail": null,
		"provider_data": {}, "metadata": null,
		"created_at": "2026-08-01T00:00:00Z", "updated_at": "2026-08-01T00:00:00Z"
	}}`)

	intent, err := client.CancelIntent(context.Background(), intentID)
	if err != nil {
		t.Fatalf("CancelIntent: %v", err)
	}
	if intent.Status != IntentStatusCancelled {
		t.Fatalf("status = %s, want cancelled", intent.Status)
	}
	if r := getReq(); r.Method != http.MethodPost || r.URL.Path != "/intents/"+intentID.String()+"/cancel" {
		t.Fatalf("got %s %s, want POST /intents/%s/cancel", r.Method, r.URL.Path, intentID)
	}
}
