package main

import (
	"TriePayments/initialization"

	_ "github.com/lib/pq"
)

var app *initialization.TriePayments

func init() {
	app = initialization.TriePaymentsSetup()
}
