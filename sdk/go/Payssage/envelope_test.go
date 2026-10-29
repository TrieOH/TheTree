package payssage

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestDecodeReadsBothShapes(t *testing.T) {
	tests := map[string]string{
		"enveloped": `{"code":200,"message":"","module":"Payssage","timestamp":"2026-10-29T00:00:00Z","data":{"name":"w"}}`,
		"bare":      `{"name":"w"}`,
	}
	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			var w Wallet
			err := decode([]byte(body), &w)
			if err != nil {
				t.Fatal(err)
			}
			if w.Name != "w" {
				t.Fatalf("name = %q, want w", w.Name)
			}
		})
	}
}

func TestListReadsItemsOrABareArray(t *testing.T) {
	sellerID := uuid.New()
	tests := map[string]string{
		"enveloped": `{"code":200,"module":"Payssage","data":[{"id":"` + sellerID.String() + `"}]}`,
		"items":     `{"items":[{"id":"` + sellerID.String() + `"}]}`,
	}
	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			client, _, _ := newTestClient(t, http.StatusOK, body)
			sellers, err := client.ListWalletSellers(context.Background(), uuid.New())
			if err != nil {
				t.Fatal(err)
			}
			if len(sellers) != 1 || sellers[0].ID != sellerID {
				t.Fatalf("sellers = %+v", sellers)
			}
		})
	}
}

func TestNoContent(t *testing.T) {
	client, _, _ := newTestClient(t, http.StatusNoContent, "")
	err := client.SetWalletFee(context.Background(), uuid.New(), 100)
	if err != nil {
		t.Fatalf("SetWalletFee: %v", err)
	}
}

func TestAPIErrorReadsABareErrorObject(t *testing.T) {
	apiErr := &APIError{}
	_ = json.Unmarshal([]byte(`{"code":"NOT_FOUND","message":"wallet not found"}`), apiErr)
	if apiErr.Message != "wallet not found" {
		t.Fatalf("message = %q, want wallet not found", apiErr.Message)
	}
}
