package authn

import (
	"IdentityX/models"
	"context"
	"encoding/json"
	"io"
	"lib/crypto"
	"lib/oauth"
	"lib/telemetry"
	"net/http"
	"time"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func (o *Operations) OAuthCallback(ctx context.Context, provider, code string) (*models.UserTokensOutput, error) {
	ctx, span := telemetry.StartSpan(ctx, "OAuthCallback")
	defer span.End()

	info, providerToken, err := o.fetchUserInfo(ctx, provider, code)
	if err != nil {
		return nil, err
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

	identity, err := o.externalIdentities.GetByProviderAndSubject(ctx, provider, info.SubString())
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}

	var actor *models.Actor
	if identity != nil {
		_, err = o.externalIdentities.UpdateTokens(ctx, models.ActorExternalIdentities{
			Provider:              models.OAuthProvider(provider),
			Subject:               info.SubString(),
			EncryptedAccessToken:  &encryptedAccess,
			EncryptedRefreshToken: encryptedRefresh,
			TokenExpiresAt:        tokenExpiresAt,
		})
		if err != nil {
			return nil, err
		}
		actor, err = o.actors.GetByID(ctx, identity.ActorID)
		if err != nil {
			return nil, err
		}
	} else {
		actor, err = o.actors.Register(ctx, models.Actor{
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
	}

	return o.issueTokens(ctx, actor)
}

func (o *Operations) fetchUserInfo(ctx context.Context, provider, code string) (*oauth.UserInfo, *oauth2.Token, error) {
	p, ok := oauth.Registry[provider]
	if !ok {
		return nil, nil, fun.ErrBadRequest("unsupported provider: " + provider)
	}

	providerToken, err := p.Config.Exchange(ctx, code)
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
