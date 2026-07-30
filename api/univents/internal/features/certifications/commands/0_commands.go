package commands

import (
	"lib/email"
	"univents/ports"
)

type Commands struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	certs    ports.CertificationRepo
	programs ports.ProgramRepo
	email    *email.Client
}

func NewCommands(
	events ports.EventRepo,
	editions ports.EditionRepo,
	certs ports.CertificationRepo,
	programs ports.ProgramRepo,
	email *email.Client,
) *Commands {
	return &Commands{
		events:   events,
		editions: editions,
		certs:    certs,
		programs: programs,
		email:    email,
	}
}
