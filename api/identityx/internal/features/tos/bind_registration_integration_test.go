package tos_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"IdentityX/models"
	"lib/database"

	"github.com/google/uuid"
)

var errAbortRegistration = errors.New("abort registration")

// TestBindRegistrationJoinsTheTx pins what authn's register and the OAuth
// signup rely on: the consent row is written on the registration's
// transaction, so it sees the actor inserted just before it, commits with
// it, and rolls back with it.
func TestBindRegistrationJoinsTheTx(t *testing.T) {
	f := setupTosFixture(t)
	ctx := context.Background()

	v1, err := f.ops.Create(f.ownerCtx, f.projectID, "Termos v1 em Português")
	if err != nil {
		t.Fatalf("create v1: %v", err)
	}
	tx := database.NewPGXTxRunner(f.pool)

	// register mirrors authn's register: the actor and its consent row in
	// one transaction, which abort rolls back after both are written.
	register := func(email string, abort bool) (uuid.UUID, error) {
		var actorID uuid.UUID
		err := tx.WithinTx(ctx, func(ctx context.Context) error {
			actor, err := f.repos.Actors.Register(ctx, models.Actor{
				ProjectID:    &f.projectID,
				AuthMethod:   models.PasswordAuthMethod,
				Email:        &email,
				Type:         models.HumanActorType,
				PasswordHash: new("x"),
			})
			if err != nil {
				return err
			}
			actorID = actor.ID
			err = f.ops.BindRegistration(ctx, actor.ID, f.projectID, v1)
			if err != nil {
				return err
			}
			if abort {
				return errAbortRegistration
			}
			return nil
		})
		return actorID, err
	}

	committed, err := register("joins@acme.com.br", false)
	if err != nil {
		t.Fatalf("register must commit the actor and its consent together: %v", err)
	}
	rolledBack, err := register("aborts@acme.com.br", true)
	if !errors.Is(err, errAbortRegistration) {
		t.Fatalf("aborted register must return its own error, got %v", err)
	}

	ledger, err := f.repos.Tos.ListAcceptances(ctx, f.projectID)
	if err != nil {
		t.Fatalf("list acceptances: %v", err)
	}
	holds := func(actorID uuid.UUID) bool {
		return slices.ContainsFunc(ledger, func(row models.TosAcceptanceWithActor) bool {
			return row.ActorID == actorID && row.TosVersion == v1.Version
		})
	}
	if !holds(committed) {
		t.Fatalf("committed registration must leave its v1 consent row, got %+v", ledger)
	}
	if holds(rolledBack) {
		t.Fatalf("rolled-back registration must not leave a consent row, got %+v", ledger)
	}
}
