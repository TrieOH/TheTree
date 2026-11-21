package backend

// The process lifecycle (ADR-0014): validator, telemetry, the boot span,
// serving, then shutdown in reverse order on SIGTERM/SIGINT.

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"lib/openapispec"
	"lib/telemetry"
	"lib/validator"
)

// shutdownTimeout bounds how long in-flight requests get to finish after a
// stop signal.
const shutdownTimeout = 10 * time.Second

// boot runs the process: validator, tracer, start (the modules, the router,
// Build), serving, then the shutdown start returned, in reverse order, on
// SIGTERM/SIGINT. It returns an error instead of exiting: errx.Exit stays
// at the process edge (cmd/main).
func boot(app App, ops []openapispec.Operation) error {
	validator.SetupValidator()

	ctx := context.Background()
	shutdownTelemetry, err := telemetry.Start(ctx, app.Name)
	if err != nil {
		return fmt.Errorf("%s boot: %w", app.Name, err)
	}

	bootCtx, bootSpan := startBoot(ctx, app.Name)
	handler, shutdown, err := start(bootCtx, app, ops)
	if err != nil {
		failBoot(bootSpan, err)
		// Flush before exiting so the boot error reaches Grafana; bounded so
		// an unreachable collector can't hold a dying process.
		flushCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		telemetry.Shutdown(flushCtx, shutdownTelemetry, app.Name)
		cancel()
		return fmt.Errorf("%s boot: %w", app.Name, err)
	}

	if app.Server.ProfilePort != "" {
		go servePprof(app.Server.ProfilePort, app.Name)
	}

	// SIGTERM is how the process is stopped (docker stop, and Cloudflare
	// when a container sleeps or is replaced), so stopping gracefully is what
	// lets the shutdown drain background work and telemetry flush its last
	// batch.
	srv := newServer(handler, app.Server.Port)
	stopCtx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	go func() {
		<-stopCtx.Done()
		drainCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		_ = srv.Shutdown(drainCtx)
	}()

	// The boot ends once the port is bound: from then on the platform's
	// readiness check gets an answer.
	var serveErr error
	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", srv.Addr)
	if err != nil {
		failBoot(bootSpan, err)
		serveErr = err
	} else {
		bootSpan.End()
		if telemetry.Traces.Enabled() {
			go probeEgress(bootCtx)
		}
		telemetry.Log().InfoContext(ctx, app.Name+" listening", slog.String("port", app.Server.Port))
		serveErr = srv.Serve(ln)
	}
	if errors.Is(serveErr, http.ErrServerClosed) {
		serveErr = nil
	}

	err = shutdown(ctx)
	if err != nil {
		telemetry.Log().ErrorContext(ctx, app.Name+" shutdown failed", telemetry.Err(err))
	}
	telemetry.Shutdown(ctx, shutdownTelemetry, app.Name)
	return serveErr
}

// newServer is the HTTP server for handler on port, with the timeouts every
// Backend serves with.
func newServer(handler http.Handler, port string) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// servePprof serves the pprof endpoints on port; Run starts it when
// PROFILE_PORT is set.
func servePprof(port, appName string) {
	pmux := http.NewServeMux()
	pmux.HandleFunc("/debug/pprof/", pprof.Index)
	pmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      pmux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	ctx := context.Background()
	telemetry.Log().InfoContext(ctx, appName+" pprof listening", slog.String("port", port))
	// pprof is a debugging aid: when it can't serve, the API carries on.
	err := srv.ListenAndServe()
	if err != nil {
		telemetry.Log().ErrorContext(ctx, appName+" pprof server stopped", telemetry.Err(err))
	}
}
