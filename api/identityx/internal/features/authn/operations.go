// Package authn is the authn feature: its repo, handlers and operations.
// Build wires it.
package authn

import (
	"IdentityX/errs"
	"IdentityX/internal/features/tos"
	"IdentityX/internal/tokens"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"errors"
	"lib/crypto"
	"lib/database"
	"lib/telemetry"
	"log/slog"
	"strings"
	"time"
	"uuid"
)

// Operations owns the session and account lifecycle: password login,
// registration, refresh, logout, setup, JWKS, and the verify/reset email
// flows. OAuth login lives in the oauth_providers module — authn crosses
// its Connect/Callback interface and touches no OAuth state itself.
type Deps struct {
	Actors        ports.ActorRepo
	Projects      ports.ProjectRepo
	PlatformRoles ports.PlatformRolesRepo
	// Tokens owns the token lifecycle (verify/mint/rotate/revoke); login,
	// refresh, logout cross it instead of touching keys, blacklist, or
	// token claims directly.
	Tokens *tokens.Manager
	// ActionTokens owns the single-use action-token lifecycle; verify and
	// reset links are redeemed through it instead of touching the HMAC
	// secret or the anti-replay repo directly.
	ActionTokens *tokens.ActionTokenManager
	EmailSender  ports.EmailSender
	// Tos owns the terms documents and the consent ledger; registration
	// gates on it and pins the acceptance row through it.
	Tos *tos.Operations
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// JWKS publishes the project's public signing keys as a JWKS document. The
// tokens module owns key policy: scope validation and the key read live
// there, not here.
func (o *Operations) JWKS(ctx context.Context, projectID *uuid.UUID) (map[string]any, error) {
	return o.Tokens.JWKS(ctx, projectID)
}

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
				slog.String("actor_id", actor.ID.String()),
				slog.String("project_id", actor.ProjectID.String()),
				telemetry.Err(err),
			)
		}
	}
	return o.Tokens.Mint(ctx, actor)
}

// Logout ends the session: the tokens module blacklists the access token
// and, when the refresh token is still verifiable, the refresh token too.
// A dead refresh token never fails the logout — the access token is
// already revoked, so the session is over either way.
func (o *Operations) Logout(ctx context.Context, in models.LogoutInput) error {
	return o.Tokens.Revoke(ctx, in.AccessToken, in.RefreshToken)
}

// ForgotPassword mints and dispatches a reset link for the actor's email.
// It always succeeds from the caller's perspective (anti-enumeration):
// unknown emails and passwordless (OAuth-only) accounts are silent no-ops.
func (o *Operations) ForgotPassword(ctx context.Context, in models.ForgotPasswordInput) error {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	actor, err := o.Actors.GetByEmail(ctx, email, in.ProjectID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
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
		project, err = o.Projects.GetByID(ctx, *in.ProjectID)
		if err != nil {
			return err
		}
	}
	err = o.EmailSender.SendReset(ctx, actor, project)
	if err != nil {
		telemetry.Log().Error("failed to enqueue password reset email",
			slog.String("actor_id", actor.ID.String()),
			telemetry.Err(err),
		)
	}
	return nil
}

// ResetPassword redeems a single-use reset link through the token module:
// validates the HMAC JWT (signature, purpose, expiry), consumes the jti
// anti-replay record, and replaces the actor's password hash. Unlike
// verification, reset is never idempotent — a consumed token is dead, and
// the user requests a new link if they need another reset.
func (o *Operations) ResetPassword(ctx context.Context, in models.ResetPasswordInput) error {
	actorID, err := o.ActionTokens.Redeem(ctx, models.PasswordResetActionTokenPurpose, in.Token)
	if err != nil {
		return err
	}

	hashed, err := crypto.Hash(in.NewPassword, crypto.Strong)
	if err != nil {
		return err
	}
	return o.Actors.UpdatePasswordHash(ctx, actorID, hashed)
}

// Refresh rotates the session: the tokens module verifies the old refresh
// token, blacklists it and the access token it anchors, and mints a fresh
// pair for the actor. Fail-closed — if the old pair cannot be revoked, no
// new pair is issued.
func (o *Operations) Refresh(ctx context.Context, refreshToken string) (*models.UserTokensOutput, error) {
	return o.Tokens.Rotate(ctx, refreshToken)
}

