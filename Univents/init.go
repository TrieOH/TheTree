package main

import (
	"univents/initialization"

	_ "github.com/lib/pq"
)

var app *initialization.UniventsApp

func init() {
	app = initialization.UniventsSetup()
}
