package tokens

import (
	"context"
	"errors"
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

const (
	testIssuer     = "test-issuer"
	testAccessTTL  = 15 * time.Minute
	testRefreshTTL = 24 * time.Hour
)

// fixture wires a Manager over per-test mockio mocks plus a real key pair,
// so tests mint, rotate, and revoke real signed tokens through the seam.
// Blacklist appends are collected into appended; appendErr forces a
// fail-closed path; revoked is consulted at call time, so tests can mark a
// jti revoked after minting.
type fixture struct {
	mgr         *Manager
	keys        ports.CryptoKeysRepo
	bl          ports.BlacklistRepo
	actors      ports.ActorRepo
	key         models.CryptoKey
	resolvedKey models.CryptoKey
	keyNotFound bool
	kp          *crypto.KeyPair
	actor       models.Actor
	now         time.Time
	revoked     map[string]bool
	appended    []models.BlacklistEntry
	appendErr   error
	appendCalls int
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	testEnv(t)
	mock.SetUp(t)

	kp, err := crypto.GenerateKeyPair("signing")
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	email := "actor@trieoh.com"
	actor := models.Actor{ID: uuid.New(), Email: &email, Type: models.HumanActorType}
	key := models.CryptoKey{
		ID: uuid.New(), Type: models.SigningCryptoKeyType, Status: models.CryptoKeyStatusActive,
		PublicKey: kp.Public, EncryptedPrivateKey: kp.EncryptedPrivate, Algorithm: kp.Algorithm,
	}

	f := &fixture{key: key, resolvedKey: key, kp: kp, actor: actor, revoked: map[string]bool{}}

	keys := mock.Mock[ports.CryptoKeysRepo]()
	mock.When(keys.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenAnswer(func(_ []any) []any {
			if f.keyNotFound {
				return []any{nil, fun.ErrNotFound("key gone")}
			}
			return []any{&f.resolvedKey, nil}
		})
	mock.When(keys.GetActive(mock.AnyContext(), mock.Equal(models.SigningCryptoKeyType), mock.Any[*uuid.UUID]())).
		ThenReturn(&key, nil)

	bl := mock.Mock[ports.BlacklistRepo]()
	mock.When(bl.Append(mock.AnyContext(), mock.Any[models.BlacklistEntry]())).
		ThenAnswer(func(args []any) []any {
			f.appendCalls++
			if f.appendErr != nil {
				return []any{f.appendErr}
			}
			f.appended = append(f.appended, args[1].(models.BlacklistEntry))
			return []any{nil}
		})
	mock.When(bl.GetByTargetAndType(mock.AnyContext(), mock.Any[string](), mock.Any[models.BlacklistEntryType]())).
		ThenAnswer(func(args []any) []any {
			target := args[1].(string)
			if f.revoked[target] {
				return []any{&models.BlacklistEntry{Target: target, Type: models.BlacklistEntryTypeToken}, nil}
			}
			return []any{nil, fun.ErrNotFound("not blacklisted")}
		})

	actors := mock.Mock[ports.ActorRepo]()
	mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actor.ID))).ThenReturn(&actor, nil)

	f.now = time.Now().UTC()
	f.mgr = NewManager(keys, bl, actors, mock.Mock[ports.ProjectRepo](), Config{
		Issuer:     testIssuer,
		AccessTTL:  testAccessTTL,
		RefreshTTL: testRefreshTTL,
	}, WithClock(func() time.Time { return f.now }))
	f.keys = keys
	f.bl = bl
	f.actors = actors
	return f
}

// mint signs a fresh pair through the Manager's own Mint.
func (f *fixture) mint(t *testing.T) *models.UserTokensOutput {
	t.Helper()
	out, err := f.mgr.Mint(context.Background(), &f.actor)
	if err != nil {
		t.Fatalf("Mint: %v", err)
	}
	return out
}

// stubKey replaces the resolved signing key (e.g. with a revoked one).
func (f *fixture) stubKey(key *models.CryptoKey) {
	f.resolvedKey = *key
}

