package app

import (
	"net/http"

	"lib/authz"
	"lib/errx"
	"lib/httpserver"
	spec "univents"
	"univents/internal/handlers"
	"univents/internal/handlers/webhooks"
	"univents/internal/openapi"

	"github.com/go-chi/chi/v5"
	"riverqueue.com/riverui"
)

func (app *Univents) CreateRouter(primitives authz.Primitives, h *handlers.Server, riverUIHandler *riverui.Handler) http.Handler {
	resolver, err := authz.NewResolver(spec.OpenAPISpec, primitives, authz.Options{})
	errx.Exit(err, "resolve auth chains")
	return httpserver.NewRouter(httpserver.Config{
		AppName:            app.cfg.AppName,
		CorsAllowedOrigins: app.cfg.CorsAllowedOrigins,
		CorsAllowedHeaders: app.cfg.CorsAllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		SkipLogPrefixes:    []string{"/admin/asynq"},
		Routes: func(r *chi.Mux) {
			mountStrict(r, h, resolver.Chains())

			// Raw realtime routes (split 6) — deliberately outside the strict
			// handler: the WS handshake cannot carry Authorization headers (the
			// one-time query token is the auth), and SSE must stream without the
			// fun/validate envelope machinery buffering the body. Both are
			// documented in the spec's description block, not as spec ops.
			r.Get("/editions/{edition_id}/store/stream", h.ServeStoreStream)
			r.Handle("/ws", http.HandlerFunc(h.ServeWS))

			r.Group(func(r chi.Router) {
				// River's job dashboard is an ops surface, not a spec
				// operation: the auth chains do not cover it, so it is
				// gated with the shared basic auth (SIMPLE_AUTH_* env).
				r.Use(httpserver.BasicAuth)
				r.Mount("/riverui", riverUIHandler)
			})
		},
	})
}

// mountStrict is the backend's one strict-server mount point: the generated
// package's constructors cross the seam as values, and every piece of
// policy — middleware order, fail-closed dispatch, unified param-binding
// error mapping — lives in httpserver.MountStrict. The raw-request capture
// middleware runs first so the Payssage webhook can verify the signature
// against the exact body bytes payssage POSTed (the strict server decodes
// the body afterwards; D2). Path-scoped to /webhooks/ inside the middleware.
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
			Middlewares:      []openapi.MiddlewareFunc{webhooks.RawRequestMiddleware},
			ErrorHandlerFunc: httpserver.ParamBindingErrorHandler(),
		},
		openapi.HandlerWithOptions,
	)
}
