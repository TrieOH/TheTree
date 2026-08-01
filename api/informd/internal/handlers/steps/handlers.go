// Package steps implements the StrictServerInterface methods for the
// steps feature, including the duplicated namespaced routes.
package steps

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/internal/services"
	"Informd/models"
)

const module = "Informd"

type Handlers struct {
	ops *services.Steps
}

func New(ops *services.Steps) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListSteps(ctx context.Context, req openapi.ListStepsRequestObject) (openapi.ListStepsResponseObject, error) {
	steps, err := h.ops.List(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ListSteps200JSONResponse{
		Code: 200, Data: &steps, Timestamp: time.Now(), Module: module,
	}, nil
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
	return openapi.CreateStep201JSONResponse{
		Code: 201, Data: step, Timestamp: time.Now(), Module: module,
	}, nil
}

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
	return openapi.BulkEditSteps200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListStepsNamespaced(ctx context.Context, req openapi.ListStepsNamespacedRequestObject) (openapi.ListStepsNamespacedResponseObject, error) {
	steps, err := h.ops.ListNamespaced(ctx, req.FormId, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListStepsNamespaced200JSONResponse{
		Code: 200, Data: &steps, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateStepNamespaced(ctx context.Context, req openapi.CreateStepNamespacedRequestObject) (openapi.CreateStepNamespacedResponseObject, error) {
	step, err := h.ops.CreateNamespaced(ctx, models.CreateNamespacedFormStepInput{
		NamespaceID:  req.NamespaceId,
		FormID:       req.FormId,
		Title:        req.Body.Title,
		Description:  req.Body.Description,
		PositionHint: req.Body.PositionHint,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateStepNamespaced201JSONResponse{
		Code: 201, Data: step, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) BulkEditStepsNamespaced(ctx context.Context, req openapi.BulkEditStepsNamespacedRequestObject) (openapi.BulkEditStepsNamespacedResponseObject, error) {
	payload := make([]models.UpdateNamespacedFormStepInput, 0, len(*req.Body))
	for _, s := range *req.Body {
		payload = append(payload, models.UpdateNamespacedFormStepInput{
			NamespaceID:  req.NamespaceId,
			FormID:       req.FormId,
			ID:           s.Id,
			Title:        s.Title,
			Description:  s.Description,
			PositionHint: s.PositionHint,
		})
	}
	err := h.ops.BulkEditNamespaced(ctx, req.FormId, req.NamespaceId, payload)
	if err != nil {
		return nil, err
	}
	return openapi.BulkEditStepsNamespaced200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
