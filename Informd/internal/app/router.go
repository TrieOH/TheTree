package app

import (
	"Informd/internal/features/forms/handler"
	"Informd/internal/features/keys"
	"Informd/internal/features/namespaces"
	_ "Informd/internal/shared/contracts"
	"net/http"

	"github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Deps struct {
	ProjectsHandler *namespaces.Handler
	ApiKeysHandler  *keys.Handler
	FormsHandler    *handler.Handler
	AsynqmonHandler http.Handler

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
// @title Forms API
// @version 0.0.1
// @description API for managing forms.
// @termsOfService https://github.com/Univents/Univents/blob/main/LICENSE
// @contact.name Univents Team
// @contact.url https://github.com/Univents
// @contact.email support@univents.io
// @license.name MIT License
// @license.url https://github.com/Univents/Univents/blob/main/LICENSE
// @host localhost:8080
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
// @response 400 {object} contracts.ErrorResponse "Standard error response for bad requests"
// @response 401 {object} contracts.ErrorResponse "Standard error response for unauthorized requests"
// @response 403 {object} contracts.ErrorResponse "Standard error response for forbidden requests"
// @response 404 {object} contracts.ErrorResponse "Standard error response for not found errors"
// @response 413 {object} contracts.ErrorResponse "Standard error response for payload too large 1MB"
// @response 429 {object} contracts.ErrorResponse "Standard error response for too many requests"
// @response 500 {object} contracts.ErrorResponse "Standard error response for internal server errors"
// @securityDefinitions.apikey Cookie
// @in header
// @name Cookie
// @description Type "Cookie" followed by a cookie in the format "access_token=xxx; refresh_token=yyy"
func CreateRouter(deps *Deps) http.Handler {
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
	r.Handle("/metrics", promhttp.Handler())
	r.Mount("/admin/asynq", deps.AsynqmonHandler)

	namespaces.RegisterRoutes(r, deps.ProjectsHandler, deps.Jwt)
	keys.RegisterRoutes(r, deps.ApiKeysHandler, deps.Jwt)
	handler.RegisterRoutes(r, deps.FormsHandler, deps.AnyAuth)

	r.Get("/health", handlers.Health(deps.AppName).Handle)

	return otelhttp.NewHandler(r, "http.server")
}
