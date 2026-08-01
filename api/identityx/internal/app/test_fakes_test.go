package app

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"IdentityX/models"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func testEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ENCRYPTION_KEY", strings.Repeat("ab", 32))
}

func testActor() models.Actor {
	email := "actor@trieoh.com"
	return models.Actor{ID: uuid.New(), Email: &email, Type: models.HumanActorType}
}

// ── in-memory fakes ───────────────────────────────────────────────────────

type fakeCryptoKeysRepo struct {
	keys map[uuid.UUID]models.CryptoKey
}

func (f *fakeCryptoKeysRepo) GetActive(_ context.Context, keyType models.CryptoKeyType, projectID *uuid.UUID) (*models.CryptoKey, error) {
	for _, k := range f.keys {
		if k.Type == keyType && k.Status == models.CryptoKeyStatusActive {
			return &k, nil
		}
	}
	return nil, fun.ErrNotFound("no active key")
}
func (f *fakeCryptoKeysRepo) GetByID(_ context.Context, id uuid.UUID) (*models.CryptoKey, error) {
	k, ok := f.keys[id]
	if !ok {
		return nil, fun.ErrNotFound("key not found")
	}
	return &k, nil
}
func (f *fakeCryptoKeysRepo) GetActiveSigningKeys(_ context.Context, projectID *uuid.UUID) ([]models.ActiveSigningKey, error) {
	return nil, nil
}
func (f *fakeCryptoKeysRepo) Create(_ context.Context, projectID *uuid.UUID, pair *crypto.KeyPair, keyType string) (*models.CryptoKey, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}

type fakeAPIKeysRepo struct{}

func (f *fakeAPIKeysRepo) Create(context.Context, models.APIKey) (*models.APIKey, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}
func (f *fakeAPIKeysRepo) GetByPrefix(context.Context, string) (*models.APIKey, error) {
	return nil, fun.ErrNotFound("no api key")
}

type fakeActorsRepo struct{ actor models.Actor }

func (f *fakeActorsRepo) Register(context.Context, models.Actor) (*models.Actor, error) {
	return &f.actor, nil
}
func (f *fakeActorsRepo) GetByEmail(context.Context, string, *uuid.UUID) (*models.Actor, error) {
	return &f.actor, nil
}
func (f *fakeActorsRepo) GetByID(context.Context, uuid.UUID) (*models.Actor, error) {
	return &f.actor, nil
}
func (f *fakeActorsRepo) List(context.Context, uuid.UUID) ([]models.Actor, error) { return nil, nil }
func (f *fakeActorsRepo) GetProjectServiceAccount(context.Context, uuid.UUID) (*models.Actor, error) {
	return nil, fun.ErrNotFound("no service account")
}
func (f *fakeActorsRepo) UpdateLastLoginAt(context.Context, uuid.UUID) error { return nil }

type fakeCapabilitiesRepo struct{}

func (f *fakeCapabilitiesRepo) Create(context.Context, models.Capability) (*models.Capability, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}
func (f *fakeCapabilitiesRepo) List(context.Context, uuid.UUID) ([]models.Capability, error) {
	return nil, nil
}
func (f *fakeCapabilitiesRepo) ValidateCapabilities(context.Context, *uuid.UUID, []uuid.UUID) (bool, error) {
	return true, nil
}
func (f *fakeCapabilitiesRepo) AssignToAPIKey(context.Context, uuid.UUID, []uuid.UUID, uuid.UUID) error {
	return nil
}
func (f *fakeCapabilitiesRepo) ListByAPIKeyPrefix(context.Context, string) ([]models.Capability, error) {
	return nil, nil
}

type fakeBlacklistRepo struct {
	entries []models.BlacklistEntry
}

