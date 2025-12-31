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
// @title        Greet Service API
// @version      0.6.0
// @description  This is the GoAuth IdP API
// @contact.name   TrieOH Support
// @contact.url    https://github.com/TrieOH
// @host      localhost:8080
// @BasePath  /
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
