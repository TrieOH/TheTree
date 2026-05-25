package commands

import (
	"IdentityX/models"
	"context"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (c *Commands) Logout(ctx context.Context, in models.LogoutInput) error {
	ctx, span := c.tracer.Start(ctx, "Logout")
	defer span.End()

	accessClaims := &models.AccessClaims{}
	token, err := crypto.OpenUnverified(in.AccessToken, accessClaims)
	if err != nil {
		return err
	}
	if accessClaims == nil {
		return fun.ErrBadRequest("empty access claims")
	}
	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return fun.ErrUnauthorized("missing kid")
	}
	keyID, err := uuid.Parse(kid)
	if err != nil {
		return fun.ErrUnauthorized("invalid kid")
	}
	cryptoKey, err := c.cryptoKeys.GetByID(ctx, keyID)
	if err != nil && fun.Is(err, fun.CodeNotFound) {
		return fun.ErrUnauthorized("outdated token")
	}
	if err != nil {
		return err
	}

	_, err = crypto.VerifyToken(in.AccessToken, cryptoKey.PublicKey, accessClaims)
	if err != nil {
		c.logger.Error("access token verification failed", zap.Error(err))
		return fun.ErrUnauthorized("invalid access token")
	}

	accessEntry := models.BlacklistEntry{
		CreatedByActorID: &accessClaims.Sub.ID,
		ProjectID:        accessClaims.Sub.ProjectID,
		Type:             models.BlacklistEntryTypeToken,
		Target:           accessClaims.ID,
		Reason:           new("logout"),
		Metadata:         nil,
		ExpiresAt:        &accessClaims.ExpiresAt.Time,
	}
	if err = c.blacklist.Append(ctx, accessEntry); err != nil {
		c.logger.Error("error appending access token to blacklist for "+accessClaims.Sub.ID.String(), zap.Error(err))
	}

	refreshClaims := &models.RefreshClaims{}
	_, err = crypto.VerifyToken(in.RefreshToken, cryptoKey.PublicKey, refreshClaims)
	if err != nil {
		c.logger.Error("refresh token verification failed", zap.Error(err))
		return fun.ErrUnauthorized("invalid refresh token")
	}

	refreshEntry := models.BlacklistEntry{
		CreatedByActorID: &accessClaims.Sub.ID,
		ProjectID:        accessClaims.Sub.ProjectID,
		Type:             models.BlacklistEntryTypeToken,
		Target:           refreshClaims.ID,
		Reason:           new("logout"),
		Metadata:         nil,
		ExpiresAt:        &refreshClaims.ExpiresAt.Time,
	}
	if err = c.blacklist.Append(ctx, refreshEntry); err != nil {
		c.logger.Error("error appending refresh token to blacklist for "+accessClaims.Sub.ID.String(), zap.Error(err))
	}

	return nil
}
