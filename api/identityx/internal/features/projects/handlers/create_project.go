package handlers

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) CreateProject(ctx context.Context, req openapi.CreateProjectRequestObject) (openapi.CreateProjectResponseObject, error) {
	project, err := h.ops.Create(ctx, models.CreateProjectInput{
		Name:      req.Body.Name,
		Domain:    req.Body.Domain,
		BrandSlug: req.Body.BrandSlug,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateProject201JSONResponse(*project), nil
}
