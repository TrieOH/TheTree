package authn

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

type fakeBlacklist struct {
	entries []models.BlacklistEntry
}

func (f *fakeBlacklist) Append(_ context.Context, e models.BlacklistEntry) error {
	f.entries = append(f.entries, e)
	return nil
}

func (f *fakeBlacklist) find(target string, entryType models.BlacklistEntryType) (*models.BlacklistEntry, error) {
	for i := range f.entries {
		if f.entries[i].Target == target && f.entries[i].Type == entryType {
			return &f.entries[i], nil
		}
	}
	return nil, fun.ErrNotFound("not blacklisted")
}

func (f *fakeBlacklist) GetByTarget(ctx context.Context, target string) (*models.BlacklistEntry, error) {
	return f.find(target, "")
}

func (f *fakeBlacklist) GetByTargetAndType(ctx context.Context, target string, entryType models.BlacklistEntryType) (*models.BlacklistEntry, error) {
	return f.find(target, entryType)
}

type fakeCryptoKeys struct {
	keys map[uuid.UUID]models.CryptoKey
}

func (f *fakeCryptoKeys) GetActive(_ context.Context, keyType models.CryptoKeyType, projectID *uuid.UUID) (*models.CryptoKey, error) {
	for _, k := range f.keys {
		if k.Type == keyType && k.Status == models.CryptoKeyStatusActive {
			return &k, nil
		}
	}
	return nil, fun.ErrNotFound("no active key")
}

func (f *fakeCryptoKeys) GetByID(_ context.Context, id uuid.UUID) (*models.CryptoKey, error) {
	k, ok := f.keys[id]
	if !ok {
		return nil, fun.ErrNotFound("key not found")
	}
	return &k, nil
}

func (f *fakeCryptoKeys) GetActiveSigningKeys(_ context.Context, projectID *uuid.UUID) ([]models.ActiveSigningKey, error) {
	return nil, nil
}

func (f *fakeCryptoKeys) Create(_ context.Context, projectID *uuid.UUID, pair *crypto.KeyPair, keyType string) (*models.CryptoKey, error) {
	return nil, nil
}

type fakeActors struct{ actor models.Actor }

func (f *fakeActors) Register(context.Context, models.Actor) (*models.Actor, error) {
	return &f.actor, nil
}
func (f *fakeActors) GetByEmail(context.Context, string, *uuid.UUID) (*models.Actor, error) {
	return &f.actor, nil
}
func (f *fakeActors) GetByID(context.Context, uuid.UUID) (*models.Actor, error) { return &f.actor, nil }
func (f *fakeActors) List(context.Context, uuid.UUID) ([]models.Actor, error)   { return nil, nil }
func (f *fakeActors) GetProjectServiceAccount(context.Context, uuid.UUID) (*models.Actor, error) {
	return nil, fun.ErrNotFound("no service account")
}
func (f *fakeActors) UpdateLastLoginAt(context.Context, uuid.UUID) error { return nil }

type fakeProjects struct{}

func (f *fakeProjects) Create(context.Context, models.Project) (*models.Project, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}
func (f *fakeProjects) GetByID(context.Context, uuid.UUID) (*models.Project, error) {
	return nil, fun.ErrNotFound("no project")
}
func (f *fakeProjects) ListByOrganization(context.Context, uuid.UUID) ([]models.Project, error) {
	return nil, nil
}
func (f *fakeProjects) ListJoined(context.Context, uuid.UUID) ([]models.Project, error) {
	return nil, nil
}
func (f *fakeProjects) ListOwned(context.Context, uuid.UUID) ([]models.Project, error) {
	return nil, nil
}
func (f *fakeProjects) AddMember(context.Context, models.ProjectMember) error    { return nil }
func (f *fakeProjects) RemoveMember(context.Context, uuid.UUID, uuid.UUID) error { return nil }
func (f *fakeProjects) GetMember(context.Context, uuid.UUID, uuid.UUID) (*models.ProjectMember, error) {
	return nil, fun.ErrNotFound("no member")
}
func (f *fakeProjects) GetRole(context.Context, uuid.UUID, uuid.UUID) (models.ProjectRole, error) {
	return models.ProjectRoleMember, nil
}
func (f *fakeProjects) ListMembers(context.Context, uuid.UUID) ([]models.ProjectMember, error) {
	return nil, nil
}

