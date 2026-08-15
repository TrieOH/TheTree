package authn

import (
	"context"
	"errors"
	"strings"
	"time"

	"IdentityX/internal/tokens"
	"IdentityX/models"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

// VerifyEmail redeems a single-use verify link through the token module:
// validates the HMAC JWT (signature, purpose, expiry), consumes the jti
// anti-replay record, and stamps the actor verified. A consumed link on an
// already-verified account succeeds (idempotent) — that fall-through is
// actor-state policy, so it lives here: the anti-replay table still blocks
// replaying a token against an unverified account.
func (o *Operations) VerifyEmail(ctx context.Context, in models.VerifyEmailInput) error {
	ctx, span := telemetry.StartSpan(ctx, "VerifyEmail")
	defer span.End()

	actorID, err := o.actionTokens.Redeem(ctx, models.EmailVerifyActionTokenPurpose, in.Token)
	if err != nil {
		if !errors.Is(err, tokens.ErrActionTokenUsed) {
			return err
		}
		actor, aerr := o.actors.GetByID(ctx, actorID)
		if aerr != nil {
			return aerr
		}
		if actor.VerifiedAt != nil {
			return nil
		}
		return err
	}

	return o.actors.SetVerifiedAt(ctx, actorID, time.Now())
}

// ResendVerification mints and dispatches a fresh verify link. It always
// succeeds from the caller's perspective (anti-enumeration): unknown
// emails and already-verified accounts are silent no-ops.
func (o *Operations) ResendVerification(ctx context.Context, in models.ResendVerificationInput) error {
	ctx, span := telemetry.StartSpan(ctx, "ResendVerification")
	defer span.End()

	email := strings.TrimSpace(strings.ToLower(in.Email))
	actor, err := o.actors.GetByEmail(ctx, email, in.ProjectID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil
		}
		return err
	}
	if actor.VerifiedAt != nil || actor.Email == nil {
		return nil
	}

	var project *models.Project
	if in.ProjectID != nil {
		project, err = o.projects.GetByID(ctx, *in.ProjectID)
		if err != nil {
			return err
		}
	}
	err = o.emailSender.SendVerify(ctx, actor, project)
	if err != nil {
		telemetry.Log().Error("failed to enqueue verification email",
			zap.String("actor_id", actor.ID.String()),
			zap.Error(err),
		)
	}
	return nil
}
