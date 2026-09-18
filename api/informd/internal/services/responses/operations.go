package responses

import (
	"Informd/ports"
	"lib/database"
)

type Operations struct {
	responders ports.ResponderRepo
	responses  ports.ResponseRepo
	answers    ports.AnswerRepo
	forms      ports.FormsRepo
	// tx opens the transaction the response submit runs in; threaded
	// from boot, no package-level runner.
	tx database.TxRunner
}

func NewOperations(
	responders ports.ResponderRepo,
	responses ports.ResponseRepo,
	answers ports.AnswerRepo,
	forms ports.FormsRepo,
	tx database.TxRunner,
) *Operations {
	return &Operations{
		responders: responders,
		responses:  responses,
		answers:    answers,
		forms:      forms,
		tx:         tx,
	}
}
