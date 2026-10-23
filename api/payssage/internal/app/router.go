package app

import (
	"net/http"

	libauthz "lib/authz"
	"lib/errx"
	"lib/httpserver"
	"lib/jobs/httpjobs"
	spec "payssage"
	"payssage/internal/handlers"
	"payssage/internal/handlers/webhooks"
	"payssage/internal/openapi"

	"github.com/go-chi/chi/v5"
)

// CreateRouter builds the API router. events, when non-nil, is the job
// events endpoint the queue and scheduler deliver to; it is outside the spec
// and its auth chains, so the deploy must keep it unreachable from the
// internet (the Cloudflare Worker refuses it on public requests).
func (app *Payssage) CreateRouter(primitives libauthz.Primitives, h *handlers.Server, events http.Handler) http.Handler {
	resolver, err := libauthz.NewResolver(spec.OpenAPISpec, primitives, libauthz.Options{})
	errx.Exit(err, "resolve auth chains")
	return httpserver.NewRouter(httpserver.Config{
		AppName:            app.cfg.AppName,
		CorsAllowedOrigins: app.cfg.AllowedOrigins,
		CorsAllowedHeaders: app.cfg.AllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		RateLimit:          app.cfg.ToRateLimit(),
		Routes: func(r *chi.Mux) {
			mountStrict(r, h, resolver.Chains())

			if events != nil {
				r.Method(http.MethodPost, httpjobs.Path, events)
			}
		},
	})
}

// mountStrict is the backend's one strict-server mount point: the generated
// package's constructors cross the seam as values, and every piece of
// policy — middleware order, fail-closed dispatch, unified param-binding
// error mapping — lives in httpserver.MountStrict. The raw-request capture
// middleware runs first so the provider webhook receive can verify
// signatures against the exact body bytes.
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
