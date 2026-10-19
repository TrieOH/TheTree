package main

import (
	"IdentityX/internal/app"
	"lib/errx"
)

func main() {
	errx.Exit(app.Run(), "identityx exited with error")
}
