package main

import (
	_ "Informd/docs"
	"Informd/internal/app"
)

func main() {
	app.New().Run()
}
