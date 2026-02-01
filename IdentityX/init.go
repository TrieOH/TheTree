package main

import (
	"GoAuth/initialization"

	_ "github.com/lib/pq"
)

var app *initialization.GoauthApp

func init() {
	app = initialization.GoAuthSetup()
}
