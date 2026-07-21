package commands

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"strings"

	"github.com/MintzyG/fun"
)

func (c *Commands) Login(ctx context.Context, in models.IDXLoginInput) (tokens *models.UserTokensOutput, err error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	ctx, span := c.tracer.Start(ctx, "Login")
	defer span.End()

	actor, err := c.actors.GetByEmail(ctx, in.Email, in.ProjectID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}
	if err != nil {
		return nil, err
	}
	if actor.PasswordHash == nil {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}
	err = crypto.Verify(in.Password, *actor.PasswordHash)
	if err != nil {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}
	err = c.actors.UpdateLastLoginAt(ctx, actor.ID)
	if err != nil {
		return nil, err
	}
	return c.issueTokens(ctx, actor)
}
