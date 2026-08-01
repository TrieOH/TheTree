package certifications

import (
	"lib/email"
	"univents/ports"
)

type Operations struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	certs    ports.CertificationRepo
	programs ports.ProgramRepo
	email    *email.Client
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	certs ports.CertificationRepo,
	programs ports.ProgramRepo,
	email *email.Client,
) *Operations {
	return &Operations{
		events:   events,
		editions: editions,
		certs:    certs,
		programs: programs,
		email:    email,
	}
}
