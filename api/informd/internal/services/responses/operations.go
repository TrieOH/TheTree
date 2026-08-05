package responses

import (
	"Informd/ports"
)

type Operations struct {
	responders ports.ResponderRepo
	responses  ports.ResponseRepo
	answers    ports.AnswerRepo
	forms      ports.FormsRepo
}

func NewOperations(
	responders ports.ResponderRepo,
	responses ports.ResponseRepo,
	answers ports.AnswerRepo,
	forms ports.FormsRepo,
) *Operations {
	return &Operations{
		responders: responders,
		responses:  responses,
		answers:    answers,
		forms:      forms,
	}
}
