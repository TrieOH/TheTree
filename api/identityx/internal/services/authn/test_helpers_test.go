package authn

import (
	"strings"
	"testing"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/internal/services/oauth_providers"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

// mockOAuthProviderOps wires the oauth_providers operations over fresh
// per-test mocks, for authn tests that never touch the provider repo.
func mockOAuthProviderOps(t *testing.T) *oauth_providers.Operations {
	t.Helper()
	projects := mock.Mock[ports.ProjectRepo]()
	return oauth_providers.NewOperations(
		mock.Mock[ports.ProjectOAuthProvidersRepo](),
		projects,
		authz.New(mock.Mock[ports.OrganizationRepo](), projects, mock.Mock[ports.PlatformRolesRepo]()),
	)
}

func testEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ENCRYPTION_KEY", strings.Repeat("ab", 32))
}

func testActor() models.Actor {
	email := "actor@trieoh.com"
	return models.Actor{ID: uuid.New(), Email: &email, Type: models.HumanActorType}
}

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

// testPair is a freshly signed access/refresh pair plus the key and actor
// they were minted for. Pure crypto — no mocks involved.
type testPair struct {
	accessToken, refreshToken string
	accessJTI, refreshJTI     uuid.UUID
	key                       models.CryptoKey
	kp                        *crypto.KeyPair
	actor                     models.Actor
}

func mintPair(t *testing.T) *testPair {
	t.Helper()
	testEnv(t)
	kp, err := crypto.GenerateKeyPair("signing")
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	keyID := uuid.New()
	actor := testActor()
	key := models.CryptoKey{
		ID: keyID, Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
		PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
	}

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
	return &testPair{
		accessToken: accessToken, refreshToken: refreshToken,
		accessJTI: accessJTI, refreshJTI: refreshJTI,
		key: key, kp: kp, actor: actor,
	}
}
