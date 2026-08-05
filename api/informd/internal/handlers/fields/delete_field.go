package fields

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) DeleteField(ctx context.Context, req openapi.DeleteFieldRequestObject) (openapi.DeleteFieldResponseObject, error) {
	err := h.ops.Delete(ctx, req.FormId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteField200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
