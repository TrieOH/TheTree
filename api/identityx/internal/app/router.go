package app

import (
	"IdentityX/generated/docs"
	"IdentityX/internal/features/actors"
	"IdentityX/internal/features/api_keys"
	"log"
	"net/http"
	"net/http/pprof"

	"IdentityX/internal/features/authn"
	"IdentityX/internal/features/organizations"
	"IdentityX/internal/features/projects"
	_ "IdentityX/models"

	fh "github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// CreateRouter godoc
// CreateRouter creates a new Chi router and registers all the routes.
// @title IdentityX API
// @version 0.19.0-alpha
// @description This is the API for the IdentityX, an Identity Provider (IdP) service. It provides user authentication, and project management functionalities.
// @termsOfService https://git.trieoh.com/TrieOH/TheTree/blob/main/api/identityx/LICENSE
// @contact.name TrieOH
// @contact.url https://github.com/TrieOH
// @contact.email contact@trieoh.com
// @license.name TSAL 1.2 License
// @license.url https://git.trieoh.com/TrieOH/TheTree/blob/main/api/identityx/LICENSE
// @host identityx.com.br
// @BasePath /
// @schemes http https
// @tag.name authn
// @tag.description "Operations related to user authentication"
// @tag.name organizations
// @tag.description "Operations related to organization management"
// @tag.name projects
// @tag.description "Operations related to project management"
// @tag.name apikeys
// @tag.description "Operations related to api key management"
// @produce json
// @consumes json
// @response 200 {object} fun.Response "Standard success response"
// @response 201 {object} fun.Response "Standard creation response"
// @response 400 {object} fun.Response "Standard error response for bad requests"
// @response 401 {object} fun.Response "Standard error response for unauthorized requests"
// @response 403 {object} fun.Response "Standard error response for forbidden requests"
// @response 404 {object} fun.Response "Standard error response for not found errors"
// @response 413 {object} fun.Response "Standard error response for payload too large 1MB"
// @response 429 {object} fun.Response "Standard error response for too many requests"
// @response 500 {object} fun.Response "Standard error response for internal server errors"
// @response 503 {object} fun.Response "Standard error response for service unavailable"
func (app *IdentityX) CreateRouter(middlewares middlewares, handlers handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(otelchi.Middleware(app.cfg.AppName,
		otelchi.WithChiRoutes(r),
		otelchi.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/health" && r.URL.Path != "/metrics"
		}),
	))

	//r.Use(deps.RealIP)
	//r.Use(deps.RequestID)
	r.Use(middlewares.logger)
	r.Use(middlewares.metrics)
	//r.Use(deps.Recover)
	//r.Use(deps.Timeout)
	//r.Use(deps.BodySize)
	//r.Use(deps.RateLimit)
	r.Use(middlewares.cors)

	//endpoints := riverui.NewEndpoints(app.river, nil)
	//
	//handler, err := riverui.NewHandler(&riverui.HandlerOpts{
	//	Endpoints: endpoints,
	//	Logger:    slog.Default(),
	//	Prefix:    "/riverui",
	//})
	//if err != nil {
	//	errx.Exit(err, "failed to create river handler")
	//}
	//err = handler.Start(context.Background())
	//if err != nil {
	//	errx.Exit(err, "failed to start river handler")
	//}
	//r.Mount("/riverui", handler)

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(docs.SwaggerJSON)
	})

	r.Handle("/metrics", promhttp.Handler())

	actors.RegisterRoutes(r, handlers.Actors, middlewares.jwtAuth, middlewares.clientOnly)
	api_keys.RegisterRoutes(r, handlers.ApiKeys, middlewares.jwtAuth, middlewares.clientOnly)
	authn.RegisterRoutes(r, handlers.Authn, middlewares.jwtAuth, middlewares.anyAuth)
	organizations.RegisterRoutes(r, handlers.Orgs, middlewares.jwtAuth, middlewares.clientOnly)
	projects.RegisterRoutes(r, handlers.Projects, middlewares.anyAuth, middlewares.clientOnly)

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
	log.Printf("identityx pprof listening on :%s", port)
	if err := http.ListenAndServe(":"+port, pmux); err != nil {
		log.Fatalf("identityx pprof server error: %v", err)
	}
}
