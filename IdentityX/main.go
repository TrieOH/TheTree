package main

import (
	"log"
	"net/http"

	"GoAuth/internal/router"
)

func main() {
	defer Db.Close()
	mux := router.CreateRouter(Db)

	log.Printf("GoAuth listening on :%s", Port)
	log.Fatal(http.ListenAndServe(":"+Port, mux))
}
