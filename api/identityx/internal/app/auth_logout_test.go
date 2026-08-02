package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"IdentityX/internal/handlers"
	"IdentityX/internal/services"
	"IdentityX/internal/services/authn"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/globals"
	"lib/validator"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/matchers"
	"github.com/ovechkin-dm/mockio/mock"
)

// mountLogoutServer wires the real authn service over per-test mockio
// mocks and mounts the strict server with a JWT stub that injects the
// actor's identity.
func mountLogoutServer(t *testing.T, key models.CryptoKey, actor models.Actor, bl ports.BlacklistRepo) http.Handler {
	t.Helper()
	validator.SetupValidator()

	cryptoKeys := mock.Mock[ports.CryptoKeysRepo]()
	mock.When(cryptoKeys.GetByID(mock.AnyContext(), mock.Equal(key.ID))).ThenReturn(&key, nil)

	ops := authn.NewOperations(
		mock.Mock[ports.ActorRepo](),
		mock.Mock[ports.ProjectRepo](),
		mock.Mock[ports.PlatformRolesRepo](),
		cryptoKeys, bl,
		mock.Mock[ports.ExternalIdentitiesRepo](),
	)
	server := handlers.NewServer(&services.Operations{Authn: ops})
	jwtStub := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := models.WithIdentity(r.Context(), &models.Identity{
				Sub:  models.Subject{ID: actor.ID, ProjectID: nil, Email: actor.Email, Type: actor.Type},
				Cred: models.Credential{Type: models.TokenCredentialType},
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	globals.MarkSetupComplete()
	return newTestRouter(t, server, middlewares{
		jwtAuth:    jwtStub,
		apiKeyAuth: mwJWT,
		anyAuth:    mwAnyAuth,
		clientOnly: mwClientOnly,
	})
}

// stubBlacklist creates a fresh blacklist repo mock per test: appends are
// captured, lookups report not-found.
func stubBlacklist() (ports.BlacklistRepo, matchers.ArgumentCaptor[models.BlacklistEntry]) {
	bl := mock.Mock[ports.BlacklistRepo]()
	captor := mock.Captor[models.BlacklistEntry]()
	mock.When(bl.Append(mock.AnyContext(), captor.Capture())).ThenReturn(nil)
	mock.When(bl.GetByTargetAndType(mock.AnyContext(), mock.Any[string](), mock.Any[models.BlacklistEntryType]())).
		ThenReturn(nil, fun.ErrNotFound("not blacklisted"))
	return bl, captor
}

func doLogout(r http.Handler, accessToken, refreshToken string) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/auth/logout", nil)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	if refreshToken != "" {
		req.Header.Set("Refresh-Token", refreshToken)
	}
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestLogoutHTTP(t *testing.T) {
	kp, keyID := newSigningKey(t)
	actor := testActor()
	accessJTI, refreshJTI := uuid.New(), uuid.New()
	accessToken := mintAccessToken(t, kp, keyID, accessJTI, actor)
	refreshToken := mintRefreshToken(t, kp, keyID, refreshJTI, models.RefreshClaims{
		Sub: models.RefreshSub{ID: actor.ID, ProjectID: nil, AccessJTI: accessJTI},
	})
	key := models.CryptoKey{
		ID: keyID, Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
		PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
	}

	tests := []struct {
		name         string
		accessToken  string
		refreshToken string
		wantCode     int
		wantEntries  int
	}{
		{
			name:        "missing refresh token header rejected at binding",
			accessToken: accessToken,
			wantCode:    http.StatusBadRequest,
		},
		{
			name:         "both headers blacklist both tokens",
			accessToken:  accessToken,
			refreshToken: refreshToken,
			wantCode:     http.StatusOK,
			wantEntries:  2,
		},
		{
			name:         "invalid access token rejected",
			accessToken:  "garbage",
			refreshToken: refreshToken,
			wantCode:     http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.SetUp(t)
			bl, captor := stubBlacklist()
			r := mountLogoutServer(t, key, actor, bl)

			rec := doLogout(r, tt.accessToken, tt.refreshToken)
			if rec.Code != tt.wantCode {
				t.Fatalf("want %d, got %d %s", tt.wantCode, rec.Code, rec.Body.String())
			}
			if got := len(captor.Values()); got != tt.wantEntries {
				t.Fatalf("want %d blacklist entries, got %d: %v", tt.wantEntries, got, captor.Values())
			}
		})
	}
}

// TestLogoutEndToEndBody pins the success envelope shape.
func TestLogoutEndToEndBody(t *testing.T) {
	mock.SetUp(t)
	kp, keyID := newSigningKey(t)
	actor := testActor()
	accessJTI, refreshJTI := uuid.New(), uuid.New()
	accessToken := mintAccessToken(t, kp, keyID, accessJTI, actor)
	refreshToken := mintRefreshToken(t, kp, keyID, refreshJTI, models.RefreshClaims{
		Sub: models.RefreshSub{ID: actor.ID, ProjectID: nil, AccessJTI: accessJTI},
	})
	key := models.CryptoKey{
		ID: keyID, Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
		PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
	}
	bl, _ := stubBlacklist()
	r := mountLogoutServer(t, key, actor, bl)

	rec := doLogout(r, accessToken, refreshToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"code":200`) {
		t.Fatalf("want 200 envelope, got %s", rec.Body.String())
	}
}
