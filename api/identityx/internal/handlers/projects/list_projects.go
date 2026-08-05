package projects

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListProjects(ctx context.Context, _ openapi.ListProjectsRequestObject) (openapi.ListProjectsResponseObject, error) {
	projects, err := h.ops.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListProjects200JSONResponse{
		Code: 200, Data: &projects, Timestamp: time.Now(), Module: module,
	}, nil
}
