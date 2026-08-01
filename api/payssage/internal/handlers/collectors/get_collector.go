package collectors

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) GetCollector(ctx context.Context, req openapi.GetCollectorRequestObject) (openapi.GetCollectorResponseObject, error) {
	collector, err := h.ops.GetByID(ctx, req.CollectorId)
	if err != nil {
		return nil, err
	}
	return openapi.GetCollector200JSONResponse{
		Code: 200, Data: collector, Timestamp: time.Now(), Module: module,
	}, nil
}
