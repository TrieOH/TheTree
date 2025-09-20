package main

import (
	"log"
	"net/http"
	"GoAuth/internal/handler"
	"GoAuth/internal/service"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()

	goAuthService := service.NewGoAuthService("John")
	goAuthHandler := handler.NewGoAuthHandler(goAuthService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /hi", goAuthHandler.Hi)

	corsMux := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Refresh"},
		AllowCredentials: true,
	}).Handler(mux)

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	log.Print("Started server on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, corsMux))
}