func (f *fixture) stubKeyNotFound() {
	f.keyNotFound = true
}

// signAccess mints an access token for the fixture actor with the fixture
// key, optionally without a kid header.
func (f *fixture) signAccess(t *testing.T, withoutKid bool) string {
	t.Helper()
	claims := models.AccessClaims{
		Sub: models.AccessSub{ID: f.actor.ID, ProjectID: f.actor.ProjectID, Email: f.actor.Email, Type: f.actor.Type},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(f.now.Add(testAccessTTL)),
			Issuer:    testIssuer, ID: uuid.New().String(), IssuedAt: jwt.NewNumericDate(f.now),
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

// mintExpiredRefresh signs a refresh token with the fixture key that is
// already past its expiry.
func (f *fixture) mintExpiredRefresh(t *testing.T) string {
	t.Helper()
	claims := models.RefreshClaims{
		Sub: models.RefreshSub{ID: f.actor.ID, ProjectID: nil, AccessJTI: uuid.New()},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			Issuer:    testIssuer, ID: uuid.New().String(), IssuedAt: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = f.key.ID
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

func openAccess(t *testing.T, tokenStr string) *models.AccessClaims {
	t.Helper()
	claims := &models.AccessClaims{}
	_, err := crypto.OpenUnverified(tokenStr, claims)
	if err != nil {
		t.Fatalf("OpenUnverified access: %v", err)
	}
	return claims
}

func openRefresh(t *testing.T, tokenStr string) *models.RefreshClaims {
	t.Helper()
	claims := &models.RefreshClaims{}
	_, err := crypto.OpenUnverified(tokenStr, claims)
	if err != nil {
		t.Fatalf("OpenUnverified refresh: %v", err)
	}
	return claims
}

func appendedTargets(entries []models.BlacklistEntry) []string {
	targets := make([]string, 0, len(entries))
	for _, e := range entries {
		targets = append(targets, e.Target)
	}
	return targets
}

func hasTarget(entries []models.BlacklistEntry, target string) bool {
	for _, e := range entries {
		if e.Target == target {
			return true
		}
	}
	return false
}

// ── Verify ───────────────────────────────────────────────────────────────

func TestVerifyValidToken(t *testing.T) {
	f := newFixture(t)
	out := f.mint(t)

	claims := &models.AccessClaims{}
	err := f.mgr.Verify(context.Background(), out.AccessToken, claims)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.Sub.ID != f.actor.ID {
		t.Fatalf("claims must be populated from the token, got %v", claims.Sub.ID)
	}
}

func TestVerifyRejectsBlacklistedToken(t *testing.T) {
	f := newFixture(t)
	out := f.mint(t)
	f.revoked[openAccess(t, out.AccessToken).ID] = true

	err := f.mgr.Verify(context.Background(), out.AccessToken, &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a revoked token, got %v", err)
	}
}

func TestVerifyRejectsBadSignature(t *testing.T) {
	f := newFixture(t)
	out := f.mint(t)
	tampered := out.AccessToken[:len(out.AccessToken)-2] + "xx"

	err := f.mgr.Verify(context.Background(), tampered, &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a tampered token, got %v", err)
	}
}

func TestVerifyMissingKid(t *testing.T) {
	f := newFixture(t)

	err := f.mgr.Verify(context.Background(), f.signAccess(t, true), &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a kid-less token, got %v", err)
	}
}

func TestVerifyRejectsRevokedSigningKey(t *testing.T) {
	f := newFixture(t)
	revoked := f.key
	revoked.Status = models.CryptoKeyStatusRevoked
	f.stubKey(&revoked)

	err := f.mgr.Verify(context.Background(), f.signAccess(t, false), &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a revoked signing key, got %v", err)
	}
}

func TestVerifyRejectsRetiredSigningKey(t *testing.T) {
	f := newFixture(t)
	retired := f.key
	retired.Status = models.CryptoKeyStatusRetired
	f.stubKey(&retired)

	err := f.mgr.Verify(context.Background(), f.signAccess(t, false), &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a retired signing key, got %v", err)
	}
}

func TestVerifyAcceptsRetiringSigningKey(t *testing.T) {
	f := newFixture(t)
	retiring := f.key
	retiring.Status = models.CryptoKeyStatusRetiring
	f.stubKey(&retiring)

	err := f.mgr.Verify(context.Background(), f.signAccess(t, false), &models.AccessClaims{})
	if err != nil {
		t.Fatalf("want a retiring signing key to still verify (tokens in flight), got %v", err)
	}
}

func TestVerifyOutdatedKey(t *testing.T) {
	f := newFixture(t)
	f.stubKeyNotFound()

	err := f.mgr.Verify(context.Background(), f.signAccess(t, false), &models.AccessClaims{})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for an outdated key, got %v", err)
	}
}

// ── Mint ─────────────────────────────────────────────────────────────────

func TestMintPinsLifetimesFromConfig(t *testing.T) {
	f := newFixture(t)
	out := f.mint(t)

	if !out.AccessExpiresAt.Equal(f.now.Add(testAccessTTL)) {
		t.Fatalf("access expires = %v, want %v", out.AccessExpiresAt, f.now.Add(testAccessTTL))
	}
	if !out.RefreshExpiresAt.Equal(f.now.Add(testRefreshTTL)) {
		t.Fatalf("refresh expires = %v, want %v", out.RefreshExpiresAt, f.now.Add(testRefreshTTL))
	}
	if out.Domain != testIssuer {
		t.Fatalf("domain = %q, want the configured issuer %q", out.Domain, testIssuer)
	}
}

func TestMintLinksRefreshToAccessJTI(t *testing.T) {
	f := newFixture(t)
	out := f.mint(t)

	access := openAccess(t, out.AccessToken)
	refresh := openRefresh(t, out.RefreshToken)
	if refresh.Sub.AccessJTI.String() != access.ID {
		t.Fatalf("refresh must anchor the access jti: refresh=%s access=%s", refresh.Sub.AccessJTI, access.ID)
	}
	if access.Sub.ID != f.actor.ID || refresh.Sub.ID != f.actor.ID {
		t.Fatalf("both tokens must carry the actor id")
	}
	if access.Issuer != testIssuer || refresh.Issuer != testIssuer {
		t.Fatalf("both tokens must carry the configured issuer")
	}
}

// ── Rotate ───────────────────────────────────────────────────────────────

func TestRotateBlacklistsOldPairAndMintsNew(t *testing.T) {
	f := newFixture(t)
	old := f.mint(t)
	oldAccessJTI := openAccess(t, old.AccessToken).ID
	oldRefreshJTI := openRefresh(t, old.RefreshToken).ID

	out, err := f.mgr.Rotate(context.Background(), old.RefreshToken)
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatal("want a fresh pair")
	}
	if out.RefreshToken == old.RefreshToken || out.AccessToken == old.AccessToken {
		t.Fatal("rotation must mint a new pair, not reuse the old tokens")
	}

	if !hasTarget(f.appended, oldRefreshJTI) || !hasTarget(f.appended, oldAccessJTI) {
		t.Fatalf("rotation must blacklist the old refresh and the access token it anchors, got %v", appendedTargets(f.appended))
	}
	for _, e := range f.appended {
		if e.Reason == nil || *e.Reason != "refresh" {
			t.Fatalf("rotation entries must carry reason %q, got %v", "refresh", e.Reason)
		}
		if e.CreatedByActorID == nil || *e.CreatedByActorID != f.actor.ID {
			t.Fatalf("rotation entries must be stamped with the actor, got %v", e.CreatedByActorID)
		}
	}
}

func TestRotateRejectsBlacklistedRefresh(t *testing.T) {
	f := newFixture(t)
	old := f.mint(t)
	f.revoked[openRefresh(t, old.RefreshToken).ID] = true

	_, err := f.mgr.Rotate(context.Background(), old.RefreshToken)
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for a revoked refresh token, got %v", err)
	}
}

func TestRotateRejectsGarbageRefresh(t *testing.T) {
	f := newFixture(t)

	_, err := f.mgr.Rotate(context.Background(), "garbage")
	if err == nil {
		t.Fatal("want an error for a garbage refresh token")
	}
}

func TestRotateRejectsExpiredRefresh(t *testing.T) {
	f := newFixture(t)

	_, err := f.mgr.Rotate(context.Background(), f.mintExpiredRefresh(t))
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for an expired refresh token, got %v", err)
	}
}

