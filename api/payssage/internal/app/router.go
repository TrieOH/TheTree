package app

import (
	"log"
	"net/http"
	"net/http/pprof"
	"payssage/generated/docs"
	"payssage/internal/features/collectors"
	"payssage/internal/features/intents"
	"payssage/internal/features/oauth"
	"payssage/internal/features/orgs"
	"payssage/internal/features/sellers"
	"payssage/internal/features/wallets"
	"payssage/internal/features/webhooks"

	fh "github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"riverqueue.com/riverui"
)

// CreateRouter godoc
// CreateRouter creates a new Chi router and registers all the routes.
// @title Univents API
// @version 1.0.0
// @description API for managing events, editions, and tickets.
// @termsOfService https://github.com/Univents/Univents/blob/main/LICENSE
// @contact.name Univents Team
// @contact.url https://github.com/Univents
// @contact.email support@univents.io
// @license.name MIT License
// @license.url https://github.com/Univents/Univents/blob/main/LICENSE
// @host localhost:8080
// @BasePath /
// @schemes http https
// @tag.name organizations
// @tag.description "Operations related to organizations"
// @tag.name wallets
// @tag.description "Operations related to wallets"
// @produce json
// @consumes json
// @response 200 {object} object "Standard success response"
// @response 400 {object} fun.Response "Standard error response for bad requests"
// @response 401 {object} fun.Response "Standard error response for unauthorized requests"
// @response 403 {object} fun.Response "Standard error response for forbidden requests"
// @response 404 {object} fun.Response "Standard error response for not found errors"
// @response 413 {object} fun.Response "Standard error response for payload too large 1MB"
// @response 429 {object} fun.Response "Standard error response for too many requests"
// @response 500 {object} fun.Response "Standard error response for internal server errors"
// @securityDefinitions.apikey Cookie
// @in header
// @name Cookie
// @description Type "Cookie" followed by a cookie in the format "access_token=xxx; refresh_token=yyy"
func (app *Payssage) CreateRouter(handlers handlers, middlewares middlewares, riverUIHandler *riverui.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.logger)
	r.Use(middlewares.cors)

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(docs.SwaggerJSON)
	})

	r.Handle("/metrics", promhttp.Handler())

	// Routes
	orgs.RegisterRoutes(r, handlers.orgs, middlewares.jwtAuth)
	wallets.RegisterRoutes(r, handlers.wallets, middlewares.jwtAuth)
	collectors.RegisterRoutes(r, handlers.collectors, middlewares.jwtAuth)
	sellers.RegisterRoutes(r, handlers.sellers, middlewares.jwtAuth)
	intents.RegisterRoutes(r, handlers.intents, middlewares.jwtAuth)
	oauth.RegisterRoutes(r, handlers.oauth, middlewares.jwtAuth)
	webhooks.RegisterRoutes(r, handlers.webhooks)

	// River UI — internal ops dashboard, gated behind basic auth
	// (SIMPLE_AUTH_USER / SIMPLE_AUTH_PASS), not tenant-scoped JWT auth.
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

func servePprof(port string) {
	pmux := http.NewServeMux()
	pmux.HandleFunc("/debug/pprof/", pprof.Index)
	pmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	log.Printf("payssage pprof listening on :%s", port)
	if err := http.ListenAndServe(":"+port, pmux); err != nil {
		log.Fatalf("payssage pprof server error: %v", err)
	}
}
