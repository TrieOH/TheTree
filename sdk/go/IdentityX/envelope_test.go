package idx

import (
	"encoding/json"
	"testing"
)

func TestEnvelopeReadsBothShapes(t *testing.T) {
	tests := map[string]string{
		"enveloped": `{"code":200,"message":"","module":"IdentityX","timestamp":"2026-10-29T00:00:00Z","data":{"email":"a@b.c"}}`,
		"bare":      `{"email":"a@b.c"}`,
	}
	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			var res envelope[Actor]
			err := json.Unmarshal([]byte(body), &res)
			if err != nil {
				t.Fatal(err)
			}
			if res.Data.Email == nil || *res.Data.Email != "a@b.c" {
				t.Fatalf("email = %v, want a@b.c", res.Data.Email)
			}
		})
	}
}

func TestAPIErrorReadsABareErrorObject(t *testing.T) {
	apiErr := &APIError{}
	_ = json.Unmarshal([]byte(`{"code":"NOT_FOUND","message":"actor not found"}`), apiErr)
	if apiErr.Message != "actor not found" {
		t.Fatalf("message = %q, want actor not found", apiErr.Message)
	}
}
