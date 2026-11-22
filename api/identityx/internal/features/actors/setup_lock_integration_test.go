package actors_test

import (
	"context"
	"testing"

	"IdentityX/internal/sqlc"
	"IdentityX/internal/testrepos"
	"IdentityX/models"
	"lib/testdb"
)

func TestSetupLockAndHasAnyRealDB(t *testing.T) {
	r := testrepos.New(sqlc.New(testdb.Postgres(t, "../../../db/migrations")))
	ctx := context.Background()

	err := r.Actors.LockSetup(ctx)
	if err != nil {
		t.Fatalf("LockSetup: %v", err)
	}
	exists, err := r.Actors.HasAny(ctx)
	if err != nil || exists {
		t.Fatalf("HasAny on an empty database: %v, %v; want false", exists, err)
	}

	email := "first@trieoh.com"
	_, err = r.Actors.Register(ctx, models.Actor{
		AuthMethod: models.PasswordAuthMethod, Email: &email, Type: models.HumanActorType, PasswordHash: new("x"),
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	exists, err = r.Actors.HasAny(ctx)
	if err != nil || !exists {
		t.Fatalf("HasAny after the first actor: %v, %v; want true", exists, err)
	}
}
