package service

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"strings"
	"testing"

	"IdentityX/errs"
	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	libauthz "lib/authz"
	"lib/database"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

// ── helpers ───────────────────────────────────────────────────────────────

const schemaV2 = `{"type":"object","required":["full_name","display_name"],"properties":{"full_name":{"type":"string"},"display_name":{"type":"string"}}}`

// schemaV1Full is the OLD platform profile schema (more fields: tagline,
// location, full visibility block). schemaV2New is the NEW schema with
// fewer fields — both pasted verbatim from the reported incident.
var schemaV1Full = fixture("schema_v1_full.json")

// schemaV2New is the NEW schema: fewer fields than v1 — tagline and
// location are gone, visibility only keeps hideLegalName.
var schemaV2New = fixture("schema_v2_new.json")

// fullyFilledV1Profile is a profile with every v1 field filled in,
// including the fields v2 removed (tagline, location, visibility extras).
var fullyFilledV1Profile = fixture("fully_filled_v1_profile.json")

// fixture reads a JSON fixture from the feature's testdata, shared with
// the integration tests one level up.
func fixture(name string) string {
	b, err := fs.ReadFile(os.DirFS("../testdata"), name)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(b))
}

func testOps(profiles ports.ProfileRepo, schemas ports.ProfileSchemaRepo, actors ports.ActorRepo, projects ports.ProjectRepo) *Operations {
	return New(Deps{
		Profiles: profiles,
		Schemas:  schemas,
		Actors:   actors,
		Authz:    authz.New(nil, projects, nil),
	})
}

func mockProjectActor(actors ports.ActorRepo, actorID, projectID uuid.UUID) {
	mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actorID))).
		ThenReturn(&models.Actor{ID: actorID, ProjectID: &projectID, Type: models.HumanActorType}, nil)
}

func testIdentity() context.Context {
	return testIdentityFor(uuid.New())
}

func testIdentityFor(actorID uuid.UUID) context.Context {
	email := "jane@trieoh.com"
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub:  models.Subject{ID: actorID, Email: &email, Type: models.HumanActorType},
		Cred: models.Credential{Type: models.TokenCredentialType},
	})
}

func testIdentityScopedTo(actorID, projectID uuid.UUID) context.Context {
	ctx := testIdentityFor(actorID)
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		panic(err)
	}
	ident.Sub.ProjectID = &projectID
	return ctx
}

func schemaOf(projectID *uuid.UUID, version int) *models.ProjectProfileSchema {
	return &models.ProjectProfileSchema{
		ProjectID: projectID,
		Schema:    json.RawMessage(schemaV2),
		Version:   version,
		Active:    true,
	}
}

func mockProjectAdmin(projects ports.ProjectRepo) {
	mock.When(projects.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.ProjectRoleAdmin, nil)
}

// ── on-demand migration on read ──────────────────────────────────────────

func TestGetProfileMigratesWhenNewerSchemaValid(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()
	profile := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       json.RawMessage(`{"full_name":"Jane","display_name":"JD"}`),
		SchemaVersion: 1,
	}
	migrated := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       profile.Profile,
		SchemaVersion: 2,
	}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(actorID))).ThenReturn(profile, nil)
	mock.When(profiles.SetMigrationState(mock.AnyContext(), mock.Equal(actorID), mock.Equal(2), mock.Equal(false))).ThenReturn(migrated, nil)

	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).ThenReturn(schemaOf(&projectID, 2), nil)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	got, err := testOps(profiles, schemas, actors, projects).GetProfile(testIdentity(), actorID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	_, _ = mock.Verify(profiles, mock.Once()).SetMigrationState(mock.AnyContext(), mock.Equal(actorID), mock.Equal(2), mock.Equal(false))
	if got.SchemaVersion != 2 || got.Outdated {
		t.Fatalf("want migrated profile v2, got %+v", got)
	}
}

