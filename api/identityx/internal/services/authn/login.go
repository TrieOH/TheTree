package authn

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"
	"strings"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

func (o *Operations) Login(ctx context.Context, in models.IDXLoginInput) (tokens *models.UserTokensOutput, err error) {
	ctx, span := telemetry.StartSpan(ctx, "Login")
	defer span.End()

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	actor, err := o.actors.GetByEmail(ctx, in.Email, in.ProjectID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}
	if err != nil {
		return nil, err
	}
	if actor.PasswordHash == nil {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}
	err = crypto.Verify(in.Password, *actor.PasswordHash)
	if err != nil {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}
	err = o.actors.UpdateLastLoginAt(ctx, actor.ID)
	if err != nil {
		return nil, err
	}
	// Continued-use consent stamp: after a version's effective date, the
	// first login records the ledger row (no-op before the date or when
	// already accepted). A stamp failure must not block the login — the
	// use happened either way; log it and move on.
	if actor.ProjectID != nil {
		err = o.tos.StampUse(ctx, actor.ID, *actor.ProjectID)
		if err != nil {
			telemetry.Log().Error("failed to stamp tos continued-use acceptance",
				zap.String("actor_id", actor.ID.String()),
				zap.String("project_id", actor.ProjectID.String()),
				zap.Error(err),
			)
		}
	}
	return o.tokens.Mint(ctx, actor)
}
