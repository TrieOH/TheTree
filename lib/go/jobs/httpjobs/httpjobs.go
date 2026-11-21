// Package httpjobs carries lib/jobs over plain HTTP, for when something
// outside the process queues and schedules the work:
//
//   - Enqueuer POSTs each job's Envelope as JSON to a URL. Whatever answers
//     there (a queue producer) owns delivery and retries. A delayed job
//     carries its delay in the DelayHeader, in whole seconds. A POST that
//     fails outright or answers 5xx is tried again a couple of times, so a
//     blip on the way to the queue doesn't lose the job.
//   - Handler serves Path: it receives one Envelope per POST, from the queue
//     consumer or a scheduler, and dispatches it. 2xx means done; 5xx means
//     the job failed and should be retried; 4xx means it never will succeed.
//
// The Envelope is the whole contract, so any queue that can POST JSON can sit
// in between. Mount Handler only where the outside world cannot reach it.
package httpjobs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"lib/jobs"
	"lib/telemetry"
)

// Path is where Handler is mounted. It must not collide with any API
// route: Univents has its own /events resource.
const Path = "/_jobs"

// DelayHeader carries a delayed job's delay, in whole seconds, to the queue
// producer.
const DelayHeader = "Job-Delay-Seconds"

// maxEnvelopeBytes bounds one Envelope.
const maxEnvelopeBytes = 1 << 20

// enqueueTimeout bounds one enqueue POST.
const enqueueTimeout = 10 * time.Second

// retryBackoff is the wait before each retry of a failed enqueue POST: three
// tries in all.
var retryBackoff = []time.Duration{200 * time.Millisecond, 800 * time.Millisecond}

type Enqueuer struct {
	url    string
	client *http.Client
}

var _ jobs.DelayedEnqueuer = (*Enqueuer)(nil)

// New builds an Enqueuer that POSTs to url.
func New(url string) (*Enqueuer, error) {
	if url == "" {
		return nil, errors.New("httpjobs: enqueue URL is required")
	}
	return &Enqueuer{url: url, client: &http.Client{Timeout: enqueueTimeout}}, nil
}

func (e *Enqueuer) Enqueue(ctx context.Context, job jobs.Job) error {
	return e.post(ctx, job, 0)
}

// EnqueueAfter POSTs the job with its delay, rounded up to whole seconds;
// the queue holds the message until then.
func (e *Enqueuer) EnqueueAfter(ctx context.Context, job jobs.Job, delay time.Duration) error {
	if delay > jobs.MaxDelay {
		return jobs.ErrDelayTooLong
	}
	return e.post(ctx, job, delay)
}

func (e *Enqueuer) post(ctx context.Context, job jobs.Job, delay time.Duration) error {
	env, err := jobs.Encode(job)
	if err != nil {
		return err
	}
	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("httpjobs: encode envelope: %w", err)
	}
	for attempt := 0; ; attempt++ {
		retry, err := e.send(ctx, body, delay)
		if err == nil {
			return nil
		}
		if !retry || attempt == len(retryBackoff) {
			return fmt.Errorf("httpjobs: enqueue %s: %w", env.Kind, err)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("httpjobs: enqueue %s: %w", env.Kind, errors.Join(err, ctx.Err()))
		case <-time.After(retryBackoff[attempt]):
		}
	}
}

// send makes one enqueue POST. retry reports whether a failure is worth
// another try: a transport error or a 5xx is, a 4xx never succeeds.
func (e *Enqueuer) send(ctx context.Context, body []byte, delay time.Duration) (retry bool, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.url, bytes.NewReader(body))
	if err != nil {
		return false, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if delay > 0 {
		req.Header.Set(DelayHeader, strconv.FormatInt(int64((delay+time.Second-1)/time.Second), 10))
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return ctx.Err() == nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxEnvelopeBytes))
	if resp.StatusCode/100 != 2 {
		return resp.StatusCode >= 500, fmt.Errorf("status %d", resp.StatusCode)
	}
	return false, nil
}

// Handler dispatches the posted Envelope into registry.
func Handler(registry *jobs.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		raw, err := io.ReadAll(io.LimitReader(r.Body, maxEnvelopeBytes))
		if err != nil {
			http.Error(w, "read envelope", http.StatusBadRequest)
			return
		}
		var env jobs.Envelope
		err = json.Unmarshal(raw, &env)
		if err != nil || env.Kind == "" {
			http.Error(w, "decode envelope", http.StatusBadRequest)
			return
		}
		err = registry.Dispatch(r.Context(), env)
		switch {
		case errors.Is(err, jobs.ErrUnknownKind):
			http.Error(w, "unknown job kind", http.StatusBadRequest)
		case err != nil:
			telemetry.Log().ErrorContext(r.Context(), "job failed", slog.String("kind", env.Kind), telemetry.Err(err))
			http.Error(w, "job failed", http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	})
}
