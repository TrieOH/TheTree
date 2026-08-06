package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) CreateOrganizationProject(ctx context.Context, req openapi.CreateOrganizationProjectRequestObject) (openapi.CreateOrganizationProjectResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	project, err := h.ops.CreateProject(ctx, models.CreateOrgProjectInput{
		OrganizationID: req.OrganizationId,
		Name:           req.Body.Name,
		Domain:         req.Body.Domain,
		BrandSlug:      req.Body.BrandSlug,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateOrganizationProject201JSONResponse{
		Code: 201, Data: project, Timestamp: time.Now(), Module: module,
	}, nil
}
