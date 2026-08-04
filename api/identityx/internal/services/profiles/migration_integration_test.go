package profiles

import (
	"context"
	"encoding/json"
	"testing"

	"IdentityX/internal/repos"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"lib/testdb"
)

// TestGetPlatformProfileAfterSchemaShrinksV1ToV2RealDB pins the on-demand
// auto-migration against a real Postgres (testcontainers), walking the
// exact sequence through the real repos and services:
//
//  1. active schema is v1 (old, more fields)
//  2. user profile is set fully filled out per v1
//  3. schema is updated to v2 (fewer fields)
//  4. GET the user profile
//
// Expected: the v1 document is migrated on read — the fields v2 forbids
// (tagline, location, visibility extras) are dropped, the document is
// persisted as schema_version 2 and NOT outdated. The user does nothing;
// only a missing required field falls through to admin resolution.
func TestGetPlatformProfileAfterSchemaShrinksV1ToV2RealDB(t *testing.T) {
	pool := testdb.Postgres(t, "../../../db/migrations")

	ctx := context.Background()
	q := sqlc.New(pool)
	r := repos.New(q)
	ops := testOps(r.Profiles, r.ProfileSchemas, r.Actors, r.Projects)

	// platform-scoped actor (project_id NULL), like a Univents member
	email := "jane@trieoh.com"
	actor, err := q.RegisterActor(ctx, sqlc.RegisterActorParams{
		ProjectID:  nil,
		AuthMethod: "password",
		Email:      &email,
		Type:       string(models.HumanActorType),
	})
	if err != nil {
		t.Fatalf("RegisterActor: %v", err)
	}
	actorID := actor.ID

	// step 1 — schema is v1 (old shape, more fields)
	v1, err := r.ProfileSchemas.Upsert(ctx, models.ProjectProfileSchema{
		ProjectID: nil,
		Schema:    json.RawMessage(schemaV1Full),
		Active:    true,
	})
	if err != nil {
		t.Fatalf("step 1 upsert schema v1: %v", err)
	}
	t.Logf("step 1: schema v1 active (version=%d)", v1.Version)

	// step 2 — user profile fully filled out per v1
	upserted, err := ops.UpsertPlatformProfile(testIdentityFor(actorID), models.UpsertProfileInput{
		ActorID: actorID,
		Profile: json.RawMessage(fullyFilledV1Profile),
	})
	if err != nil {
		t.Fatalf("step 2 UpsertPlatformProfile: %v", err)
	}
	t.Logf("step 2: profile fully filled per v1, stamped schema_version=%d outdated=%t",
		upserted.SchemaVersion, upserted.Outdated)

	// step 3 — schema updated to v2 (fewer fields)
	v2, err := r.ProfileSchemas.Upsert(ctx, models.ProjectProfileSchema{
		ProjectID: nil,
		Schema:    json.RawMessage(schemaV2New),
		Active:    true,
	})
	if err != nil {
		t.Fatalf("step 3 upsert schema v2: %v", err)
	}
	t.Logf("step 3: schema updated to v%d (fewer fields)", v2.Version)

	// step 4 — GET the user profile: must auto-migrate the document to v2
	got, err := ops.GetPlatformProfile(ctx, actorID)
	if err != nil {
		t.Fatalf("step 4 GetPlatformProfile: %v", err)
	}
	t.Logf("step 4: GET returned schema_version=%d outdated=%t", got.SchemaVersion, got.Outdated)
	if got.SchemaVersion != 2 || got.Outdated {
		t.Fatalf("step 4: want migrated profile v2 not outdated, got %+v", got)
	}

	// raw persisted state, straight from the repo (no on-read migration)
	raw, err := r.Profiles.Get(ctx, actorID)
	if err != nil {
		t.Fatalf("step 4 read back: %v", err)
	}
	t.Logf("        persisted: schema_version=%d outdated=%t", raw.SchemaVersion, raw.Outdated)
	if raw.SchemaVersion != 2 || raw.Outdated {
		t.Fatalf("persisted: want v2 not outdated, got %+v", raw)
	}

	// the stored document was rewritten: v1-only fields dropped, kept intact
	var doc map[string]any
	if err := json.Unmarshal(raw.Profile, &doc); err != nil {
		t.Fatalf("unmarshal persisted profile: %v", err)
	}
	if _, ok := doc["tagline"]; ok {
		t.Fatalf("migrated doc must drop tagline, got %s", raw.Profile)
	}
	if _, ok := doc["location"]; ok {
		t.Fatalf("migrated doc must drop location, got %s", raw.Profile)
	}
	vis, _ := doc["visibility"].(map[string]any)
	if len(vis) != 1 || vis["hideLegalName"] != true {
		t.Fatalf("migrated doc must keep only hideLegalName, got %s", raw.Profile)
	}
	if doc["legalName"] != "Jane Doe" {
		t.Fatalf("migrated doc must keep legalName, got %s", raw.Profile)
	}
}

