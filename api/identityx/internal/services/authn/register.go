package authn

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"
	"strings"

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

	actor, err := o.actors.Register(ctx, models.Actor{
		ProjectID:    in.ProjectID,
		AuthMethod:   models.PasswordAuthMethod,
		PasswordHash: &hashedPassword,
		Email:        &in.Email,
		Type:         models.HumanActorType,
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
