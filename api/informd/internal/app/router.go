package app

import (
	"Informd/internal/features/fields"
	"Informd/internal/features/forms"
	"Informd/internal/features/namespaces"
	"Informd/internal/features/responses"
	"Informd/internal/features/steps"
	"net/http"

	fh "github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func (app *Informd) CreateRouter(handlers handlers, middlewares middlewares) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.realIP)
	r.Use(middlewares.requestID)
	r.Use(middlewares.logger)
	r.Use(middlewares.metrics)
	r.Use(middlewares.recover)
	r.Use(middlewares.timeout)
	r.Use(middlewares.bodySize)
	r.Use(middlewares.ratelimit)
	r.Use(middlewares.cors)

	r.Handle("/metrics", promhttp.Handler())

	namespaces.RegisterRoutes(r, handlers.namespaces, middlewares.jwt)
	forms.RegisterRoutes(r, handlers.forms, middlewares.anyAuth)
	steps.RegisterRoutes(r, handlers.steps, middlewares.anyAuth)
	fields.RegisterRoutes(r, handlers.fields, middlewares.anyAuth)
	responses.RegisterRoutes(r, handlers.responses)

	r.Get("/health", fh.Health(app.cfg.AppName).Handle)

	return otelhttp.NewHandler(r, "http.server",
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/health"
		}),
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/metrics"
		}),
	)
}
