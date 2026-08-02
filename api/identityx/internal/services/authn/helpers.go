package authn

import (
	"IdentityX/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"lib/crypto"
	"lib/env"
	"lib/oauth"
	"lib/telemetry"
	"net/http"
	"os"
	"time"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func (o *Operations) cryptoKeyFromToken(ctx context.Context, token *jwt.Token) (*models.CryptoKey, error) {
	ctx, span := telemetry.StartSpan(ctx, "cryptoKeyFromToken")
	defer span.End()

	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, fun.ErrUnauthorized("missing kid")
	}
	keyID, err := uuid.Parse(kid)
	if err != nil {
		return nil, fun.ErrUnauthorized("invalid kid")
	}
	cryptoKey, err := o.cryptoKeys.GetByID(ctx, keyID)
	if err != nil && fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrUnauthorized("outdated token")
	}
	if err != nil {
		return nil, err
	}
	if cryptoKey.Status == "revoked" {
		return nil, fun.ErrUnauthorized("token signing key revoked")
	}
	return cryptoKey, nil
}

func (o *Operations) issueTokens(ctx context.Context, actor *models.Actor) (*models.UserTokensOutput, error) {
	ctx, span := telemetry.StartSpan(ctx, "issueTokens")
	defer span.End()

	activeKeyPair, err := o.cryptoKeys.GetActive(ctx, models.SigningCryptoKeyType, actor.ProjectID)
	if err != nil {
		return nil, err
	}
	accessJTI := uuid.New()
	refreshJTI := uuid.New()
	accessExpiresAt := time.Now().Add(env.Get[time.Duration]("ACCESS_TOKEN_LIFETIME", time.ParseDuration, 15*time.Minute))
	refreshExpiresAt := time.Now().Add(env.Get[time.Duration]("REFRESH_TOKEN_LIFETIME", time.ParseDuration, 7*24*time.Hour))
	accessPayload, err := o.newAccessToken(*actor, accessJTI, activeKeyPair.ID, accessExpiresAt)
	if err != nil {
		return nil, err
	}
	refreshPayload, err := o.newIDXRefreshToken(actor, refreshJTI, accessJTI, activeKeyPair.ID, refreshExpiresAt)
	if err != nil {
		return nil, err
	}
	kp := activeKeyPair.ToKeyPair()
	accessToken, err := crypto.SignToken(accessPayload, kp)
	if err != nil {
		return nil, err
	}
	refreshToken, err := crypto.SignToken(refreshPayload, kp)
	if err != nil {
		return nil, err
	}
	return &models.UserTokensOutput{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		Domain:           os.Getenv("ISSUER"),
	}, nil
}

