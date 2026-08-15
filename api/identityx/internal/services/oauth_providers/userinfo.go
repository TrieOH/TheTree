package oauth_providers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"IdentityX/models"
	"lib/oauth"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

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
	_, err = o.loginStates.CreateState(ctx, models.OAuthLoginState{
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
	_, err := o.external.UpdateTokens(ctx, models.ActorExternalIdentities{
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
	_, err = o.external.Create(ctx, models.ActorExternalIdentities{
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
// credentials and resolves the provider's userinfo. The code exchange and
// the userinfo call both go through the injected client, so tests point it
// at a local server instead of touching the network.
func (o *Operations) fetchUserInfo(ctx context.Context, provider string, creds oauth.Credentials, code string) (*oauth.UserInfo, *oauth2.Token, error) {
	p, ok := o.meta[provider]
	if !ok {
		return nil, nil, fun.ErrBadRequest("unsupported provider: " + provider)
	}

	cfg := p.Config(creds)
	exchangeCtx := context.WithValue(ctx, oauth2.HTTPClient, o.client.Client())
	providerToken, err := cfg.Exchange(exchangeCtx, code)
	if err != nil {
		telemetry.Log().Error("oauth code exchange failed", zap.Error(err))
		return nil, nil, fun.ErrUnauthorized("failed to exchange code")
	}
	if providerToken == nil {
		return nil, nil, fun.ErrUnauthorized("empty token from provider")
	}

	resp, err := o.client.R().SetContext(ctx).
		SetHeader("Authorization", "Bearer "+providerToken.AccessToken).
		Get(p.Userinfo)
	if err != nil {
		return nil, nil, err
	}
	if resp.IsStatusFailure() {
		return nil, nil, fun.ErrUnauthorized("userinfo request failed")
	}

	var info oauth.UserInfo
	err = json.Unmarshal(resp.Bytes(), &info)
	if err != nil {
		return nil, nil, err
	}

	if info.SubString() == "" {
		return nil, nil, fun.ErrUnauthorized("incomplete userinfo from provider")
	}

	if info.Email == "" && provider == "github" {
		info.Email, err = o.fetchGitHubEmail(ctx, providerToken.AccessToken)
		if err != nil {
			return nil, nil, fun.ErrUnauthorized("could not fetch github email")
		}
	}

	if info.Email == "" {
		return nil, nil, fun.ErrUnauthorized("incomplete userinfo from provider")
	}

	return &info, providerToken, nil
}

// fetchGitHubEmail resolves a GitHub user's verified primary email. GitHub
// does not include email in the userinfo payload, so it must be fetched
// separately; it uses the same injected client as the rest of the flow.
func (o *Operations) fetchGitHubEmail(ctx context.Context, accessToken string) (string, error) {
	resp, err := o.client.R().SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Accept", "application/vnd.github+json").
		Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	if resp.IsStatusFailure() {
		return "", fun.ErrUnauthorized("could not fetch github email")
	}
	var emails []oauth.GitHubEmail
	err = json.Unmarshal(resp.Bytes(), &emails)
	if err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", errors.New("no verified primary email found")
}
