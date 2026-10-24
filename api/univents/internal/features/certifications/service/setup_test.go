package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	libjobs "lib/jobs"
	idx "sdk/identityx"
	"univents/internal/authz"
	certifications "univents/internal/features/certifications/service"
	"univents/ports"
)

func TestMain(m *testing.M) {
	fun.SetConfig(fun.Config{
		DefaultModule:        "test",
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
	})
	os.Exit(m.Run())
}

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
