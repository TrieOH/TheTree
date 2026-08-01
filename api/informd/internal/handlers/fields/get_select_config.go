package fields

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) GetSelectConfig(ctx context.Context, req openapi.GetSelectConfigRequestObject) (openapi.GetSelectConfigResponseObject, error) {
	config, err := h.ops.GetSelectConfig(ctx, req.FormId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSelectConfig200JSONResponse{
		Code: 200, Data: config, Timestamp: time.Now(), Module: module,
	}, nil
}
