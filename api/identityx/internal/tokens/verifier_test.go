package tokens

import (
	"context"
	"strings"
	"testing"
	"time"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

func testEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ENCRYPTION_KEY", strings.Repeat("ab", 32))
}

// fixture wires a verifier over per-test mockio mocks plus a real key
// pair, so tests mint and verify real signed tokens through the seam.
type fixture struct {
	verifier *Verifier
	keys     ports.CryptoKeysRepo
	bl       ports.BlacklistRepo
	key      models.CryptoKey
	kp       *crypto.KeyPair
}

func newFixture(t *testing.T, blacklisted bool) *fixture {
	t.Helper()
	testEnv(t)
	mock.SetUp(t)

	kp, err := crypto.GenerateKeyPair("signing")
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	key := models.CryptoKey{
		ID: uuid.New(), Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
		PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
	}
	keys := mock.Mock[ports.CryptoKeysRepo]()

	bl := mock.Mock[ports.BlacklistRepo]()
	if blacklisted {
		mock.When(bl.GetByTargetAndType(mock.AnyContext(), mock.Any[string](), mock.Any[models.BlacklistEntryType]())).
			ThenReturn(&models.BlacklistEntry{Type: models.BlacklistEntryTypeToken}, nil)
	} else {
		mock.When(bl.GetByTargetAndType(mock.AnyContext(), mock.Any[string](), mock.Any[models.BlacklistEntryType]())).
			ThenReturn(nil, fun.ErrNotFound("not blacklisted"))
	}

	return &fixture{verifier: NewVerifier(keys, bl), keys: keys, bl: bl, key: key, kp: kp}
}

// stubKey makes GetByID return the given key; the default per-test stubs
// are set explicitly so each test controls key resolution.
func (f *fixture) stubKey(key *models.CryptoKey) {
	mock.When(f.keys.GetByID(mock.AnyContext(), mock.Equal(f.key.ID))).ThenReturn(key, nil)
}

func (f *fixture) stubKeyNotFound() {
	mock.When(f.keys.GetByID(mock.AnyContext(), mock.Equal(f.key.ID))).
		ThenReturn(nil, fun.ErrNotFound("key gone"))
}

// mintAccess signs a real access token with the fixture's key.
func (f *fixture) mintAccess(t *testing.T, withoutKid bool) string {
	t.Helper()
	claims := models.AccessClaims{
		Sub: models.AccessSub{ID: uuid.New()},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "test-issuer", ID: uuid.New().String(), IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	if !withoutKid {
		token.Header["kid"] = f.key.ID
	}
	payload, err := token.SigningString()
	if err != nil {
		t.Fatalf("SigningString: %v", err)
	}
	signed, err := crypto.SignToken([]byte(payload), f.kp)
	if err != nil {
		t.Fatalf("SignToken: %v", err)
	}
	return signed
}

func TestVerifyValidToken(t *testing.T) {
	f := newFixture(t, false)
	f.stubKey(&f.key)

	claims := &models.AccessClaims{}
	_, key, err := f.verifier.Verify(context.Background(), f.mintAccess(t, false), claims)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.Sub.ID == uuid.Nil {
		t.Fatal("claims must be populated")
	}
	if key == nil || key.ID != f.key.ID {
		t.Fatalf("want the signing key back, got %+v", key)
	}
}

func TestVerifyRejectsBlacklistedToken(t *testing.T) {
	f := newFixture(t, true)
	f.stubKey(&f.key)

	_, _, err := f.verifier.Verify(context.Background(), f.mintAccess(t, false), &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a revoked token, got %v", err)
	}
}

func TestVerifyRejectsBadSignature(t *testing.T) {
	f := newFixture(t, false)
	f.stubKey(&f.key)

	token := f.mintAccess(t, false)
	tampered := token[:len(token)-2] + "xx"
	_, _, err := f.verifier.Verify(context.Background(), tampered, &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a tampered token, got %v", err)
	}
}

func TestVerifyMissingKid(t *testing.T) {
	f := newFixture(t, false)
	f.stubKey(&f.key)

	_, _, err := f.verifier.Verify(context.Background(), f.mintAccess(t, true), &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a kid-less token, got %v", err)
	}
}

func TestVerifyRejectsRevokedSigningKey(t *testing.T) {
	f := newFixture(t, false)
	revoked := f.key
	revoked.Status = "revoked"
	f.stubKey(&revoked)

	_, _, err := f.verifier.Verify(context.Background(), f.mintAccess(t, false), &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a revoked signing key, got %v", err)
	}
}

func TestVerifyOutdatedKey(t *testing.T) {
	f := newFixture(t, false)
	f.stubKeyNotFound()

	_, _, err := f.verifier.Verify(context.Background(), f.mintAccess(t, false), &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for an outdated key, got %v", err)
	}
}

func TestKeyForTokenResolvesKey(t *testing.T) {
	f := newFixture(t, false)
	f.stubKey(&f.key)

	claims := &models.AccessClaims{}
	token, err := crypto.OpenUnverified(f.mintAccess(t, false), claims)
	if err != nil {
		t.Fatalf("OpenUnverified: %v", err)
	}
	key, err := f.verifier.KeyForToken(context.Background(), token)
	if err != nil {
		t.Fatalf("KeyForToken: %v", err)
	}
	if key == nil || key.ID != f.key.ID {
		t.Fatalf("want the signing key, got %+v", key)
	}
}

func TestKeyForTokenMissingKid(t *testing.T) {
	f := newFixture(t, false)
	f.stubKey(&f.key)

	claims := &models.AccessClaims{}
	token, err := crypto.OpenUnverified(f.mintAccess(t, true), claims)
	if err != nil {
		t.Fatalf("OpenUnverified: %v", err)
	}
	_, err = f.verifier.KeyForToken(context.Background(), token)
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a kid-less token, got %v", err)
	}
}
