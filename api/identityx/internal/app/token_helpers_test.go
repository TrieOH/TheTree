package app

import (
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

// newSigningKey mints a fresh signing key pair and returns it plus its ID.
func newSigningKey(t *testing.T) (*crypto.KeyPair, uuid.UUID) {
	t.Helper()
	testEnv(t)
	kp, err := crypto.GenerateKeyPair("signing")
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	return kp, uuid.New()
}

// mintAccessToken builds a signed access token for the actor with the
// given jti, signed by the key pair.
func mintAccessToken(t *testing.T, kp *crypto.KeyPair, keyID, jti uuid.UUID, actor models.Actor) string {
	t.Helper()
	claims := models.AccessClaims{
		Sub: models.AccessSub{ID: actor.ID, ProjectID: actor.ProjectID, Email: actor.Email, Type: actor.Type},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "test-issuer", ID: jti.String(), IssuedAt: jwt.NewNumericDate(time.Now()),
		},
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

// mintRefreshToken builds a signed refresh token for the given claims.
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
