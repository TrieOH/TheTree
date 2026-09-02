package certifications_test

import (
	"context"
	"os"
	"testing"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/internal/services/certifications"
	"univents/ports"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
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

// recordingRiver captures Insert calls (the emit endpoint's only river use);
// satisfied by the same shape as *river.Client[pgx.Tx].
type recordingRiver struct {
	inserted []rivertype.JobArgs
}

func (r *recordingRiver) Insert(_ context.Context, args river.JobArgs, _ *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	r.inserted = append(r.inserted, args)
	return &rivertype.JobInsertResult{}, nil
}

// newOps builds a certifications Operations for the emit path: mocked
// programs/editions/event repos (authz), the given river seam, and nil for
// everything the emit path never touches (events, certs, email).
func newOps(
	programs ports.ProgramRepo,
	editions ports.EditionRepo,
	authzEvents ports.EventRepo,
	river certifications.River,
) *certifications.Operations {
	return certifications.NewOperations(nil, editions, nil, programs, nil, river, authz.New(authzEvents))
}
