package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"IdentityX/models"

	"github.com/google/uuid"
)

// newAuthMW wires the real JWT middleware over in-memory fakes.
func newAuthMW(t *testing.T, bl *fakeBlacklistRepo, kp models.CryptoKey, actor models.Actor) func(http.Handler) http.Handler {
	t.Helper()
	app := &IdentityX{}
	mw := app.SetupAuthMiddlewares(
		&fakeCryptoKeysRepo{keys: map[uuid.UUID]models.CryptoKey{kp.ID: kp}},
		&fakeAPIKeysRepo{},
		&fakeActorsRepo{actor: actor},
		&fakeCapabilitiesRepo{},
		bl,
	)
	return mw.JWT()
}

func serveThroughJWT(t *testing.T, jwt func(http.Handler) http.Handler, token string) *httptest.ResponseRecorder {
	t.Helper()
	var gotIdentity *models.Identity
	handler := jwt(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ident, err := models.RequireIdentity(r.Context())
		if err != nil {
			http.Error(w, "identity missing", http.StatusInternalServerError)
			return
		}
		gotIdentity = ident
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if gotIdentity != nil {
		rec.Header().Set("X-Test-Identity", gotIdentity.Sub.ID.String())
	}
	return rec
}

func TestJWTAuthAcceptsCleanToken(t *testing.T) {
	kp, keyID := newSigningKey(t)
	actor := testActor()
	jti := uuid.New()
	accessToken := mintAccessToken(t, kp, keyID, jti, actor)

	jwt := newAuthMW(t, &fakeBlacklistRepo{}, models.CryptoKey{
		ID: keyID, Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
		PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
	}, actor)
	rec := serveThroughJWT(t, jwt, accessToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200 for clean token, got %d %s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("X-Test-Identity") != actor.ID.String() {
		t.Fatalf("want identity %s, got %s", actor.ID, rec.Header().Get("X-Test-Identity"))
	}
}

// TestJWTAuthRejectsBlacklistedToken pins the logout contract: a token
// whose jti was blacklisted at logout must no longer authenticate.
func TestJWTAuthRejectsBlacklistedToken(t *testing.T) {
	kp, keyID := newSigningKey(t)
	actor := testActor()
	jti := uuid.New()
	accessToken := mintAccessToken(t, kp, keyID, jti, actor)

	bl := &fakeBlacklistRepo{}
	expiresAt := time.Now().Add(time.Hour)
	_ = bl.Append(t.Context(), models.BlacklistEntry{
		Type: models.BlacklistEntryTypeToken, Target: jti.String(), ExpiresAt: &expiresAt,
	})

	jwt := newAuthMW(t, bl, models.CryptoKey{
		ID: keyID, Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
		PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
	}, actor)
	rec := serveThroughJWT(t, jwt, accessToken)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 for blacklisted token, got %d %s", rec.Code, rec.Body.String())
	}
}
