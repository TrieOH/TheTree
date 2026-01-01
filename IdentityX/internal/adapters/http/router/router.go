//go:build !test

package router

import (
	"GoAuth/internal/adapters/http/middleware"
	"database/sql"
	"net/http"
	"strings"
	"time"

	_ "GoAuth/docs"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// CreateRouter godoc
// @title GoAuth API
// @version 0.6.0
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
// @response 400 {object} http.ErrorResponse "Standard error response for bad requests"
// @response 401 {object} http.ErrorResponse "Standard error response for unauthorized requests"
// @response 404 {object} http.ErrorResponse "Standard error response for not found errors"
// @response 500 {object} http.ErrorResponse "Standard error response for internal server errors"
func CreateRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ","),
		AllowedMethods:   strings.Split(viper.GetString("CORS_ALLOWED_METHODS"), ","),
		AllowedHeaders:   strings.Split(viper.GetString("CORS_ALLOWED_HEADERS"), ","),
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.Logs)
	r.Use(middleware.Metrics)

	r.Handle("/swagger/*", httpSwagger.WrapHandler)
	r.Handle("/metrics", middleware.Handler())

	r = registerRoutes(db, r)

	return otelhttp.NewHandler(r, "http.server")
}
