package steps

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/models"
)

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
