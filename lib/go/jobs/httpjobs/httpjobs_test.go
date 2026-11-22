package httpjobs_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func post(t *testing.T, method, body string, headers map[string]string) (int, httpjobs.Outcome) {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(context.Background(), method, httpjobs.Path, strings.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	httpjobs.Handler(registry(), "").ServeHTTP(rec, req)
	var out httpjobs.Outcome
	if rec.Code == http.StatusOK {
		err := json.Unmarshal(rec.Body.Bytes(), &out)
		if err != nil {
			t.Fatalf("body %q: %v", rec.Body.String(), err)
		}
	}
	return rec.Code, out
}

func TestHandlerOutcomes(t *testing.T) {
	last := strconv.Itoa(jobs.DefaultPolicy.MaxAttempts)
	cases := []struct {
		name    string
		method  string
		body    string
		headers map[string]string
		status  int
		want    httpjobs.Outcome
	}{
		{"success", http.MethodPost, `{"kind":"work","args":{"fail":false}}`, nil, http.StatusOK, httpjobs.Outcome{Outcome: jobs.Done}},
		{"failure retries after the backoff", http.MethodPost, `{"kind":"work","args":{"fail":true}}`, map[string]string{httpjobs.AttemptHeader: "2"}, http.StatusOK, httpjobs.Outcome{Outcome: jobs.Retry, RetryAfterSeconds: 60}},
		{"last attempt drops", http.MethodPost, `{"kind":"work","args":{"fail":true}}`, map[string]string{httpjobs.AttemptHeader: last}, http.StatusOK, httpjobs.Outcome{Outcome: jobs.Drop}},
		{"dead only gives up", http.MethodPost, `{"kind":"work","args":{"fail":true}}`, map[string]string{httpjobs.DeadHeader: "true"}, http.StatusOK, httpjobs.Outcome{Outcome: jobs.Drop}},
		{"unknown kind drops", http.MethodPost, `{"kind":"nope"}`, nil, http.StatusOK, httpjobs.Outcome{Outcome: jobs.Drop}},
		{"missing kind", http.MethodPost, `{}`, nil, http.StatusBadRequest, httpjobs.Outcome{}},
		{"bad json", http.MethodPost, `{`, nil, http.StatusBadRequest, httpjobs.Outcome{}},
		{"wrong method", http.MethodGet, ``, nil, http.StatusMethodNotAllowed, httpjobs.Outcome{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			status, out := post(t, c.method, c.body, c.headers)
			if status != c.status || out != c.want {
				t.Fatalf("got %d %+v, want %d %+v", status, out, c.status, c.want)
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
		{"right secret", "jobs-secret", http.StatusOK},
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
