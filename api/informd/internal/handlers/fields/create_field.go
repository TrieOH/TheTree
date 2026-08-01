package fields

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) CreateField(ctx context.Context, req openapi.CreateFieldRequestObject) (openapi.CreateFieldResponseObject, error) {
	field, err := h.ops.Create(ctx, models.CreateStepFieldInput{
		FormID:       req.FormId,
		StepID:       req.StepId,
		Key:          req.Body.Key,
		Title:        req.Body.Title,
		Description:  req.Body.Description,
		PositionHint: req.Body.PositionHint,
		Required:     req.Body.Required,
		Type:         req.Body.Type,
		Placeholder:  mustMarshal(req.Body.Placeholder),
		DefaultValue: mustMarshal(req.Body.DefaultValue),
		Config:       mustMarshal(req.Body.Config),
		SelectConfig: mapSelectConfig(req.Body.SelectConfig),
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateField201JSONResponse{
		Code: 201, Data: field, Timestamp: time.Now(), Module: module,
	}, nil
}
