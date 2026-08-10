package payssage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// newTestClient spins an httptest server that records the last request and
// responds with a canned envelope, and returns the client pointed at it.
func newTestClient(t *testing.T, status int, envelope string) (*Client, func() *http.Request, func() string) {
	t.Helper()
	var captured *http.Request
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		b := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(b)
		body = string(b)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(envelope))
	}))
	t.Cleanup(srv.Close)
	return New(Config{BaseURL: srv.URL, APIKey: "test-api-key"}), func() *http.Request { return captured }, func() string { return body }
}

func TestCreateWallet(t *testing.T) {
	walletID := uuid.New()
	ownerID := uuid.New()
	client, getReq, getBody := newTestClient(t, http.StatusCreated, `{"data": {
		"id": "`+walletID.String()+`", "owner_id": "`+ownerID.String()+`",
		"organization_id": null, "name": "event-slug", "sandbox": false,
		"fee_bps": 0, "collector_id": null, "created_at": "2026-08-01T00:00:00Z"
	}}`)

	wallet, err := client.CreateWallet(context.Background(), CreateWalletRequest{Name: "event-slug"})
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}
	if wallet.ID != walletID {
		t.Fatalf("wallet id = %s, want %s", wallet.ID, walletID)
	}
	if wallet.Name != "event-slug" {
		t.Fatalf("wallet name = %q, want event-slug", wallet.Name)
	}
	if r := getReq(); r.Method != http.MethodPost || r.URL.Path != "/wallets" {
		t.Fatalf("got %s %s, want POST /wallets", r.Method, r.URL.Path)
	}
	if r := getReq(); r.Header.Get("X-API-Key") != "test-api-key" {
		t.Fatalf("X-API-Key header missing")
	}
	if !strings.Contains(getBody(), `"name":"event-slug"`) {
		t.Fatalf("body missing name: %s", getBody())
	}
}

func TestCreateWalletWithOrganization(t *testing.T) {
	orgID := uuid.New()
	client, _, getBody := newTestClient(t, http.StatusCreated, `{"data": {}}`)
	_, err := client.CreateWallet(context.Background(), CreateWalletRequest{Name: "org-wallet", OrganizationID: &orgID})
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}
	if !strings.Contains(getBody(), `"organization_id":"`+orgID.String()+`"`) {
		t.Fatalf("body missing organization_id: %s", getBody())
	}
}

func TestSetWalletFee(t *testing.T) {
	walletID := uuid.New()
	client, getReq, getBody := newTestClient(t, http.StatusOK, `{"code":200,"message":"ok","timestamp":"2026-08-01T00:00:00Z"}`)
	err := client.SetWalletFee(context.Background(), walletID, 500)
	if err != nil {
		t.Fatalf("SetWalletFee: %v", err)
	}
	if r := getReq(); r.Method != http.MethodPatch || r.URL.Path != "/wallets/"+walletID.String()+"/fee" {
		t.Fatalf("got %s %s, want PATCH /wallets/%s/fee", r.Method, r.URL.Path, walletID)
	}
	if !strings.Contains(getBody(), `"fee_bps":500`) {
		t.Fatalf("body missing fee_bps: %s", getBody())
	}
}

func TestListWalletSellers(t *testing.T) {
	walletID := uuid.New()
	sellerID := uuid.New()
	client, getReq, _ := newTestClient(t, http.StatusOK, `{"data": [{
		"id": "`+sellerID.String()+`", "wallet_id": "`+walletID.String()+`",
		"provider": "mercado_pago", "provider_user_id": "12345",
		"created_at": "2026-08-01T00:00:00Z", "revoked_at": null
	}]}`)

	sellers, err := client.ListWalletSellers(context.Background(), walletID)
	if err != nil {
		t.Fatalf("ListWalletSellers: %v", err)
	}
	if len(sellers) != 1 || sellers[0].ID != sellerID || sellers[0].Provider != "mercado_pago" {
		t.Fatalf("unexpected sellers: %+v", sellers)
	}
	if r := getReq(); r.Method != http.MethodGet || r.URL.Path != "/wallets/"+walletID.String()+"/sellers" {
		t.Fatalf("got %s %s, want GET /wallets/%s/sellers", r.Method, r.URL.Path, walletID)
	}
}

func TestConnectProvider(t *testing.T) {
	walletID := uuid.New()
	client, getReq, getBody := newTestClient(t, http.StatusOK, `{"data":"https://auth.mercadopago.com/consent"}`)

	url, err := client.ConnectProvider(context.Background(), "mercado_pago", ConnectProviderRequest{
		Flow:             OAuthFlowSeller,
		WalletID:         &walletID,
		FinalRedirectURL: "https://events.example/events/abc/payssage/oauth/callback",
	})
	if err != nil {
		t.Fatalf("ConnectProvider: %v", err)
	}
	if url != "https://auth.mercadopago.com/consent" {
		t.Fatalf("url = %q", url)
	}
	if r := getReq(); r.Method != http.MethodPost || r.URL.Path != "/providers/mercado_pago/connect" {
		t.Fatalf("got %s %s", r.Method, r.URL.Path)
	}
	var payload map[string]any
	err = json.Unmarshal([]byte(getBody()), &payload)
	if err != nil {
		t.Fatalf("body not json: %v", err)
	}
	if payload["flow"] != "seller" || payload["wallet_id"] != walletID.String() {
		t.Fatalf("unexpected payload: %s", getBody())
	}
	if _, ok := payload["provider_redirect_url"]; ok {
		t.Fatalf("provider_redirect_url must not be sent (D7): %s", getBody())
	}
	if payload["final_redirect_url"] != "https://events.example/events/abc/payssage/oauth/callback" {
		t.Fatalf("unexpected final_redirect_url: %s", getBody())
	}
}

func TestAPIError(t *testing.T) {
	client, _, _ := newTestClient(t, http.StatusNotFound, `{"code":404,"message":"wallet not found","error":null,"timestamp":"2026-08-01T00:00:00Z","module":"payssage"}`)
	_, err := client.CreateWallet(context.Background(), CreateWalletRequest{Name: "x"})
	if err == nil {
		t.Fatal("want error")
	}
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false", err)
	}
}
