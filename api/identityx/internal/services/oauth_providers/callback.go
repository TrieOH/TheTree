package oauth_providers

import (
	"context"

	"IdentityX/models"
	"lib/crypto"
	"lib/telemetry"
	"time"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

// Callback completes an OAuth login: it validates the one-time state,
// exchanges the code for provider tokens, resolves or links the identity,
// and mints the session pair through the Token-lifecycle module. The whole
// flow — state lifecycle, provider resolution, userinfo fetch, identity
// linking — lives behind this interface.
func (o *Operations) Callback(ctx context.Context, provider, code, state string) (*models.UserTokensOutput, error) {
	ctx, span := telemetry.StartSpan(ctx, "Callback")
	defer span.End()

	loginState, err := o.loginStates.GetByState(ctx, state)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.ErrBadRequest("invalid or expired OAuth state; start a new login")
		}
		return nil, err
	}
	// The state is one-time-use: consume it regardless of the outcome so it
	// can never be replayed.
	err = o.loginStates.DeleteState(ctx, loginState.ID)
	if err != nil {
		telemetry.Log().Warn("failed to delete oauth login state", zap.String("state_id", loginState.ID.String()))
	}
	if string(loginState.Provider) != provider {
		return nil, fun.ErrBadRequest("invalid OAuth state")
	}
	if time.Now().After(loginState.ExpiresAt) {
		return nil, fun.ErrBadRequest("OAuth state expired; start a new login")
	}

	resolved, err := o.resolveLoginProvider(ctx, string(loginState.Provider), loginState.ProjectID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.ErrBadRequest(
				"this project has disabled " + string(loginState.Provider) + " login; contact the project to enable it",
			)
		}
		return nil, err
	}

	info, providerToken, err := o.fetchUserInfo(ctx, provider, resolved.Creds, code)
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

	identity, err := o.external.GetByProviderAndSubject(ctx, provider, info.SubString(), loginState.ProjectID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}

	var actor *models.Actor
	if identity != nil {
		actor, err = o.updateExistingIdentity(ctx, provider, info, identity, encryptedAccess, encryptedRefresh, tokenExpiresAt, loginState.ProjectID)
	} else {
		// A disabled provider allows existing identities to log back in,
		// but never new sign-ups.
		if resolved.Disabled {
			return nil, fun.ErrForbidden(
				"this project has disabled " + provider + " login; contact the project to enable it",
			)
		}
		actor, err = o.registerNewIdentity(ctx, provider, info, encryptedAccess, encryptedRefresh, tokenExpiresAt, loginState.ProjectID)
	}
	if err != nil {
		return nil, err
	}

	return o.tokens.Mint(ctx, actor)
}
