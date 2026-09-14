package authn

import (
	"testing"

	"IdentityX/internal/authz"
	"IdentityX/internal/emails"
	"IdentityX/internal/services/tos"
	"IdentityX/models"
	"IdentityX/ports"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

func testActor() models.Actor {
	email := "actor@trieoh.com"
	return models.Actor{ID: uuid.New(), Email: &email, Type: models.HumanActorType}
}

// newTestTosOps builds a tos operations over per-test mockio repos. Tests
// that never touch registration flow through it untouched: unstubbed
// mockio calls resolve to "no terms" (v0).
func newTestTosOps(t *testing.T) *tos.Operations {
	t.Helper()
	return tos.NewOperations(
		mock.Mock[ports.TosRepo](),
		mock.Mock[ports.ProjectRepo](),
		mock.Mock[ports.ActorRepo](),
		authz.New(mock.Mock[ports.OrganizationRepo](), mock.Mock[ports.ProjectRepo](), mock.Mock[ports.PlatformRolesRepo]()),
		emails.NewTosNotifier(mock.Mock[emails.Enqueuer]()),
	)
}
