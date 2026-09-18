package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	payssage "sdk/payssage"

	"github.com/google/uuid"
)

// TestVerifyPayssageWallet_Resolves pins the happy path of the D6 boot
// check: a wallet that resolves must not stop the boot.
func TestVerifyPayssageWallet_Resolves(t *testing.T) {
	walletID := uuid.New()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"data":{
			"id":"` + walletID.String() + `","owner_id":"` + uuid.New().String() + `",
			"organization_id":null,"name":"platform","sandbox":false,
			"fee_bps":500,"collector_id":null,"created_at":"2026-08-01T00:00:00Z"
		},"timestamp":"2026-08-01T00:00:00Z"}`))
	}))
	t.Cleanup(srv.Close)

	client := payssage.New(payssage.Config{BaseURL: srv.URL, APIKey: "test"})

	// Must not exit — the resolved wallet passes the boot gate.
	err := VerifyPayssageWallet(context.Background(), client, walletID)
	if err != nil {
		t.Fatalf("resolved wallet must not stop the boot: %v", err)
	}
}

// TestVerifyPayssageWallet_FailsOnUnresolvableWallet pins the fail-fast
// side: a wallet that does not resolve (wrong id / unreachable Payssage)
// surfaces as a boot error, so Boot stops the process before serving.
func TestVerifyPayssageWallet_FailsOnUnresolvableWallet(t *testing.T) {
	client := payssage.New(payssage.Config{BaseURL: "http://127.0.0.1:1", APIKey: "test"})
	err := VerifyPayssageWallet(context.Background(), client, uuid.New())
	if err == nil {
		t.Fatal("want a boot error from an unresolvable wallet, got nil")
	}
}
