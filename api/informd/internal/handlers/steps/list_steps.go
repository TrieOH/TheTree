package steps

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListSteps(ctx context.Context, req openapi.ListStepsRequestObject) (openapi.ListStepsResponseObject, error) {
	steps, err := h.ops.List(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ListSteps200JSONResponse{
		Code: 200, Data: &steps, Timestamp: time.Now(), Module: module,
	}, nil
}
