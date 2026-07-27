package app

import (
	"net/http"
	"payssage/internal/features/collectors"
	"payssage/internal/features/intents"
	"payssage/internal/features/oauth"
	"payssage/internal/features/orgs"
	"payssage/internal/features/sellers"
	"payssage/internal/features/wallets"
	"payssage/internal/features/webhook_deliveries"
	"payssage/internal/features/webhook_endpoints"
	"payssage/internal/features/webhook_events"
	"payssage/internal/features/webhooks"

	fh "github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"riverqueue.com/riverui"
)

func (app *Payssage) CreateRouter(handlers handlers, middlewares middlewares, riverUIHandler *riverui.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.logger)
	r.Use(middlewares.cors)

	r.Handle("/metrics", promhttp.Handler())

	orgs.RegisterRoutes(r, handlers.orgs, middlewares.jwtAuth)
	wallets.RegisterRoutes(r, handlers.wallets, middlewares.jwtAuth)
	collectors.RegisterRoutes(r, handlers.collectors, middlewares.jwtAuth)
	sellers.RegisterRoutes(r, handlers.sellers, middlewares.jwtAuth)
	intents.RegisterRoutes(r, handlers.intents, middlewares.jwtAuth)
	oauth.RegisterRoutes(r, handlers.oauth, middlewares.jwtAuth)
	webhooks.RegisterRoutes(r, handlers.webhooks)
	webhook_endpoints.RegisterRoutes(r, handlers.endpoints, middlewares.jwtAuth)
	webhook_events.RegisterRoutes(r, handlers.events, middlewares.jwtAuth)
	webhook_deliveries.RegisterRoutes(r, handlers.deliveries, middlewares.jwtAuth)

	r.Group(func(r chi.Router) {
		r.Use(basicAuth)
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
