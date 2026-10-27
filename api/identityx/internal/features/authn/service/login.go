package service

import (
	"IdentityX/errs"
	"IdentityX/models"
	"context"
	"errors"
	"lib/crypto"
	"lib/database"
	"lib/telemetry"
	"strings"

	"go.uber.org/zap"
)

func (o *Operations) Login(ctx context.Context, in models.IDXLoginInput) (tokens *models.UserTokensOutput, err error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	actor, err := o.Actors.GetByEmail(ctx, in.Email, in.ProjectID)
	if errors.Is(err, database.ErrNotFound) {
		return nil, errs.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if actor.PasswordHash == nil {
		return nil, errs.ErrInvalidCredentials
	}
	err = crypto.Verify(in.Password, *actor.PasswordHash)
	if err != nil {
		return nil, errs.ErrInvalidCredentials
	}
	err = o.Actors.UpdateLastLoginAt(ctx, actor.ID)
	if err != nil {
		return nil, err
	}
	// Continued-use consent stamp: after a version's effective date, the
	// first login records the ledger row (no-op before the date or when
	// already accepted). A stamp failure must not block the login — the
	// use happened either way; log it and move on.
	if actor.ProjectID != nil {
		err = o.Tos.StampUse(ctx, actor.ID, *actor.ProjectID)
		if err != nil {
			telemetry.Log().Error("failed to stamp tos continued-use acceptance",
				zap.String("actor_id", actor.ID.String()),
				zap.String("project_id", actor.ProjectID.String()),
				zap.Error(err),
			)
		}
	}
	return o.Tokens.Mint(ctx, actor)
}
