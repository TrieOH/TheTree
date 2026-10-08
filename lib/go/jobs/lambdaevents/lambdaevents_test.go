package lambdaevents_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"lib/jobs"
	"lib/jobs/lambdaevents"
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

func post(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, lambdaevents.Path, strings.NewReader(body))
	lambdaevents.Handler(registry()).ServeHTTP(rec, req)
	return rec
}

func TestSQSBatchReportsOnlyFailures(t *testing.T) {
	rec := post(t, `{"Records":[
		{"messageId":"ok","eventSource":"aws:sqs","body":"{\"kind\":\"work\",\"args\":{\"fail\":false}}"},
		{"messageId":"bad","eventSource":"aws:sqs","body":"{\"kind\":\"work\",\"args\":{\"fail\":true}}"},
		{"messageId":"unknown","eventSource":"aws:sqs","body":"{\"kind\":\"nope\"}"}
	]}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var resp struct {
		BatchItemFailures []struct {
			ItemIdentifier string `json:"itemIdentifier"`
		} `json:"batchItemFailures"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, 0, len(resp.BatchItemFailures))
	for _, f := range resp.BatchItemFailures {
		got = append(got, f.ItemIdentifier)
	}
	if strings.Join(got, ",") != "bad,unknown" {
		t.Fatalf("failures = %v, want [bad unknown]", got)
	}
}

func TestScheduledEnvelope(t *testing.T) {
	if rec := post(t, `{"kind":"work"}`); rec.Code != http.StatusNoContent {
		t.Fatalf("ok status = %d, want 204", rec.Code)
	}
	if rec := post(t, `{"kind":"work","args":{"fail":true}}`); rec.Code != http.StatusInternalServerError {
		t.Fatalf("failing status = %d, want 500", rec.Code)
	}
	if rec := post(t, `{"detail":"something else"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("unrecognized status = %d, want 400", rec.Code)
	}
}
