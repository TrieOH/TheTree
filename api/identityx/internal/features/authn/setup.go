package authn

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/database"
	"lib/telemetry"
	"strings"
)

func (o *Operations) Setup(ctx context.Context, in models.SetupInput) error {
	ctx, span := telemetry.StartSpan(ctx, "Setup")
	defer span.End()

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	hashedPassword, err := crypto.Hash(in.Password, crypto.Strong)
	if err != nil {
		return err
	}

	var actor *models.Actor
	err = database.RunTx(ctx, func(ctx context.Context) error {
		actor, err = o.actors.Register(ctx, models.Actor{
			AuthMethod:   models.PasswordAuthMethod,
			PasswordHash: &hashedPassword,
			Email:        &in.Email,
			Type:         models.HumanActorType,
		})
		if err != nil {
			return err
		}

		_, err = o.platformRoles.Give(ctx, actor.ID, models.PlatformRoleSuperAdmin, nil)
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
