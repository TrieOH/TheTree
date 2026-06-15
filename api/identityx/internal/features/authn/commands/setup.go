package commands

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"strings"
)

func (c *Commands) Setup(ctx context.Context, in models.SetupInput) error {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	hashedPassword, err := crypto.Hash(in.Password, crypto.Strong)
	if err != nil {
		return err
	}

	var actor *models.Actor
	if err = c.tx.WithinTx(ctx, func(ctx context.Context) error {
		actor, err = c.actors.Register(ctx, models.Actor{
			AuthMethod:   models.PasswordAuthMethod,
			PasswordHash: &hashedPassword,
			Email:        &in.Email,
			Type:         models.HumanActorType,
		})
		if err != nil {
			return err
		}

		_, err = c.platformRoles.Give(ctx, actor.ID, models.PlatformRoleSuperAdmin, nil)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
