package service

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"strings"
)

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
