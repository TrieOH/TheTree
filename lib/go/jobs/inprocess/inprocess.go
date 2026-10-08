// Package inprocess carries jobs inside a long-lived process: a bounded
// queue drained by a worker pool, with retries and backoff, plus a ticker
// per periodic job. It holds work in memory only — work still queued when
// the process dies is lost — which is the accepted trade for hosted mode.
package inprocess

import (
	"context"
	"errors"
	"sync"
	"time"

	"lib/jobs"
	"lib/telemetry"

	"go.uber.org/zap"
)

// Options tunes the runner; zero values take the defaults.
type Options struct {
	Workers     int           // default 4
	QueueSize   int           // default 1024
	MaxAttempts int           // default 5
	BaseBackoff time.Duration // default 1s, doubled per attempt
}

func (o Options) withDefaults() Options {
	if o.Workers <= 0 {
		o.Workers = 4
	}
	if o.QueueSize <= 0 {
		o.QueueSize = 1024
	}
	if o.MaxAttempts <= 0 {
		o.MaxAttempts = 5
	}
	if o.BaseBackoff <= 0 {
		o.BaseBackoff = time.Second
	}
	return o
}

// ErrStopped is returned by Enqueue once the runner is stopping.
var ErrStopped = errors.New("inprocess: runner stopped")

// ErrQueueFull is returned by Enqueue when the queue is at capacity.
var ErrQueueFull = errors.New("inprocess: queue full")

type Runner struct {
	registry *jobs.Registry
	opts     Options

	mu       sync.RWMutex
	stopped  bool
	queue    chan jobs.Envelope
	tickStop chan struct{}

	// ctx is what handlers run under. It is cancelled only when Stop's
	// drain window runs out, so queued work gets to finish normally.
	ctx     context.Context
	cancel  context.CancelFunc
	workers sync.WaitGroup
	tickers sync.WaitGroup
}

var _ jobs.Enqueuer = (*Runner)(nil)

func New(registry *jobs.Registry, opts Options) *Runner {
	opts = opts.withDefaults()
	ctx, cancel := context.WithCancel(context.Background())
	return &Runner{
		registry: registry,
		opts:     opts,
		queue:    make(chan jobs.Envelope, opts.QueueSize),
		tickStop: make(chan struct{}),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Enqueue queues a job without blocking the request path.
func (r *Runner) Enqueue(_ context.Context, job jobs.Job) error {
	env, err := jobs.Encode(job)
	if err != nil {
		return err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.stopped {
		return ErrStopped
	}
	select {
	case r.queue <- env:
		return nil
	default:
		return ErrQueueFull
	}
}

// Start launches the workers and one ticker per periodic job.
func (r *Runner) Start(periodic ...jobs.Periodic) {
	for range r.opts.Workers {
		r.workers.Add(1)
		go r.work()
	}
	for _, p := range periodic {
		r.tickers.Add(1)
		go r.tick(p)
	}
}

// Stop stops the tickers, refuses new work, and lets queued work drain
// until ctx is done; then it cancels whatever is still running and waits
// for it to return.
func (r *Runner) Stop(ctx context.Context) error {
	r.mu.Lock()
	if r.stopped {
		r.mu.Unlock()
		return nil
	}
	r.stopped = true
	close(r.tickStop)
	close(r.queue)
	r.mu.Unlock()

	r.tickers.Wait()
	drained := make(chan struct{})
	go func() {
		r.workers.Wait()
		close(drained)
	}()
	select {
	case <-drained:
		r.cancel()
		return nil
	case <-ctx.Done():
		r.cancel()
		<-drained
		return ctx.Err()
	}
}

func (r *Runner) work() {
	defer r.workers.Done()
	for env := range r.queue {
		r.run(env)
	}
}

// run dispatches env, retrying failures with exponential backoff. A job
// that exhausts its attempts is logged and dropped.
func (r *Runner) run(env jobs.Envelope) {
	backoff := r.opts.BaseBackoff
	for attempt := 1; ; attempt++ {
		err := r.registry.Dispatch(r.ctx, env)
		if err == nil {
			return
		}
		if errors.Is(err, jobs.ErrUnknownKind) || attempt >= r.opts.MaxAttempts || r.ctx.Err() != nil {
			telemetry.Log().Error("job failed",
				zap.String("kind", env.Kind), zap.Int("attempts", attempt), zap.Error(err))
			return
		}
		telemetry.Log().Warn("job failed, retrying",
			zap.String("kind", env.Kind), zap.Int("attempt", attempt), zap.Duration("backoff", backoff), zap.Error(err))
		select {
		case <-time.After(backoff):
			backoff *= 2
		case <-r.ctx.Done():
		}
	}
}

func (r *Runner) tick(p jobs.Periodic) {
	defer r.tickers.Done()
	if p.RunOnStart {
		r.enqueuePeriodic(p)
	}
	ticker := time.NewTicker(p.Every)
	defer ticker.Stop()
	for {
		select {
		case <-r.tickStop:
			return
		case <-ticker.C:
			r.enqueuePeriodic(p)
		}
	}
}

func (r *Runner) enqueuePeriodic(p jobs.Periodic) {
	err := r.Enqueue(r.ctx, p.Job)
	if err != nil && !errors.Is(err, ErrStopped) {
		telemetry.Log().Error("periodic job not enqueued", zap.String("kind", p.Job.Kind()), zap.Error(err))
	}
}
