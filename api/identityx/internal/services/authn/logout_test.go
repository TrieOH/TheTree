package authn

import (
	"context"
	"testing"

	"IdentityX/models"
	"IdentityX/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/matchers"
	"github.com/ovechkin-dm/mockio/mock"
)

// stubBlacklist creates a fresh blacklist repo mock per test: appends are
// captured, and the given target jtis (strings) are reported as revoked.
func stubBlacklist(revoked map[string]bool) (ports.BlacklistRepo, matchers.ArgumentCaptor[models.BlacklistEntry]) {
	bl := mock.Mock[ports.BlacklistRepo]()
	captor := mock.Captor[models.BlacklistEntry]()
	mock.When(bl.Append(mock.AnyContext(), captor.Capture())).ThenReturn(nil)
	mock.When(bl.GetByTargetAndType(mock.AnyContext(), mock.Any[string](), mock.Any[models.BlacklistEntryType]())).
		ThenAnswer(func(args []any) []any {
			target := args[1].(string)
			if revoked[target] {
				return []any{&models.BlacklistEntry{Target: target, Type: models.BlacklistEntryTypeToken}, nil}
			}
			return []any{nil, fun.ErrNotFound("not blacklisted")}
		})
	return bl, captor
}

// stubCryptoKeys creates a fresh crypto-keys repo mock per test that
// resolves the pair's signing key by id and as the active key.
func stubCryptoKeys(pair *testPair) ports.CryptoKeysRepo {
	keys := mock.Mock[ports.CryptoKeysRepo]()
	mock.When(keys.GetByID(mock.AnyContext(), mock.Equal(pair.key.ID))).ThenReturn(&pair.key, nil)
	mock.When(keys.GetActive(mock.AnyContext(), mock.Equal(models.SigningCryptoKeyType), mock.Any[*uuid.UUID]())).
		ThenReturn(&pair.key, nil)
	return keys
}

func appendTargets(captor matchers.ArgumentCaptor[models.BlacklistEntry]) []string {
	targets := make([]string, 0, len(captor.Values()))
	for _, e := range captor.Values() {
		targets = append(targets, e.Target)
	}
	return targets
}

func hasTarget(captor matchers.ArgumentCaptor[models.BlacklistEntry], target string) bool {
	for _, e := range captor.Values() {
		if e.Target == target {
			return true
		}
	}
	return false
}

func TestLogout(t *testing.T) {
	pair := mintPair(t)
	expiredRefresh := mintExpiredRefresh(t, pair)

	tests := []struct {
		name         string
		ctx          context.Context
		accessToken  string
		refreshToken string
		wantErr      bool
		wantAccess   bool // access jti expected in the blacklist
		wantRefresh  bool // refresh jti expected in the blacklist
	}{
		{
			name:         "valid pair appends both tokens",
			ctx:          ctxWithIdentity(),
			accessToken:  pair.accessToken,
			refreshToken: pair.refreshToken,
			wantAccess:   true,
			wantRefresh:  true,
		},
		{
			name:         "invalid access token rejected without side effects",
			ctx:          ctxWithIdentity(),
			accessToken:  "garbage",
			refreshToken: pair.refreshToken,
			wantErr:      true,
		},
		{
			name:         "missing identity rejected",
			ctx:          context.Background(),
			accessToken:  pair.accessToken,
			refreshToken: pair.refreshToken,
			wantErr:      true,
		},
		{
			name:         "invalid refresh token still logs out the access token",
			ctx:          ctxWithIdentity(),
			accessToken:  pair.accessToken,
			refreshToken: "garbage-refresh-token",
			wantAccess:   true,
		},
		{
			name:         "expired refresh token still logs out the access token",
			ctx:          ctxWithIdentity(),
			accessToken:  pair.accessToken,
			refreshToken: expiredRefresh,
			wantAccess:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.SetUp(t)

			keys := stubCryptoKeys(pair)
			bl, captor := stubBlacklist(nil)
			ops := NewOperations(
				mock.Mock[ports.ActorRepo](),
				mock.Mock[ports.ProjectRepo](),
				mock.Mock[ports.PlatformRolesRepo](),
				keys, bl,
				mock.Mock[ports.ExternalIdentitiesRepo](),
				mockOAuthProviderOps(t),
				mock.Mock[ports.OAuthLoginStatesRepo](),
				mock.Mock[ports.ActionTokenRepo](),
				mock.Mock[ports.EmailSender](),
				[]byte("test-hmac"),
			)

			err := ops.Logout(tt.ctx, models.LogoutInput{
				AccessToken:  tt.accessToken,
				RefreshToken: tt.refreshToken,
			})
			if tt.wantErr && err == nil {
				t.Fatal("want error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want no error, got %v", err)
			}
			if got := hasTarget(captor, pair.accessJTI.String()); got != tt.wantAccess {
				t.Fatalf("access token blacklisted = %v, want %v (entries: %v)", got, tt.wantAccess, appendTargets(captor))
			}
			if got := hasTarget(captor, pair.refreshJTI.String()); got != tt.wantRefresh {
				t.Fatalf("refresh token blacklisted = %v, want %v (entries: %v)", got, tt.wantRefresh, appendTargets(captor))
			}
		})
	}
}