func TestGetProfileFlagsWhenNewerSchemaInvalid(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()
	profile := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       json.RawMessage(`{"full_name":"Jane"}`),
		SchemaVersion: 1,
	}
	flagged := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       profile.Profile,
		SchemaVersion: 1,
		Outdated:      true,
	}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(actorID))).ThenReturn(profile, nil)
	mock.When(profiles.SetMigrationState(mock.AnyContext(), mock.Equal(actorID), mock.Equal(1), mock.Equal(true))).ThenReturn(flagged, nil)

	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).ThenReturn(schemaOf(&projectID, 2), nil)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	got, err := testOps(profiles, schemas, actors, projects).GetProfile(testIdentity(), actorID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	// keep the current version, flag as outdated
	_, _ = mock.Verify(profiles, mock.Once()).SetMigrationState(mock.AnyContext(), mock.Equal(actorID), mock.Equal(1), mock.Equal(true))
	if !got.Outdated {
		t.Fatalf("want flagged profile, got %+v", got)
	}
}

func TestGetProfileSkipsMigrationWhenVersionsMatch(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()
	profile := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       json.RawMessage(`{"full_name":"Jane"}`),
		SchemaVersion: 2,
	}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(actorID))).ThenReturn(profile, nil)

	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).ThenReturn(schemaOf(&projectID, 2), nil)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	got, err := testOps(profiles, schemas, actors, projects).GetProfile(testIdentity(), actorID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	_, _ = mock.Verify(profiles, mock.Never()).SetMigrationState(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[int](), mock.Any[bool]())
	if got.SchemaVersion != 2 {
		t.Fatalf("want unchanged profile, got %+v", got)
	}
}

func TestGetProfileSkipsMigrationWithoutSchema(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()
	profile := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       json.RawMessage(`{"full_name":"Jane"}`),
		SchemaVersion: 1,
	}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(actorID))).ThenReturn(profile, nil)

	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).ThenReturn(nil, database.ErrNotFound)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	got, err := testOps(profiles, schemas, actors, projects).GetProfile(testIdentity(), actorID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	_, _ = mock.Verify(profiles, mock.Never()).SetMigrationState(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[int](), mock.Any[bool]())
	if got.SchemaVersion != 1 {
		t.Fatalf("want unchanged profile, got %+v", got)
	}
}

func TestGetProfileSkipsMigrationWhenSchemaInactive(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()
	profile := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       json.RawMessage(`{"full_name":"Jane"}`),
		SchemaVersion: 1,
	}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(actorID))).ThenReturn(profile, nil)

	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	inactive := schemaOf(&projectID, 2)
	inactive.Active = false
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).ThenReturn(inactive, nil)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	got, err := testOps(profiles, schemas, actors, projects).GetProfile(testIdentity(), actorID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	_, _ = mock.Verify(profiles, mock.Never()).SetMigrationState(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[int](), mock.Any[bool]())
	if got.SchemaVersion != 1 {
		t.Fatalf("want unchanged profile, got %+v", got)
	}
}

func TestGetProfileSkipsWriteWhenAlreadyFlagged(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()
	profile := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       json.RawMessage(`{"full_name":"Jane"}`),
		SchemaVersion: 1,
		Outdated:      true,
	}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(actorID))).ThenReturn(profile, nil)

	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).ThenReturn(schemaOf(&projectID, 2), nil)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	got, err := testOps(profiles, schemas, actors, projects).GetProfile(testIdentity(), actorID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	// state (version 1, outdated) is already persisted: a read must not churn writes
	_, _ = mock.Verify(profiles, mock.Never()).SetMigrationState(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[int](), mock.Any[bool]())
	if !got.Outdated {
		t.Fatalf("want flagged profile, got %+v", got)
	}
}

