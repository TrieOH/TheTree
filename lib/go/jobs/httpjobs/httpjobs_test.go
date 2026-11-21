package httpjobs_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

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
	httpjobs.Handler(registry(), "").ServeHTTP(rec, req)
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

func TestEnqueueRetriesA5xx(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) < 3 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	e, err := httpjobs.New(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	err = e.Enqueue(context.Background(), work{})
	if err != nil {
		t.Fatalf("third try should succeed: %v", err)
	}
	if calls.Load() != 3 {
		t.Fatalf("calls = %d, want 3", calls.Load())
	}
}

func TestEnqueueDoesNotRetryA4xx(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	e, err := httpjobs.New(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if e.Enqueue(context.Background(), work{}) == nil {
		t.Fatal("expected an error")
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1", calls.Load())
	}
}

func TestEnqueueAfterSendsTheDelay(t *testing.T) {
	headers := make(chan string, 2)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers <- r.Header.Get(httpjobs.DelayHeader)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	e, err := httpjobs.New(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	err = e.EnqueueAfter(context.Background(), work{}, 1500*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if got := <-headers; got != "2" {
		t.Fatalf("delay header = %q, want 2 (rounded up)", got)
	}
	err = e.Enqueue(context.Background(), work{})
	if err != nil {
		t.Fatal(err)
	}
	if got := <-headers; got != "" {
		t.Fatalf("delay header = %q on a plain enqueue, want none", got)
	}
}

func TestHandlerRequiresTheSecretWhenSet(t *testing.T) {
	h := httpjobs.Handler(registry(), "jobs-secret")
	for _, tc := range []struct {
		name, sent string
		want       int
	}{
		{"right secret", "jobs-secret", http.StatusNoContent},
		{"wrong secret", "nope", http.StatusUnauthorized},
		{"no secret", "", http.StatusUnauthorized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, httpjobs.Path, strings.NewReader(`{"kind":"work"}`))
			if tc.sent != "" {
				req.Header.Set(httpjobs.SecretHeader, tc.sent)
			}
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d", rec.Code, tc.want)
			}
		})
	}
}