// TestGetPlatformProfileFlagsWhenV2AddsRequiredFieldRealDB pins the other
// side of on-demand migration against a real Postgres: when the new schema
// requires a field the old document cannot provide, pruning cannot help, so
// the profile keeps its version and is flagged outdated for admin
// resolution instead of being stamped v2.
func TestGetPlatformProfileFlagsWhenV2AddsRequiredFieldRealDB(t *testing.T) {
	pool := testdb.Postgres(t, "../../../db/migrations")

	ctx := context.Background()
	q := sqlc.New(pool)
	r := repos.New(q)
	ops := testOps(r.Profiles, r.ProfileSchemas, r.Actors, r.Projects)

	email := "jane@trieoh.com"
	actor, err := q.RegisterActor(ctx, sqlc.RegisterActorParams{
		ProjectID:  nil,
		AuthMethod: "password",
		Email:      &email,
		Type:       string(models.HumanActorType),
	})
	if err != nil {
		t.Fatalf("RegisterActor: %v", err)
	}
	actorID := actor.ID

	// v1 does not require display_name
	if _, err := r.ProfileSchemas.Upsert(ctx, models.ProjectProfileSchema{
		ProjectID: nil,
		Schema:    json.RawMessage(`{"type":"object","properties":{"full_name":{"type":"string"}},"additionalProperties":false}`),
		Active:    true,
	}); err != nil {
		t.Fatalf("upsert schema v1: %v", err)
	}
	// user saved a profile without display_name
	if _, err := ops.UpsertPlatformProfile(testIdentityFor(actorID), models.UpsertProfileInput{
		ActorID: actorID,
		Profile: json.RawMessage(`{"full_name":"Jane"}`),
	}); err != nil {
		t.Fatalf("UpsertPlatformProfile: %v", err)
	}
	// v2 adds display_name as required
	if _, err := r.ProfileSchemas.Upsert(ctx, models.ProjectProfileSchema{
		ProjectID: nil,
		Schema:    json.RawMessage(`{"type":"object","required":["display_name"],"properties":{"full_name":{"type":"string"},"display_name":{"type":"string"}},"additionalProperties":false}`),
		Active:    true,
	}); err != nil {
		t.Fatalf("upsert schema v2: %v", err)
	}

	got, err := ops.GetPlatformProfile(ctx, actorID)
	if err != nil {
		t.Fatalf("GetPlatformProfile: %v", err)
	}
	if got.SchemaVersion != 1 || !got.Outdated {
		t.Fatalf("want flagged v1 profile (version=1 outdated=true) for missing required field, got %+v", got)
	}
	// document must NOT have been rewritten
	var doc map[string]any
	if err := json.Unmarshal(got.Profile, &doc); err != nil {
		t.Fatalf("unmarshal flagged profile: %v", err)
	}
	if len(doc) != 1 || doc["full_name"] != "Jane" {
		t.Fatalf("flagged profile document must stay untouched, got %s", got.Profile)
	}
}
