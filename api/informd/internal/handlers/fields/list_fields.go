package fields

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListFields(ctx context.Context, req openapi.ListFieldsRequestObject) (openapi.ListFieldsResponseObject, error) {
	fields, err := h.ops.List(ctx, req.FormId, req.StepId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFields200JSONResponse{
		Code: 200, Data: &fields, Timestamp: time.Now(), Module: module,
	}, nil
}
