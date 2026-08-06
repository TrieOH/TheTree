package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) ListOrganizationProjectActors(ctx context.Context, req openapi.ListOrganizationProjectActorsRequestObject) (openapi.ListOrganizationProjectActorsResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	actors, err := h.ops.ListProjectActors(ctx, req.OrgId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationProjectActors200JSONResponse{
		Code: 200, Data: &actors, Timestamp: time.Now(), Module: module,
	}, nil
}
