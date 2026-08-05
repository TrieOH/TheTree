package collectors

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) ListOrganizationCollectors(ctx context.Context, req openapi.ListOrganizationCollectorsRequestObject) (openapi.ListOrganizationCollectorsResponseObject, error) {
	collectors, err := h.ops.ListByOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationCollectors200JSONResponse{
		Code: 200, Data: &collectors, Timestamp: time.Now(), Module: module,
	}, nil
}
