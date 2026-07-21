package commands

import (
	"IdentityX/models"
	"context"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

func (c *Commands) Refresh(ctx context.Context, refreshToken string) (*models.UserTokensOutput, error) {
	ctx, span := c.tracer.Start(ctx, "Refresh")
	defer span.End()

	refreshClaims := &models.RefreshClaims{}
	token, err := crypto.OpenUnverified(refreshToken, refreshClaims)
	if err != nil {
		return nil, err
	}
	cryptoKey, err := c.cryptoKeyFromToken(ctx, token)
	if err != nil {
		return nil, err
	}

	_, err = crypto.VerifyToken(refreshToken, cryptoKey.PublicKey, refreshClaims)
	if err != nil {
		c.logger.Error("refresh token verification failed", zap.Error(err))
		return nil, fun.ErrUnauthorized("invalid access token")
	}

	err = c.blacklist.Append(ctx, refreshClaims.ToRefreshBlacklistEntry())
	if err != nil {
		c.logger.Error("error appending refresh token to blacklist", zap.Error(err))
	}

	err = c.blacklist.Append(ctx, refreshClaims.ToAccessBlacklistEntry())
	if err != nil {
		c.logger.Error("error appending access token to blacklist", zap.Error(err))
	}

	actor, err := c.actors.GetByID(ctx, refreshClaims.Sub.ID)
	if err != nil {
		return nil, err
	}

	return c.issueTokens(ctx, actor)
}
