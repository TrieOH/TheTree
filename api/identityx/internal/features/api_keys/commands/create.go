package commands

import (
	"IdentityX/models"
	"context"
	"lib/api_keys"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

func (c *Commands) Create(ctx context.Context, payload models.CreateAPIKeyInput) (*models.APIKey, string, error) {
	ctx, span := c.tracer.Start(ctx, "Create")
	defer span.End()

	project, err := c.projects.GetByID(ctx, *payload.ProjectID)
	if err != nil {
		telemetry.Log().Info("Create api key", zap.String("project_id", payload.ProjectID.String()))
		return nil, "", err
	}

	var created *models.APIKey
	var generated *api_keys.GeneratedAPIKey
	if err = c.tx.WithinTx(ctx, func(ctx context.Context) error {
		created, generated, err = c.createInternal(ctx, *project, payload)
		return err
	}); err != nil {
		return nil, "", err
	}

	return created, generated.Raw, nil
}

func (c *Commands) createInternal(ctx context.Context, project models.Project, payload models.CreateAPIKeyInput) (*models.APIKey, *api_keys.GeneratedAPIKey, error) {
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

	apiKey := models.APIKey{
		SubjectID:     actorID,
		Name:          payload.Name,
		DisplayPrefix: generated.DisplayPrefix,
		KeyHash:       generated.Hash,
		ExpiresAt:     payload.ExpiresAt,
		CreatedBy:     ident.Sub.ID,
	}

	var created *models.APIKey
	created, err = c.apiKeys.Create(ctx, apiKey)
	if err != nil {
		return nil, nil, err
	}

	if len(payload.Capabilities) > 0 {
		err = c.capabilities.AssignToAPIKey(ctx, created.ID, payload.Capabilities, ident.Sub.ID)
		if err != nil {
			return nil, nil, err
		}
	}

	return created, generated, nil
}
