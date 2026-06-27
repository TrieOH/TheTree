package app

import (
	"log"
	"net/http"
	"net/http/pprof"
	idx "sdk/identityx"

	"Informd/generated/docs"
	"Informd/internal/features/fields"
	"Informd/internal/features/forms"
	"Informd/internal/features/namespaces"
	"Informd/internal/features/responses"
	"Informd/internal/features/steps"
	_ "Informd/models"

	_ "Informd/generated/docs"

	_ "github.com/MintzyG/fun"
	fh "github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// CreateRouter
// @title Informd
// @version 0.0.1
// @description API for managing forms.
// @termsOfService https://git.trieoh.com/TrieOH/TheTree/src/branch/main/api/informd/LICENSE
// @contact.name TrieOH
// @contact.url https://trieoh.com
// @contact.email support@trieoh.com
// @license.name TSAL License
// @license.url https://git.trieoh.com/TrieOH/TheTree/src/branch/main/api/informd/LICENSE
// @host informd.trieoh.com
// @BasePath /
// @schemes http https
// @tag.name forms
// @tag.description "Operations related to form creation"
// @tag.name api_keys
// @tag.description "Operations related to api keys"
// @tag.name namespaces
// @tag.description "Operations related to namespaces"
// @tag.name steps
// @tag.description "Operations related to steps"
// @produce json
// @consumes json
// @response 200 {object} fun.Response "Standard success response"
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

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(docs.SwaggerJSON)
	})

	r.Handle("/metrics", promhttp.Handler())

	namespaces.RegisterRoutes(r, handlers.namespaces, middlewares.jwt)
	forms.RegisterRoutes(r, handlers.forms, middlewares.anyAuth)
	steps.RegisterRoutes(r, handlers.steps, middlewares.anyAuth)
	fields.RegisterRoutes(r, handlers.fields, middlewares.anyAuth)
	responses.RegisterRoutes(r, handlers.responses)

	r.Get("/health", fh.Health(app.cfg.AppName).Handle)

	r.Mount("/idx/setup", idx.NewSetupHandler(app.idxClient))

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
	log.Printf("informd pprof listening on :%s", port)
	if err := http.ListenAndServe(":"+port, pmux); err != nil {
		log.Fatalf("informd pprof server error: %v", err)
	}
}
