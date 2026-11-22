package blacklist_test

import (
	"context"
	"testing"
	"time"

	"IdentityX/internal/sqlc"
	"IdentityX/internal/testrepos"
	"IdentityX/models"
	"lib/testdb"
)

func entry(target string, typ models.BlacklistEntryType) models.BlacklistEntry {
	expires := time.Now().Add(time.Hour)
	return models.BlacklistEntry{Type: typ, Target: target, ExpiresAt: &expires}
}

// Append is the claim a rotation makes on a refresh token: only the first
// one adds the entry.
func TestAppendReportsWhetherItAddedTheEntryRealDB(t *testing.T) {
	r := testrepos.New(sqlc.New(testdb.Postgres(t, "../../../db/migrations")))
	ctx := context.Background()

	added, err := r.Blacklist.Append(ctx, entry("jti-1", models.BlacklistEntryTypeToken))
	if err != nil || !added {
		t.Fatalf("first Append: added=%v err=%v, want added", added, err)
	}
	added, err = r.Blacklist.Append(ctx, entry("jti-1", models.BlacklistEntryTypeToken))
	if err != nil || added {
		t.Fatalf("second Append: added=%v err=%v, want not added and no error", added, err)
	}
}

// Put restarts an entry that is already there, even with a null project.
func TestPutRestartsAnExistingEntryRealDB(t *testing.T) {
	r := testrepos.New(sqlc.New(testdb.Postgres(t, "../../../db/migrations")))
	ctx := context.Background()

	err := r.Blacklist.Put(ctx, entry("actor-1", models.BlacklistEntryTypeActor))
	if err != nil {
		t.Fatalf("first Put: %v", err)
	}
	first, err := r.Blacklist.GetByTargetAndType(ctx, "actor-1", models.BlacklistEntryTypeActor)
	if err != nil {
		t.Fatalf("read first: %v", err)
	}

	time.Sleep(20 * time.Millisecond)
	err = r.Blacklist.Put(ctx, entry("actor-1", models.BlacklistEntryTypeActor))
	if err != nil {
		t.Fatalf("second Put: %v", err)
	}
	second, err := r.Blacklist.GetByTargetAndType(ctx, "actor-1", models.BlacklistEntryTypeActor)
	if err != nil {
		t.Fatalf("read second: %v", err)
	}
	if !second.CreatedAt.After(first.CreatedAt) {
		t.Fatalf("want created_at moved forward, got %v then %v", first.CreatedAt, second.CreatedAt)
	}
}
