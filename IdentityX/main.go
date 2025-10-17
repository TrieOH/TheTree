package main

import (
	"log"
	"net/http"

	"GoAuth/internal/router"
)

func main() {
	defer DB.Close()
	defer scheduler.Shutdown()
	mux := router.CreateRouter(DB)

	log.Printf("GoAuth listening on :%s", Port)
	log.Fatal(http.ListenAndServe(":"+Port, mux))
}
