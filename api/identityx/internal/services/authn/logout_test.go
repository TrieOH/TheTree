package authn

import (
	"context"
	"strings"
	"testing"
	"time"

	"IdentityX/models"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func blacklistTargets(bl *fakeBlacklist) map[string]bool {
	targets := make(map[string]bool)
	for _, e := range bl.entries {
		targets[e.Target] = true
	}
	return targets
}

func TestLogoutAppendsAccessAndRefresh(t *testing.T) {
	bl := &fakeBlacklist{}
	ops, pair := newTestOps(t, bl)

	err := ops.Logout(ctxWithIdentity(), models.LogoutInput{
		AccessToken:  pair.accessToken,
		RefreshToken: pair.refreshToken,
	})
	if err != nil {
		t.Fatalf("Logout: %v", err)
	}
	targets := blacklistTargets(bl)
	if len(bl.entries) != 2 {
		t.Fatalf("want 2 blacklist entries, got %d", len(bl.entries))
	}
	if !targets[pair.accessJTI.String()] {
		t.Fatal("access token jti not blacklisted")
	}
	if !targets[pair.refreshJTI.String()] {
		t.Fatal("refresh token jti not blacklisted")
	}
}

func TestLogoutInvalidAccessRejected(t *testing.T) {
	bl := &fakeBlacklist{}
	ops, pair := newTestOps(t, bl)

	err := ops.Logout(ctxWithIdentity(), models.LogoutInput{
		AccessToken:  "garbage",
		RefreshToken: pair.refreshToken,
	})
	if err == nil {
		t.Fatal("want error for invalid access token, got nil")
	}
	if len(bl.entries) != 0 {
		t.Fatalf("nothing must be blacklisted on invalid access, got %d entries", len(bl.entries))
	}
}

func TestLogoutMissingIdentityRejected(t *testing.T) {
	bl := &fakeBlacklist{}
	ops, pair := newTestOps(t, bl)

	err := ops.Logout(context.Background(), models.LogoutInput{
		AccessToken:  pair.accessToken,
		RefreshToken: pair.refreshToken,
	})
	if err == nil {
		t.Fatal("want error without identity in context, got nil")
	}
}

func TestLogoutInvalidRefreshStillLogsOut(t *testing.T) {
	bl := &fakeBlacklist{}
	ops, pair := newTestOps(t, bl)

	err := ops.Logout(ctxWithIdentity(), models.LogoutInput{
		AccessToken:  pair.accessToken,
		RefreshToken: "garbage-refresh-token",
	})
	if err != nil {
		t.Fatalf("logout must succeed even with an invalid refresh token, got %v", err)
	}
	targets := blacklistTargets(bl)
	if !targets[pair.accessJTI.String()] {
		t.Fatal("access token jti must still be blacklisted")
	}
	if targets[pair.refreshJTI.String()] {
		t.Fatal("unverifiable refresh token must not be blacklisted")
	}
}

func TestLogoutExpiredRefreshStillLogsOut(t *testing.T) {
	bl := &fakeBlacklist{}
	ops, pair := newTestOps(t, bl)

	expiredPayload := mintPayload(t, models.RefreshClaims{
		Sub: models.RefreshSub{ID: uuid.New(), ProjectID: nil, AccessJTI: pair.accessJTI},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			Issuer:    "test-issuer", ID: uuid.New().String(), IssuedAt: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}, pair.keyID)
	expiredRefresh, err := crypto.SignToken(expiredPayload, pair.kp)
	if err != nil {
		t.Fatalf("SignToken: %v", err)
	}

	err = ops.Logout(ctxWithIdentity(), models.LogoutInput{
		AccessToken:  pair.accessToken,
		RefreshToken: expiredRefresh,
	})
	if err != nil {
		t.Fatalf("logout must succeed even with an expired refresh token, got %v", err)
	}
	if !blacklistTargets(bl)[pair.accessJTI.String()] {
		t.Fatal("access token jti must still be blacklisted")
	}
}

func TestLogoutErrorIsUnauthorized(t *testing.T) {
	bl := &fakeBlacklist{}
	ops, pair := newTestOps(t, bl)

	err := ops.Logout(ctxWithIdentity(), models.LogoutInput{AccessToken: "garbage", RefreshToken: pair.refreshToken})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized error, got %v", err)
	}
}

func TestRefreshRejectsBlacklistedRefreshToken(t *testing.T) {
	bl := &fakeBlacklist{}
	ops, pair := newTestOps(t, bl)

	// simulate a previous logout/rotation blacklisting the refresh token
	if err := bl.Append(context.Background(), models.BlacklistEntry{
		Type: models.BlacklistEntryTypeToken, Target: pair.refreshJTI.String(),
		ExpiresAt: func() *time.Time { t := time.Now().Add(time.Hour); return &t }(),
	}); err != nil {
		t.Fatalf("seed blacklist: %v", err)
	}

	_, err := ops.Refresh(context.Background(), pair.refreshToken)
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized for blacklisted refresh token, got %v", err)
	}
}

func TestRefreshAcceptsCleanRefreshToken(t *testing.T) {
	bl := &fakeBlacklist{}
	ops, pair := newTestOps(t, bl)

	out, err := ops.Refresh(context.Background(), pair.refreshToken)
	if err != nil {
		t.Fatalf("Refresh with clean token: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatal("want a fresh token pair")
	}
	if strings.TrimSpace(out.RefreshToken) == pair.refreshToken {
		t.Fatal("refresh must rotate: new refresh token must differ")
	}
}
