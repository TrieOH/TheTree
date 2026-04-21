//go:build !test

package router

import (
	"IdentityX/internal/features/account"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/auth"
	"IdentityX/internal/features/projects"
	"IdentityX/internal/features/sessions"
	"IdentityX/internal/interfaces/http/middleware"
	"IdentityX/internal/interfaces/http/system"
	"log"
	"net/http"
	"strings"
	"time"

	_ "IdentityX/internal/shared/contracts"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Handlers struct {
	System   *system.Handler
	ApiKeys  *api_keys.Handler
	Users    *auth.Handler
	Accounts *account.Handler
	Sessions *sessions.Handler
	Projects *projects.Handler
	AuthMW   middleware.AuthMiddleware
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
func CreateRouter(handlers Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	if !viper.GetBool("DISABLE_RATE_LIMIT") {
		r.Use(httprate.Limit(
			400,
			1*time.Minute,
			httprate.WithKeyFuncs(httprate.KeyByRealIP),
		))
	}

	r.Use(middleware.MaxBodySize(1 << 20)) // 1 MB

	r.Use(cors.Handler(GetCORSOptions()))

	r.Use(middleware.RequestID)
	r.Use(middleware.Logs)
	r.Use(middleware.Metrics)

	r.Handle("/swagger/*", httpSwagger.WrapHandler)
	r.Handle("/metrics", middleware.Handler())

	r = registerRoutes(handlers, r)

	return otelhttp.NewHandler(r, "http.server")
}

func splitAndCleanCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}

	if len(out) == 0 {
		return nil
	}

	return out
}

// GetCORSOptions builds a safe cors.Options configuration from environment
// variables, applying sensible defaults and preventing invalid empty values.
func GetCORSOptions() cors.Options {
	allowedOrigins := splitAndCleanCSV(viper.GetString("CORS_ALLOWED_ORIGINS"))
	allowedMethods := splitAndCleanCSV(viper.GetString("CORS_ALLOWED_METHODS"))
	allowedHeaders := splitAndCleanCSV(viper.GetString("CORS_ALLOWED_HEADERS"))

	// Never default origins to "*"
	if allowedOrigins == nil {
		log.Fatalf("No AllowedOrigins set in CORS_ALLOWED_ORIGINS")
	}

	if allowedMethods == nil {
		allowedMethods = []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		}
	}

	if allowedHeaders == nil {
		allowedHeaders = []string{
			"Accept",
			"Authorization",
			"Content-Type",
		}
	}

	return cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		AllowCredentials: true,
		MaxAge:           300,
	}
}
