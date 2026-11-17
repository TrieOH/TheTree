package service_test

import (
	"context"
	"testing"
	"time"
	"uuid"

	"github.com/ovechkin-dm/mockio/mock"

	libjobs "lib/jobs"

	idx "sdk/identityx"
	"univents/internal/authz"
	badges "univents/internal/features/badges/service"
	"univents/models"
	"univents/ports"
)

// newOps builds a badges Operations with mocked ports. The caller wires the
// behaviors it cares about; un-mocked methods return zero values. jobs may be
// nil: a recording fakeJobs is used instead.
func newOps(
	t *testing.T,
	templates ports.BadgeTemplateRepo,
	emissions ports.BadgeEmissionRepo,
	registrations ports.RegistrationRepo,
	editions ports.EditionRepo,
	events ports.EventRepo,
	jobs badges.Jobs,
) *badges.Operations {
	t.Helper()
	return newOpsWithDeps(t, templates, emissions, registrations, editions, events, &fakeActors{}, jobs, mock.Mock[ports.EventRepo]())
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
	jobs badges.Jobs,
	authzEvents ports.EventRepo,
) *badges.Operations {
	t.Helper()
	if jobs == nil {
		jobs = &fakeJobs{}
	}
	return badges.New(badges.Deps{
		Repo:          templates,
		Emissions:     emissions,
		Registrations: registrations,
		Editions:      editions,
		Events:        events,
		Actors:        actors,
		Jobs:          jobs,
		AppURL:        "https://events.test",
		Authz:         authz.New(authzEvents),
	})
}

// fakeJobs records every enqueued job.
type fakeJobs struct{ enqueued []libjobs.Job }

func (f *fakeJobs) EnqueueAfter(_ context.Context, job libjobs.Job, _ time.Duration) error {
	f.enqueued = append(f.enqueued, job)
	return nil
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
		ID: id, EditionID: editionID, UserID: userID,
		Origin: origin, Status: models.BadgeEmissionStatusActive,
		EmittedAt:    time.Now(),
		EditionName:  editionName,
		EndsAt:       endsAt,
		EventName:    eventName,
		TicketTypeID: ticketTypeID,
		TicketName:   ticketName,
	}
}
