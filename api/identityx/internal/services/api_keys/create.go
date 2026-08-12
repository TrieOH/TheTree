package api_keys

import (
	"IdentityX/models"
	"context"
	"lib/api_keys"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateAPIKeyInput) (*models.APIKey, string, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	project, err := o.projects.GetByID(ctx, *payload.ProjectID)
	if err != nil {
		telemetry.Log().Info("Create api key", zap.String("project_id", payload.ProjectID.String()))
		return nil, "", err
	}

	created, generated, err := o.createInternal(ctx, *project, payload)
	if err != nil {
		return nil, "", err
	}

	return created, generated.Raw, nil
}

func (o *Operations) createInternal(ctx context.Context, project models.Project, payload models.CreateAPIKeyInput) (*models.APIKey, *api_keys.GeneratedAPIKey, error) {
	ctx, span := telemetry.StartSpan(ctx, "createInternal")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, nil, err
	}

	actorID := ident.Sub.ID
	if payload.SubjectID != nil {
		actorID = *payload.SubjectID
		err = o.authz.CheckProject(ctx, ident.Sub.ID, project.ID, nil, models.ProjectRoleAdmin)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(payload.Capabilities) > 0 {
		valid, err := o.capabilities.ValidateCapabilities(ctx, &project.ID, payload.Capabilities)
		if err != nil {
			return nil, nil, err
		}
		if !valid {
			return nil, nil, fun.ErrBadRequest("invalid capabilities")
		}
	}

	var generated *api_keys.GeneratedAPIKey
	generated, err = api_keys.GenerateAPIKey(project.BrandSlug, payload.Env, o.hmacSecret)
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
	created, err = o.apiKeys.Create(ctx, apiKey)
	if err != nil {
		return nil, nil, err
	}

	if len(payload.Capabilities) > 0 {
		err = o.capabilities.AssignToAPIKey(ctx, created.ID, payload.Capabilities, ident.Sub.ID)
		if err != nil {
			return nil, nil, err
		}
	}

	return created, generated, nil
}