func TestGetPlatformProfileMigratesWithPlatformSchema(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	profile := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       json.RawMessage(`{"full_name":"Jane","display_name":"JD"}`),
		SchemaVersion: 1,
	}
	migrated := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       profile.Profile,
		SchemaVersion: 2,
	}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(actorID))).ThenReturn(profile, nil)
	mock.When(profiles.SetMigrationState(mock.AnyContext(), mock.Equal(actorID), mock.Equal(2), mock.Equal(false))).ThenReturn(migrated, nil)

	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	captor := mock.Captor[*uuid.UUID]()
	mock.When(schemas.Get(mock.AnyContext(), captor.Capture())).ThenReturn(schemaOf(nil, 2), nil)

	got, err := testOps(profiles, schemas, mock.Mock[ports.ActorRepo](), mock.Mock[ports.ProjectRepo]()).GetPlatformProfile(testIdentity(), actorID)
	if err != nil {
		t.Fatalf("GetPlatformProfile: %v", err)
	}
	_, _ = mock.Verify(schemas, mock.Once()).Get(mock.AnyContext(), mock.Any[*uuid.UUID]())
	if captor.Values()[0] != nil {
		t.Fatalf("platform scope must query schema with nil project_id, got %v", captor.Values()[0])
	}
	_, _ = mock.Verify(profiles, mock.Once()).SetMigrationState(mock.AnyContext(), mock.Equal(actorID), mock.Equal(2), mock.Equal(false))
	if got.SchemaVersion != 2 {
		t.Fatalf("want migrated profile v2, got %+v", got)
	}
}

// ── project user vs project member ────────────────────────────────────────

func TestGetProfileAllowsSelfReadWithoutMembership(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()
	profile := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       json.RawMessage(`{"full_name":"Jane"}`),
		SchemaVersion: 1,
	}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(actorID))).ThenReturn(profile, nil)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	// no GetRole stub: a project user has no project_members row, so the
	// member gate must never be consulted for a self-read
	projects := mock.Mock[ports.ProjectRepo]()

	// no schema configured: nothing to migrate
	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).
		ThenReturn(nil, database.ErrNotFound)

	got, err := testOps(profiles, schemas, actors, projects).
		GetProfile(testIdentityFor(actorID), actorID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	_, _ = mock.Verify(projects, mock.Never()).GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())
	if got.ActorID != actorID {
		t.Fatalf("want own profile, got %+v", got)
	}
}

func TestGetProfileAllowsProjectUserReadingOtherProjectUsers(t *testing.T) {
	mock.SetUp(t)
	requesterID := uuid.New()
	targetID := uuid.New()
	projectID := uuid.New()
	profile := &models.ActorProfile{
		ActorID:       targetID,
		Profile:       json.RawMessage(`{"full_name":"Jane"}`),
		SchemaVersion: 1,
	}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(targetID))).ThenReturn(profile, nil)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, targetID, projectID)

	// requester is a project user (scoped actor, no membership row): the
	// member gate must not be consulted for same-project reads
	projects := mock.Mock[ports.ProjectRepo]()

	// no schema configured: nothing to migrate
	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).
		ThenReturn(nil, database.ErrNotFound)

	got, err := testOps(profiles, schemas, actors, projects).
		GetProfile(testIdentityScopedTo(requesterID, projectID), targetID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	_, _ = mock.Verify(projects, mock.Never()).GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())
	if got.ActorID != targetID {
		t.Fatalf("want target profile, got %+v", got)
	}
}

func TestGetProfileAllowsAnonymousPublicRead(t *testing.T) {
	mock.SetUp(t)
	targetID := uuid.New()
	projectID := uuid.New()
	profile := &models.ActorProfile{
		ActorID:       targetID,
		Profile:       json.RawMessage(`{"full_name":"Jane"}`),
		SchemaVersion: 1,
	}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(targetID))).ThenReturn(profile, nil)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, targetID, projectID)

	// public read: no identity, no membership consulted
	projects := mock.Mock[ports.ProjectRepo]()

	// no schema configured: nothing to migrate
	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).
		ThenReturn(nil, database.ErrNotFound)

	got, err := testOps(profiles, schemas, actors, projects).
		GetProfile(context.Background(), targetID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	_, _ = mock.Verify(projects, mock.Never()).GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())
	if got.ActorID != targetID {
		t.Fatalf("want target profile, got %+v", got)
	}
}

