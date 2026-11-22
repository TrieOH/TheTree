package jobs

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"lib/telemetry"
)

// Delivery is one hand-over of an Envelope by a carrier. Attempt counts the
// hand-overs, 1 on the first; Dead marks a job the carrier itself gave up
// redelivering (a dead-letter queue), which only gets its give-up run.
type Delivery struct {
	Envelope

	Attempt int
	Dead    bool
}

// Outcome is what the carrier does with a delivered job.
type Outcome string

const (
	// Done acks the job.
	Done Outcome = "done"
	// Retry hands the job over again after Result.After.
	Retry Outcome = "retry"
	// Drop acks a job that did not succeed: it never can, or it was given
	// up and its give-up has run.
	Drop Outcome = "drop"
)

// Result is Handle's verdict on a delivery. Carriers map it onto their own
// primitives (an ack, a delayed redelivery, a terminate) and decide nothing
// themselves.
type Result struct {
	Outcome Outcome
	After   time.Duration // Retry only
}

// Policy is how often and how patiently a failing job is retried.
type Policy struct {
	// MaxAttempts is the delivery after which a failing job is given up.
	MaxAttempts int
	// Backoff is the wait before the delivery after attempt.
	Backoff func(attempt int) time.Duration
}

// DefaultPolicy gives a job 10 attempts over a few hours: 30s, doubling,
// capped at an hour. A tenant's webhook endpoint gets that long to come back.
var DefaultPolicy = Policy{MaxAttempts: 10, Backoff: Exponential(30*time.Second, time.Hour)}

// Exponential waits base after the first attempt, doubling per attempt up to
// limit.
func Exponential(base, limit time.Duration) func(attempt int) time.Duration {
	return func(attempt int) time.Duration {
		wait := base
		for range max(attempt, 1) - 1 {
			wait *= 2
			if wait >= limit {
				return limit
			}
		}
		return min(wait, limit)
	}
}

// permanent marks an error no retry can fix.
type permanent struct{ err error }

func (p permanent) Error() string { return p.err.Error() }
func (p permanent) Unwrap() error { return p.err }

// Permanent marks err as one no retry can fix: the job is dropped (without
// its give-up) instead of retried.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return permanent{err}
}

// IsPermanent reports whether err, or an error it wraps, came from
// Permanent or is ErrUnknownKind.
func IsPermanent(err error) bool {
	var p permanent
	return errors.As(err, &p) || errors.Is(err, ErrUnknownKind)
}

// Handle runs a delivered job and decides what the carrier does with it:
// done on success, a drop for a permanent error, a retry after the policy's
// backoff, and once the attempts are spent (or the carrier delivers it
// dead) the kind's give-up and a drop.
func (r *Registry) Handle(ctx context.Context, d Delivery) Result {
	attempt := max(d.Attempt, 1)
	if d.Dead {
		return r.giveUp(ctx, d.Envelope, attempt, nil)
	}
	err := r.Dispatch(ctx, d.Envelope)
	switch {
	case err == nil:
		return Result{Outcome: Done}
	case IsPermanent(err):
		telemetry.Log().ErrorContext(ctx, "job dropped", slog.String("kind", d.Kind), slog.Int("attempt", attempt), telemetry.Err(err))
		return Result{Outcome: Drop}
	case attempt >= r.policy.MaxAttempts:
		return r.giveUp(ctx, d.Envelope, attempt, err)
	default:
		after := r.policy.Backoff(attempt)
		telemetry.Log().WarnContext(ctx, "job failed, retrying",
			slog.String("kind", d.Kind), slog.Int("attempt", attempt), slog.Duration("backoff", after), telemetry.Err(err))
		return Result{Outcome: Retry, After: after}
	}
}

// giveUp runs env's give-up job, if its kind has one, and drops env. A
// give-up that fails retries env, which gives up again on its next delivery.
func (r *Registry) giveUp(ctx context.Context, env Envelope, attempt int, cause error) Result {
	attrs := []any{slog.String("kind", env.Kind), slog.Int("attempts", attempt)}
	if cause != nil {
		attrs = append(attrs, telemetry.Err(cause))
	}
	telemetry.Log().ErrorContext(ctx, "job given up", attrs...)
	e, ok := r.kinds[env.Kind]
	if !ok || e.giveUp == nil {
		return Result{Outcome: Drop}
	}
	g, err := e.giveUp(env.Args)
	if err == nil {
		g.Trace = env.Trace
		err = r.Dispatch(ctx, g)
	}
	if err != nil && !IsPermanent(err) {
		after := r.policy.Backoff(attempt)
		telemetry.Log().ErrorContext(ctx, "job give-up failed, retrying", slog.String("kind", env.Kind), slog.Duration("backoff", after), telemetry.Err(err))
		return Result{Outcome: Retry, After: after}
	}
	return Result{Outcome: Drop}
}
