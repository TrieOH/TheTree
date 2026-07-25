package commands

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"
	"strings"
)

func (c *Commands) Register(ctx context.Context, in models.IDXRegisterInput) error {
	ctx, span := telemetry.StartSpan(ctx, "Register")
	defer span.End()

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	hashedPassword, err := crypto.Hash(in.Password, crypto.Strong)
	if err != nil {
		return err
	}

	if in.ProjectID != nil {
		_, err := c.projects.GetByID(ctx, *in.ProjectID)
		if err != nil {
			return err
		}
	}

	_, err = c.actors.Register(ctx, models.Actor{
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
