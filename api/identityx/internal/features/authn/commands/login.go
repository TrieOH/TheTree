package commands

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/errx"
	"os"
	"strings"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) Login(ctx context.Context, in models.IDXLoginInput) (tokens *models.UserTokensOutput, err error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	ctx, span := c.tracer.Start(ctx, "Login")
	defer span.End()

	actor, err := c.actors.GetByEmail(ctx, in.Email, nil)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}
	if err != nil {
		return nil, err
	}
	if actor.PasswordHash == nil {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}
	if err = crypto.Verify(in.Password, *actor.PasswordHash); err != nil {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}

	if err = c.actors.UpdateLastLoginAt(ctx, actor.ID); err != nil {
		return nil, err
	}

	activeKeyPair, err := c.cryptoKeys.GetActive(ctx, models.SigningCryptoKeyType, nil)
	if err != nil {
		return nil, err
	}

	accessJTI := uuid.New()
	refreshJTI := uuid.New()
	accessExpiresAt := time.Now().Add(errx.Env[time.Duration]("ACCESS_TOKEN_EXPIRATION", time.ParseDuration, 15*time.Minute))
	refreshExpiresAt := time.Now().Add(errx.Env[time.Duration]("REFRESH_TOKEN_EXPIRATION", time.ParseDuration, 15*time.Minute))

	accessPayload, err := c.newIDXAccessToken(*actor, accessJTI, activeKeyPair.ID, accessExpiresAt)
	if err != nil {
		return nil, err
	}

	refreshPayload, err := c.newIDXRefreshToken(refreshJTI, accessJTI, activeKeyPair.ID, refreshExpiresAt)
	if err != nil {
		return nil, err
	}

	accessToken, err := crypto.SignToken(accessPayload, activeKeyPair.ToKeyPair())
	if err != nil {
		return nil, err
	}

	refreshToken, err := crypto.SignToken(refreshPayload, activeKeyPair.ToKeyPair())
	if err != nil {
		return nil, err
	}

	return &models.UserTokensOutput{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		Domain:           os.Getenv("ISSUER"),
	}, nil
}
