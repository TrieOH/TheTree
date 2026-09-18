package main

import (
	"lib/errx"
	"payssage/internal/app"
)

func main() {
	errx.Exit(app.Run(), "payssage exited with error")
}
