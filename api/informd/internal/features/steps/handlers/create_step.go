package handlers

import (
	"context"

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
	return openapi.CreateStep201JSONResponse(*step), nil
}
