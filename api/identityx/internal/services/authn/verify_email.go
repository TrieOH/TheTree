package authn

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"
	"strings"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// VerifyEmail redeems a single-use verify link: validates the HMAC JWT
// (signature, purpose, expiry), consumes the jti anti-replay record, and
// stamps the actor verified. Re-clicking an already-consumed link on an
// already-verified account succeeds (idempotent) — the anti-replay table
// still blocks replaying a token against an unverified account.
func (o *Operations) VerifyEmail(ctx context.Context, in models.VerifyEmailInput) error {
	ctx, span := telemetry.StartSpan(ctx, "VerifyEmail")
	defer span.End()

	actorID, jti, err := o.parseActionToken(ctx, in.Token, models.EmailVerifyActionTokenPurpose)
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
		actor, aerr := o.actors.GetByID(ctx, actorID)
		if aerr != nil {
			return aerr
		}
		if actor.VerifiedAt != nil {
			return nil
		}
		return fun.ErrBadRequest("token already used")
	}
	if time.Now().After(record.ExpiresAt) {
		return fun.ErrBadRequest("token expired")
	}

	_, err = o.actionTokens.Consume(ctx, jti)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			// consumed concurrently; fall through to the used-token path
			actor, aerr := o.actors.GetByID(ctx, actorID)
			if aerr != nil {
				return aerr
			}
			if actor.VerifiedAt != nil {
				return nil
			}
			return fun.ErrBadRequest("token already used")
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

// parseActionToken verifies the HMAC signature and expiry, scopes the
// purpose, and extracts the actor id + jti. Every failure collapses to the
// same generic message so the endpoint never leaks token internals.
func (o *Operations) parseActionToken(ctx context.Context, tokenStr string, purpose models.ActionTokenPurpose) (uuid.UUID, uuid.UUID, error) {
	_, span := telemetry.StartSpan(ctx, "parseActionToken")
	defer span.End()

	claims := &models.ActionTokenClaims{}
	_, err := crypto.ParseHMACJWT(tokenStr, claims, o.hmacSecret)
	if err != nil {
		return uuid.Nil, uuid.Nil, fun.ErrBadRequest("invalid or expired token")
	}
	if claims.Purpose != string(purpose) {
		return uuid.Nil, uuid.Nil, fun.ErrBadRequest("invalid or expired token")
	}
	actorID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, uuid.Nil, fun.ErrBadRequest("invalid or expired token")
	}
	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return uuid.Nil, uuid.Nil, fun.ErrBadRequest("invalid or expired token")
	}
	return actorID, jti, nil
}
