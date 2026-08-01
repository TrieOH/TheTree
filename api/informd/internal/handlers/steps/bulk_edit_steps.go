package steps

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/models"
)

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
