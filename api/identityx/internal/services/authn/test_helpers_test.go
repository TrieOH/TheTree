package authn

import (
	"context"
	"strings"
	"testing"
	"time"

	"IdentityX/models"
	"lib/crypto"

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

// mintExpiredRefresh builds a refresh token signed with the pair's key
// that is already expired.
func mintExpiredRefresh(t *testing.T, pair *testPair) string {
	t.Helper()
	payload := mintPayload(t, models.RefreshClaims{
		Sub: models.RefreshSub{ID: pair.actor.ID, ProjectID: nil, AccessJTI: pair.accessJTI},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			Issuer:    "test-issuer", ID: uuid.New().String(), IssuedAt: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}, pair.key.ID)
	token, err := crypto.SignToken(payload, pair.kp)
	if err != nil {
		t.Fatalf("SignToken: %v", err)
	}
	return token
}

func ctxWithIdentity() context.Context {
	actor := testActor()
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub:  models.Subject{ID: actor.ID, ProjectID: nil, Email: actor.Email, Type: actor.Type},
		Cred: models.Credential{Type: models.TokenCredentialType},
	})
}
