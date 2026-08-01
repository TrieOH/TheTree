package authn

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"
	"strings"
)

func (o *Operations) Register(ctx context.Context, in models.IDXRegisterInput) error {
	ctx, span := telemetry.StartSpan(ctx, "Register")
	defer span.End()

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	hashedPassword, err := crypto.Hash(in.Password, crypto.Strong)
	if err != nil {
		return err
	}

	if in.ProjectID != nil {
		_, err := o.projects.GetByID(ctx, *in.ProjectID)
		if err != nil {
			return err
		}
	}

	_, err = o.actors.Register(ctx, models.Actor{
		ProjectID:    in.ProjectID,
		AuthMethod:   models.PasswordAuthMethod,
		PasswordHash: &hashedPassword,
		Email:        &in.Email,
		Type:         models.HumanActorType,
	})
	if err != nil {
		return err
	}

	return nil
}
