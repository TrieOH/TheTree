package main

import (
	"GoAuth/internal/adapters/http/router"
	"GoAuth/internal/infrastructure/telemetry"
	"context"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	shutdown, err := telemetry.InitTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func(ctx context.Context) {
		err := shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(ctx)

	defer func() {
		err := DB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	defer func() {
		err := scheduler.StopJobs()
		if err != nil {
			log.Printf("Error stopping jobs: %v", err)
		}
		err = scheduler.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	mux := router.CreateRouter(DB)

	log.Printf("GoAuth listening on :%s", Port)
	log.Fatal(http.ListenAndServe(":"+Port, mux))
}