func (o *Operations) newAccessToken(actor models.Actor, jti, kid uuid.UUID, expiresAt time.Time) ([]byte, error) {
	claims := models.AccessClaims{
		Sub: models.AccessSub{
			ID:           actor.ID,
			ProjectID:    actor.ProjectID,
			Email:        actor.Email,
			Type:         actor.Type,
			Capabilities: nil,
			Metadata:     nil,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    os.Getenv("ISSUER"),
			ID:        jti.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = kid

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func (o *Operations) newIDXRefreshToken(actor *models.Actor, jti, accessJTI, kid uuid.UUID, expiresAt time.Time) ([]byte, error) {
	claims := models.RefreshClaims{
		Sub: models.RefreshSub{
			ID:        actor.ID,
			ProjectID: actor.ProjectID,
			AccessJTI: accessJTI,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    os.Getenv("ISSUER"),
			ID:        jti.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = kid

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

// ── OAuth helpers ────────────────────────────────────────────────────────

// oauthStateTTL bounds how long a connect attempt stays valid. It must
// outlive the provider consent round trip; it is deliberately short.
const oauthStateTTL = 10 * time.Minute

// resolvedCredentials is the credential scope for one OAuth flow plus
// whether the provider is disabled for that scope. Connect proceeds either
// way (existing identities may re-login); the callback enforces the rule.
type resolvedCredentials struct {
	creds    oauth.Credentials
	disabled bool
}

// platformCredentials returns the env-configured credentials for
// IdentityX itself (no project scope).
func platformCredentials(provider string) (resolvedCredentials, error) {
	creds, ok := oauth.EnvCredentials(provider)
	if !ok {
		return resolvedCredentials{}, fun.ErrBadRequest("provider not configured: " + provider)
	}
	return resolvedCredentials{creds: creds}, nil
}

// projectCredentials resolves a project's configured provider row and
// decrypts its client secret. A missing row surfaces as NotFound — callers
// decide how to phrase it.
func (o *Operations) projectCredentials(ctx context.Context, provider string, projectID uuid.UUID) (resolvedCredentials, error) {
	row, err := o.oauthProviders.GetByProjectAndProvider(ctx, projectID, models.OAuthProvider(provider))
	if err != nil {
		return resolvedCredentials{}, err
	}
	secret, err := crypto.DecryptPrivateKey(row.EncryptedClientSecret)
	if err != nil {
		return resolvedCredentials{}, err
	}
	return resolvedCredentials{
		creds:    oauth.Credentials{ClientID: row.ClientID, ClientSecret: string(secret)},
		disabled: !row.Enabled,
	}, nil
}

// resolveCredentials picks the scope's credentials: the project's
// configured provider row, or the platform env credentials when the flow
// targets IdentityX itself.
func (o *Operations) resolveCredentials(ctx context.Context, provider string, projectID *uuid.UUID) (resolvedCredentials, error) {
	if projectID == nil {
		return platformCredentials(provider)
	}
	_, err := o.projects.GetByID(ctx, *projectID)
	if err != nil {
		return resolvedCredentials{}, err
	}
	creds, err := o.projectCredentials(ctx, provider, *projectID)
	if err != nil && fun.Is(err, fun.CodeNotFound) {
		return resolvedCredentials{}, fun.ErrBadRequest("provider not configured for this project: " + provider)
	}
	return creds, err
}

// createLoginState records the connect attempt and returns the opaque state
// token the provider will hand back on the callback. Deleted there.
func (o *Operations) createLoginState(ctx context.Context, provider string, projectID *uuid.UUID) (string, error) {
	token, err := newStateToken()
	if err != nil {
		return "", err
	}
	_, err = o.oauthLoginStates.CreateState(ctx, models.OAuthLoginState{
		State:     token,
		Provider:  models.OAuthProvider(provider),
		ProjectID: projectID,
		ExpiresAt: time.Now().Add(oauthStateTTL),
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func newStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// updateExistingIdentity refreshes the stored provider tokens of an already
// linked identity and returns its actor. Allowed even when the provider is
// disabled — the user already has an account.
func (o *Operations) updateExistingIdentity(
	ctx context.Context,
	provider string,
	info *oauth.UserInfo,
	identity *models.ActorExternalIdentities,
	encryptedAccess string,
	encryptedRefresh *string,
	tokenExpiresAt *time.Time,
) (*models.Actor, error) {
	_, err := o.externalIdentities.UpdateTokens(ctx, models.ActorExternalIdentities{
		Provider:              models.OAuthProvider(provider),
		Subject:               info.SubString(),
		EncryptedAccessToken:  &encryptedAccess,
		EncryptedRefreshToken: encryptedRefresh,
		TokenExpiresAt:        tokenExpiresAt,
	})
	if err != nil {
		return nil, err
	}
	return o.actors.GetByID(ctx, identity.ActorID)
}

// registerNewIdentity creates the actor and its external identity for a
// first-time OAuth login.
func (o *Operations) registerNewIdentity(
	ctx context.Context,
	provider string,
	info *oauth.UserInfo,
	encryptedAccess string,
	encryptedRefresh *string,
	tokenExpiresAt *time.Time,
) (*models.Actor, error) {
	actor, err := o.actors.Register(ctx, models.Actor{
		AuthMethod: models.AuthMethod(provider),
		Email:      &info.Email,
		Type:       models.HumanActorType,
	})
	if err != nil {
		return nil, err
	}
	_, err = o.externalIdentities.Create(ctx, models.ActorExternalIdentities{
		ActorID:               actor.ID,
		Provider:              models.OAuthProvider(provider),
		Subject:               info.SubString(),
		Email:                 &info.Email,
		EncryptedAccessToken:  &encryptedAccess,
		EncryptedRefreshToken: encryptedRefresh,
		TokenExpiresAt:        tokenExpiresAt,
	})
	if err != nil {
		return nil, err
	}
	return actor, nil
}

// resolveCallbackCredentials resolves the credentials the state was created
// for. When the project's provider row is gone (deleted between connect and
// callback) the flow cannot proceed at all: there are no keys to exchange
// the code with, so the user is told the project disabled the auth.
func (o *Operations) resolveCallbackCredentials(ctx context.Context, state models.OAuthLoginState) (resolvedCredentials, error) {
	if state.ProjectID == nil {
		return platformCredentials(string(state.Provider))
	}
	creds, err := o.projectCredentials(ctx, string(state.Provider), *state.ProjectID)
	if err != nil && fun.Is(err, fun.CodeNotFound) {
		return resolvedCredentials{}, fun.ErrBadRequest(
			"this project has disabled " + string(state.Provider) + " login; contact the project to enable it",
		)
	}
	return creds, err
}

// fetchUserInfo exchanges the authorization code with the scope's
// credentials and resolves the provider's userinfo.
func (o *Operations) fetchUserInfo(ctx context.Context, provider string, creds oauth.Credentials, code string) (*oauth.UserInfo, *oauth2.Token, error) {
	p, ok := oauth.Registry[provider]
	if !ok {
		return nil, nil, fun.ErrBadRequest("unsupported provider: " + provider)
	}

	providerToken, err := p.Config(creds).Exchange(ctx, code)
	if err != nil {
		telemetry.Log().Error("oauth code exchange failed", zap.Error(err))
		return nil, nil, fun.ErrUnauthorized("failed to exchange code")
	}
	if providerToken == nil {
		return nil, nil, fun.ErrUnauthorized("empty token from provider")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.Userinfo, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+providerToken.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			telemetry.Log().Warn("failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	telemetry.Log().Info("userinfo response", zap.String("provider", provider), zap.String("body", string(body)))

	var info oauth.UserInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, nil, err
	}

	if info.SubString() == "" {
		return nil, nil, fun.ErrUnauthorized("incomplete userinfo from provider")
	}

	if info.Email == "" && provider == "github" {
		info.Email, err = oauth.FetchGitHubEmail(ctx, providerToken.AccessToken)
		if err != nil {
			return nil, nil, fun.ErrUnauthorized("could not fetch github email")
		}
	}

	if info.Email == "" {
		return nil, nil, fun.ErrUnauthorized("incomplete userinfo from provider")
	}

	return &info, providerToken, nil
}
