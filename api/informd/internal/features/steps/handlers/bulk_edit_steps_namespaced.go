package handlers

import (
	"context"

	"Informd/internal/openapi"
	"Informd/models"
)

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
	return openapi.BulkEditStepsNamespaced204Response{}, nil
}
