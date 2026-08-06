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
	return newOpsWithAuthz(t, templates, emissions, registrations, editions, events, email, mock.Mock[ports.EventRepo]())
}

// newOpsWithAuthz is newOps with control over the event repo used by authz
// (stub GetRole to a role that passes feature checks).
func newOpsWithAuthz(
	t *testing.T,
	templates ports.BadgeTemplateRepo,
	emissions ports.BadgeEmissionRepo,
	registrations ports.RegistrationRepo,
	editions ports.EditionRepo,
	events ports.EventRepo,
	email *email.Client,
	authzEvents ports.EventRepo,
) *badges.Operations {
	t.Helper()
	return badges.NewOperations(templates, emissions, registrations, editions, events, email, authz.New(authzEvents))
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
