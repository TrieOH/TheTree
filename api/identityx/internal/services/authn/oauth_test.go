package authn

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"
	"lib/oauth"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
	"golang.org/x/oauth2"
)

// stubOAuthFlow points the google registry entry at a local provider
// serving the token exchange and userinfo calls, so tests never touch the
// network. The original entry is restored after the test.
func stubOAuthFlow(t *testing.T) {
	t.Helper()
	orig := oauth.Registry["google"]
	t.Cleanup(func() { oauth.Registry["google"] = orig })

	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.FormValue("grant_type") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"access_token": "test-access-token",
			"token_type":   "bearer",
			"expires_in":   "3600",
		})
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-access-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": "subject-1", "email": "user@example.com",
		})
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	oauth.Registry["google"] = oauth.Provider{
		Endpoint:    oauth2.Endpoint{TokenURL: srv.URL + "/token", AuthURL: srv.URL + "/auth"},
		Scopes:      []string{"email"},
		Userinfo:    srv.URL + "/userinfo",
		RedirectURL: "http://localhost/callback",
	}
}

// oauthRepos bundles the per-test mockio mocks backing the authn
// operations for the OAuth flow. Each test creates its own via newOAuthOps
// and stubs them inline.
type oauthRepos struct {
	actors    ports.ActorRepo
	projects  ports.ProjectRepo
	keys      ports.CryptoKeysRepo
	blacklist ports.BlacklistRepo
	external  ports.ExternalIdentitiesRepo
	providers ports.ProjectOAuthProvidersRepo
	states    ports.OAuthLoginStatesRepo
}

// newOAuthOps creates a fresh set of per-test mocks and wires the authn
// operations over them.
func newOAuthOps(t *testing.T) (*Operations, *oauthRepos) {
	t.Helper()
	mock.SetUp(t)
	r := &oauthRepos{
		actors:    mock.Mock[ports.ActorRepo](),
		projects:  mock.Mock[ports.ProjectRepo](),
		keys:      mock.Mock[ports.CryptoKeysRepo](),
		blacklist: mock.Mock[ports.BlacklistRepo](),
		external:  mock.Mock[ports.ExternalIdentitiesRepo](),
		providers: mock.Mock[ports.ProjectOAuthProvidersRepo](),
		states:    mock.Mock[ports.OAuthLoginStatesRepo](),
	}
	ops := NewOperations(
		r.actors, r.projects,
		mock.Mock[ports.PlatformRolesRepo](),
		r.keys, r.blacklist, r.external, r.providers, r.states,
		mock.Mock[ports.ActionTokenRepo](),
		mock.Mock[ports.EmailSender](),
		[]byte("test-hmac"),
	)
	return ops, r
}

func envState() models.OAuthLoginState {
	return models.OAuthLoginState{
		ID: uuid.New(), State: "state-token", Provider: models.GoogleIdentityProvider,
		ProjectID: nil, ExpiresAt: time.Now().Add(time.Hour),
	}
}

func projectState(projectID uuid.UUID) models.OAuthLoginState {
	return models.OAuthLoginState{
		ID: uuid.New(), State: "state-token", Provider: models.GoogleIdentityProvider,
		ProjectID: &projectID, ExpiresAt: time.Now().Add(time.Hour),
	}
}

