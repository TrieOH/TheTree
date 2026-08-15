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
			VerifiedAt:   actor.VerifiedAt,
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
// disabled — the user already has an account. The project keeps the update
// scoped: only the identity row whose actor lives in that project is touched.
func (o *Operations) updateExistingIdentity(
	ctx context.Context,
	provider string,
	info *oauth.UserInfo,
	identity *models.ActorExternalIdentities,
	encryptedAccess string,
	encryptedRefresh *string,
	tokenExpiresAt *time.Time,
	projectID *uuid.UUID,
) (*models.Actor, error) {
	_, err := o.externalIdentities.UpdateTokens(ctx, models.ActorExternalIdentities{
		Provider:              models.OAuthProvider(provider),
		Subject:               info.SubString(),
		EncryptedAccessToken:  &encryptedAccess,
		EncryptedRefreshToken: encryptedRefresh,
		TokenExpiresAt:        tokenExpiresAt,
	}, projectID)
	if err != nil {
		return nil, err
	}
	return o.actors.GetByID(ctx, identity.ActorID)
}

// registerNewIdentity creates the actor and its external identity for a
// first-time OAuth login. When the flow was started from a project, the new
// actor is scoped to that project so its tokens carry the project id.
func (o *Operations) registerNewIdentity(
	ctx context.Context,
	provider string,
	info *oauth.UserInfo,
	encryptedAccess string,
	encryptedRefresh *string,
	tokenExpiresAt *time.Time,
	projectID *uuid.UUID,
) (*models.Actor, error) {
	actor, err := o.actors.Register(ctx, models.Actor{
		ProjectID:  projectID,
		AuthMethod: models.AuthMethod(provider),
		Email:      &info.Email,
		Type:       models.HumanActorType,
		// The provider already verified the email (Google/GitHub), so the
		// account ships verified and never needs the verify link.
		VerifiedAt: new(time.Now()),
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
