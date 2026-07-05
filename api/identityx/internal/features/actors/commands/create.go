package commands

import (
	"IdentityX/models"
	"context"
	"time"

	"github.com/MintzyG/fun"
)

func (c *Commands) Create(ctx context.Context, payload models.CreateActorInput) (*models.Actor, error) {
	ctx, span := c.tracer.Start(ctx, "Create")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var project *models.Project
	project, err = c.projects.GetByID(ctx, *payload.ProjectID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != project.OwnerID {
		member, err := c.projects.GetMember(ctx, ident.Sub.ID, project.ID)
		if err != nil {
			return nil, err
		}
		if member.Role != models.ProjectRoleAdmin {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
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
