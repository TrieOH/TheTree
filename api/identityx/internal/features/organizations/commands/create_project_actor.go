package commands

import (
	"IdentityX/models"
	"context"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) CreateProjectActor(ctx context.Context, orgID uuid.UUID, payload models.CreateActorInput) (*models.Actor, error) {
	ctx, span := c.tracer.Start(ctx, "CreateActor")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := c.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	project, err := c.projects.GetByID(ctx, *payload.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.OrganizationID != nil && *project.OrganizationID != orgID {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	if ident.Sub.ID != org.OwnerID {
		_, err = c.orgs.GetMember(ctx, ident.Sub.ID, orgID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			_, err = c.projects.GetMember(ctx, ident.Sub.ID, project.ID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return nil, err
			}
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
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
