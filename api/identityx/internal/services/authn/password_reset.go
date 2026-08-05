package authn

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"
	"strings"
	"time"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

// ForgotPassword mints and dispatches a reset link for the actor's email.
// It always succeeds from the caller's perspective (anti-enumeration):
// unknown emails and passwordless (OAuth-only) accounts are silent no-ops.
func (o *Operations) ForgotPassword(ctx context.Context, in models.ForgotPasswordInput) error {
	ctx, span := telemetry.StartSpan(ctx, "ForgotPassword")
	defer span.End()

	email := strings.TrimSpace(strings.ToLower(in.Email))
	actor, err := o.actors.GetByEmail(ctx, email, in.ProjectID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil
		}
		return err
	}
	// An OAuth-only account has no password to reset; sending a link would
	// only let an attacker turn it into a password account.
	if actor.PasswordHash == nil {
		return nil
	}

	var project *models.Project
	if in.ProjectID != nil {
		project, err = o.projects.GetByID(ctx, *in.ProjectID)
		if err != nil {
			return err
		}
	}
	err = o.emailSender.SendReset(ctx, actor, project)
	if err != nil {
		telemetry.Log().Error("failed to enqueue password reset email",
			zap.String("actor_id", actor.ID.String()),
			zap.Error(err),
		)
	}
	return nil
}

// ResetPassword redeems a single-use reset link: validates the HMAC JWT
// (signature, purpose, expiry), consumes the jti anti-replay record, and
// replaces the actor's password hash. Unlike verification, reset is never
// idempotent — a consumed token is dead, and the user requests a new link
// if they need another reset.
func (o *Operations) ResetPassword(ctx context.Context, in models.ResetPasswordInput) error {
	ctx, span := telemetry.StartSpan(ctx, "ResetPassword")
	defer span.End()

	actorID, jti, err := o.parseActionToken(ctx, in.Token, models.PasswordResetActionTokenPurpose)
	if err != nil {
		return err
	}

	record, err := o.actionTokens.GetByJTI(ctx, jti)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return fun.ErrBadRequest("invalid or expired token")
		}
		return err
	}
	if record.UsedAt != nil {
		return fun.ErrBadRequest("token already used")
	}
	if time.Now().After(record.ExpiresAt) {
		return fun.ErrBadRequest("token expired")
	}

	_, err = o.actionTokens.Consume(ctx, jti)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return fun.ErrBadRequest("token already used")
		}
		return err
	}

	hashed, err := crypto.Hash(in.NewPassword, crypto.Strong)
	if err != nil {
		return err
	}
	return o.actors.UpdatePasswordHash(ctx, actorID, hashed)
}
