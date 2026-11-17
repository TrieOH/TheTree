package steps

import (
	"context"

	"Informd/internal/features/steps/service"
	"Informd/internal/openapi"
	"Informd/models"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) BulkEditSteps(ctx context.Context, req openapi.BulkEditStepsRequestObject) (openapi.BulkEditStepsResponseObject, error) {
	payload := make([]models.UpdateFormStepInput, 0, len(*req.Body))
	for _, s := range *req.Body {
		payload = append(payload, models.UpdateFormStepInput{
			FormID:       req.FormId,
			ID:           s.Id,
			Title:        s.Title,
			Description:  s.Description,
			PositionHint: s.PositionHint,
		})
	}
	err := h.ops.BulkEdit(ctx, req.FormId, payload)
	if err != nil {
		return nil, err
	}
	return openapi.BulkEditSteps204Response{}, nil
}

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

func (h *Handlers) ListSteps(ctx context.Context, req openapi.ListStepsRequestObject) (openapi.ListStepsResponseObject, error) {
	steps, err := h.ops.List(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ListSteps200JSONResponse{Items: steps}, nil
}
