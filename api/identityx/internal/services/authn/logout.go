package authn

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

func (o *Operations) Logout(ctx context.Context, in models.LogoutInput) error {
	ctx, span := telemetry.StartSpan(ctx, "Logout")
	defer span.End()

	accessClaims := &models.AccessClaims{}
	token, err := crypto.OpenUnverified(in.AccessToken, accessClaims)
	if err != nil {
		telemetry.Log().Error("access token verification failed", zap.Error(err))
		return fun.ErrUnauthorized("invalid access token")
	}
	// the signing key resolves even though the access token is deliberately
	// left unverified: a dead token must not fail logout, but its key is
	// needed to verify the refresh token below.
	key, err := o.verifier.KeyForToken(ctx, token)
	if err != nil {
		return err
	}

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	accessEntry := models.BlacklistEntry{
		CreatedByActorID: &ident.Sub.ID,
		ProjectID:        ident.Sub.ProjectID,
		Type:             models.BlacklistEntryTypeToken,
		Target:           accessClaims.ID,
		Reason:           new("logout"),
		Metadata:         nil,
		ExpiresAt:        &accessClaims.ExpiresAt.Time,
	}
	err = o.blacklist.Append(ctx, accessEntry)
	if err != nil {
		telemetry.Log().Error("error appending access token to blacklist for "+accessClaims.Sub.ID.String(), zap.Error(err))
	}

	// A dead refresh token (expired, revoked, garbage) must not fail the
	// logout: the access token is already blacklisted, so the session is
	// over either way. Only a verified refresh token gets blacklisted.
	refreshClaims := &models.RefreshClaims{}
	_, err = crypto.VerifyToken(in.RefreshToken, key.PublicKey, refreshClaims)
	if err != nil {
		telemetry.Log().Warn("refresh token not blacklisted at logout (unverifiable)", zap.Error(err))
		return nil
	}

	refreshEntry := models.BlacklistEntry{
		CreatedByActorID: &ident.Sub.ID,
		ProjectID:        ident.Sub.ProjectID,
		Type:             models.BlacklistEntryTypeToken,
		Target:           refreshClaims.ID,
		Reason:           new("logout"),
		Metadata:         nil,
		ExpiresAt:        &refreshClaims.ExpiresAt.Time,
	}
	err = o.blacklist.Append(ctx, refreshEntry)
	if err != nil {
		telemetry.Log().Error("error appending refresh token to blacklist for "+accessClaims.Sub.ID.String(), zap.Error(err))
	}

	return nil
}