type fakePlatformRoles struct{}

func (f *fakePlatformRoles) Give(context.Context, uuid.UUID, models.PlatformRole, *json.RawMessage) (*models.PlatformRoleRelation, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}

type fakeExternalIdentities struct{}

func (f *fakeExternalIdentities) GetByProviderAndSubject(context.Context, string, string) (*models.ActorExternalIdentities, error) {
	return nil, fun.ErrNotFound("no identity")
}
func (f *fakeExternalIdentities) Create(context.Context, models.ActorExternalIdentities) (*models.ActorExternalIdentities, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}
func (f *fakeExternalIdentities) UpdateTokens(context.Context, models.ActorExternalIdentities) (*models.ActorExternalIdentities, error) {
	return nil, fun.ErrNotImplemented("unused in tests")
}

// ── token minting ──────────────────────────────────────────────────────────

// mintPayload builds the signing string (header.payload) of a token, the
// same way the service's newAccessToken/newIDXRefreshToken do.
func mintPayload(t *testing.T, claims jwt.Claims, kid uuid.UUID) []byte {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = kid
	payload, err := token.SigningString()
	if err != nil {
		t.Fatalf("SigningString: %v", err)
	}
	return []byte(payload)
}

type signedPair struct {
	accessToken, refreshToken string
	accessJTI, refreshJTI     uuid.UUID
	keyID                     uuid.UUID
	kp                        *crypto.KeyPair
}

// newTestOps wires an authn Operations over in-memory fakes with a fresh
// signing key, and mints a valid access/refresh pair signed by it.
func newTestOps(t *testing.T, bl *fakeBlacklist) (*Operations, *signedPair) {
	t.Helper()
	testEnv(t)
	kp, err := crypto.GenerateKeyPair("signing")
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	keyID := uuid.New()
	actor := testActor()
	cryptoKeys := &fakeCryptoKeys{keys: map[uuid.UUID]models.CryptoKey{
		keyID: {
			ID: keyID, Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
			PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
		},
	}}
	if bl == nil {
		bl = &fakeBlacklist{}
	}
	ops := NewOperations(&fakeActors{actor: actor}, &fakeProjects{}, &fakePlatformRoles{}, cryptoKeys, bl, &fakeExternalIdentities{})

	accessJTI, refreshJTI := uuid.New(), uuid.New()
	accessPayload := mintPayload(t, models.AccessClaims{
		Sub: models.AccessSub{ID: actor.ID, ProjectID: actor.ProjectID, Email: actor.Email, Type: actor.Type},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "test-issuer", ID: accessJTI.String(), IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, keyID)
	refreshPayload := mintPayload(t, models.RefreshClaims{
		Sub: models.RefreshSub{ID: actor.ID, ProjectID: actor.ProjectID, AccessJTI: accessJTI},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "test-issuer", ID: refreshJTI.String(), IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, keyID)
	accessToken, err := crypto.SignToken(accessPayload, kp)
	if err != nil {
		t.Fatalf("SignToken access: %v", err)
	}
	refreshToken, err := crypto.SignToken(refreshPayload, kp)
	if err != nil {
		t.Fatalf("SignToken refresh: %v", err)
	}
	return ops, &signedPair{
		accessToken: accessToken, refreshToken: refreshToken,
		accessJTI: accessJTI, refreshJTI: refreshJTI, keyID: keyID, kp: kp,
	}
}

func ctxWithIdentity() context.Context {
	actor := testActor()
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub:  models.Subject{ID: actor.ID, ProjectID: nil, Email: actor.Email, Type: actor.Type},
		Cred: models.Credential{Type: models.TokenCredentialType},
	})
}
