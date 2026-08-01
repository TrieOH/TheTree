package steps

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) CreateStep(ctx context.Context, req openapi.CreateStepRequestObject) (openapi.CreateStepResponseObject, error) {
	step, err := h.ops.Create(ctx, models.CreateFormStepInput{
		FormID:       req.FormId,
		Title:        req.Body.Title,
		Description:  req.Body.Description,
		PositionHint: req.Body.PositionHint,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateStep201JSONResponse{
		Code: 201, Data: step, Timestamp: time.Now(), Module: module,
	}, nil
}
