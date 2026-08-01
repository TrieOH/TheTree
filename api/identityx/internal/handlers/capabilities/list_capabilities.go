package capabilities

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListCapabilities(ctx context.Context, req openapi.ListCapabilitiesRequestObject) (openapi.ListCapabilitiesResponseObject, error) {
	caps, err := h.ops.List(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListCapabilities200JSONResponse{
		Code: 200, Data: &caps, Timestamp: time.Now(), Module: module,
	}, nil
}
