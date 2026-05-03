package app

import (
	"IdentityX/internal/features/account"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/auth"
	"IdentityX/internal/features/projects"
	"IdentityX/internal/features/sessions"
	"IdentityX/internal/interfaces/http/middleware"
	_ "IdentityX/internal/shared/contracts"
	"fmt"
	"net/http"

	"github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Handlers struct {
	ApiKeys  *api_keys.Handler
	Users    *auth.Handler
	Accounts *account.Handler
	Sessions *sessions.Handler
	Projects *projects.Handler

	Logger    func(http.Handler) http.Handler
	RequestID func(http.Handler) http.Handler
	BodySize  func(http.Handler) http.Handler
	Metrics   func(http.Handler) http.Handler
	CORS      func(http.Handler) http.Handler
	RealIP    func(http.Handler) http.Handler
	Recover   func(http.Handler) http.Handler
	Timeout   func(http.Handler) http.Handler
	RateLimit func(http.Handler) http.Handler
	Jwt       func(http.Handler) http.Handler
	ApiKey    func(http.Handler) http.Handler
	AnyAuth   func(http.Handler) http.Handler

	AppName string
}

// CreateRouter godoc
// CreateRouter creates a new Chi router and registers all the routes.
// @title GoAuth API
// @version 0.17.10
// @description This is the API for the GoAuth Identity Provider (IdP) service. It provides user authentication, authorization, and project management functionalities.
// @termsOfService https://github.com/TrieOH/GoAuth/blob/main/LICENSE
// @contact.name TrieOH
// @contact.url https://github.com/TrieOH
// @contact.email trieoh@trieoh.com
// @license.name MIT License
// @license.url https://github.com/TrieOH/GoAuth/blob/main/LICENSE
// @host localhost:8080
// @BasePath /
// @schemes http https
// @tag.name auth
// @tag.description "Operations related to user authentication and authorization"
// @tag.name projects
// @tag.description "Operations related to project management"
// @produce json
// @consumes json
// @response 200 {object} object "Standard success response"
// @response 400 {object} contracts.ErrorResponse "Standard error response for bad requests"
// @response 401 {object} contracts.ErrorResponse "Standard error response for unauthorized requests"
// @response 403 {object} contracts.ErrorResponse "Standard error response for forbidden requests"
// @response 404 {object} contracts.ErrorResponse "Standard error response for not found errors"
// @response 413 {object} contracts.ErrorResponse "Standard error response for payload too large 1MB"
// @response 429 {object} contracts.ErrorResponse "Standard error response for too many requests"
// @response 500 {object} contracts.ErrorResponse "Standard error response for internal server errors"
func CreateRouter(deps Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(deps.RealIP)
	r.Use(deps.RequestID)
	r.Use(deps.Logger)
	r.Use(deps.Metrics)
	r.Use(deps.Recover)
	r.Use(deps.Timeout)
	r.Use(deps.BodySize)
	r.Use(deps.RateLimit)
	r.Use(deps.CORS)

	r.Handle("/swagger/*", httpSwagger.WrapHandler)
	r.Handle("/metrics", middleware.Handler())

	auth.RegisterAuthRoutes(r, deps.Users, deps.Jwt)
	account.RegisterRoutes(r, deps.Accounts, deps.Jwt)
	sessions.RegisterRoutes(r, deps.Sessions, deps.Jwt)
	projects.RegisterRoutes(r, deps.Projects, deps.AnyAuth)
	api_keys.RegisterRoutes(r, deps.ApiKeys, deps.Jwt)

	r.Get("/health", handlers.Health("IdentityX-API").Handle)

	if viper.GetBool("DEBUG_MODE") {
		_ = chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			fmt.Printf("[%s] %s\n", method, route)
			return nil
		})
	}

	return otelhttp.NewHandler(r, "http.server")
}