func (f *fakeBlacklistRepo) Append(_ context.Context, e models.BlacklistEntry) error {
	f.entries = append(f.entries, e)
	return nil
}
func (f *fakeBlacklistRepo) GetByTarget(_ context.Context, target string) (*models.BlacklistEntry, error) {
	for i := range f.entries {
		if f.entries[i].Target == target {
			return &f.entries[i], nil
		}
	}
	return nil, fun.ErrNotFound("not blacklisted")
}
func (f *fakeBlacklistRepo) GetByTargetAndType(_ context.Context, target string, entryType models.BlacklistEntryType) (*models.BlacklistEntry, error) {
	for i := range f.entries {
		if f.entries[i].Target == target && f.entries[i].Type == entryType {
			return &f.entries[i], nil
		}
	}
	return nil, fun.ErrNotFound("not blacklisted")
}

type fakeProjectsRepo struct{}

func (f *fakeProjectsRepo) Create(context.Context, models.Project) (*models.Project, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}
func (f *fakeProjectsRepo) GetByID(context.Context, uuid.UUID) (*models.Project, error) {
	return nil, fun.ErrNotFound("no project")
}
func (f *fakeProjectsRepo) ListByOrganization(context.Context, uuid.UUID) ([]models.Project, error) {
	return nil, nil
}
func (f *fakeProjectsRepo) ListJoined(context.Context, uuid.UUID) ([]models.Project, error) {
	return nil, nil
}
func (f *fakeProjectsRepo) ListOwned(context.Context, uuid.UUID) ([]models.Project, error) {
	return nil, nil
}
func (f *fakeProjectsRepo) AddMember(context.Context, models.ProjectMember) error    { return nil }
func (f *fakeProjectsRepo) RemoveMember(context.Context, uuid.UUID, uuid.UUID) error { return nil }
func (f *fakeProjectsRepo) GetMember(context.Context, uuid.UUID, uuid.UUID) (*models.ProjectMember, error) {
	return nil, fun.ErrNotFound("no member")
}
func (f *fakeProjectsRepo) GetRole(context.Context, uuid.UUID, uuid.UUID) (models.ProjectRole, error) {
	return models.ProjectRoleMember, nil
}
func (f *fakeProjectsRepo) ListMembers(context.Context, uuid.UUID) ([]models.ProjectMember, error) {
	return nil, nil
}

type fakePlatformRolesRepo struct{}

func (f *fakePlatformRolesRepo) Give(context.Context, uuid.UUID, models.PlatformRole, *json.RawMessage) (*models.PlatformRoleRelation, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}

type fakeExternalIdentitiesRepo struct{}

func (f *fakeExternalIdentitiesRepo) GetByProviderAndSubject(context.Context, string, string) (*models.ActorExternalIdentities, error) {
	return nil, fun.ErrNotFound("no identity")
}
func (f *fakeExternalIdentitiesRepo) Create(context.Context, models.ActorExternalIdentities) (*models.ActorExternalIdentities, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}
func (f *fakeExternalIdentitiesRepo) UpdateTokens(context.Context, models.ActorExternalIdentities) (*models.ActorExternalIdentities, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}

// ── token minting ──────────────────────────────────────────────────────────

func mintAccessToken(t *testing.T, kp *crypto.KeyPair, keyID, jti uuid.UUID, actor models.Actor) string {
	t.Helper()
	claims := models.AccessClaims{
		Sub: models.AccessSub{ID: actor.ID, ProjectID: actor.ProjectID, Email: actor.Email, Type: actor.Type},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "test-issuer", ID: jti.String(), IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = keyID
	payload, err := token.SigningString()
	if err != nil {
		t.Fatalf("SigningString: %v", err)
	}
	signed, err := crypto.SignToken([]byte(payload), kp)
	if err != nil {
		t.Fatalf("SignToken: %v", err)
	}
	return signed
}

// newSigningKey mints a fresh signing key pair and returns it plus its ID.
func newSigningKey(t *testing.T) (*crypto.KeyPair, uuid.UUID) {
	t.Helper()
	testEnv(t)
	kp, err := crypto.GenerateKeyPair("signing")
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	return kp, uuid.New()
}
