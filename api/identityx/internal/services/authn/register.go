package authn

import (
	"context"

	"IdentityX/models"
	"lib/crypto"
	"lib/database"
	"lib/telemetry"
	"strings"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

func (o *Operations) Register(ctx context.Context, in models.IDXRegisterInput) error {
	ctx, span := telemetry.StartSpan(ctx, "Register")
	defer span.End()

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	hashedPassword, err := crypto.Hash(in.Password, crypto.Strong)
	if err != nil {
		return err
	}

	var project *models.Project
	if in.ProjectID != nil {
		project, err = o.projects.GetByID(ctx, *in.ProjectID)
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
		currentTos, hasTos, err = o.tos.CurrentTos(ctx, *in.ProjectID)
		if err != nil {
			return err
		}
		if hasTos && !in.AcceptedTos {
			return fun.ErrValidation("you must accept the terms of service to register")
		}
	}

	// The account and its consent row commit together: a compliance
	// artifact must never exist without the ledger entry that proves it.
	// A failure here rolls back the actor and the registration 500s —
	// the user retries, nothing half-created.
	var actor *models.Actor
	err = database.RunTx(ctx, func(ctx context.Context) error {
		actor, err = o.actors.Register(ctx, models.Actor{
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
			return o.tos.BindRegistration(ctx, actor.ID, *in.ProjectID, currentTos)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Always dispatch a verification email. A failure to enqueue is not
	// the user's fault (the account exists); log it and let the user
	// recover via resend-verification.
	err = o.emailSender.SendVerify(ctx, actor, project)
	if err != nil {
		telemetry.Log().Error("failed to enqueue verification email",
			zap.String("actor_id", actor.ID.String()),
			zap.Error(err),
		)
	}

	return nil
}
