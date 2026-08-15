package oauth_providers

import (
	"context"
	"testing"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

// testOps bundles the per-test mockio mocks backing the operations.
type testOps struct {
	providers ports.ProjectOAuthProvidersRepo
	projects  ports.ProjectRepo
	orgs      ports.OrganizationRepo
	ops       *Operations
}

// newTestOps creates a fresh set of per-test mocks and wires the
// operations over them, with the project role resolved by authz.
func newTestOps(t *testing.T, role models.ProjectRole) *testOps {
	t.Helper()
	testEnv(t)
	mock.SetUp(t)
	o := &testOps{
		providers: mock.Mock[ports.ProjectOAuthProvidersRepo](),
		projects:  mock.Mock[ports.ProjectRepo](),
		orgs:      mock.Mock[ports.OrganizationRepo](),
	}
	mock.When(o.projects.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(role, nil)
	o.ops = NewOperations(o.providers, o.projects, authz.New(o.orgs, o.projects))
	return o
}

func testEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ENCRYPTION_KEY", "abababababababababababababababab")
}

func ctxWithMember() context.Context {
	email := "member@trieoh.com"
	actor := models.Actor{ID: uuid.New(), Email: &email, Type: models.HumanActorType}
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub:  models.Subject{ID: actor.ID, ProjectID: nil, Email: actor.Email, Type: actor.Type},
		Cred: models.Credential{Type: models.TokenCredentialType},
	})
}

