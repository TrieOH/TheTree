package main

import (
	_ "payssage/generated/docs"
	"payssage/internal/app"
)

func main() {
	app.New().Run()
}
