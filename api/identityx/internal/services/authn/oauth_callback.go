package authn

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"
	"time"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

func (o *Operations) OAuthCallback(ctx context.Context, provider, code, state string) (*models.UserTokensOutput, error) {
	ctx, span := telemetry.StartSpan(ctx, "OAuthCallback")
	defer span.End()

	loginState, err := o.oauthLoginStates.GetByState(ctx, state)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.ErrBadRequest("invalid or expired OAuth state; start a new login")
		}
		return nil, err
	}
	// The state is one-time-use: consume it regardless of the outcome so it
	// can never be replayed.
	err = o.oauthLoginStates.DeleteState(ctx, loginState.ID)
	if err != nil {
		telemetry.Log().Warn("failed to delete oauth login state", zap.String("state_id", loginState.ID.String()))
	}
	if string(loginState.Provider) != provider {
		return nil, fun.ErrBadRequest("invalid OAuth state")
	}
	if time.Now().After(loginState.ExpiresAt) {
		return nil, fun.ErrBadRequest("OAuth state expired; start a new login")
	}

	creds, err := o.resolveCallbackCredentials(ctx, *loginState)
	if err != nil {
		return nil, err
	}

	info, providerToken, err := o.fetchUserInfo(ctx, provider, creds.creds, code)
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

	identity, err := o.externalIdentities.GetByProviderAndSubject(ctx, provider, info.SubString(), loginState.ProjectID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}

	var actor *models.Actor
	if identity != nil {
		actor, err = o.updateExistingIdentity(ctx, provider, info, identity, encryptedAccess, encryptedRefresh, tokenExpiresAt, loginState.ProjectID)
	} else {
		// A disabled provider allows existing identities to log back in,
		// but never new sign-ups.
		if creds.disabled {
			return nil, fun.ErrForbidden(
				"this project has disabled " + provider + " login; contact the project to enable it",
			)
		}
		actor, err = o.registerNewIdentity(ctx, provider, info, encryptedAccess, encryptedRefresh, tokenExpiresAt, loginState.ProjectID)
	}
	if err != nil {
		return nil, err
	}

	return o.issueTokens(ctx, actor)
}
