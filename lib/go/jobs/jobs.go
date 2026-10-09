// Package jobs is the provider-neutral background-work seam. A Backend
// declares its jobs (a Job is a JSON-serializable value with a Kind),
// registers one handler per kind, and enqueues through an Enqueuer — it never
// learns what carries the work:
//
//   - inprocess: a bounded worker pool plus tickers, for a long-lived
//     process (hosted mode).
//   - sqsjobs: an SQS queue, drained by a worker function through
//     lambdaevents (managed mode). Periodic jobs are then fired by a
//     scheduler that invokes the same worker with an Envelope.
//
// The Envelope is the only wire format, so every carrier delivers the same
// bytes to the same Registry.
package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"lib/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Job is one unit of background work. Implementations are plain structs
// that round-trip through encoding/json.
type Job interface {
	Kind() string
}

// Enqueuer hands a Job to whatever carries background work.
type Enqueuer interface {
	Enqueue(ctx context.Context, job Job) error
}

// Envelope is the wire format every carrier transports.
type Envelope struct {
	Kind string          `json:"kind"`
	Args json.RawMessage `json:"args,omitempty"`
}

// Encode wraps a Job in its Envelope.
func Encode(job Job) (Envelope, error) {
	args, err := json.Marshal(job)
	if err != nil {
		return Envelope{}, fmt.Errorf("encode job %s: %w", job.Kind(), err)
	}
	return Envelope{Kind: job.Kind(), Args: args}, nil
}

// Periodic declares a job to run on an interval. Hosted mode runs it on a
// ticker; managed mode turns the same declaration into a schedule at deploy
// time, so the cadence lives in one place.
type Periodic struct {
	Job        Job
	Every      time.Duration
	RunOnStart bool
}

// ErrUnknownKind is returned for an envelope no handler is registered for.
var ErrUnknownKind = errors.New("unknown job kind")

type handler func(ctx context.Context, args json.RawMessage) error

// Registry maps job kinds to their handlers.
type Registry struct {
	handlers map[string]handler
}

func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]handler)}
}

// Register binds work to the kind of T. Registering a kind twice panics:
// it is a wiring bug, caught at boot.
func Register[T Job](r *Registry, work func(ctx context.Context, job T) error) {
	var zero T
	kind := zero.Kind()
	if _, dup := r.handlers[kind]; dup {
		panic("jobs: kind registered twice: " + kind)
	}
	r.handlers[kind] = func(ctx context.Context, args json.RawMessage) error {
		var job T
		if len(args) > 0 {
			err := json.Unmarshal(args, &job)
			if err != nil {
				return fmt.Errorf("decode job %s: %w", kind, err)
			}
		}
		return work(ctx, job)
	}
}

// Kinds lists the registered kinds.
func (r *Registry) Kinds() []string {
	kinds := make([]string, 0, len(r.handlers))
	for kind := range r.handlers {
		kinds = append(kinds, kind)
	}
	return kinds
}

// Dispatch runs the handler for env inside a span named after its kind.
func (r *Registry) Dispatch(ctx context.Context, env Envelope) error {
	h, ok := r.handlers[env.Kind]
	if !ok {
		return fmt.Errorf("%w: %q", ErrUnknownKind, env.Kind)
	}
	ctx, span := telemetry.StartSpan(ctx, "job "+env.Kind)
	defer span.End()
	span.SetAttributes(attribute.String("job.kind", env.Kind))
	err := h(ctx, env.Args)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}
