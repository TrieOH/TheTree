package commands

import (
	"IdentityX/models"
	"context"
	"encoding/json"
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

	googleToken, err := p.Config.Exchange(ctx, code)
	if err != nil {
		c.logger.Error("oauth code exchange failed", zap.Error(err))
		return nil, fun.ErrUnauthorized("failed to exchange code")
	}
	if googleToken == nil {
		return nil, fun.ErrUnauthorized("empty token from provider")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.Userinfo, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+googleToken.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info oauth.UserInfo
	if err = json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	if info.Sub == "" || info.Email == "" {
		return nil, fun.ErrUnauthorized("incomplete userinfo from provider")
	}

	encryptedAccess, err := crypto.EncryptPrivateKey([]byte(googleToken.AccessToken))
	if err != nil {
		return nil, err
	}
	var encryptedRefresh *string
	if googleToken.RefreshToken != "" {
		e, err := crypto.EncryptPrivateKey([]byte(googleToken.RefreshToken))
		if err != nil {
			return nil, err
		}
		encryptedRefresh = &e
	}
	var tokenExpiresAt *time.Time
	if !googleToken.Expiry.IsZero() {
		tokenExpiresAt = &googleToken.Expiry
	}

	identity, err := c.externalIdentities.GetByProviderAndSubject(ctx, provider, info.Sub)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}

	var actor *models.Actor
	if identity != nil {
		_, err = c.externalIdentities.UpdateTokens(ctx, models.ActorExternalIdentities{
			Provider:              models.OAuthProvider(provider),
			Subject:               info.Sub,
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
			Subject:               info.Sub,
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
