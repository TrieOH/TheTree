package commands

import (
	"Informd/ports"
)

type Commands struct {
	responders ports.ResponderRepo
	responses  ports.ResponseRepo
	answers    ports.AnswerRepo
	forms      ports.FormsRepo
}

func NewCommands(
	responders ports.ResponderRepo,
	responses ports.ResponseRepo,
	answers ports.AnswerRepo,
	forms ports.FormsRepo,
) *Commands {
	return &Commands{
		responders: responders,
		responses:  responses,
		answers:    answers,
		forms:      forms,
	}
}
