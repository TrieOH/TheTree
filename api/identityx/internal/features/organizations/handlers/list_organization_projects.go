package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListOrganizationProjects(ctx context.Context, req openapi.ListOrganizationProjectsRequestObject) (openapi.ListOrganizationProjectsResponseObject, error) {
	projects, err := h.ops.ListOrgProjects(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationProjects200JSONResponse{Items: projects}, nil
}
