package idx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestAPIErrorReadsTheProblem(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"type":"https://api.identityx.test/problems/not-found","status":404,` +
			`"title":"Not Found","detail":"actor not found","instance":"/actors/x"}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(Config{BaseURL: srv.URL, APIKey: "test-api-key", ProjectID: uuid.New()})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.Actors.GetByID(context.Background(), uuid.New())
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false", err)
	}
	apiErr, _ := errors.AsType[*APIError](err)
	if apiErr.Type != "https://api.identityx.test/problems/not-found" || apiErr.Title != "Not Found" ||
		apiErr.Message != "actor not found" {
		t.Fatalf("APIError = %+v", apiErr)
	}
}
