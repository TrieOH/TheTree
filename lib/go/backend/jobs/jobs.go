// Package jobs is the kit module for background work (ADR-0011, ADR-0013).
// The carrier is picked by the jobs URL: empty runs jobs in-process; set
// (on Cloudflare, by the Worker) POSTs them to the queue producer, and the
// queue delivers them back to the /_jobs endpoint the kit mounts.
//
// The carrier is in Deps.Jobs as soon as the module opens, so services can
// hold it; the registry depends on those services, so Build hands it over
// later with Serve.
package jobs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"lib/backend"
	libjobs "lib/jobs"
	"lib/jobs/httpjobs"
	"lib/jobs/inprocess"
)

// drainTimeout bounds how long shutdown waits for in-process jobs.
const drainTimeout = 15 * time.Second

// ErrNotServing is returned by an in-process enqueue before Build has
// called Serve: there is no runner to queue onto yet.
var ErrNotServing = errors.New("jobs: enqueue before Serve")

// Module fills Deps.Jobs with the carrier for jobsURL (JOBS_URL).
func Module(jobsURL string) backend.Module {
	return backend.Module{
		Name: "jobs",
		Open: func(_ context.Context, d *backend.Deps) (backend.Close, error) {
			c, err := newCarrier(jobsURL)
			if err != nil {
				return nil, err
			}
			d.Jobs = c
			return c.close, nil
		},
	}
}

// carrier implements backend.Jobs over httpjobs or the in-process runner.
type carrier struct {
	http *httpjobs.Enqueuer // nil when in-process

	mu      sync.RWMutex
	runner  *inprocess.Runner // in-process only, set by Serve
	handler http.Handler      // httpjobs only, set by Serve
	served  bool
}

var _ backend.Jobs = (*carrier)(nil)

func newCarrier(jobsURL string) (*carrier, error) {
	if jobsURL == "" {
		return &carrier{}, nil
	}
	e, err := httpjobs.New(jobsURL)
	if err != nil {
		return nil, fmt.Errorf("jobs carrier: %w", err)
	}
	return &carrier{http: e}, nil
}

func (c *carrier) Enqueue(ctx context.Context, job libjobs.Job) error {
	if c.http != nil {
		return c.http.Enqueue(ctx, job)
	}
	runner, err := c.inprocess()
	if err != nil {
		return err
	}
	return runner.Enqueue(ctx, job)
}

func (c *carrier) EnqueueAfter(ctx context.Context, job libjobs.Job, delay time.Duration) error {
	if c.http != nil {
		return c.http.EnqueueAfter(ctx, job, delay)
	}
	runner, err := c.inprocess()
	if err != nil {
		return err
	}
	return runner.EnqueueAfter(ctx, job, delay)
}

func (c *carrier) inprocess() (*inprocess.Runner, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.runner == nil {
		return nil, ErrNotServing
	}
	return c.runner, nil
}

// Serve hands over the registry. Over HTTP it builds the /_jobs handler
// (the deploy's scheduler fires the periodic jobs there, so periodic is
// unused); in-process it starts the runner, with a ticker per periodic job.
// Calling it twice is a wiring bug and panics.
func (c *carrier) Serve(registry *libjobs.Registry, periodic ...libjobs.Periodic) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.served {
		panic("jobs: Serve called twice")
	}
	c.served = true
	if c.http != nil {
		c.handler = httpjobs.Handler(registry)
		return
	}
	c.runner = inprocess.New(registry, inprocess.Options{})
	c.runner.Start(periodic...)
}

func (c *carrier) Served() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.served
}

func (c *carrier) Handler() http.Handler {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.handler
}

// close drains the in-process runner; the HTTP carrier holds nothing.
func (c *carrier) close(ctx context.Context) error {
	c.mu.RLock()
	runner := c.runner
	c.mu.RUnlock()
	if runner == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, drainTimeout)
	defer cancel()
	return runner.Stop(ctx)
}
