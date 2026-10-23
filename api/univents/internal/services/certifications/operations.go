package certifications

import (
	"context"

	"lib/email"
	libjobs "lib/jobs"
	"univents/internal/authz"
	"univents/ports"
)

// Jobs is the cert-grant scheduling surface: a plain enqueue outside any tx
// (the emit endpoint is a fire-and-forget enqueue of the grant job; the job
// itself re-checks existing certs, so re-emission is safe). Satisfied by
// the lib/jobs carriers.
type Jobs interface {
	Enqueue(ctx context.Context, job libjobs.Job) error
}

type Operations struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	certs    ports.CertificationRepo
	programs ports.ProgramRepo
	email    *email.Client
	jobs     Jobs
	authz    *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	certs ports.CertificationRepo,
	programs ports.ProgramRepo,
	email *email.Client,
	jobs Jobs,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events:   events,
		editions: editions,
		certs:    certs,
		programs: programs,
		email:    email,
		jobs:     jobs,
		authz:    authz,
	}
}
