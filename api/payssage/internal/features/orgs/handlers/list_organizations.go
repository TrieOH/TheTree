package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) ListOrganizations(ctx context.Context, _ openapi.ListOrganizationsRequestObject) (openapi.ListOrganizationsResponseObject, error) {
	orgs, err := h.ops.ListOrgs(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizations200JSONResponse{Items: orgs}, nil
}