func TestRotateFailsClosedOnAppendError(t *testing.T) {
	f := newFixture(t)
	old := f.mint(t)
	f.appendErr = errors.New("blacklist down")

	_, err := f.mgr.Rotate(context.Background(), old.RefreshToken)
	if err == nil {
		t.Fatal("rotation must fail when the old pair cannot be blacklisted")
	}
	if f.appendCalls < 1 {
		t.Fatal("expected the blacklist append to be attempted before the failure")
	}
}

// ── Revoke ───────────────────────────────────────────────────────────────

func TestRevokeBlacklistsBothTokens(t *testing.T) {
	f := newFixture(t)
	pair := f.mint(t)

	err := f.mgr.Revoke(context.Background(), pair.AccessToken, pair.RefreshToken)
	if err != nil {
		t.Fatalf("Revoke: %v", err)
	}

	wantAccess := openAccess(t, pair.AccessToken).ID
	wantRefresh := openRefresh(t, pair.RefreshToken).ID
	if !hasTarget(f.appended, wantAccess) || !hasTarget(f.appended, wantRefresh) {
		t.Fatalf("logout must blacklist both tokens, got %v", appendedTargets(f.appended))
	}
	for _, e := range f.appended {
		if e.Reason == nil || *e.Reason != "logout" {
			t.Fatalf("logout entries must carry reason %q, got %v", "logout", e.Reason)
		}
		// identity is derived from the claims, not from request context
		if e.CreatedByActorID == nil || *e.CreatedByActorID != f.actor.ID {
			t.Fatalf("logout entries must be stamped with the actor from the claims, got %v", e.CreatedByActorID)
		}
	}
}

