package identityx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MintzyG/fun"
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
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("err = %v, want forbidden", err)
	}
}

func TestAPIKeysRefusedByDefault(t *testing.T) {
	hook := apiKeyHook("http://unused.invalid", options{})
	_, err := hook(context.Background(), "any")
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("err = %v, want forbidden", err)
	}
}
