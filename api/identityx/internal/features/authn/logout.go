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
	cryptoKey, err := o.cryptoKeyFromToken(ctx, token)
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

	refreshClaims := &models.RefreshClaims{}
	_, err = crypto.VerifyToken(in.RefreshToken, cryptoKey.PublicKey, refreshClaims)
	if err != nil {
		telemetry.Log().Error("refresh token verification failed", zap.Error(err))
		return fun.ErrUnauthorized("invalid refresh token")
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
