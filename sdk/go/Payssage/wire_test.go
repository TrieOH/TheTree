package payssage

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestNoContent(t *testing.T) {
	client, _, _ := newTestClient(t, http.StatusNoContent, "")
	err := client.SetWalletFee(context.Background(), uuid.New(), 100)
	if err != nil {
		t.Fatalf("SetWalletFee: %v", err)
	}
}

func TestAPIErrorReadsTheProblem(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"type":"https://api.payssage.test/problems/not-found","status":404,` +
			`"title":"Not Found","detail":"wallet not found","instance":"/wallets"}`))
	}))
	t.Cleanup(srv.Close)

	client := New(Config{BaseURL: srv.URL, APIKey: "test-api-key"})
	_, err := client.CreateWallet(context.Background(), CreateWalletRequest{Name: "x"})
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false", err)
	}
	apiErr, _ := errors.AsType[*APIError](err)
	if apiErr.Type != "https://api.payssage.test/problems/not-found" || apiErr.Title != "Not Found" ||
		apiErr.Message != "wallet not found" {
		t.Fatalf("APIError = %+v", apiErr)
	}
}
