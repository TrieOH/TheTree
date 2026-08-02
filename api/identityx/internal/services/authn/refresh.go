package authn

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

func (o *Operations) Refresh(ctx context.Context, refreshToken string) (*models.UserTokensOutput, error) {
	ctx, span := telemetry.StartSpan(ctx, "Refresh")
	defer span.End()

	refreshClaims := &models.RefreshClaims{}
	token, err := crypto.OpenUnverified(refreshToken, refreshClaims)
	if err != nil {
		return nil, err
	}
	cryptoKey, err := o.cryptoKeyFromToken(ctx, token)
	if err != nil {
		return nil, err
	}

	_, err = crypto.VerifyToken(refreshToken, cryptoKey.PublicKey, refreshClaims)
	if err != nil {
		telemetry.Log().Error("refresh token verification failed", zap.Error(err))
		return nil, fun.ErrUnauthorized("invalid access token")
	}

	// a refresh token blacklisted at logout must not issue a new pair
	_, err = o.blacklist.GetByTargetAndType(ctx, refreshClaims.ID, models.BlacklistEntryTypeToken)
	if err == nil {
		return nil, fun.ErrUnauthorized("refresh token has been revoked")
	}
	if !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}

	err = o.blacklist.Append(ctx, refreshClaims.ToRefreshBlacklistEntry())
	if err != nil {
		telemetry.Log().Error("error appending refresh token to blacklist", zap.Error(err))
	}

	err = o.blacklist.Append(ctx, refreshClaims.ToAccessBlacklistEntry())
	if err != nil {
		telemetry.Log().Error("error appending access token to blacklist", zap.Error(err))
	}

	actor, err := o.actors.GetByID(ctx, refreshClaims.Sub.ID)
	if err != nil {
		return nil, err
	}

	return o.issueTokens(ctx, actor)
}
