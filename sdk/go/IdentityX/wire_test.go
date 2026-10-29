package idx

import (
	"encoding/json"
	"testing"
)

func TestAPIErrorReadsTheErrorObject(t *testing.T) {
	apiErr := &APIError{}
	err := json.Unmarshal([]byte(`{"code":"NOT_FOUND","message":"actor not found"}`), apiErr)
	if err != nil {
		t.Fatal(err)
	}
	if apiErr.Code != "NOT_FOUND" || apiErr.Message != "actor not found" {
		t.Fatalf("APIError = %+v", apiErr)
	}
}
