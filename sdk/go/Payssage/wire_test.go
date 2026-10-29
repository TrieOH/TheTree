package payssage

import (
	"context"
	"encoding/json"
	"net/http"
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

func TestAPIErrorReadsTheErrorObject(t *testing.T) {
	apiErr := &APIError{}
	err := json.Unmarshal([]byte(`{"code":"NOT_FOUND","message":"wallet not found"}`), apiErr)
	if err != nil {
		t.Fatal(err)
	}
	if apiErr.Code != "NOT_FOUND" || apiErr.Message != "wallet not found" {
		t.Fatalf("APIError = %+v", apiErr)
	}
}
