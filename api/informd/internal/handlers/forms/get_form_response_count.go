package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) GetFormResponseCount(ctx context.Context, req openapi.GetFormResponseCountRequestObject) (openapi.GetFormResponseCountResponseObject, error) {
	count, err := h.ops.GetResponseCount(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFormResponseCount200JSONResponse{
		Code: 200,
		Data: &struct {
			Count int `json:"count"`
		}{Count: count},
		Timestamp: time.Now(), Module: module,
	}, nil
}
