package main

import (
	"Informd/internal/app"
	"lib/errx"
)

func main() {
	errx.Exit(app.Run(), "informd exited with error")
}
