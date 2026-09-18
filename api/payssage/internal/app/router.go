package app

import (
	"net/http"

	libauthz "lib/authz"
	"lib/errx"
	"lib/httpserver"
	spec "payssage"
	"payssage/internal/handlers"
	"payssage/internal/handlers/webhooks"
	"payssage/internal/openapi"

	"github.com/go-chi/chi/v5"
	"riverqueue.com/riverui"
)

func (app *Payssage) CreateRouter(primitives libauthz.Primitives, h *handlers.Server, riverUIHandler *riverui.Handler) http.Handler {
	resolver, err := libauthz.NewResolver(spec.OpenAPISpec, primitives, libauthz.Options{})
	errx.Exit(err, "resolve auth chains")
	return httpserver.NewRouter(httpserver.Config{
		AppName:            app.cfg.AppName,
		CorsAllowedOrigins: app.cfg.CorsAllowedOrigins,
		CorsAllowedHeaders: app.cfg.CorsAllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		Routes: func(r *chi.Mux) {
			mountStrict(r, h, resolver.Chains())

			r.Group(func(r chi.Router) {
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
