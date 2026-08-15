package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListOrganizationProjectActors(ctx context.Context, req openapi.ListOrganizationProjectActorsRequestObject) (openapi.ListOrganizationProjectActorsResponseObject, error) {
	actors, err := h.ops.ListProjectActors(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationProjectActors200JSONResponse{
		Code: 200, Data: &actors, Timestamp: time.Now(), Module: module,
	}, nil
}
