package app

import (
	"log"
	"net/http"
	"net/http/pprof"

	_ "univents/contracts"
	"univents/internal/features/activities"
	"univents/internal/features/checkpoints"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/products"
	"univents/internal/features/purchases"
	"univents/internal/features/tickets"

	_ "univents/generated/docs"

	fh "github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// CreateRouter godoc
// @title Univents API
// @version 1.0.0
// @description API for managing events, editions, and tickets.
// @termsOfService https://git.trieoh.com/TrieOH/TheTree/src/branch/main/api/univents/LICENSE
// @contact.name TrieOH
// @contact.url https://trieoh.com
// @contact.email support@trieoh.com
// @license.name TSAL License
// @license.url https://git.trieoh.com/TrieOH/TheTree/src/branch/main/api/univents/LICENSE
// @host univents.com.br
// @BasePath /
// @schemes http https
// @tag.name auth
// @tag.description "Operations related to user authentication and authorization"
// @tag.name events
// @tag.description "Operations related to event management"
// @tag.name editions
// @tag.description "Operations related to edition management"
// @tag.name tickets
// @tag.description "Operations related to ticket management"
// @tag.name system
// @tag.description "System operations"
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the access token
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY
// @description API key for service-to-service authentication
func (app *Univents) CreateRouter(middlewares middlewares, handlers handlers) http.Handler {
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

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(docs.SwaggerJSON)
	})

	r.Handle("/metrics", promhttp.Handler())

	//r.With(middlewares.jwt).Get("/ws/token", deps.Security.WSAuth)
	events.Routes(r, handlers.Events, middlewares.jwt)
	editions.Routes(r, handlers.Editions, middlewares.jwt)
	tickets.Routes(r, handlers.Tickets, middlewares.jwt)
	activities.Routes(r, handlers.Activities, middlewares.jwt)
	checkpoints.Routes(r, handlers.Checkpoints, middlewares.jwt)
	products.Routes(r, handlers.Products, middlewares.jwt)
	purchases.Routes(r, handlers.Purchases, middlewares.jwt)

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
	log.Printf("univents pprof listening on :%s", port)
	if err := http.ListenAndServe(":"+port, pmux); err != nil {
		log.Fatalf("univents pprof server error: %v", err)
	}
}
