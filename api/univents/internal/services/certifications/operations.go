package certifications

import (
	"lib/email"
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	certs    ports.CertificationRepo
	programs ports.ProgramRepo
	email    *email.Client
	authz    *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	certs ports.CertificationRepo,
	programs ports.ProgramRepo,
	email *email.Client,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events:   events,
		editions: editions,
		certs:    certs,
		programs: programs,
		email:    email,
		authz:    authz,
	}
}
