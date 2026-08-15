package authn

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"go.uber.org/zap"
)

func (o *Operations) Refresh(ctx context.Context, refreshToken string) (*models.UserTokensOutput, error) {
	ctx, span := telemetry.StartSpan(ctx, "Refresh")
	defer span.End()

	refreshClaims := &models.RefreshClaims{}
	_, _, err := o.verifier.Verify(ctx, refreshToken, refreshClaims)
	if err != nil {
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
