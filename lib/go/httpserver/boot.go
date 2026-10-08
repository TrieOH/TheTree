package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"lib/telemetry"
)

// StartFunc opens everything the backend varies on — its storage (a Postgres
// pool today; a Mongo client or nothing tomorrow — Boot never learns what a
// database is), its clients, its workers and jobs — and returns the root
// handler to serve plus the shutdown to run once serving has stopped. The
// shutdown runs before the tracer shuts down, so spans from stopping
// workers are kept; stopping the backend's components in the right order is
// the backend's concern, localized here.
type StartFunc func(ctx context.Context) (handler http.Handler, shutdown func(context.Context) error, err error)

// Boot is the shared process lifecycle every Backend runs through: FUN
// runtime, tracer, the backend's components (via start), serving, then
// shutdown in reverse order on SIGTERM/SIGINT. The sequence lives here once — the per-backend
// "call this before X" contracts it replaces (FUN before any request,
// tracer before any span, constraint messages before constraint
// validation, tx runner before the first query) are either owned here,
// became explicit data (the constraint registry travels into SetupDB as an
// argument), or became explicit dependencies (the TxRunner threads into
// the services that open transactions).
//
// Boot returns an error instead of exiting: errx.Exit stays at the process
// edge (cmd/main), and the boot path stays testable through this interface.
func Boot(cfg Config, start StartFunc) error {
	SetupFUN(cfg.AppName)

	ctx := context.Background()
	shutdownTelemetry, err := telemetry.Start(ctx, cfg.AppName)
	if err != nil {
		return fmt.Errorf("%s boot: %w", cfg.AppName, err)
	}

	handler, shutdown, err := start(ctx)
	if err != nil {
		telemetry.Shutdown(ctx, shutdownTelemetry, cfg.AppName)
		return fmt.Errorf("%s boot: %w", cfg.AppName, err)
	}

	if cfg.ProfilePort != "" {
		go servePprof(cfg.ProfilePort, cfg.AppName)
	}

	// SIGTERM is how both runtimes stop the process (docker stop, and Lambda
	// when it reclaims an instance — the Lambda Web Adapter extension makes
	// Lambda deliver it), so stopping gracefully is what lets the shutdown
	// hook drain background work and telemetry flush its last batch.
	srv := newServer(handler, cfg.Port)
	stopCtx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	go func() {
		<-stopCtx.Done()
		drainCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		_ = srv.Shutdown(drainCtx)
	}()

	log.Printf("%s listening on :%s", cfg.AppName, cfg.Port)
	serveErr := srv.ListenAndServe()
	if errors.Is(serveErr, http.ErrServerClosed) {
		serveErr = nil
	}

	if shutdown != nil {
		err := shutdown(ctx)
		if err != nil {
			log.Printf("%s shutdown: %v", cfg.AppName, err)
		}
	}
	telemetry.Shutdown(ctx, shutdownTelemetry, cfg.AppName)
	return serveErr
}

// shutdownTimeout bounds how long in-flight requests get to finish after a
// stop signal.
const shutdownTimeout = 10 * time.Second