// ── project-scope enforcement ─────────────────────────────────────────────

func TestGetProfileDeniesActorFromAnotherProject(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()
	otherProject := uuid.New()

	actors := mock.Mock[ports.ActorRepo]()
	mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actorID))).
		ThenReturn(&models.Actor{ID: actorID, ProjectID: &otherProject, Type: models.HumanActorType}, nil)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	_, err := testOps(mock.Mock[ports.ProfileRepo](), mock.Mock[ports.ProfileSchemaRepo](), actors, projects).
		GetProfile(testIdentity(), actorID, projectID)
	if !errors.Is(err, errs.ErrActorOutsideProject) {
		t.Fatalf("want forbidden for cross-project actor, got %v", err)
	}
}

func TestGetProfileDeniesPlatformActorThroughProjectRoute(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()

	actors := mock.Mock[ports.ActorRepo]()
	mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actorID))).
		ThenReturn(&models.Actor{ID: actorID, ProjectID: nil, Type: models.HumanActorType}, nil)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	_, err := testOps(mock.Mock[ports.ProfileRepo](), mock.Mock[ports.ProfileSchemaRepo](), actors, projects).
		GetProfile(testIdentity(), actorID, projectID)
	if !errors.Is(err, errs.ErrActorOutsideProject) {
		t.Fatalf("want forbidden for platform actor via project route, got %v", err)
	}
}

func TestUpsertProfileDeniesCrossProjectActor(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()
	otherProject := uuid.New()

	actors := mock.Mock[ports.ActorRepo]()
	mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actorID))).
		ThenReturn(&models.Actor{ID: actorID, ProjectID: &otherProject, Type: models.HumanActorType}, nil)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	_, err := testOps(mock.Mock[ports.ProfileRepo](), mock.Mock[ports.ProfileSchemaRepo](), actors, projects).
		UpsertProfile(testIdentity(), models.UpsertProfileInput{
			ActorID: actorID,
			Profile: json.RawMessage(`{"full_name":"Jane"}`),
		}, projectID)
	if !errors.Is(err, errs.ErrActorOutsideProject) {
		t.Fatalf("want forbidden for cross-project actor, got %v", err)
	}
}

func TestUpsertProfileDeniesMemberUpdatingOthers(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	projects := mock.Mock[ports.ProjectRepo]()
	mock.When(projects.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Equal(projectID))).
		ThenReturn(models.ProjectRoleMember, nil)

	_, err := testOps(mock.Mock[ports.ProfileRepo](), mock.Mock[ports.ProfileSchemaRepo](), actors, projects).
		UpsertProfile(testIdentity(), models.UpsertProfileInput{
			ActorID: actorID,
			Profile: json.RawMessage(`{"full_name":"Jane"}`),
		}, projectID)
	if !errors.Is(err, libauthz.ErrInsufficientRole) {
		t.Fatalf("want forbidden for member updating another member's profile, got %v", err)
	}
}

func TestUpsertProfileAllowsAdminUpdatingOthers(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()

	profiles := mock.Mock[ports.ProfileRepo]()
	captor := mock.Captor[models.ActorProfile]()
	mock.When(profiles.Upsert(mock.AnyContext(), captor.Capture())).ThenReturn(nil, nil)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).ThenReturn(schemaOf(&projectID, 2), nil)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	_, err := testOps(profiles, schemas, actors, projects).
		UpsertProfile(testIdentity(), models.UpsertProfileInput{
			ActorID: actorID,
			Profile: json.RawMessage(`{"full_name":"Jane","display_name":"JD"}`),
		}, projectID)
	if err != nil {
		t.Fatalf("UpsertProfile: %v", err)
	}
	_, _ = mock.Verify(profiles, mock.Once()).Upsert(mock.AnyContext(), mock.Any[models.ActorProfile]())
}

// ── upsert stamps the active schema version ──────────────────────────────

