package app

import (
	"IdentityX/internal/features/authn"
	_ "IdentityX/models"
	"context"
	"fmt"
	"lib/errx"
	"log/slog"
	"net/http"

	_ "IdentityX/generated/docs"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"github.com/swaggo/swag/v2"
	"riverqueue.com/riverui"
)

type RouterDeps struct {
	AppName string

	CORS       func(http.Handler) http.Handler
	Logger     func(http.Handler) http.Handler
	JwtAuth    func(http.Handler) http.Handler
	ApiKeyAuth func(http.Handler) http.Handler
	AnyAuth    func(http.Handler) http.Handler

	Authn *authn.Handlers
}

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
// @tag.name projects
// @tag.description "Operations related to project management"
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
func (app *IdentityX) CreateRouter(deps RouterDeps, debugMode, disableRateLimit bool) http.Handler {
	r := chi.NewRouter()

	r.Use(otelchi.Middleware(deps.AppName,
		otelchi.WithChiRoutes(r),
		otelchi.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/health" && r.URL.Path != "/metrics"
		}),
	))

	//r.Use(deps.RealIP)
	//r.Use(deps.RequestID)
	r.Use(deps.Logger)
	//r.Use(deps.Metrics)
	//r.Use(deps.Recover)
	//r.Use(deps.Timeout)
	//r.Use(deps.BodySize)
	//r.Use(deps.RateLimit)
	r.Use(deps.CORS)

	endpoints := riverui.NewEndpoints(app.river, nil)

	handler, err := riverui.NewHandler(&riverui.HandlerOpts{
		Endpoints: endpoints,
		Logger:    slog.Default(),
		Prefix:    "/riverui",
	})
	if err != nil {
		errx.Exit(err, "failed to create river handler")
	}
	err = handler.Start(context.Background())
	if err != nil {
		errx.Exit(err, "failed to start river handler")
	}
	r.Mount("/riverui", handler)

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc()
		if fun.Bail(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(doc))
	})

	r.Handle("/metrics", promhttp.Handler())

	authn.RegisterRoutes(r, deps.Authn, deps.JwtAuth)
	//account.RegisterRoutes(r, deps.Accounts, deps.Jwt)
	//sessions.RegisterRoutes(r, deps.Sessions, deps.Jwt)
	//projects.RegisterRoutes(r, deps.Projects, deps.AnyAuth)
	//api_keys.RegisterRoutes(r, deps.ApiKeys, deps.Jwt)

	r.Get("/health", handlers.Health("IdentityX-API").Handle)

	if debugMode {
		_ = chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			fmt.Printf("[%s] %s\n", method, route)
			return nil
		})
	}

	return r
}
