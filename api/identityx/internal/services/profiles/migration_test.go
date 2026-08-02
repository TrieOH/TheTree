package profiles

import (
	"context"
	"encoding/json"
	"testing"

	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

// ── helpers ───────────────────────────────────────────────────────────────

const schemaV2 = `{"type":"object","required":["full_name","display_name"],"properties":{"full_name":{"type":"string"},"display_name":{"type":"string"}}}`

func testOps(profiles ports.ProfileRepo, schemas ports.ProfileSchemaRepo, actors ports.ActorRepo, projects ports.ProjectRepo) *Operations {
	return NewOperations(profiles, schemas, actors, authz.New(nil, projects))
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
	mock.Verify(profiles, mock.Once()).SetMigrationState(mock.AnyContext(), mock.Equal(actorID), mock.Equal(2), mock.Equal(false))
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
	mock.Verify(profiles, mock.Once()).SetMigrationState(mock.AnyContext(), mock.Equal(actorID), mock.Equal(1), mock.Equal(true))
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
	mock.Verify(profiles, mock.Never()).SetMigrationState(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[int](), mock.Any[bool]())
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
	mock.When(schemas.Get(mock.AnyContext(), mock.Equal(&projectID))).ThenReturn(nil, fun.ErrNotFound("project_profile_schema not found"))

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, actorID, projectID)

	projects := mock.Mock[ports.ProjectRepo]()
	mockProjectAdmin(projects)

	got, err := testOps(profiles, schemas, actors, projects).GetProfile(testIdentity(), actorID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	mock.Verify(profiles, mock.Never()).SetMigrationState(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[int](), mock.Any[bool]())
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
	mock.Verify(profiles, mock.Never()).SetMigrationState(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[int](), mock.Any[bool]())
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
	mock.Verify(profiles, mock.Never()).SetMigrationState(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[int](), mock.Any[bool]())
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
	mock.Verify(schemas, mock.Once()).Get(mock.AnyContext(), mock.Any[*uuid.UUID]())
	if captor.Values()[0] != nil {
		t.Fatalf("platform scope must query schema with nil project_id, got %v", captor.Values()[0])
	}
	mock.Verify(profiles, mock.Once()).SetMigrationState(mock.AnyContext(), mock.Equal(actorID), mock.Equal(2), mock.Equal(false))
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
		ThenReturn(nil, fun.ErrNotFound("project_profile_schema not found"))

	got, err := testOps(profiles, schemas, actors, projects).
		GetProfile(testIdentityFor(actorID), actorID, projectID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	mock.Verify(projects, mock.Never()).GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())
	if got.ActorID != actorID {
		t.Fatalf("want own profile, got %+v", got)
	}
}

func TestGetProfileDeniesProjectUserReadingOthers(t *testing.T) {
	mock.SetUp(t)
	requesterID := uuid.New()
	targetID := uuid.New()
	projectID := uuid.New()

	actors := mock.Mock[ports.ActorRepo]()
	mockProjectActor(actors, targetID, projectID)

	// requester is a project user with no membership row -> GetRole not found
	projects := mock.Mock[ports.ProjectRepo]()
	mock.When(projects.GetRole(mock.AnyContext(), mock.Equal(requesterID), mock.Equal(projectID))).
		ThenReturn(models.ProjectRole(""), fun.ErrNotFound("project member not found"))

	_, err := testOps(mock.Mock[ports.ProfileRepo](), mock.Mock[ports.ProfileSchemaRepo](), actors, projects).
		GetProfile(testIdentityFor(requesterID), targetID, projectID)
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden for project user reading another's profile, got %v", err)
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
	if !fun.Is(err, fun.CodeForbidden) {
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
	if !fun.Is(err, fun.CodeForbidden) {
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
	if !fun.Is(err, fun.CodeForbidden) {
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
	if !fun.Is(err, fun.CodeForbidden) {
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
	mock.Verify(profiles, mock.Once()).Upsert(mock.AnyContext(), mock.Any[models.ActorProfile]())
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
	mock.Verify(profiles, mock.Once()).Upsert(mock.AnyContext(), mock.Any[models.ActorProfile]())
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
	if !fun.Is(err, fun.CodeForbidden) {
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
