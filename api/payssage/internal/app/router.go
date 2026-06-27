package app

import (
	"net/http"
	"time"

	"payssage/internal/features/api_keys"
	"payssage/internal/features/intents"
	"payssage/internal/features/oauth"
	"payssage/internal/features/webhooks"
	"payssage/internal/features/workspaces"

	fh "github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HTTPDeps struct {
	IntentsHandler    *intents.Handler
	WorkspacesHandler *workspaces.Handler
	ApiKeysHandler    *api_keys.Handler
	WebhooksHandler   *webhooks.Handler
	OauthHandler      *oauth.Handler
}

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
// @securityDefinitions.apikey Cookie
// @in header
// @name Cookie
// @description Type "Cookie" followed by a cookie in the format "access_token=xxx; refresh_token=yyy"
func CreateRouter(deps *HTTPDeps) http.Handler {
	r := chi.NewRouter()

	if !viper.GetBool("DISABLE_RATE_LIMIT") {
		r.Use(httprate.Limit(
			400,
			1*time.Minute,
			httprate.WithKeyFuncs(httprate.KeyByRealIP),
		))
	}

	r.Handle("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/health", fh.Health("payssage").Handle)

	registerRoutes(r, deps)
	return otelhttp.NewHandler(r, "http.server")
}
