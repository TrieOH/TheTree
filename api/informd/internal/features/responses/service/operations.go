package service

import (
	"Informd/ports"
	"lib/database"
)

type Deps struct {
	Responders ports.ResponderRepo
	Responses  ports.ResponseRepo
	Answers    ports.AnswerRepo
	Forms      ports.FormsRepo
	// Tx opens the transaction the response submit runs in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
