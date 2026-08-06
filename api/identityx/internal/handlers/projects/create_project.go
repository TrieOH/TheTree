package projects

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) CreateProject(ctx context.Context, req openapi.CreateProjectRequestObject) (openapi.CreateProjectResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	project, err := h.ops.Create(ctx, models.CreateProjectInput{
		Name:      req.Body.Name,
		Domain:    req.Body.Domain,
		BrandSlug: req.Body.BrandSlug,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateProject201JSONResponse{
		Code: 201, Data: project, Timestamp: time.Now(), Module: module,
	}, nil
}
