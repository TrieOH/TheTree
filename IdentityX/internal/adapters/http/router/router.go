//go:build !test

package router

import (
	"GoAuth/internal/adapters/http/middleware"
	"database/sql"
	"net/http"
	"strings"

	_ "GoAuth/docs"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// CreateRouter godoc
// @title        Greet Service API
// @version      0.1
// @description  This is the GreetService API that handles user greetings.
// @contact.name   TrieOH Support
// @contact.url    https://github.com/TrieOH
// @host      localhost:8080
// @BasePath  /
func CreateRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux = registerRoutes(db, mux)

	mux.Handle("GET /metrics", middleware.Handler())
	withMetrics := middleware.Metrics(mux)
	withLogging := middleware.Logs(withMetrics)
	withID := middleware.RequestID(withLogging)

	withCors := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ","),
		AllowedMethods:   strings.Split(viper.GetString("CORS_ALLOWED_METHODS"), ","),
		AllowedHeaders:   strings.Split(viper.GetString("CORS_ALLOWED_HEADERS"), ","),
		AllowCredentials: true,
	}).Handler(withID)

	otelHandler := otelhttp.NewHandler(
		withCors,
		"http.server",
	)

	return otelHandler
}