func TestUpsertProfileStampsSchemaVersion(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()

	profiles := mock.Mock[ports.ProfileRepo]()
	captor := mock.Captor[models.ActorProfile]()
	mock.When(profiles.Upsert(mock.AnyContext(), captor.Capture())).ThenReturn(nil, nil)

	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).ThenReturn(schemaOf(&projectID, 3), nil)

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	_, err := testOps(profiles, schemas, actors, projects).UpsertProfile(testIdentity(), models.UpsertProfileInput{
		ActorID: actorID,
		Profile: json.RawMessage(`{"full_name":"Jane","display_name":"JD"}`),
	}, projectID)
	if err != nil {
		t.Fatalf("UpsertProfile: %v", err)
	}
	_, _ = mock.Verify(profiles, mock.Once()).Upsert(mock.AnyContext(), mock.Any[models.ActorProfile]())
	upserted := captor.Values()[0]
	if upserted.SchemaVersion != 3 || upserted.Outdated {
		t.Fatalf("want profile stamped v3, not outdated, got %+v", upserted)
	}
}

// ── list outdated profiles ────────────────────────────────────────────────

func TestListOutdatedProjectProfiles(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	want := []models.ActorProfile{{ActorID: uuid.New(), Profile: json.RawMessage(`{}`), SchemaVersion: 1, Outdated: true}}

	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.ListOutdated(mock.AnyContext(), mock.Equal(&projectID))).ThenReturn(want, nil)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	got, err := testOps(profiles, mock.Mock[ports.ProfileSchemaRepo](), mock.Mock[ports.ActorRepo](), projects).ListOutdatedProfiles(testIdentity(), &projectID)
	if err != nil {
		t.Fatalf("ListOutdatedProfiles: %v", err)
	}
	if len(got) != 1 || got[0].ActorID != want[0].ActorID {
		t.Fatalf("want %+v, got %+v", want, got)
	}
}

func TestListOutdatedProjectProfilesRequiresAdmin(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()

	projects := mock.Mock[ports.ProjectRepo]()
	mock.When(projects.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Equal(projectID))).
		ThenReturn(models.ProjectRoleMember, nil)

	_, err := testOps(mock.Mock[ports.ProfileRepo](), mock.Mock[ports.ProfileSchemaRepo](), mock.Mock[ports.ActorRepo](), projects).
		ListOutdatedProfiles(testIdentity(), &projectID)
	if !errors.Is(err, libauthz.ErrInsufficientRole) {
		t.Fatalf("want forbidden for member, got %v", err)
	}
}

func TestListOutdatedPlatformProfiles(t *testing.T) {
	mock.SetUp(t)
	want := []models.ActorProfile{{ActorID: uuid.New(), Profile: json.RawMessage(`{}`), SchemaVersion: 1, Outdated: true}}

	profiles := mock.Mock[ports.ProfileRepo]()
	captor := mock.Captor[*uuid.UUID]()
	mock.When(profiles.ListOutdated(mock.AnyContext(), captor.Capture())).ThenReturn(want, nil)

	got, err := testOps(profiles, mock.Mock[ports.ProfileSchemaRepo](), mock.Mock[ports.ActorRepo](), mock.Mock[ports.ProjectRepo]()).
		ListOutdatedProfiles(testIdentity(), nil)
	if len(captor.Values()) != 1 || captor.Values()[0] != nil {
		t.Fatalf("want platform scope (nil), got %v", captor.Values())
	}
	if err != nil {
		t.Fatalf("ListOutdatedProfiles: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 profile, got %+v", got)
	}
}

// ── reproduction: schema shrinks from v1 to v2 ──────────────────────────

