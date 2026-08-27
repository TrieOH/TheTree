package badges_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"lib/email"

	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/internal/services/badges"
	"univents/models"
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

// newOps builds a badges Operations with mocked ports. The caller wires the
// behaviors it cares about; un-mocked methods return zero values. email may be
// nil for tests that never trigger the badge email goroutine.
func newOps(
	t *testing.T,
	templates ports.BadgeTemplateRepo,
	emissions ports.BadgeEmissionRepo,
	registrations ports.RegistrationRepo,
	editions ports.EditionRepo,
	events ports.EventRepo,
	email *email.Client,
) *badges.Operations {
	t.Helper()
	return newOpsWithDeps(t, templates, emissions, registrations, editions, events, &fakeActors{}, email, mock.Mock[ports.EventRepo]())
}

// newOpsWithAuthz is newOps with control over the event repo used by authz
// (stub GetRole to a role that passes feature checks).
func newOpsWithAuthz(
	t *testing.T,
	templates ports.BadgeTemplateRepo,
	emissions ports.BadgeEmissionRepo,
	editions ports.EditionRepo,
	authzEvents ports.EventRepo,
) *badges.Operations {
	t.Helper()
	return newOpsWithDeps(t, templates, emissions, nil, editions, nil, &fakeActors{}, nil, authzEvents)
}

// newOpsWithDeps is the full constructor: everything newOps provides plus the
// actor resolver (the profile-badges gift claim needs the account email) and
// the authz event repo.
func newOpsWithDeps(
	t *testing.T,
	templates ports.BadgeTemplateRepo,
	emissions ports.BadgeEmissionRepo,
	registrations ports.RegistrationRepo,
	editions ports.EditionRepo,
	events ports.EventRepo,
	actors ports.ActorResolver,
	email *email.Client,
	authzEvents ports.EventRepo,
) *badges.Operations {
	t.Helper()
	return badges.NewOperations(templates, emissions, registrations, editions, events, actors, email, authz.New(authzEvents))
}

// fakeActors is a no-account actor resolver: unknown id/email →
// ports.ErrActorNotFound, so the claim short-circuits unless seeded.
// Seed with fakeActors{byID: ...} for claim tests.
type fakeActors struct {
	byID map[uuid.UUID]*idx.Actor
}

func (f *fakeActors) GetByID(_ context.Context, id uuid.UUID) (*idx.Actor, error) {
	if f.byID != nil {
		if a, ok := f.byID[id]; ok {
			return a, nil
		}
	}
	return nil, ports.ErrActorNotFound
}

func (f *fakeActors) GetByEmail(_ context.Context, _ string) (*idx.Actor, error) {
	return nil, ports.ErrActorNotFound
}

// ownerCtx returns a context with an authenticated identity for feature tests
// that pass through RequireIdentity.
func ownerCtx() context.Context {
	return idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: uuid.New()},
	})
}

func view(
	id, editionID, userID uuid.UUID,
	origin models.BadgeEmissionOrigin,
	editionName string,
	endsAt time.Time,
	eventName string,
	ticketTypeID *uuid.UUID,
	ticketName *string,
) models.BadgeEmissionView {
	return models.BadgeEmissionView{
		BadgeEmission: models.BadgeEmission{
			ID: id, EditionID: editionID, UserID: userID,
			Origin: origin, Status: models.BadgeEmissionStatusActive,
			EmittedAt: time.Now(),
		},
		EditionName:  editionName,
		EndsAt:       endsAt,
		EventName:    eventName,
		TicketTypeID: ticketTypeID,
		TicketName:   ticketName,
	}
}
