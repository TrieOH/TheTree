package main

import (
	"lib/errx"
	"univents/internal/app"
)

func main() {
	errx.Exit(app.Run(), "univents exited with error")
}
