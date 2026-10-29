package identityx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntrospectReturnsTheIdentity(t *testing.T) {
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("X-API-KEY")
		if r.URL.Path != "/auth/introspect" {
			t.Errorf("path = %s, want /auth/introspect", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"data":{"subject":{"id":"019a0000-0000-7000-8000-000000000001","type":"service"},"cred":{"type":"api_key"}}}`))
	}))
	defer srv.Close()

	hook := apiKeyHook(srv.URL+"/", options{acceptAPIKeys: true, httpClient: srv.Client()})
	ctx, err := hook(context.Background(), "key-1")
	if err != nil {
		t.Fatal(err)
	}
	if gotKey != "key-1" {
		t.Fatalf("X-API-KEY = %q, want key-1", gotKey)
	}
	if ctx == nil {
		t.Fatal("no context returned")
	}
}

func TestIntrospectRejectsAnInvalidKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	hook := apiKeyHook(srv.URL, options{acceptAPIKeys: true, httpClient: srv.Client()})
	_, err := hook(context.Background(), "bad")
	if !errors.Is(err, ErrInvalidAPIKey) {
		t.Fatalf("err = %v, want ErrInvalidAPIKey", err)
	}
}

func TestAPIKeysRefusedByDefault(t *testing.T) {
	hook := apiKeyHook("http://unused.invalid", options{})
	_, err := hook(context.Background(), "any")
	if !errors.Is(err, ErrAPIKeysNotAccepted) {
		t.Fatalf("err = %v, want ErrAPIKeysNotAccepted", err)
	}
}

func TestIntrospectReadsBothShapes(t *testing.T) {
	const identity = `{"subject":{"id":"019a0000-0000-7000-8000-000000000001","type":"service"},"cred":{"type":"api_key"}}`
	bodies := map[string]string{
		"enveloped": `{"code":200,"message":"","module":"IdentityX","timestamp":"2026-10-29T00:00:00Z","data":` + identity + `}`,
		"bare":      identity,
	}
	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(body))
			}))
			defer srv.Close()

			got, err := introspect(context.Background(), srv.Client(), srv.URL, "key-1")
			if err != nil {
				t.Fatal(err)
			}
			if got.Sub.ID.String() != "019a0000-0000-7000-8000-000000000001" {
				t.Fatalf("subject id = %s, want 019a0000-0000-7000-8000-000000000001", got.Sub.ID)
			}
		})
	}
}