// encryptedProjectSecret is a real encrypted "proj-secret" so the
// connect/callback path exercises decryption.
func encryptedProjectSecret(t *testing.T) string {
	t.Helper()
	enc, err := crypto.EncryptPrivateKey([]byte("proj-secret"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	return enc
}

// ── Connect ──────────────────────────────────────────────────────────────

func TestOAuthConnectPlatformUsesEnvCredentials(t *testing.T) {
	testEnv(t)
	t.Setenv("GOOGLE_CLIENT_ID", "platform-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "platform-secret")
	stubOAuthFlow(t)
	ops, r := newOAuthOps(t)
	var created []models.OAuthLoginState
	mock.When(r.states.CreateState(mock.AnyContext(), mock.Any[models.OAuthLoginState]())).
		ThenAnswer(func(args []any) []any {
			s := args[1].(models.OAuthLoginState)
			created = append(created, s)
			return []any{&s, nil}
		})

	connectURL, err := ops.OAuthConnect(context.Background(), "google", nil)
	if err != nil {
		t.Fatalf("OAuthConnect: %v", err)
	}
	parsed, err := url.Parse(connectURL)
	if err != nil {
		t.Fatalf("parse connect url: %v", err)
	}
	if got := parsed.Query().Get("client_id"); got != "platform-id" {
		t.Fatalf("client_id = %q, want %q", got, "platform-id")
	}
	state := parsed.Query().Get("state")
	if state == "" {
		t.Fatal("state param must be present")
	}
	if len(created) != 1 {
		t.Fatalf("state rows created = %d, want 1", len(created))
	}
	if created[0].ProjectID != nil {
		t.Fatalf("platform state must not carry a project, got %v", created[0].ProjectID)
	}
	if created[0].State != state {
		t.Fatalf("stored state %q != URL state %q", created[0].State, state)
	}
}

func TestOAuthConnectProjectUsesProjectCredentials(t *testing.T) {
	testEnv(t)
	stubOAuthFlow(t)
	projectID := uuid.New()
	ops, r := newOAuthOps(t)
	mock.When(r.projects.GetByID(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(&models.Project{ID: projectID}, nil)
	row := models.ProjectOAuthProviders{
		ID: uuid.New(), ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "proj-client-id", EncryptedClientSecret: encryptedProjectSecret(t), Enabled: true,
		CallbackURL: "https://myapp.example.com/auth/callback",
	}
	mock.When(r.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(&row, nil)
	var created []models.OAuthLoginState
	mock.When(r.states.CreateState(mock.AnyContext(), mock.Any[models.OAuthLoginState]())).
		ThenAnswer(func(args []any) []any {
			s := args[1].(models.OAuthLoginState)
			created = append(created, s)
			return []any{&s, nil}
		})

	connectURL, err := ops.OAuthConnect(context.Background(), "google", &projectID)
	if err != nil {
		t.Fatalf("OAuthConnect: %v", err)
	}
	parsed, err := url.Parse(connectURL)
	if err != nil {
		t.Fatalf("parse connect url: %v", err)
	}
	if got := parsed.Query().Get("client_id"); got != "proj-client-id" {
		t.Fatalf("client_id = %q, want %q", got, "proj-client-id")
	}
	if got := parsed.Query().Get("redirect_uri"); got != row.CallbackURL {
		t.Fatalf("redirect_uri = %q, want the project callback_url %q", got, row.CallbackURL)
	}
	if len(created) != 1 || created[0].ProjectID == nil || *created[0].ProjectID != projectID {
		t.Fatalf("state must carry the project, got %+v", created)
	}
}

func TestOAuthConnectDisabledProviderStillReturnsURL(t *testing.T) {
	testEnv(t)
	stubOAuthFlow(t)
	projectID := uuid.New()
	ops, r := newOAuthOps(t)
	mock.When(r.projects.GetByID(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(&models.Project{ID: projectID}, nil)
	row := models.ProjectOAuthProviders{
		ID: uuid.New(), ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "proj-client-id", EncryptedClientSecret: encryptedProjectSecret(t), Enabled: false,
	}
	mock.When(r.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(&row, nil)
	mock.When(r.states.CreateState(mock.AnyContext(), mock.Any[models.OAuthLoginState]())).
		ThenAnswer(func(args []any) []any {
			s := args[1].(models.OAuthLoginState)
			return []any{&s, nil}
		})

	_, err := ops.OAuthConnect(context.Background(), "google", &projectID)
	if err != nil {
		t.Fatalf("disabled provider must still issue a connect URL, got %v", err)
	}
}

func TestOAuthConnectProviderNotConfigured(t *testing.T) {
	testEnv(t)
	projectID := uuid.New()
	ops, r := newOAuthOps(t)
	mock.When(r.projects.GetByID(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(&models.Project{ID: projectID}, nil)
	mock.When(r.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(nil, fun.ErrNotFound("not configured"))

	_, err := ops.OAuthConnect(context.Background(), "google", &projectID)
	if !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request, got %v", err)
	}
}

func TestOAuthConnectUnknownProject(t *testing.T) {
	testEnv(t)
	projectID := uuid.New()
	ops, r := newOAuthOps(t)
	mock.When(r.projects.GetByID(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(nil, fun.ErrNotFound("project not found"))

	_, err := ops.OAuthConnect(context.Background(), "google", &projectID)
	if !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("want not found, got %v", err)
	}
}

func TestOAuthConnectUnsupportedProvider(t *testing.T) {
	testEnv(t)
	ops, _ := newOAuthOps(t)

	_, err := ops.OAuthConnect(context.Background(), "x", nil)
	if !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request, got %v", err)
	}
}

// ── Callback ─────────────────────────────────────────────────────────────

func TestOAuthCallbackSignupNewIdentity(t *testing.T) {
	testEnv(t)
	t.Setenv("GOOGLE_CLIENT_ID", "platform-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "platform-secret")
	stubOAuthFlow(t)
	pair := mintPair(t)
	ops, r := newOAuthOps(t)
	state := envState()
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Equal(state.State))).ThenReturn(&state, nil)
	mock.When(r.states.DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)
	mock.When(r.external.GetByProviderAndSubject(mock.AnyContext(), mock.Equal("google"), mock.Equal("subject-1"), mock.Equal(state.ProjectID))).
		ThenReturn(nil, fun.ErrNotFound("no identity"))
	captor := mock.Captor[models.Actor]()
	mock.When(r.actors.Register(mock.AnyContext(), captor.Capture())).
		ThenAnswer(func(args []any) []any {
			a := args[1].(models.Actor)
			a.ID = uuid.New()
			return []any{&a, nil}
		})
	mock.When(r.external.Create(mock.AnyContext(), mock.Any[models.ActorExternalIdentities]())).
		ThenAnswer(func(args []any) []any {
			e := args[1].(models.ActorExternalIdentities)
			return []any{&e, nil}
		})
	mock.When(r.keys.GetActive(mock.AnyContext(), mock.Equal(models.SigningCryptoKeyType), mock.Any[*uuid.UUID]())).
		ThenReturn(&pair.key, nil)

	out, err := ops.OAuthCallback(context.Background(), "google", "code", "state-token")
	if err != nil {
		t.Fatalf("OAuthCallback: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatal("want a token pair")
	}
	actor := captor.Last()
	if actor.ProjectID != nil {
		t.Fatalf("platform signup must stay platform-scoped, got %v", actor.ProjectID)
	}
	_, _ = mock.Verify(r.external, mock.Times(1)).Create(mock.AnyContext(), mock.Any[models.ActorExternalIdentities]())
	_ = mock.Verify(r.states, mock.Times(1)).DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestOAuthCallbackProjectSignupScopesActorToProject(t *testing.T) {
	testEnv(t)
	stubOAuthFlow(t)
	projectID := uuid.New()
	pair := mintPair(t)
	ops, r := newOAuthOps(t)
	state := projectState(projectID)
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Equal(state.State))).ThenReturn(&state, nil)
	mock.When(r.states.DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)
	mock.When(r.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(&models.ProjectOAuthProviders{
			ID: uuid.New(), ProjectID: projectID, Provider: models.GoogleIdentityProvider,
			ClientID: "proj-client-id", EncryptedClientSecret: encryptedProjectSecret(t), Enabled: true,
		}, nil)
	mock.When(r.external.GetByProviderAndSubject(mock.AnyContext(), mock.Equal("google"), mock.Equal("subject-1"), mock.Equal(state.ProjectID))).
		ThenReturn(nil, fun.ErrNotFound("no identity"))
	captor := mock.Captor[models.Actor]()
	mock.When(r.actors.Register(mock.AnyContext(), captor.Capture())).
		ThenAnswer(func(args []any) []any {
			a := args[1].(models.Actor)
			a.ID = uuid.New()
			return []any{&a, nil}
		})
	mock.When(r.external.Create(mock.AnyContext(), mock.Any[models.ActorExternalIdentities]())).
		ThenAnswer(func(args []any) []any {
			e := args[1].(models.ActorExternalIdentities)
			return []any{&e, nil}
		})
	mock.When(r.keys.GetActive(mock.AnyContext(), mock.Equal(models.SigningCryptoKeyType), mock.Any[*uuid.UUID]())).
		ThenReturn(&pair.key, nil)

	out, err := ops.OAuthCallback(context.Background(), "google", "code", "state-token")
	if err != nil {
		t.Fatalf("OAuthCallback: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatal("want a token pair")
	}
	actor := captor.Last()
	if actor.ProjectID == nil || *actor.ProjectID != projectID {
		t.Fatalf("want actor scoped to project %s, got %v", projectID, actor.ProjectID)
	}
}

func TestOAuthCallbackExistingIdentityLogsIn(t *testing.T) {
	testEnv(t)
	t.Setenv("GOOGLE_CLIENT_ID", "platform-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "platform-secret")
	stubOAuthFlow(t)
	pair := mintPair(t)
	ops, r := newOAuthOps(t)
	state := envState()
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Equal(state.State))).ThenReturn(&state, nil)
	mock.When(r.states.DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)
	mock.When(r.external.GetByProviderAndSubject(mock.AnyContext(), mock.Equal("google"), mock.Equal("subject-1"), mock.Equal(state.ProjectID))).
		ThenReturn(&models.ActorExternalIdentities{ID: uuid.New(), ActorID: pair.actor.ID}, nil)
	mock.When(r.external.UpdateTokens(mock.AnyContext(), mock.Any[models.ActorExternalIdentities](), mock.Equal(state.ProjectID))).
		ThenAnswer(func(args []any) []any {
			e := args[1].(models.ActorExternalIdentities)
			return []any{&e, nil}
		})
	mock.When(r.actors.GetByID(mock.AnyContext(), mock.Equal(pair.actor.ID))).ThenReturn(&pair.actor, nil)
	mock.When(r.keys.GetActive(mock.AnyContext(), mock.Equal(models.SigningCryptoKeyType), mock.Any[*uuid.UUID]())).
		ThenReturn(&pair.key, nil)

	out, err := ops.OAuthCallback(context.Background(), "google", "code", "state-token")
	if err != nil {
		t.Fatalf("OAuthCallback: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatal("want a token pair")
	}
	_, _ = mock.Verify(r.external, mock.Times(1)).UpdateTokens(mock.AnyContext(), mock.Any[models.ActorExternalIdentities](), mock.Equal(state.ProjectID))
	_, _ = mock.Verify(r.actors, mock.Times(0)).Register(mock.AnyContext(), mock.Any[models.Actor]())
}

func TestOAuthCallbackInvalidState(t *testing.T) {
	testEnv(t)
	ops, r := newOAuthOps(t)
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Any[string]())).
		ThenReturn(nil, fun.ErrNotFound("no state"))

	_, err := ops.OAuthCallback(context.Background(), "google", "code", "nope")
	if !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request, got %v", err)
	}
}

func TestOAuthCallbackProviderMismatch(t *testing.T) {
	testEnv(t)
	ops, r := newOAuthOps(t)
	state := models.OAuthLoginState{
		ID: uuid.New(), State: "state-token", Provider: models.GithubIdentityProvider,
		ProjectID: nil, ExpiresAt: time.Now().Add(time.Hour),
	}
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Equal(state.State))).ThenReturn(&state, nil)
	mock.When(r.states.DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)

	_, err := ops.OAuthCallback(context.Background(), "google", "code", "state-token")
	if !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request, got %v", err)
	}
}

func TestOAuthCallbackExpiredState(t *testing.T) {
	testEnv(t)
	ops, r := newOAuthOps(t)
	state := models.OAuthLoginState{
		ID: uuid.New(), State: "state-token", Provider: models.GoogleIdentityProvider,
		ProjectID: nil, ExpiresAt: time.Now().Add(-time.Hour),
	}
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Equal(state.State))).ThenReturn(&state, nil)
	mock.When(r.states.DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)

	_, err := ops.OAuthCallback(context.Background(), "google", "code", "state-token")
	if !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request, got %v", err)
	}
}

func TestOAuthCallbackProviderRowDeletedMidFlight(t *testing.T) {
	testEnv(t)
	projectID := uuid.New()
	ops, r := newOAuthOps(t)
	state := projectState(projectID)
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Equal(state.State))).ThenReturn(&state, nil)
	mock.When(r.states.DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)
	mock.When(r.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(nil, fun.ErrNotFound("deleted"))

	_, err := ops.OAuthCallback(context.Background(), "google", "code", "state-token")
	if !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request, got %v", err)
	}
	if !strings.Contains(err.Error(), "contact the project") {
		t.Fatalf("want friendly disabled message, got %q", err.Error())
	}
}

func TestOAuthCallbackDisabledBlocksNewSignup(t *testing.T) {
	testEnv(t)
	stubOAuthFlow(t)
	projectID := uuid.New()
	ops, r := newOAuthOps(t)
	state := projectState(projectID)
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Equal(state.State))).ThenReturn(&state, nil)
	mock.When(r.states.DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)
	row := models.ProjectOAuthProviders{
		ID: uuid.New(), ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "proj-client-id", EncryptedClientSecret: encryptedProjectSecret(t), Enabled: false,
	}
	mock.When(r.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(&row, nil)
	mock.When(r.external.GetByProviderAndSubject(mock.AnyContext(), mock.Equal("google"), mock.Equal("subject-1"), mock.Equal(state.ProjectID))).
		ThenReturn(nil, fun.ErrNotFound("no identity"))

	_, err := ops.OAuthCallback(context.Background(), "google", "code", "state-token")
	if !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
	_, _ = mock.Verify(r.actors, mock.Times(0)).Register(mock.AnyContext(), mock.Any[models.Actor]())
}

