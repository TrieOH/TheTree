package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"IdentityX/models"
	"IdentityX/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

// newAuthMW wires the real JWT middleware over per-test mockio mocks.
func newAuthMW(t *testing.T, key models.CryptoKey, blacklisted func(jti string) bool) func(http.Handler) http.Handler {
	t.Helper()
	cryptoKeys := mock.Mock[ports.CryptoKeysRepo]()
	mock.When(cryptoKeys.GetByID(mock.AnyContext(), mock.Equal(key.ID))).ThenReturn(&key, nil)

	bl := mock.Mock[ports.BlacklistRepo]()
	mock.When(bl.GetByTargetAndType(mock.AnyContext(), mock.Any[string](), mock.Any[models.BlacklistEntryType]())).
		ThenAnswer(func(args []any) []any {
			target := args[1].(string)
			if blacklisted(target) {
				return []any{&models.BlacklistEntry{Target: target, Type: models.BlacklistEntryTypeToken}, nil}
			}
			return []any{nil, fun.ErrNotFound("not blacklisted")}
		})

	mw := (&IdentityX{}).SetupAuthMiddlewares(
		cryptoKeys,
		mock.Mock[ports.APIKeysRepo](),
		mock.Mock[ports.ActorRepo](),
		mock.Mock[ports.CapabilityRepo](),
		bl,
	)
	return mw.JWT()
}

func serveThroughJWT(t *testing.T, jwt func(http.Handler) http.Handler, token string) (*httptest.ResponseRecorder, string) {
	t.Helper()
	var identityID string
	handler := jwt(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ident, err := models.RequireIdentity(r.Context())
		if err != nil {
			http.Error(w, "identity missing", http.StatusInternalServerError)
			return
		}
		identityID = ident.Sub.ID.String()
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec, identityID
}

func TestJWTAuth(t *testing.T) {
	kp, keyID := newSigningKey(t)
	actor := testActor()
	jti := uuid.New()
	accessToken := mintAccessToken(t, kp, keyID, jti, actor)
	key := models.CryptoKey{
		ID: keyID, Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
		PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
	}

	tests := []struct {
		name        string
		token       string
		blacklisted bool
		wantCode    int
		wantActor   bool
	}{
		{
			name:      "clean token authenticates",
			token:     accessToken,
			wantCode:  http.StatusOK,
			wantActor: true,
		},
		{
			name:        "blacklisted token rejected",
			token:       accessToken,
			blacklisted: true,
			wantCode:    http.StatusUnauthorized,
		},
		{
			name:     "missing token rejected",
			wantCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.SetUp(t)
			jwt := newAuthMW(t, key, func(got string) bool {
				return tt.blacklisted && got == jti.String()
			})

			rec, identityID := serveThroughJWT(t, jwt, tt.token)
			if rec.Code != tt.wantCode {
				t.Fatalf("want %d, got %d %s", tt.wantCode, rec.Code, rec.Body.String())
			}
			if tt.wantActor && identityID != actor.ID.String() {
				t.Fatalf("want identity %s, got %s", actor.ID, identityID)
			}
		})
	}
}
