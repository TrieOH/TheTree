package river

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"lib/httpserver"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"riverqueue.com/riverui"
)

// Start creates and starts the river client in ctx. Queues default to the
// shared default queue when nil. Call Migrate beforehand; migration is the
// caller's deployment decision, not an unrequested side effect.
func Start(ctx context.Context, dbPool *pgxpool.Pool, workers *river.Workers, queues map[string]river.QueueConfig, periodicJobs []*river.PeriodicJob) (*river.Client[pgx.Tx], error) {
	client := NewClient(dbPool, workers, queues, periodicJobs)
	err := client.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("start river client: %w", err)
	}
	return client, nil
}

// Dashboard returns the riverui job-dashboard handler, started in ctx. The
// ops policy lives here, once: production mode, the /riverui prefix, and
// job args hidden by default. Mount it with MountDashboard.
func Dashboard(ctx context.Context, client *river.Client[pgx.Tx]) (*riverui.Handler, error) {
	handler, err := riverui.NewHandler(&riverui.HandlerOpts{
		DevMode:                  false,
		Endpoints:                riverui.NewEndpoints[pgx.Tx](client, nil),
		Logger:                   slog.Default(),
		Prefix:                   "/riverui",
		JobListHideArgsByDefault: true,
	})
	if err != nil {
		return nil, fmt.Errorf("create river ui handler: %w", err)
	}
	err = handler.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("start river ui handler: %w", err)
	}
	return handler, nil
}

// MountDashboard mounts the dashboard under /riverui behind the shared
// basic auth (SIMPLE_AUTH_* env). The dashboard is an ops surface, not a
// spec operation: the auth chains do not cover it.
func MountDashboard(r chi.Router, handler http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(httpserver.BasicAuth)
		r.Mount("/riverui", handler)
	})
}
