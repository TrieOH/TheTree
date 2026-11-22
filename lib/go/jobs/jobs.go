// Package jobs is the provider-neutral background-work seam. A Backend
// declares its jobs (a Job is a JSON-serializable value with a Kind),
// registers one handler per kind, and enqueues through an Enqueuer — it never
// learns what carries the work:
//
//   - inprocess: a bounded worker pool plus tickers, for a long-lived
//     process (local compose).
//   - httpjobs: the Envelope POSTed to an external queue, which delivers it
//     back to the Backend's job endpoint (/_jobs); an external scheduler fires
//     the periodic jobs there (Cloudflare Queues and Cron Triggers).
//
// The Envelope is the only wire format, so every carrier delivers the same
// bytes to the same Registry. The job semantics live here too: a carrier
// hands each delivery to Registry.Handle, which runs it and says what to do
// next (done, retry after a backoff, drop, give up), and the carrier only
// maps that onto its own ack, delayed redelivery or terminate. A message
// broker (NATS JetStream, RabbitMQ) is one more carrier of the same shape.
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
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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

// DelayedEnqueuer also hands over a Job to run no sooner than delay from
// now. Both carriers implement it: the in-process runner holds a timer, and
// httpjobs passes the delay to the queue, which holds the message.
type DelayedEnqueuer interface {
	Enqueuer
	EnqueueAfter(ctx context.Context, job Job, delay time.Duration) error
}

// MaxDelay is the longest delay every carrier honors: Cloudflare Queues
// hold a message for at most 12 hours.
const MaxDelay = 12 * time.Hour

// ErrDelayTooLong is returned for a delay over MaxDelay.
var ErrDelayTooLong = errors.New("job delay over the maximum")

// Envelope is the wire format every carrier transports.
type Envelope struct {
	Kind string          `json:"kind"`
	Args json.RawMessage `json:"args,omitempty"`
	// Trace is the W3C trace context (traceparent, tracestate) of the code
	// that enqueued the job, so the job's span joins that trace: a request,
	// the job it queued and the email the job sent read as one trace. Empty
	// for a job a scheduler fired.
	Trace map[string]string `json:"trace,omitempty"`
}

// traceFormat is the wire format of Envelope.Trace. It's fixed here rather
// than read from the global propagator: the envelope is a contract between
// processes.
var traceFormat = propagation.TraceContext{}

// Encode wraps a Job in its Envelope, with the trace ctx is in.
func Encode(ctx context.Context, job Job) (Envelope, error) {
	args, err := json.Marshal(job)
	if err != nil {
		return Envelope{}, fmt.Errorf("encode job %s: %w", job.Kind(), err)
	}
	env := Envelope{Kind: job.Kind(), Args: args}
	carrier := propagation.MapCarrier{}
	traceFormat.Inject(ctx, carrier)
	if len(carrier) > 0 {
		env.Trace = carrier
	}
	return env, nil
}

// Periodic declares a job to run on an interval. The in-process runner runs
// it on a ticker; with httpjobs the deploy's scheduler fires it instead.
type Periodic struct {
	Job        Job
	Every      time.Duration
	RunOnStart bool
}

// ErrUnknownKind is returned for an envelope no handler is registered for.
var ErrUnknownKind = errors.New("unknown job kind")

type handler func(ctx context.Context, args json.RawMessage) error

// kind is one registered kind: its work, and the job that records its
// give-up.
type kind struct {
	work   handler
	giveUp func(args json.RawMessage) (Envelope, error)
}

// Registry maps job kinds to their handlers, and holds the retry policy
// Handle applies to them.
type Registry struct {
	kinds  map[string]*kind
	policy Policy
}

// RegistryOption tunes a Registry.
type RegistryOption func(*Registry)

// WithPolicy replaces DefaultPolicy.
func WithPolicy(p Policy) RegistryOption {
	return func(r *Registry) { r.policy = p }
}

func NewRegistry(opts ...RegistryOption) *Registry {
	r := &Registry{kinds: make(map[string]*kind), policy: DefaultPolicy}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Option tunes one kind at Register.
type Option[T Job] func(*kind)

// GiveUpWith declares the job that records the give-up of a T whose retries
// ran out: Handle runs giveUp(job) in its place before dropping it. The
// give-up's kind must be registered too.
func GiveUpWith[T Job](giveUp func(job T) Job) Option[T] {
	return func(k *kind) {
		k.giveUp = func(args json.RawMessage) (Envelope, error) {
			job, err := decode[T](args)
			if err != nil {
				return Envelope{}, err
			}
			return Encode(context.Background(), giveUp(job))
		}
	}
}

// Register binds work to the kind of T. Registering a kind twice panics:
// it is a wiring bug, caught at boot.
func Register[T Job](r *Registry, work func(ctx context.Context, job T) error, opts ...Option[T]) {
	var zero T
	name := zero.Kind()
	if _, dup := r.kinds[name]; dup {
		panic("jobs: kind registered twice: " + name)
	}
	k := &kind{work: func(ctx context.Context, args json.RawMessage) error {
		job, err := decode[T](args)
		if err != nil {
			return err
		}
		return work(ctx, job)
	}}
	for _, opt := range opts {
		opt(k)
	}
	r.kinds[name] = k
}

// decode reads a T's args. Args that don't decode never will, so the error
// is permanent.
func decode[T Job](args json.RawMessage) (T, error) {
	var job T
	if len(args) > 0 {
		err := json.Unmarshal(args, &job)
		if err != nil {
			return job, Permanent(fmt.Errorf("decode job %s: %w", job.Kind(), err))
		}
	}
	return job, nil
}

// Kinds lists the registered kinds.
func (r *Registry) Kinds() []string {
	kinds := make([]string, 0, len(r.kinds))
	for name := range r.kinds {
		kinds = append(kinds, name)
	}
	return kinds
}

// Dispatch runs the handler for env inside a span named after its kind. With
// a trace in env the span is a child of the code that enqueued the job, and
// links to whatever delivered it (the /_jobs request).
func (r *Registry) Dispatch(ctx context.Context, env Envelope) error {
	k, ok := r.kinds[env.Kind]
	if !ok {
		return fmt.Errorf("%w: %q", ErrUnknownKind, env.Kind)
	}
	var opts []trace.SpanStartOption
	parent := trace.SpanContextFromContext(traceFormat.Extract(context.Background(), propagation.MapCarrier(env.Trace)))
	if parent.IsValid() {
		delivery := trace.SpanContextFromContext(ctx)
		if delivery.IsValid() {
			opts = append(opts, trace.WithLinks(trace.Link{SpanContext: delivery}))
		}
		ctx = trace.ContextWithRemoteSpanContext(ctx, parent)
	}
	ctx, span := telemetry.StartSpan(ctx, "job "+env.Kind, opts...)
	defer span.End()
	span.SetAttributes(attribute.String("job.kind", env.Kind))
	err := k.work(ctx, env.Args)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}
