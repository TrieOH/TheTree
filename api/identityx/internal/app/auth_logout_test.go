package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"IdentityX/internal/handlers"
	"IdentityX/internal/services"
	"IdentityX/internal/services/authn"
	"IdentityX/models"
	"lib/crypto"
	"lib/globals"
	"lib/validator"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func mintRefreshToken(t *testing.T, kp *crypto.KeyPair, keyID, jti uuid.UUID, claims models.RefreshClaims) string {
	t.Helper()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    "test-issuer", ID: jti.String(), IssuedAt: jwt.NewNumericDate(time.Now()),
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

// mountLogoutServer wires the real authn service over in-memory fakes and
// mounts the strict server with a JWT stub that injects the identity.
func mountLogoutServer(t *testing.T, bl *fakeBlacklistRepo, kp models.CryptoKey, actor models.Actor) http.Handler {
	t.Helper()
	validator.SetupValidator()
	ops := authn.NewOperations(
		&fakeActorsRepo{actor: actor}, &fakeProjectsRepo{}, &fakePlatformRolesRepo{},
		&fakeCryptoKeysRepo{keys: map[uuid.UUID]models.CryptoKey{kp.ID: kp}},
		bl, &fakeExternalIdentitiesRepo{},
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

func doLogout(t *testing.T, r http.Handler, accessToken, refreshToken string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
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

// TestLogoutMissingRefreshTokenHeader pins the wire contract: the spec
// requires the Refresh-Token header; a client that omits it gets a 400 at
// the generated parameter binding.
func TestLogoutMissingRefreshTokenHeader(t *testing.T) {
	kp, keyID := newSigningKey(t)
	actor := testActor()
	bl := &fakeBlacklistRepo{}
	r := mountLogoutServer(t, bl, models.CryptoKey{
		ID: keyID, Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
		PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
	}, actor)

	jti := uuid.New()
	accessToken := mintAccessToken(t, kp, keyID, jti, actor)
	rec := doLogout(t, r, accessToken, "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 when Refresh-Token header is missing, got %d %s", rec.Code, rec.Body.String())
	}
}

// TestLogoutEndToEnd pins the full contract: with both headers the
// request succeeds and both token jtis land in the blacklist.
func TestLogoutEndToEnd(t *testing.T) {
	kp, keyID := newSigningKey(t)
	actor := testActor()
	bl := &fakeBlacklistRepo{}
	r := mountLogoutServer(t, bl, models.CryptoKey{
		ID: keyID, Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
		PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
	}, actor)

	jti := uuid.New()
	accessToken := mintAccessToken(t, kp, keyID, jti, actor)
	refreshJTI := uuid.New()
	refreshToken := mintRefreshToken(t, kp, keyID, refreshJTI, models.RefreshClaims{
		Sub: models.RefreshSub{ID: actor.ID, ProjectID: nil, AccessJTI: jti},
	})

	rec := doLogout(t, r, accessToken, refreshToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"code":200`) {
		t.Fatalf("want 200 envelope, got %s", rec.Body.String())
	}
	var targets []string
	for _, e := range bl.entries {
		targets = append(targets, e.Target)
	}
	if len(bl.entries) != 2 {
		t.Fatalf("want 2 blacklist entries, got %d: %v", len(bl.entries), targets)
	}
}