func TestRevokeBlacklistsAccessWhenRefreshDead(t *testing.T) {
	f := newFixture(t)
	pair := f.mint(t)

	err := f.mgr.Revoke(context.Background(), pair.AccessToken, "garbage")
	if err != nil {
		t.Fatalf("a dead refresh token must not fail the logout, got %v", err)
	}
	wantAccess := openAccess(t, pair.AccessToken).ID
	if !hasTarget(f.appended, wantAccess) {
		t.Fatalf("the access token must be blacklisted, got %v", appendedTargets(f.appended))
	}
	if len(f.appended) != 1 {
		t.Fatalf("only the access token may be blacklisted when the refresh is dead, got %v", appendedTargets(f.appended))
	}
}

func TestRevokeBlacklistsAccessWhenRefreshExpired(t *testing.T) {
	f := newFixture(t)
	pair := f.mint(t)

	err := f.mgr.Revoke(context.Background(), pair.AccessToken, f.mintExpiredRefresh(t))
	if err != nil {
		t.Fatalf("an expired refresh token must not fail the logout, got %v", err)
	}
	if len(f.appended) != 1 {
		t.Fatalf("only the access token may be blacklisted, got %v", appendedTargets(f.appended))
	}
}

func TestRevokeRejectsInvalidAccess(t *testing.T) {
	f := newFixture(t)

	err := f.mgr.Revoke(context.Background(), "garbage", "x")
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for an invalid access token, got %v", err)
	}
	if len(f.appended) != 0 {
		t.Fatal("an invalid access token must not blacklist anything")
	}
}

func TestRevokeFailsClosedOnAppendError(t *testing.T) {
	f := newFixture(t)
	pair := f.mint(t)
	f.appendErr = errors.New("blacklist down")

	err := f.mgr.Revoke(context.Background(), pair.AccessToken, pair.RefreshToken)
	if err == nil {
		t.Fatal("logout must fail when the access token cannot be blacklisted")
	}
}
