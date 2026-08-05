package collectors

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) ListCollectors(ctx context.Context, _ openapi.ListCollectorsRequestObject) (openapi.ListCollectorsResponseObject, error) {
	collectors, err := h.ops.ListOwned(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListCollectors200JSONResponse{
		Code: 200, Data: &collectors, Timestamp: time.Now(), Module: module,
	}, nil
}
