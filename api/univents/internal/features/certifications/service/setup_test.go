package service_test

import (
	"context"

	"github.com/google/uuid"

	libjobs "lib/jobs"
	idx "sdk/identityx"
	"univents/internal/authz"
	certifications "univents/internal/features/certifications/service"
	"univents/ports"
)

// ownerCtx returns a context with an authenticated identity for feature tests
// that pass through RequireIdentity.
func ownerCtx() context.Context {
	return idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: uuid.New()},
	})
}

// recordingJobs captures Enqueue calls (the emit endpoint's only job use).
type recordingJobs struct {
	inserted []libjobs.Job
}

func (r *recordingJobs) Enqueue(_ context.Context, job libjobs.Job) error {
	r.inserted = append(r.inserted, job)
	return nil
}

// newOps builds a certifications Operations for the emit path: mocked
// programs/editions/event repos (authz), the given jobs seam, and nil for
// everything the emit path never touches (events, certs, email).
func newOps(
	programs ports.ProgramRepo,
	editions ports.EditionRepo,
	authzEvents ports.EventRepo,
	enq certifications.Jobs,
) *certifications.Operations {
	return certifications.New(certifications.Deps{
		Editions: editions,
		Programs: programs,
		Jobs:     enq,
		Authz:    authz.New(authzEvents),
	})
}
