package authn

import (
	"context"
	"errors"
	"testing"

	"IdentityX/errs"
	"IdentityX/models"
	"IdentityX/ports"

	"github.com/ovechkin-dm/mockio/mock"
)

// Two first setups can both pass the handler's check. The one that takes
// the lock second must find the first account and create nothing.
func TestSetupRefusesWhenAnActorAppearedUnderTheLock(t *testing.T) {
	mock.SetUp(t)
	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.LockSetup(mock.AnyContext())).ThenReturn(nil)
	_ = mock.When(actors.HasAny(mock.AnyContext())).ThenReturn(true, nil)

	ops := emailOps(t, actors, nil)
	err := ops.Setup(context.Background(), models.SetupInput{Email: "admin@trieoh.com", Password: "S3curePassw0rd!"})
	if !errors.Is(err, errs.ErrSetupComplete) {
		t.Fatalf("want setup complete, got %v", err)
	}
	_, _ = mock.Verify(actors, mock.Never()).Register(mock.AnyContext(), mock.Any[models.Actor]())
}
