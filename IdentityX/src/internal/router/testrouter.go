package router

import (
	"log"
	"net/http"

	"GreetService/internal/handler"
	"GreetService/internal/metrics"
	"GreetService/internal/logs"
	"GreetService/internal/repository"
	"GreetService/internal/service"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func CreateTestRouter(db *sql.DB) http.Handler {
	queries := repository.New(db)
	service := service.NewGreetService(queries)
	handler := handler.NewGreetHandler(service)

	mux := http.NewServeMux()
	sgsu := viper.GetString("SPECIAL_GREETING_SERVICE_URL")
	if sgsu != "" {
		mux.HandleFunc("POST /greet/{id}/{greeting_id}", handler.SpecialGreetById)
	} else {
		log.Println("SPECIAL_GREETING_SERVICE_URL -> Unavailable")
		log.Println("Not serving special greet routes")
	}

	mux.Handle("GET /metrics", metrics.Handler())
	withMetrics := metrics.MetricsMW(mux)
	withLogging := logs.LogsMW(withMetrics)
	withID := logs.RequestIDMW(withLogging)

	withCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Refresh"},
		AllowCredentials: true,
	}).Handler(withID)

	return withCors
}
