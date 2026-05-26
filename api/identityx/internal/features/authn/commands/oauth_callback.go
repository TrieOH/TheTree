package commands

import (
	"IdentityX/models"
	"context"
	"encoding/json"
	"io"
	"lib/crypto"
	"lib/oauth"
	"net/http"
	"time"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

func (c *Commands) OAuthCallback(ctx context.Context, provider, code string) (*models.UserTokensOutput, error) {
	ctx, span := c.tracer.Start(ctx, "OAuthCallback")
	defer span.End()

	p, ok := oauth.Registry[provider]
	if !ok {
		return nil, fun.ErrBadRequest("unsupported provider: " + provider)
	}

	providerToken, err := p.Config.Exchange(ctx, code)
	if err != nil {
		c.logger.Error("oauth code exchange failed", zap.Error(err))
		return nil, fun.ErrUnauthorized("failed to exchange code")
	}
	if providerToken == nil {
		return nil, fun.ErrUnauthorized("empty token from provider")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.Userinfo, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+providerToken.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.logger.Info("userinfo response", zap.String("provider", provider), zap.String("body", string(body)))

	var info oauth.UserInfo
	if err = json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	if info.SubString() == "" {
		return nil, fun.ErrUnauthorized("incomplete userinfo from provider")
	}

	if info.Email == "" && provider == "github" {
		info.Email, err = oauth.FetchGitHubEmail(ctx, providerToken.AccessToken)
		if err != nil {
			return nil, fun.ErrUnauthorized("could not fetch github email")
		}
	}

	if info.Email == "" {
		return nil, fun.ErrUnauthorized("incomplete userinfo from provider")
	}

	encryptedAccess, err := crypto.EncryptPrivateKey([]byte(providerToken.AccessToken))
	if err != nil {
		return nil, err
	}
	var encryptedRefresh *string
	if providerToken.RefreshToken != "" {
		e, err := crypto.EncryptPrivateKey([]byte(providerToken.RefreshToken))
		if err != nil {
			return nil, err
		}
		encryptedRefresh = &e
	}
	var tokenExpiresAt *time.Time
	if !providerToken.Expiry.IsZero() {
		tokenExpiresAt = &providerToken.Expiry
	}

	identity, err := c.externalIdentities.GetByProviderAndSubject(ctx, provider, info.SubString())
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}

	var actor *models.Actor
	if identity != nil {
		_, err = c.externalIdentities.UpdateTokens(ctx, models.ActorExternalIdentities{
			Provider:              models.OAuthProvider(provider),
			Subject:               info.SubString(),
			EncryptedAccessToken:  &encryptedAccess,
			EncryptedRefreshToken: encryptedRefresh,
			TokenExpiresAt:        tokenExpiresAt,
		})
		if err != nil {
			return nil, err
		}
		actor, err = c.actors.GetByID(ctx, identity.ActorID)
		if err != nil {
			return nil, err
		}
	} else {
		actor, err = c.actors.Register(ctx, models.Actor{
			AuthMethod: models.AuthMethod(provider),
			Email:      &info.Email,
			Type:       models.HumanActorType,
		})
		if err != nil {
			return nil, err
		}
		_, err = c.externalIdentities.Create(ctx, models.ActorExternalIdentities{
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
	}

	return c.issueTokens(ctx, actor)
}
