package httpjobs_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"lib/jobs"
	"lib/jobs/httpjobs"
)

type work struct {
	Fail bool `json:"fail"`
}

func (work) Kind() string { return "work" }

func registry() *jobs.Registry {
	reg := jobs.NewRegistry()
	jobs.Register(reg, func(_ context.Context, w work) error {
		if w.Fail {
			return errors.New("boom")
		}
		return nil
	})
	return reg
}

func post(t *testing.T, method, body string) int {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(context.Background(), method, httpjobs.Path, strings.NewReader(body))
	httpjobs.Handler(registry()).ServeHTTP(rec, req)
	return rec.Code
}

func TestHandlerStatuses(t *testing.T) {
	cases := []struct {
		name   string
		method string
		body   string
		want   int
	}{
		{"success", http.MethodPost, `{"kind":"work","args":{"fail":false}}`, http.StatusNoContent},
		{"failure is retryable", http.MethodPost, `{"kind":"work","args":{"fail":true}}`, http.StatusInternalServerError},
		{"unknown kind is final", http.MethodPost, `{"kind":"nope"}`, http.StatusBadRequest},
		{"missing kind", http.MethodPost, `{}`, http.StatusBadRequest},
		{"bad json", http.MethodPost, `{`, http.StatusBadRequest},
		{"wrong method", http.MethodGet, ``, http.StatusMethodNotAllowed},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := post(t, c.method, c.body)
			if got != c.want {
				t.Fatalf("status = %d, want %d", got, c.want)
			}
		})
	}
}

func TestEnqueuePostsTheEnvelope(t *testing.T) {
	var got jobs.Envelope
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &got)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	e, err := httpjobs.New(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	err = e.Enqueue(context.Background(), work{Fail: true})
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "work" || string(got.Args) != `{"fail":true}` {
		t.Fatalf("envelope = %+v", got)
	}
}

func TestEnqueueFailsOnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	e, err := httpjobs.New(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if e.Enqueue(context.Background(), work{}) == nil {
		t.Fatal("expected an error")
	}
}
