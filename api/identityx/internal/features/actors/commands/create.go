package commands

import (
	"IdentityX/models"
	"context"
	"time"

	"lib/telemetry"

	"github.com/MintzyG/fun"
)

func (c *Commands) Create(ctx context.Context, payload models.CreateActorInput) (*models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = c.authz.CheckProject(ctx, ident.Sub.ID, *payload.ProjectID, nil, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	if payload.Type != models.HumanActorType && payload.AuthMethod == models.PasswordAuthMethod {
		return nil, fun.ErrBadRequest("only human actors can authenticate via password")
	}

	actor := models.Actor{
		ProjectID:  payload.ProjectID,
		AuthMethod: payload.AuthMethod,
		VerifiedAt: new(time.Now()),
		Email:      payload.Email,
		Type:       payload.Type,
	}

	return c.actors.Register(ctx, actor)
}
