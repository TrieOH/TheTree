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

	project, err := c.projects.GetByID(ctx, *payload.ProjectID)
	if err != nil {
		return nil, "", err
	}

	var created *models.ApiKey
	var generated *api_keys.GeneratedAPIKey
	if err = c.tx.WithinTx(ctx, func(ctx context.Context) error {
		created, generated, err = c.createInternal(ctx, *project, payload)
		return err
	}); err != nil {
		return nil, "", err
	}

	return created, generated.Raw, nil
}

func (c *Commands) createInternal(ctx context.Context, project models.Project, payload models.CreateApiKeyInput) (*models.ApiKey, *api_keys.GeneratedAPIKey, error) {
	ctx, span := c.tracer.Start(ctx, "createInternal")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, nil, err
	}

	actorID := ident.Sub.ID
	if payload.SubjectID != nil {
		actorID = *payload.SubjectID
		if ident.Sub.ID != project.OwnerID {
			member, err := c.projects.GetMember(ctx, ident.Sub.ID, project.ID)
			if err != nil {
				return nil, nil, err
			}
			if member.Role != models.ProjectRoleAdmin {
				return nil, nil, fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	_, err = c.projects.GetMember(ctx, actorID, project.ID)
	if err != nil {
		return nil, nil, err
	}

	if len(payload.Capabilities) > 0 {
		valid, err := c.capabilities.ValidateCapabilities(ctx, &project.ID, payload.Capabilities)
		if err != nil {
			return nil, nil, err
		}
		if !valid {
			return nil, nil, fun.ErrBadRequest("invalid capabilities")
		}
	}

	var generated *api_keys.GeneratedAPIKey
	generated, err = api_keys.GenerateAPIKey(project.BrandSlug, payload.Env, c.hmacSecret)
	if err != nil {
		return nil, nil, fun.ErrInternal(err.Error())
	}

	apiKey := models.ApiKey{
		SubjectID:     actorID,
		Name:          payload.Name,
		DisplayPrefix: generated.DisplayPrefix,
		KeyHash:       generated.Hash,
		ExpiresAt:     payload.ExpiresAt,
		CreatedBy:     ident.Sub.ID,
	}

	var created *models.ApiKey
	created, err = c.apiKeys.Create(ctx, apiKey)
	if err != nil {
		return nil, nil, err
	}

	if len(payload.Capabilities) > 0 {
		if err = c.capabilities.AssignToApiKey(ctx, created.ID, payload.Capabilities, ident.Sub.ID); err != nil {
			return nil, nil, err
		}
	}

	return created, generated, nil
}
