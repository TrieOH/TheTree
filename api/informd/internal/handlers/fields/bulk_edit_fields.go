package fields

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) BulkEditFields(ctx context.Context, req openapi.BulkEditFieldsRequestObject) (openapi.BulkEditFieldsResponseObject, error) {
	payload := make([]models.UpdateStepFieldInput, 0, len(*req.Body))
	for _, f := range *req.Body {
		payload = append(payload, models.UpdateStepFieldInput{
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
	err := h.ops.BulkEdit(ctx, req.FormId, payload)
	if err != nil {
		return nil, err
	}
	return openapi.BulkEditFields200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