// TestGetProfileAfterSchemaShrinksV1ToV2 pins the on-demand auto-migration
// for the reported incident:
//
//  1. active schema is v1 (old, more fields)
//  2. user profile is set fully filled out per v1
//  3. schema is updated to v2 (fewer fields)
//  4. GET the user profile
//
// Expected: the v1 document is migrated on read — the fields v2 forbids
// (tagline, location, visibility extras) are dropped, the document is
// persisted as schema_version 2 and NOT outdated. The user does nothing;
// only a missing required field falls through to admin resolution
// (outdated=true).
func TestGetProfileAfterSchemaShrinksV1ToV2(t *testing.T) {
	mock.SetUp(t)
	actorID := uuid.New()
	projectID := uuid.New()

	// step 1 — schema is v1 (old shape, more fields)
	activeSchema := schemaOf(&projectID, 1)
	activeSchema.Schema = json.RawMessage(schemaV1Full)

	schemas := mock.Mock[ports.ProfileSchemaRepo]()
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).ThenAnswer(func(_ []any) []any {
		return []any{activeSchema, nil}
	})

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	// step 2 — user profile fully filled out per v1 (single Upsert stub
	// captures every write: the step-2 stamp and the step-4 migration)
	var upserts []models.ActorProfile
	profiles := mock.Mock[ports.ProfileRepo]()
	mock.When(profiles.Upsert(mock.AnyContext(), mock.Any[models.ActorProfile]())).ThenAnswer(func(args []any) []any {
		p := args[1].(models.ActorProfile)
		upserts = append(upserts, p)
		return []any{&upserts[len(upserts)-1], nil}
	})

	ops := testOps(profiles, schemas, actors, mock.Mock[ports.ProjectRepo]())
	_, err := ops.UpsertProfile(testIdentityFor(actorID), models.UpsertProfileInput{
		ActorID: actorID,
		Profile: json.RawMessage(fullyFilledV1Profile),
	}, projectID)
	if err != nil {
		t.Fatalf("step 2 UpsertProfile: %v", err)
	}
	if len(upserts) != 1 || upserts[0].SchemaVersion != 1 {
		t.Fatalf("step 2: want stamped schema_version=1, got %+v", upserts)
	}

	// step 3 — schema updated to v2 (fewer fields)
	activeSchema = schemaOf(&projectID, 2)
	activeSchema.Schema = json.RawMessage(schemaV2New)

	// step 4 — GET the profile: must auto-migrate the document to v2
	stored := &models.ActorProfile{
		ActorID:       actorID,
		Profile:       json.RawMessage(fullyFilledV1Profile),
		SchemaVersion: 1,
	}
	mock.When(profiles.Get(mock.AnyContext(), mock.Equal(actorID))).ThenReturn(stored, nil)

	got, err := ops.GetProfile(testIdentityFor(actorID), actorID, projectID)
	if err != nil {
		t.Fatalf("step 4 GetProfile: %v", err)
	}
	if got.SchemaVersion != 2 || got.Outdated {
		t.Fatalf("step 4: want migrated profile v2 not outdated, got %+v", got)
	}

	// the document was rewritten on read: v1-only fields dropped, kept intact
	if len(upserts) != 2 {
		t.Fatalf("step 4: want the migration to persist the pruned doc, upserts=%+v", upserts)
	}
	migrated := upserts[1]
	if migrated.SchemaVersion != 2 || migrated.Outdated {
		t.Fatalf("step 4: want persisted v2 not outdated, got %+v", migrated)
	}
	var doc map[string]any
	err = json.Unmarshal(migrated.Profile, &doc)
	if err != nil {
		t.Fatalf("unmarshal migrated profile: %v", err)
	}
	if _, ok := doc["tagline"]; ok {
		t.Fatalf("migrated doc must drop tagline, got %s", migrated.Profile)
	}
	if _, ok := doc["location"]; ok {
		t.Fatalf("migrated doc must drop location, got %s", migrated.Profile)
	}
	vis, _ := doc["visibility"].(map[string]any)
	if _, ok := vis["hideSocials"]; ok {
		t.Fatalf("migrated doc must drop hideSocials, got %s", migrated.Profile)
	}
	if vis["hideLegalName"] != true {
		t.Fatalf("migrated doc must keep hideLegalName, got %s", migrated.Profile)
	}
	if doc["legalName"] != "Jane Doe" {
		t.Fatalf("migrated doc must keep legalName, got %s", migrated.Profile)
	}
}
