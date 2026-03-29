package main

import (
	"TrieForms/initialization"

	_ "github.com/lib/pq"
)

var app *initialization.TrieForms

func init() {
	app = initialization.TrieFormsSetup()
}
