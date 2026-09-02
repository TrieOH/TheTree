package certifications

import (
	"context"

	"lib/email"
	"univents/internal/authz"
	"univents/ports"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// River is the cert-grant scheduling surface: river.Insert outside any tx
// (the emit endpoint is a fire-and-forget enqueue of the grant job; the job
// itself re-checks existing certs, so re-emission is safe). Satisfied by
// *river.Client[pgx.Tx].
type River interface {
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

type Operations struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	certs    ports.CertificationRepo
	programs ports.ProgramRepo
	email    *email.Client
	river    River
	authz    *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	certs ports.CertificationRepo,
	programs ports.ProgramRepo,
	email *email.Client,
	river River,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events:   events,
		editions: editions,
		certs:    certs,
		programs: programs,
		email:    email,
		river:    river,
		authz:    authz,
	}
}
