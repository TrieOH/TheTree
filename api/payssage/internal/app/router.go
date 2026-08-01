package app

import (
	"net/http"

	spec "payssage"
	collectorsHandlers "payssage/internal/features/collectors/handlers"
	intentsHandlers "payssage/internal/features/intents/handlers"
	oauthHandlers "payssage/internal/features/oauth/handlers"
	orgsHandlers "payssage/internal/features/orgs/handlers"
	sellersHandlers "payssage/internal/features/sellers/handlers"
	walletsHandlers "payssage/internal/features/wallets/handlers"
	webhookDeliveriesHandlers "payssage/internal/features/webhook_deliveries/handlers"
	webhookEndpointsHandlers "payssage/internal/features/webhook_endpoints/handlers"
	webhookEventsHandlers "payssage/internal/features/webhook_events/handlers"
	webhooksHandlers "payssage/internal/features/webhooks/handlers"

	"lib/httpserver"

	"github.com/go-chi/chi/v5"
	"riverqueue.com/riverui"
)

func (app *Payssage) CreateRouter(handlers handlers, middlewares middlewares, riverUIHandler *riverui.Handler) http.Handler {
	return httpserver.NewRouter(httpserver.Config{
		AppName:     app.cfg.AppName,
		OpenAPISpec: spec.OpenAPISpec,
		Routes: func(r *chi.Mux) {
			registerRoutes(r, middlewares, handlers, riverUIHandler)
		},
	})
}

// registerRoutes wires every feature's routes onto r. Kept package-level so
// the router-parity test can walk the same registration the app serves.
func registerRoutes(r *chi.Mux, middlewares middlewares, handlers handlers, riverUIHandler *riverui.Handler) {
	orgsHandlers.RegisterRoutes(r, handlers.orgs, middlewares.jwtAuth)
	walletsHandlers.RegisterRoutes(r, handlers.wallets, middlewares.jwtAuth)
	collectorsHandlers.RegisterRoutes(r, handlers.collectors, middlewares.jwtAuth)
	sellersHandlers.RegisterRoutes(r, handlers.sellers, middlewares.jwtAuth)
	intentsHandlers.RegisterRoutes(r, handlers.intents, middlewares.jwtAuth)
	oauthHandlers.RegisterRoutes(r, handlers.oauth, middlewares.jwtAuth)
	webhooksHandlers.RegisterRoutes(r, handlers.webhooks)
	webhookEndpointsHandlers.RegisterRoutes(r, handlers.endpoints, middlewares.jwtAuth)
	webhookEventsHandlers.RegisterRoutes(r, handlers.events, middlewares.jwtAuth)
	webhookDeliveriesHandlers.RegisterRoutes(r, handlers.deliveries, middlewares.jwtAuth)

	r.Group(func(r chi.Router) {
		r.Use(basicAuth)
		r.Mount("/riverui", riverUIHandler)
	})
}
