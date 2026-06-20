package commands

import (
	"IdentityX/models"
	"context"
	"lib/api_keys"

	"github.com/MintzyG/fun"
)

func (c *Commands) Create(ctx context.Context, payload models.CreateApiKeyInput) (*models.ApiKey, string, error) {
	ctx, span := c.tracer.Start(ctx, "Create")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, "", err
	}

	var project *models.Project
	project, err = c.projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return nil, "", err
	}

	if ident.Sub.ID != project.OwnerID {
		member, err := c.projects.GetMember(ctx, ident.Sub.ID, payload.ProjectID)
		if err != nil {
			return nil, "", err
		}
		if member.Role != models.ProjectRoleAdmin {
			return nil, "", fun.ErrForbidden("insufficient permissions")
		}
	}

	actorID := ident.Sub.ID
	if payload.CreateForServiceAccount {
		svcAcc, err := c.actors.GetProjectServiceAccount(ctx, payload.ProjectID)
		if err != nil {
			return nil, "", err
		}
		actorID = svcAcc.ID
	}

	generated, err := api_keys.GenerateAPIKey("live")
	if err != nil {
		return nil, "", fun.ErrInternal(err.Error())
	}

	apiKey := models.ApiKey{
		ActorID:   actorID,
		ProjectID: &payload.ProjectID,
		Name:      payload.Name,
		KeyPrefix: generated.Prefix,
		KeyHash:   generated.Hash,
		ExpiresAt: payload.ExpiresAt,
	}

	created, err := c.apiKeys.Create(ctx, apiKey)
	if err != nil {
		return nil, "", err
	}

	return created, generated.Raw, nil
}