func TestOAuthCallbackDisabledAllowsExistingLogin(t *testing.T) {
	testEnv(t)
	stubOAuthFlow(t)
	projectID := uuid.New()
	pair := mintPair(t)
	ops, r := newOAuthOps(t)
	state := projectState(projectID)
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Equal(state.State))).ThenReturn(&state, nil)
	mock.When(r.states.DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)
	row := models.ProjectOAuthProviders{
		ID: uuid.New(), ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "proj-client-id", EncryptedClientSecret: encryptedProjectSecret(t), Enabled: false,
	}
	mock.When(r.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(&row, nil)
	mock.When(r.external.GetByProviderAndSubject(mock.AnyContext(), mock.Equal("google"), mock.Equal("subject-1"), mock.Equal(state.ProjectID))).
		ThenReturn(&models.ActorExternalIdentities{ID: uuid.New(), ActorID: pair.actor.ID}, nil)
	mock.When(r.external.UpdateTokens(mock.AnyContext(), mock.Any[models.ActorExternalIdentities](), mock.Equal(state.ProjectID))).
		ThenAnswer(func(args []any) []any {
			e := args[1].(models.ActorExternalIdentities)
			return []any{&e, nil}
		})
	mock.When(r.actors.GetByID(mock.AnyContext(), mock.Equal(pair.actor.ID))).ThenReturn(&pair.actor, nil)
	mock.When(r.keys.GetActive(mock.AnyContext(), mock.Equal(models.SigningCryptoKeyType), mock.Any[*uuid.UUID]())).
		ThenReturn(&pair.key, nil)

	out, err := ops.OAuthCallback(context.Background(), "google", "code", "state-token")
	if err != nil {
		t.Fatalf("existing identity must be able to log in on a disabled provider, got %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatal("want a token pair")
	}
	_, _ = mock.Verify(r.actors, mock.Times(0)).Register(mock.AnyContext(), mock.Any[models.Actor]())
}

// TestOAuthCallbackProjectLoginCreatesProjectActorWhenPlatformIdentityExists
// pins the scope bug: a Google account that only exists as a platform-level
// identity must NOT be reused by a project login. The scoped lookup misses
// (the platform row belongs to the NULL scope), so the callback registers a
// brand-new actor scoped to the project instead of hijacking the platform
// actor.
func TestOAuthCallbackProjectLoginCreatesProjectActorWhenPlatformIdentityExists(t *testing.T) {
	testEnv(t)
	stubOAuthFlow(t)
	projectID := uuid.New()
	pair := mintPair(t)
	ops, r := newOAuthOps(t)
	state := projectState(projectID)
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Equal(state.State))).ThenReturn(&state, nil)
	mock.When(r.states.DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)
	row := models.ProjectOAuthProviders{
		ID: uuid.New(), ProjectID: projectID, Provider: models.GoogleIdentityProvider,
		ClientID: "proj-client-id", EncryptedClientSecret: encryptedProjectSecret(t), Enabled: true,
	}
	mock.When(r.providers.GetByProjectAndProvider(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.GoogleIdentityProvider))).
		ThenReturn(&row, nil)
	// The platform-scoped identity exists in the DB, but the project-scoped
	// lookup must not see it.
	mock.When(r.external.GetByProviderAndSubject(mock.AnyContext(), mock.Equal("google"), mock.Equal("subject-1"), mock.Equal(state.ProjectID))).
		ThenReturn(nil, fun.ErrNotFound("no identity in this scope"))
	captor := mock.Captor[models.Actor]()
	mock.When(r.actors.Register(mock.AnyContext(), captor.Capture())).
		ThenAnswer(func(args []any) []any {
			a := args[1].(models.Actor)
			a.ID = uuid.New()
			return []any{&a, nil}
		})
	mock.When(r.external.Create(mock.AnyContext(), mock.Any[models.ActorExternalIdentities]())).
		ThenAnswer(func(args []any) []any {
			e := args[1].(models.ActorExternalIdentities)
			return []any{&e, nil}
		})
	mock.When(r.keys.GetActive(mock.AnyContext(), mock.Equal(models.SigningCryptoKeyType), mock.Any[*uuid.UUID]())).
		ThenReturn(&pair.key, nil)

	out, err := ops.OAuthCallback(context.Background(), "google", "code", "state-token")
	if err != nil {
		t.Fatalf("OAuthCallback: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatal("want a token pair")
	}
	actor := captor.Last()
	if actor.ProjectID == nil || *actor.ProjectID != projectID {
		t.Fatalf("project login must create a project-scoped actor, got %v", actor.ProjectID)
	}
	_, _ = mock.Verify(r.actors, mock.Times(1)).Register(mock.AnyContext(), mock.Any[models.Actor]())
	_, _ = mock.Verify(r.external, mock.Times(1)).Create(mock.AnyContext(), mock.Any[models.ActorExternalIdentities]())
}

