package service

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

type Deps struct {
	Events   ports.EventRepo
	Editions ports.EditionRepo
	Certs    ports.CertificationRepo
	Programs ports.ProgramRepo
	Email    *email.Client
	Jobs     Jobs
	Authz    *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
