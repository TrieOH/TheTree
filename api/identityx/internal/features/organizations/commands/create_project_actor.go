package commands

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) CreateProjectActor(ctx context.Context, orgID uuid.UUID, payload models.CreateActorInput) (*models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "CreateProjectActor")
	defer span.End()
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = c.authz.CheckProject(ctx, ident.Sub.ID, *payload.ProjectID, &orgID, models.ProjectRoleMember)
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
