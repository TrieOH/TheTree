package app

import (
	"net/http"

	spec "Informd"
	"Informd/internal/handlers"
	"Informd/internal/openapi"
	libauthz "lib/authz"
	"lib/errx"
	"lib/httpserver"

	"github.com/go-chi/chi/v5"
)

func (app *Informd) CreateRouter(h *handlers.Server, primitives libauthz.Primitives) http.Handler {
	resolver, err := libauthz.NewResolver(spec.OpenAPISpec, primitives, libauthz.Options{})
	errx.Exit(err, "resolve auth chains")
	return httpserver.NewRouter(httpserver.Config{
		AppName:            app.cfg.AppName,
		CorsAllowedOrigins: app.cfg.AllowedOrigins,
		CorsAllowedHeaders: app.cfg.AllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		Routes: func(r *chi.Mux) {
			mountStrict(r, h, resolver.Chains())
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
