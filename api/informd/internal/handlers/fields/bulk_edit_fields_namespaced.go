package fields

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) BulkEditFieldsNamespaced(ctx context.Context, req openapi.BulkEditFieldsNamespacedRequestObject) (openapi.BulkEditFieldsNamespacedResponseObject, error) {
	payload := make([]models.UpdateNamespacedStepFieldInput, 0, len(*req.Body))
	for _, f := range *req.Body {
		payload = append(payload, models.UpdateNamespacedStepFieldInput{
			NamespaceID:  req.NamespaceId,
			FormID:       req.FormId,
			StepID:       req.StepId,
			ID:           f.Id,
			Key:          f.Key,
			Title:        f.Title,
			Description:  f.Description,
			PositionHint: f.PositionHint,
			Required:     f.Required,
			Type:         f.Type,
			Placeholder:  mustMarshal(f.Placeholder),
			DefaultValue: mustMarshal(f.DefaultValue),
			Config:       mustMarshal(f.Config),
			SelectConfig: mapSelectConfig(f.SelectConfig),
		})
	}
	err := h.ops.BulkEditNamespaced(ctx, req.FormId, req.NamespaceId, payload)
	if err != nil {
		return nil, err
	}
	return openapi.BulkEditFieldsNamespaced200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
