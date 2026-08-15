package api_keys

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/api_keys"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

// authTestOps bundles the per-test mocks backing Authenticate. The HMAC
// secret is shared with the generated key, so verification is real crypto,
// not another mock.
type authTestOps struct {
	apiKeys  ports.APIKeysRepo
	actors   ports.ActorRepo
	caps     ports.CapabilityRepo
	projects ports.ProjectRepo
	orgs     ports.OrganizationRepo
	ops      *Operations
}

func newAuthTestOps(t *testing.T, secret []byte) *authTestOps {
	t.Helper()
	mock.SetUp(t)
	o := &authTestOps{
		apiKeys:  mock.Mock[ports.APIKeysRepo](),
		actors:   mock.Mock[ports.ActorRepo](),
		caps:     mock.Mock[ports.CapabilityRepo](),
		projects: mock.Mock[ports.ProjectRepo](),
		orgs:     mock.Mock[ports.OrganizationRepo](),
	}
	o.ops = NewOperations(
		secret,
		o.actors,
		o.apiKeys,
		o.caps,
		o.projects,
		authz.New(o.orgs, o.projects, mock.Mock[ports.PlatformRolesRepo]()),
	)
	return o
}

// mintedKey generates a real key with the given secret and returns the raw
// string plus the stored row shape Authenticate looks up by prefix.
func mintedKey(t *testing.T, secret []byte, expiresAt *time.Time) (string, models.APIKey, models.Actor) {
	t.Helper()
	generated, err := api_keys.GenerateAPIKey("acme", "prod", secret)
	if err != nil {
		t.Fatalf("GenerateAPIKey: %v", err)
	}
	projectID := uuid.New()
	email := "svc@acme.example"
	actor := models.Actor{
		ID: projectID, ProjectID: &projectID, Email: &email,
		Type: models.MachineActorType,
	}
	apiKey := models.APIKey{
		ID: uuid.New(), SubjectID: actor.ID, DisplayPrefix: generated.DisplayPrefix,
		KeyHash: generated.Hash, ExpiresAt: expiresAt,
	}
	return generated.Raw, apiKey, actor
}

func TestAuthenticateValidKeyShapesIdentity(t *testing.T) {
	secret := []byte("test-hmac-secret")
	raw, apiKey, actor := mintedKey(t, secret, nil)
	o := newAuthTestOps(t, secret)
	mock.When(o.apiKeys.GetByPrefix(mock.AnyContext(), mock.Equal(apiKey.DisplayPrefix))).ThenReturn(&apiKey, nil)
	mock.When(o.actors.GetByID(mock.AnyContext(), mock.Equal(actor.ID))).ThenReturn(&actor, nil)
	mock.When(o.caps.ListByAPIKeyPrefix(mock.AnyContext(), mock.Equal(apiKey.DisplayPrefix))).ThenReturn([]models.Capability{
		{Resource: "projects", Action: "read"},
		{Resource: "profiles", Action: "write"},
	}, nil)

	ident, err := o.ops.Authenticate(context.Background(), raw)
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if ident.Sub.ID != actor.ID {
		t.Fatalf("sub.id = %s, want %s", ident.Sub.ID, actor.ID)
	}
	if ident.Sub.ProjectID == nil || *ident.Sub.ProjectID != *actor.ProjectID {
		t.Fatalf("sub.project_id = %v, want %v", ident.Sub.ProjectID, actor.ProjectID)
	}
	if ident.Cred.Type != models.APIKeyCredentialType {
		t.Fatalf("cred type = %q, want %q", ident.Cred.Type, models.APIKeyCredentialType)
	}
	if ident.Cred.ID == nil || *ident.Cred.ID != apiKey.ID {
		t.Fatalf("cred id = %v, want %v", ident.Cred.ID, apiKey.ID)
	}
	if ident.Cred.Raw != raw {
		t.Fatalf("cred raw = %q, want the presented key", ident.Cred.Raw)
	}
	wantCaps, _ := json.Marshal([]string{"projects:read", "profiles:write"})
	if string(ident.Sub.Capabilities) != string(wantCaps) {
		t.Fatalf("capabilities = %s, want %s", ident.Sub.Capabilities, wantCaps)
	}
}

