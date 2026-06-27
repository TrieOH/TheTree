package app

import (
	"net/http"
	idx "sdk/identityx"

	"Informd/internal/features/fields"
	"Informd/internal/features/forms"
	"Informd/internal/features/namespaces"
	"Informd/internal/features/responses"
	"Informd/internal/features/steps"
	_ "Informd/models"

	_ "Informd/generated/docs"

	"github.com/MintzyG/fun"
	_ "github.com/MintzyG/fun"
	"github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/swag/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Deps struct {
	NamespacesHandler *namespaces.Handlers
	FormsHandler      *forms.Handlers
	StepsHandler      *steps.Handlers
	FieldsHandler     *fields.Handlers
	ResponsesHandler  *responses.Handlers

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

// CreateRouter
// @title Informd
// @version 0.0.1
// @description API for managing forms.
// @termsOfService https://git.trieoh.com/TrieOH/TheTree/src/branch/main/api/Informd/LICENSE
// @contact.name TrieOH
// @contact.url https://trieoh.com
// @contact.email support@trieoh.com
// @license.name TSAL License
// @license.url https://git.trieoh.com/TrieOH/TheTree/src/branch/main/api/Informd/LICENSE
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
func (app *Informd) CreateRouter(deps *Deps) http.Handler {
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

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc()
		if fun.Bail(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(doc))
	})

	r.Handle("/metrics", promhttp.Handler())

	namespaces.RegisterRoutes(r, deps.NamespacesHandler, deps.Jwt)
	forms.RegisterRoutes(r, deps.FormsHandler, deps.AnyAuth)
	steps.RegisterRoutes(r, deps.StepsHandler, deps.AnyAuth)
	fields.RegisterRoutes(r, deps.FieldsHandler, deps.AnyAuth)
	responses.RegisterRoutes(r, deps.ResponsesHandler)

	r.Get("/health", handlers.Health(deps.AppName).Handle)

	r.Mount("/", idx.NewSetupHandler(app.idxClient))

	return otelhttp.NewHandler(r, "http.server",
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/health"
		}),
	)
}
