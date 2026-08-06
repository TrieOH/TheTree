package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) ListOrganizations(ctx context.Context, _ openapi.ListOrganizationsRequestObject) (openapi.ListOrganizationsResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	orgs, err := h.ops.ListOrgs(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizations200JSONResponse{
		Code: 200, Data: &orgs, Timestamp: time.Now(), Module: module,
	}, nil
}
