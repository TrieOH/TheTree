package badges

import (
	"context"
	"encoding/json"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) UpdateBadgeTemplate(ctx context.Context, req openapi.UpdateBadgeTemplateRequestObject) (openapi.UpdateBadgeTemplateResponseObject, error) {
	input := models.UpdateBadgeTemplateInput{TemplateID: req.TemplateId}

	if req.Body.Name != nil {
		name := *req.Body.Name
		input.Name = &name
	}
	if req.Body.DesignData != nil {
		design, err := json.Marshal(req.Body.DesignData)
		if err != nil {
			return nil, err
		}
		input.DesignData = design
	}

	template, err := h.ops.UpdateTemplate(ctx, input)
	if err != nil {
		return nil, err
	}
	return openapi.UpdateBadgeTemplate200JSONResponse{
		Code: 200, Data: template, Timestamp: time.Now(), Module: module,
	}, nil
}