// TestOAuthCallbackPlatformLoginCreatesPlatformActorWhenProjectIdentityExists
// is the mirror case: a platform login must never reuse an identity that
// belongs to a project.
func TestOAuthCallbackPlatformLoginCreatesPlatformActorWhenProjectIdentityExists(t *testing.T) {
	testEnv(t)
	t.Setenv("GOOGLE_CLIENT_ID", "platform-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "platform-secret")
	stubOAuthFlow(t)
	pair := mintPair(t)
	ops, r := newOAuthOps(t)
	state := envState()
	mock.When(r.states.GetByState(mock.AnyContext(), mock.Equal(state.State))).ThenReturn(&state, nil)
	mock.When(r.states.DeleteState(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)
	mock.When(r.external.GetByProviderAndSubject(mock.AnyContext(), mock.Equal("google"), mock.Equal("subject-1"), mock.Equal(state.ProjectID))).
		ThenReturn(nil, fun.ErrNotFound("no identity in this scope"))
	captor := mock.Captor[models.Actor]()
	mock.When(r.actors.Register(mock.AnyContext(), captor.Capture())).
		ThenAnswer(func(args []any) []any {
			a := args[1].(models.Actor)
			a.ID = uuid.New()
			return []any{&a, nil}
		})
	mock.When(r.external.Create(mock.AnyContext(), mock.Any[models.ActorExternalIdentities]())).
		ThenAnswer(func(args []any) []any {
			e := args[1].(models.ActorExternalIdentities)
			return []any{&e, nil}
		})
	mock.When(r.keys.GetActive(mock.AnyContext(), mock.Equal(models.SigningCryptoKeyType), mock.Any[*uuid.UUID]())).
		ThenReturn(&pair.key, nil)

	out, err := ops.OAuthCallback(context.Background(), "google", "code", "state-token")
	if err != nil {
		t.Fatalf("OAuthCallback: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatal("want a token pair")
	}
	actor := captor.Last()
	if actor.ProjectID != nil {
		t.Fatalf("platform login must create a platform-scoped actor, got %v", actor.ProjectID)
	}
	_, _ = mock.Verify(r.actors, mock.Times(1)).Register(mock.AnyContext(), mock.Any[models.Actor]())
}
