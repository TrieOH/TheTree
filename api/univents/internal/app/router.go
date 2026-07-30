package app

import (
	"net/http"
	"univents/internal/features/certifications"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/products"
	"univents/internal/features/programs"
	"univents/internal/features/signatures"
	"univents/internal/features/ticket_types"

	fh "github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"riverqueue.com/riverui"
)

func (app *Univents) CreateRouter(middlewares middlewares, handlers handlers, riverUIHandler *riverui.Handler) http.Handler {
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

	events.RegisterRoutes(r, handlers.events, middlewares.jwt)
	editions.RegisterRoutes(r, handlers.editions, middlewares.jwt)
	ticket_types.RegisterRoutes(r, handlers.ticketTypes, middlewares.jwt)
	products.RegisterRoutes(r, handlers.products, middlewares.jwt)
	programs.RegisterRoutes(r, handlers.programs, middlewares.jwt)
	signatures.RegisterRoutes(r, handlers.signatures, middlewares.jwt)
	certifications.RegisterRoutes(r, handlers.certs, middlewares.jwt)

	r.Group(func(r chi.Router) {
		r.Mount("/riverui", riverUIHandler)
	})

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