func providerRow(id, projectID uuid.UUID, enabled bool) models.ProjectOAuthProviders {
	return models.ProjectOAuthProviders{
		ID: id, ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "client-id", Enabled: enabled,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

// ── CRUD ────────────────────────────────────────────────────────────────

func TestCreateEncryptsSecret(t *testing.T) {
	projectID := uuid.New()
	o := newTestOps(t, models.ProjectRoleAdmin)
	mock.When(o.providers.Create(mock.AnyContext(), mock.Any[models.ProjectOAuthProviders]())).
		ThenAnswer(func(args []any) []any {
			p := args[1].(models.ProjectOAuthProviders)
			p.ID = uuid.New()
			return []any{&p, nil}
		})

	out, err := o.ops.Create(ctxWithMember(), models.CreateOAuthProviderInput{
		ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "client-id", ClientSecret: "super-secret",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if out.ClientID != "client-id" || out.Provider != models.GoogleIdentityProvider {
		t.Fatalf("unexpected output: %+v", out)
	}

	captor := mock.Captor[models.ProjectOAuthProviders]()
	_, _ = mock.Verify(o.providers, mock.Times(1)).Create(mock.AnyContext(), captor.Capture())
	if len(captor.Values()) != 1 {
		t.Fatalf("created = %d, want 1", len(captor.Values()))
	}
	stored := captor.Values()[0]
	if stored.EncryptedClientSecret == "super-secret" {
		t.Fatal("secret must not be stored in plaintext")
	}
	decrypted, err := crypto.DecryptPrivateKey(stored.EncryptedClientSecret)
	if err != nil {
		t.Fatalf("decrypt stored secret: %v", err)
	}
	if string(decrypted) != "super-secret" {
		t.Fatalf("round-trip mismatch: %q", decrypted)
	}
	if !stored.Enabled {
		t.Fatal("created provider must default to enabled")
	}
}

func TestCreateRequiresAdmin(t *testing.T) {
	projectID := uuid.New()
	o := newTestOps(t, models.ProjectRoleMember)

	_, err := o.ops.Create(ctxWithMember(), models.CreateOAuthProviderInput{
		ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "client-id", ClientSecret: "super-secret",
	})
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden for member, got %v", err)
	}
	_, _ = mock.Verify(o.providers, mock.Times(0)).Create(mock.AnyContext(), mock.Any[models.ProjectOAuthProviders]())
}

func TestCreateUnknownProjectIsForbidden(t *testing.T) {
	projectID := uuid.New()
	o := newTestOps(t, models.ProjectRoleMember)

	_, err := o.ops.Create(ctxWithMember(), models.CreateOAuthProviderInput{
		ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "client-id", ClientSecret: "super-secret",
	})
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden for non-member, got %v", err)
	}
}

func TestListByProjectReturnsNoSecret(t *testing.T) {
	projectID := uuid.New()
	o := newTestOps(t, models.ProjectRoleMember)
	row := providerRow(uuid.New(), projectID, true)
	row.EncryptedClientSecret = "encrypted-blob"
	mock.When(o.providers.ListByProject(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn([]models.ProjectOAuthProviders{row}, nil)

	out, err := o.ops.ListByProject(ctxWithMember(), projectID)
	if err != nil {
		t.Fatalf("ListByProject: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("list = %d, want 1", len(out))
	}
	if out[0].ClientID != "client-id" || !out[0].Enabled {
		t.Fatalf("unexpected output: %+v", out[0])
	}
}

func TestUpdateClientSecretEncrypts(t *testing.T) {
	id, projectID := uuid.New(), uuid.New()
	o := newTestOps(t, models.ProjectRoleAdmin)
	row := providerRow(id, projectID, true)
	row.EncryptedClientSecret = "old-blob"
	mock.When(o.providers.GetByID(mock.AnyContext(), mock.Equal(id))).ThenReturn(&row, nil)
	mock.When(o.providers.UpdateClientSecret(mock.AnyContext(), mock.Equal(id), mock.Any[string]())).
		ThenAnswer(func(args []any) []any {
			encrypted := args[2].(string)
			updated := providerRow(id, projectID, true)
			updated.EncryptedClientSecret = encrypted
			return []any{&updated, nil}
		})

	newSecretValue := "new-secret"
	out, err := o.ops.Update(ctxWithMember(), models.UpdateOAuthProviderInput{
		ID: id, ClientSecret: &newSecretValue,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if out.ClientID != "client-id" {
		t.Fatalf("client id must be untouched, got %q", out.ClientID)
	}

	captor := mock.Captor[string]()
	_, _ = mock.Verify(o.providers, mock.Times(1)).UpdateClientSecret(mock.AnyContext(), mock.Equal(id), captor.Capture())
	encrypted := captor.Values()[0]
	if encrypted == "new-secret" {
		t.Fatal("secret must be encrypted before storage")
	}
	decrypted, err := crypto.DecryptPrivateKey(encrypted)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if string(decrypted) != "new-secret" {
		t.Fatalf("round-trip mismatch: %q", decrypted)
	}
}

func TestUpdateRequiresAdmin(t *testing.T) {
	id, projectID := uuid.New(), uuid.New()
	o := newTestOps(t, models.ProjectRoleMember)
	row := providerRow(id, projectID, true)
	mock.When(o.providers.GetByID(mock.AnyContext(), mock.Equal(id))).ThenReturn(&row, nil)

	clientID := "x"
	_, err := o.ops.Update(ctxWithMember(), models.UpdateOAuthProviderInput{ID: id, ClientID: &clientID})
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden for member, got %v", err)
	}
}

func TestSetEnabledToggles(t *testing.T) {
	id, projectID := uuid.New(), uuid.New()
	o := newTestOps(t, models.ProjectRoleAdmin)
	row := providerRow(id, projectID, true)
	mock.When(o.providers.GetByID(mock.AnyContext(), mock.Equal(id))).ThenReturn(&row, nil)
	mock.When(o.providers.SetEnabled(mock.AnyContext(), mock.Equal(id), mock.Any[bool]())).
		ThenAnswer(func(args []any) []any {
			enabled := args[2].(bool)
			updated := providerRow(id, projectID, enabled)
			return []any{&updated, nil}
		})

	disabled, err := o.ops.SetEnabled(ctxWithMember(), id, false)
	if err != nil {
		t.Fatalf("SetEnabled(false): %v", err)
	}
	if disabled.Enabled {
		t.Fatal("want disabled")
	}
	enabled, err := o.ops.SetEnabled(ctxWithMember(), id, true)
	if err != nil {
		t.Fatalf("SetEnabled(true): %v", err)
	}
	if !enabled.Enabled {
		t.Fatal("want enabled")
	}
}

func TestDeleteHardDeletes(t *testing.T) {
	id, projectID := uuid.New(), uuid.New()
	o := newTestOps(t, models.ProjectRoleAdmin)
	row := providerRow(id, projectID, true)
	mock.When(o.providers.GetByID(mock.AnyContext(), mock.Equal(id))).ThenReturn(&row, nil)

	err := o.ops.Delete(ctxWithMember(), id)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_ = mock.Verify(o.providers, mock.Times(1)).Delete(mock.AnyContext(), mock.Equal(id))
}

func TestDeleteRequiresAdmin(t *testing.T) {
	id, projectID := uuid.New(), uuid.New()
	o := newTestOps(t, models.ProjectRoleMember)
	row := providerRow(id, projectID, true)
	mock.When(o.providers.GetByID(mock.AnyContext(), mock.Equal(id))).ThenReturn(&row, nil)

	err := o.ops.Delete(ctxWithMember(), id)
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden for member, got %v", err)
	}
	_ = mock.Verify(o.providers, mock.Times(0)).Delete(mock.AnyContext(), mock.Equal(id))
}

// ── Discovery ────────────────────────────────────────────────────────────

func TestListEnabledProvidersProject(t *testing.T) {
	projectID := uuid.New()
	o := newTestOps(t, models.ProjectRoleMember)
	mock.When(o.projects.GetByID(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(&models.Project{ID: projectID}, nil)
	mock.When(o.providers.ListByProject(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn([]models.ProjectOAuthProviders{
			providerRow(uuid.New(), projectID, true),
			{ID: uuid.New(), ProjectID: projectID, Provider: models.GithubIdentityProvider, Enabled: false},
		}, nil)

	out, err := o.ops.ListEnabledProviders(ctxWithMember(), &projectID)
	if err != nil {
		t.Fatalf("ListEnabledProviders: %v", err)
	}
	if len(out) != 1 || out[0] != models.GoogleIdentityProvider {
		t.Fatalf("want [google], got %v", out)
	}
}

func TestListEnabledProvidersUnknownProject(t *testing.T) {
	projectID := uuid.New()
	o := newTestOps(t, models.ProjectRoleMember)
	mock.When(o.projects.GetByID(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(nil, fun.ErrNotFound("project not found"))

	_, err := o.ops.ListEnabledProviders(ctxWithMember(), &projectID)
	if !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("want not found, got %v", err)
	}
}

func TestListEnabledProvidersPlatform(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "platform-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "platform-secret")
	t.Setenv("GITHUB_CLIENT_ID", "")
	t.Setenv("GITHUB_CLIENT_SECRET", "")
	o := newTestOps(t, models.ProjectRoleMember)

	out, err := o.ops.ListEnabledProviders(ctxWithMember(), nil)
	if err != nil {
		t.Fatalf("ListEnabledProviders: %v", err)
	}
	if len(out) != 1 || out[0] != models.GoogleIdentityProvider {
		t.Fatalf("want [google] (only env-configured), got %v", out)
	}
}

// ── Login resolution ─────────────────────────────────────────────────────

func encryptedSecret(t *testing.T) string {
	t.Helper()
	enc, err := crypto.EncryptPrivateKey([]byte("proj-secret"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	return enc
}

func TestResolveLoginProviderPlatformUsesEnvCredentials(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "platform-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "platform-secret")
	o := newTestOps(t, models.ProjectRoleMember)

	out, err := o.ops.ResolveLoginProvider(context.Background(), "google", nil)
	if err != nil {
		t.Fatalf("ResolveLoginProvider: %v", err)
	}
	if out.Creds.ClientID != "platform-id" || out.Creds.ClientSecret != "platform-secret" {
		t.Fatalf("unexpected credentials: %+v", out.Creds)
	}
	if out.Disabled {
		t.Fatal("platform credentials are never disabled")
	}
	if out.Provider != models.GoogleIdentityProvider {
		t.Fatalf("provider = %v, want google", out.Provider)
	}
}

func TestResolveLoginProviderPlatformNotConfigured(t *testing.T) {
	o := newTestOps(t, models.ProjectRoleMember)

	_, err := o.ops.ResolveLoginProvider(context.Background(), "google", nil)
	if !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request for missing env credentials, got %v", err)
	}
}

func TestResolveLoginProviderProjectRow(t *testing.T) {
	projectID := uuid.New()
	o := newTestOps(t, models.ProjectRoleMember)
	row := models.ProjectOAuthProviders{
		ID: uuid.New(), ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "proj-client-id", EncryptedClientSecret: encryptedSecret(t), Enabled: true,
		CallbackURL: "https://app.example.com/callback",
	}
	mock.When(o.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(&row, nil)

	out, err := o.ops.ResolveLoginProvider(context.Background(), "google", &projectID)
	if err != nil {
		t.Fatalf("ResolveLoginProvider: %v", err)
	}
	if out.Creds.ClientID != "proj-client-id" || out.Creds.ClientSecret != "proj-secret" {
		t.Fatalf("unexpected credentials: %+v", out.Creds)
	}
	if out.Creds.RedirectURL != "https://app.example.com/callback" {
		t.Fatalf("redirect = %q, want the project callback_url", out.Creds.RedirectURL)
	}
	if out.Disabled {
		t.Fatal("enabled row must resolve as enabled")
	}
}

func TestResolveLoginProviderProjectDisabled(t *testing.T) {
	projectID := uuid.New()
	o := newTestOps(t, models.ProjectRoleMember)
	row := models.ProjectOAuthProviders{
		ID: uuid.New(), ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "proj-client-id", EncryptedClientSecret: encryptedSecret(t), Enabled: false,
	}
	mock.When(o.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(&row, nil)

	out, err := o.ops.ResolveLoginProvider(context.Background(), "google", &projectID)
	if err != nil {
		t.Fatalf("ResolveLoginProvider: %v", err)
	}
	if !out.Disabled {
		t.Fatal("disabled row must resolve as disabled")
	}
}

func TestResolveLoginProviderMissingRow(t *testing.T) {
	projectID := uuid.New()
	o := newTestOps(t, models.ProjectRoleMember)
	mock.When(o.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(nil, fun.ErrNotFound("not configured"))

	_, err := o.ops.ResolveLoginProvider(context.Background(), "google", &projectID)
	if !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("want not found for missing row, got %v", err)
	}
}