func TestLogoutErrorIsUnauthorized(t *testing.T) {
	mock.SetUp(t)
	pair := mintPair(t)

	keys := stubCryptoKeys(pair)
	bl, _ := stubBlacklist(nil)
	ops := NewOperations(
		mock.Mock[ports.ActorRepo](),
		mock.Mock[ports.ProjectRepo](),
		mock.Mock[ports.PlatformRolesRepo](),
		keys, bl,
		mock.Mock[ports.ExternalIdentitiesRepo](),
		mockOAuthProviderOps(t),
		mock.Mock[ports.OAuthLoginStatesRepo](),
		mock.Mock[ports.ActionTokenRepo](),
		mock.Mock[ports.EmailSender](),
		[]byte("test-hmac"),
	)

	err := ops.Logout(ctxWithIdentity(), models.LogoutInput{
		AccessToken:  "garbage",
		RefreshToken: pair.refreshToken,
	})
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("want unauthorized error, got %v", err)
	}
}

func TestRefresh(t *testing.T) {
	pair := mintPair(t)
	expiredRefresh := mintExpiredRefresh(t, pair)

	tests := []struct {
		name         string
		refreshToken string
		revoked      map[string]bool
		wantErr      bool
		wantRotated  bool
	}{
		{
			name:         "clean token issues a fresh pair",
			refreshToken: pair.refreshToken,
			wantRotated:  true,
		},
		{
			name:         "blacklisted token rejected",
			refreshToken: pair.refreshToken,
			revoked:      map[string]bool{pair.refreshJTI.String(): true},
			wantErr:      true,
		},
		{
			name:         "garbage token rejected",
			refreshToken: "garbage",
			wantErr:      true,
		},
		{
			name:         "expired token rejected",
			refreshToken: expiredRefresh,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.SetUp(t)

			keys := stubCryptoKeys(pair)
			bl, _ := stubBlacklist(tt.revoked)
			actors := mock.Mock[ports.ActorRepo]()
			mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(pair.actor.ID))).ThenReturn(&pair.actor, nil)

			ops := NewOperations(
				actors,
				mock.Mock[ports.ProjectRepo](),
				mock.Mock[ports.PlatformRolesRepo](),
				keys, bl,
				mock.Mock[ports.ExternalIdentitiesRepo](),
				mockOAuthProviderOps(t),
				mock.Mock[ports.OAuthLoginStatesRepo](),
				mock.Mock[ports.ActionTokenRepo](),
				mock.Mock[ports.EmailSender](),
				[]byte("test-hmac"),
			)

			out, err := ops.Refresh(context.Background(), tt.refreshToken)
			if tt.wantErr && err == nil {
				t.Fatal("want error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want no error, got %v", err)
			}
			if !tt.wantErr && tt.wantRotated {
				if out.AccessToken == "" || out.RefreshToken == "" {
					t.Fatal("want a fresh token pair")
				}
				if out.RefreshToken == tt.refreshToken {
					t.Fatal("refresh must rotate: new refresh token must differ")
				}
			}
		})
	}
}