func (o *Operations) Register(ctx context.Context, in models.IDXRegisterInput) error {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	hashedPassword, err := crypto.Hash(in.Password, crypto.Strong)
	if err != nil {
		return err
	}

	var project *models.Project
	if in.ProjectID != nil {
		project, err = o.Projects.GetByID(ctx, *in.ProjectID)
		if err != nil {
			return err
		}
	}

	// Consent gate: a project with terms only accepts registrations that
	// explicitly agreed to them. Projects without terms (v0) skip the
	// gate entirely — there is nothing to agree to.
	var currentTos *models.TermsOfService
	if in.ProjectID != nil {
		var hasTos bool
		currentTos, hasTos, err = o.Tos.CurrentTos(ctx, *in.ProjectID)
		if err != nil {
			return err
		}
		if hasTos && !in.AcceptedTos {
			return errs.ErrTermsNotAccepted
		}
	}

	// The account and its consent row commit together: a compliance
	// artifact must never exist without the ledger entry that proves it.
	// A failure here rolls back the actor and the registration 500s —
	// the user retries, nothing half-created.
	var actor *models.Actor
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		actor, err = o.Actors.Register(ctx, models.Actor{
			ProjectID:    in.ProjectID,
			AuthMethod:   models.PasswordAuthMethod,
			PasswordHash: &hashedPassword,
			Email:        &in.Email,
			Type:         models.HumanActorType,
		})
		if err != nil {
			return err
		}
		if currentTos != nil {
			return o.Tos.BindRegistration(ctx, actor.ID, *in.ProjectID, currentTos)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Always dispatch a verification email. A failure to enqueue is not
	// the user's fault (the account exists); log it and let the user
	// recover via resend-verification.
	err = o.EmailSender.SendVerify(ctx, actor, project)
	if err != nil {
		telemetry.Log().Error("failed to enqueue verification email",
			slog.String("actor_id", actor.ID.String()),
			telemetry.Err(err),
		)
	}

	return nil
}

func (o *Operations) Setup(ctx context.Context, in models.SetupInput) error {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	hashedPassword, err := crypto.Hash(in.Password, crypto.Strong)
	if err != nil {
		return err
	}

	var actor *models.Actor
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		actor, err = o.Actors.Register(ctx, models.Actor{
			AuthMethod:   models.PasswordAuthMethod,
			PasswordHash: &hashedPassword,
			Email:        &in.Email,
			Type:         models.HumanActorType,
		})
		if err != nil {
			return err
		}

		_, err = o.PlatformRoles.Give(ctx, actor.ID, models.PlatformRoleSuperAdmin, nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// VerifyEmail redeems a single-use verify link through the token module:
// validates the HMAC JWT (signature, purpose, expiry), consumes the jti
// anti-replay record, and stamps the actor verified. A consumed link on an
// already-verified account succeeds (idempotent) — that fall-through is
// actor-state policy, so it lives here: the anti-replay table still blocks
// replaying a token against an unverified account.
func (o *Operations) VerifyEmail(ctx context.Context, in models.VerifyEmailInput) error {
	actorID, err := o.ActionTokens.Redeem(ctx, models.EmailVerifyActionTokenPurpose, in.Token)
	if err != nil {
		if !errors.Is(err, errs.ErrActionTokenUsed) {
			return err
		}
		actor, aerr := o.Actors.GetByID(ctx, actorID)
		if aerr != nil {
			return aerr
		}
		if actor.VerifiedAt != nil {
			return nil
		}
		return err
	}

	return o.Actors.SetVerifiedAt(ctx, actorID, time.Now())
}

// ResendVerification mints and dispatches a fresh verify link. It always
// succeeds from the caller's perspective (anti-enumeration): unknown
// emails and already-verified accounts are silent no-ops.
func (o *Operations) ResendVerification(ctx context.Context, in models.ResendVerificationInput) error {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	actor, err := o.Actors.GetByEmail(ctx, email, in.ProjectID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil
		}
		return err
	}
	if actor.VerifiedAt != nil || actor.Email == nil {
		return nil
	}

	var project *models.Project
	if in.ProjectID != nil {
		project, err = o.Projects.GetByID(ctx, *in.ProjectID)
		if err != nil {
			return err
		}
	}
	err = o.EmailSender.SendVerify(ctx, actor, project)
	if err != nil {
		telemetry.Log().Error("failed to enqueue verification email",
			slog.String("actor_id", actor.ID.String()),
			telemetry.Err(err),
		)
	}
	return nil
}
