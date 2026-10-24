// Package backend is the kit every Backend boots from (ADR-0014). A Backend
// declares what it is — name, spec, server settings, the shared modules it
// uses, and a Build function — and Run does the rest: the process lifecycle
// (through httpserver.Boot), the modules in the listed order, the router
// frame, Build, and the shutdown of every module in reverse order.
//
//	errx.Exit(backend.Run(backend.App{
//		Name:    "payssage",
//		Spec:    spec.OpenAPISpec,
//		Server:  cfg.Server(),
//		Modules: []backend.Module{postgres.Module(...), identityx.Module(...), jobs.Module(...)},
//		Build:   build(cfg),
//	}), "payssage exited with error")
//
// There is no dependency sorting and no reflection: the module list is the
// boot order, and Deps is a plain struct.
package backend

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"lib/authz"
	"lib/database"
	"lib/email"
	"lib/httpserver"
	libjobs "lib/jobs"
	"lib/jobs/httpjobs"

	idx "sdk/identityx"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App is what a Backend declares about itself.
type App struct {
	Name    string
	Spec    []byte
	Server  Server
	Modules []Module
	Build   BuildFunc
}

// Server is the process-level HTTP settings.
type Server struct {
	Port        string
	ProfilePort string
	// CorsAllowedOrigins and CorsAllowedHeaders are comma-separated.
	CorsAllowedOrigins string
	CorsAllowedHeaders string
	RateLimit          httpserver.RateLimit
}

// BuildFunc wires the Backend's own code: repos, services, handlers,
// Backend-only infrastructure and boot checks. It receives the router frame
// (CORS, rate limit, /health, /docs already in place) and mounts its strict
// server and any raw routes on it. An error stops the boot.
type BuildFunc func(ctx context.Context, d *Deps, r chi.Router) error

// Close releases what a module opened. Run calls it once, in reverse module
// order, after serving stops (or when a later boot step fails).
type Close func(ctx context.Context) error

// NoClose is the Close of a module that holds nothing to release.
func NoClose(context.Context) error { return nil }

// Module opens one piece of shared infrastructure into Deps.
type Module struct {
	Name string
	Open func(ctx context.Context, d *Deps) (Close, error)
}

// Deps is everything the modules opened. A field stays nil when its module
// is not listed.
type Deps struct {
	// postgres
	DB *pgxpool.Pool
	Tx database.TxRunner

	// identityx: the SDK client and the auth primitives built on it.
	IDX *IdentityX

	// email: the SMTP client (it satisfies email.Sender).
	Email *email.Client

	// jobs: enqueue right away; hand over the registry with Serve in Build.
	Jobs Jobs

	spec []byte
}

// IdentityX is what the identityx module provides: the SDK client and the
// auth primitives every IdentityX consumer authenticates requests with.
type IdentityX struct {
	Client *idx.Client
	Auth   authz.Primitives
}

// Jobs is the background-work carrier the jobs module provides. Enqueue
// works from the moment the module opens; Serve hands over the registry
// (built in Build, from the Backend's repos and services) and starts
// delivery — the in-process runner, or the /_jobs endpoint the kit mounts.
type Jobs interface {
	libjobs.DelayedEnqueuer
	Serve(registry *libjobs.Registry, periodic ...libjobs.Periodic)
	// Served reports whether Serve ran; the kit refuses to start otherwise.
	Served() bool
	// Handler is the /_jobs endpoint, nil when jobs run in-process.
	Handler() http.Handler
}

// Chains resolves every spec operation's auth chain from the operation's
// security, enforced by primitives: d.IDX.Auth for IdentityX's consumers,
// IdentityX's own primitives for IdentityX itself. Pass the result to
// httpserver.MountStrict.
func (d *Deps) Chains(primitives authz.Primitives) (map[string][]func(http.Handler) http.Handler, error) {
	resolver, err := authz.NewResolver(d.spec, primitives, authz.Options{})
	if err != nil {
		return nil, fmt.Errorf("resolve auth chains: %w", err)
	}
	return resolver.Chains(), nil
}

// Run boots the Backend and serves until SIGTERM/SIGINT.
func Run(app App) error {
	if app.Build == nil {
		return fmt.Errorf("%s: backend.App.Build is required", app.Name)
	}
	return httpserver.Boot(httpserver.Config{
		AppName:            app.Name,
		Port:               app.Server.Port,
		ProfilePort:        app.Server.ProfilePort,
		CorsAllowedOrigins: app.Server.CorsAllowedOrigins,
		CorsAllowedHeaders: app.Server.CorsAllowedHeaders,
		OpenAPISpec:        app.Spec,
		RateLimit:          app.Server.RateLimit,
	}, func(ctx context.Context) (http.Handler, func(context.Context) error, error) {
		return start(ctx, app)
	})
}

// start opens the modules, builds the router through Build, and returns the
// handler with the shutdown that closes every module in reverse order.
func start(ctx context.Context, app App) (http.Handler, func(context.Context) error, error) {
	d := &Deps{spec: app.Spec}
	closeAll, err := openModules(ctx, d, app.Modules)
	if err != nil {
		return nil, nil, err
	}

	var buildErr error
	handler := httpserver.NewRouter(httpserver.Config{
		AppName:            app.Name,
		CorsAllowedOrigins: app.Server.CorsAllowedOrigins,
		CorsAllowedHeaders: app.Server.CorsAllowedHeaders,
		OpenAPISpec:        app.Spec,
		RateLimit:          app.Server.RateLimit,
		Routes: func(r *chi.Mux) {
			buildErr = app.Build(ctx, d, r)
			if buildErr == nil {
				buildErr = mountJobs(d, r)
			}
		},
	})
	if buildErr != nil {
		_ = closeAll(ctx)
		return nil, nil, fmt.Errorf("build: %w", buildErr)
	}
	return handler, closeAll, nil
}

// openModules opens each module in the listed order. If one fails, the ones
// already open are closed (in reverse) before the error returns. The
// returned Close closes them all, in reverse, joining their errors.
func openModules(ctx context.Context, d *Deps, modules []Module) (func(context.Context) error, error) {
	var closes []namedClose
	closeAll := func(ctx context.Context) error {
		var errs []error
		for i := len(closes) - 1; i >= 0; i-- {
			err := closes[i].close(ctx)
			if err != nil {
				errs = append(errs, fmt.Errorf("close %s: %w", closes[i].name, err))
			}
		}
		return errors.Join(errs...)
	}
	for _, m := range modules {
		closeFn, err := m.Open(ctx, d)
		if err != nil {
			_ = closeAll(ctx)
			return nil, fmt.Errorf("open %s: %w", m.Name, err)
		}
		if closeFn != nil {
			closes = append(closes, namedClose{name: m.Name, close: closeFn})
		}
	}
	return closeAll, nil
}

type namedClose struct {
	name  string
	close Close
}

// mountJobs mounts the jobs endpoint when the jobs module is listed. Build
// must have called Serve: a Backend that lists the module and never hands
// over its registry would silently drop every job.
func mountJobs(d *Deps, r chi.Router) error {
	if d.Jobs == nil {
		return nil
	}
	if !d.Jobs.Served() {
		return errors.New("jobs module listed but Build never called d.Jobs.Serve")
	}
	if h := d.Jobs.Handler(); h != nil {
		r.Method(http.MethodPost, httpjobs.Path, h)
	}
	return nil
}
