package service

import (
	"testing"

	"IdentityX/internal/authz"
	"IdentityX/internal/emails"
	tos "IdentityX/internal/features/tos/service"
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
	return tos.New(tos.Deps{
		Tos:      mock.Mock[ports.TosRepo](),
		Projects: mock.Mock[ports.ProjectRepo](),
		Actors:   mock.Mock[ports.ActorRepo](),
		Authz:    authz.New(mock.Mock[ports.OrganizationRepo](), mock.Mock[ports.ProjectRepo](), mock.Mock[ports.PlatformRolesRepo]()),
		Notifier: emails.NewTosNotifier(mock.Mock[emails.Enqueuer]()),
	})
}
