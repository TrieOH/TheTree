package app

import (
	"net/http"

	spec "IdentityX"
	"IdentityX/internal/handlers"
	"IdentityX/internal/openapi"
	libauthz "lib/authz"
	"lib/errx"
	"lib/httpserver"

	"github.com/go-chi/chi/v5"
	"riverqueue.com/riverui"
)

func (app *IdentityX) CreateRouter(primitives libauthz.Primitives, h *handlers.Server, riverUIHandler *riverui.Handler) http.Handler {
	// The setup guard and its op list are validated against the spec at
	// construction; a mismatch fails boot, never production. Platform-vs-
	// project scope is also a chain concern, derived from each operation's
	// x-scope annotation: the resolver validates every declared scope
	// against the registered checkers (a miss fails boot) and runs the
	// scope middleware after authn, so a forgotten annotation cannot
	// silently widen the surface.
	resolver, err := libauthz.NewResolver(spec.OpenAPISpec, primitives, libauthz.Options{
		SetupGuard:     setupGuard(),
		SkipSetupGuard: []string{"getSetup", "postSetup"},
	})
	errx.Exit(err, "resolve auth chains")
	return httpserver.NewRouter(httpserver.Config{
		AppName:            app.cfg.AppName,
		CorsAllowedOrigins: app.cfg.CorsAllowedOrigins,
		CorsAllowedHeaders: app.cfg.CorsAllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		Routes: func(r *chi.Mux) {
			mountStrict(r, h, resolver.Chains())

			// nil in tests; wired in run via initRiver
			if riverUIHandler != nil {
				r.Group(func(r chi.Router) {
					// River's job dashboard is an ops surface, not a spec
					// operation: the auth chains do not cover it, so it is
					// gated with the shared basic auth (SIMPLE_AUTH_* env).
					r.Use(httpserver.BasicAuth)
					r.Mount("/riverui", riverUIHandler)
				})
			}
		},
	})
}

// mountStrict is the backend's one strict-server mount point: the generated
// package's constructors cross the seam as values, and every piece of
// policy — middleware order, fail-closed dispatch, unified param-binding
// error mapping — lives in httpserver.MountStrict.
func mountStrict(r *chi.Mux, h openapi.StrictServerInterface, chains map[string][]func(http.Handler) http.Handler) {
	httpserver.MountStrict[openapi.StrictHandlerFunc](
		h,
		chains,
		openapi.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  httpserver.StrictRequestErrorHandler(),
			ResponseErrorHandlerFunc: httpserver.StrictResponseErrorHandler(),
		},
		openapi.NewStrictHandlerWithOptions,
		openapi.ChiServerOptions{
			BaseRouter:       r,
			ErrorHandlerFunc: httpserver.ParamBindingErrorHandler(),
		},
		openapi.HandlerWithOptions,
	)
}