func TestAuthenticateNoCapabilitiesYieldsEmptyList(t *testing.T) {
	secret := []byte("test-hmac-secret")
	raw, apiKey, actor := mintedKey(t, secret, nil)
	o := newAuthTestOps(t, secret)
	mock.When(o.apiKeys.GetByPrefix(mock.AnyContext(), mock.Equal(apiKey.DisplayPrefix))).ThenReturn(&apiKey, nil)
	mock.When(o.actors.GetByID(mock.AnyContext(), mock.Equal(actor.ID))).ThenReturn(&actor, nil)
	mock.When(o.caps.ListByAPIKeyPrefix(mock.AnyContext(), mock.Equal(apiKey.DisplayPrefix))).ThenReturn([]models.Capability{}, nil)

	ident, err := o.ops.Authenticate(context.Background(), raw)
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if string(ident.Sub.Capabilities) != "[]" {
		t.Fatalf("capabilities = %s, want []", ident.Sub.Capabilities)
	}
}

func TestAuthenticateHmacMismatchRejected(t *testing.T) {
	secret := []byte("test-hmac-secret")
	raw, _, actor := mintedKey(t, secret, nil)
	o := newAuthTestOps(t, secret)
	parsed, err := api_keys.ParseAPIKey(raw)
	if err != nil {
		t.Fatalf("ParseAPIKey: %v", err)
	}
	// The stored row is keyed by the raw key's own prefix, but its hash was
	// computed with a different secret, so verification must fail.
	other, err := api_keys.GenerateAPIKey("acme", "prod", []byte("other-secret"))
	if err != nil {
		t.Fatalf("GenerateAPIKey: %v", err)
	}
	mock.When(o.apiKeys.GetByPrefix(mock.AnyContext(), mock.Equal(parsed.DisplayPrefix))).
		ThenReturn(&models.APIKey{ID: uuid.New(), SubjectID: actor.ID, DisplayPrefix: parsed.DisplayPrefix, KeyHash: other.Hash}, nil)

	_, err = o.ops.Authenticate(context.Background(), raw)
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestAuthenticateExpiredKeyRejected(t *testing.T) {
	secret := []byte("test-hmac-secret")
	expiresAt := time.Now().Add(-time.Minute)
	raw, apiKey, actor := mintedKey(t, secret, &expiresAt)
	o := newAuthTestOps(t, secret)
	mock.When(o.apiKeys.GetByPrefix(mock.AnyContext(), mock.Equal(apiKey.DisplayPrefix))).ThenReturn(&apiKey, nil)
	mock.When(o.actors.GetByID(mock.AnyContext(), mock.Equal(actor.ID))).ThenReturn(&actor, nil)

	_, err := o.ops.Authenticate(context.Background(), raw)
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestAuthenticateUnknownPrefixRejected(t *testing.T) {
	secret := []byte("test-hmac-secret")
	o := newAuthTestOps(t, secret)
	mock.When(o.apiKeys.GetByPrefix(mock.AnyContext(), mock.Any[string]())).
		ThenReturn(nil, fun.ErrNotFound("no key"))

	_, err := o.ops.Authenticate(context.Background(), "idk-anything")
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestAuthenticateInvalidRawKeyRejected(t *testing.T) {
	o := newAuthTestOps(t, []byte("test-hmac-secret"))

	_, err := o.ops.Authenticate(context.Background(), "not-a-key")
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestAuthenticateUnknownActorRejected(t *testing.T) {
	secret := []byte("test-hmac-secret")
	raw, apiKey, _ := mintedKey(t, secret, nil)
	o := newAuthTestOps(t, secret)
	mock.When(o.apiKeys.GetByPrefix(mock.AnyContext(), mock.Equal(apiKey.DisplayPrefix))).ThenReturn(&apiKey, nil)
	mock.When(o.actors.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(nil, fun.ErrNotFound("actor gone"))

	_, err := o.ops.Authenticate(context.Background(), raw)
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestAuthenticateCapabilityErrorSurfacesRaw(t *testing.T) {
	secret := []byte("test-hmac-secret")
	raw, apiKey, actor := mintedKey(t, secret, nil)
	o := newAuthTestOps(t, secret)
	mock.When(o.apiKeys.GetByPrefix(mock.AnyContext(), mock.Equal(apiKey.DisplayPrefix))).ThenReturn(&apiKey, nil)
	mock.When(o.actors.GetByID(mock.AnyContext(), mock.Equal(actor.ID))).ThenReturn(&actor, nil)
	mock.When(o.caps.ListByAPIKeyPrefix(mock.AnyContext(), mock.Equal(apiKey.DisplayPrefix))).
		ThenReturn(nil, fun.ErrInternal("db down"))

	_, err := o.ops.Authenticate(context.Background(), raw)
	if !fun.Is(err, fun.CodeInternal) {
		t.Fatalf("want internal, got %v", err)
	}
}
