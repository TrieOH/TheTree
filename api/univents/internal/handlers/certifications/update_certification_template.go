package certifications

import (
	"encoding/json"

	"context"
	"time"

	"github.com/MintzyG/fun"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) UpdateCertificationTemplate(ctx context.Context, req openapi.UpdateCertificationTemplateRequestObject) (openapi.UpdateCertificationTemplateResponseObject, error) {
	design, err := json.Marshal(req.Body.DesignData)
	if err != nil {
		return nil, fun.ErrBadRequest("invalid design_data")
	}
	template, err := h.ops.UpdateTemplate(ctx, models.UpdateCertificationTemplateInput{
		TemplateID:  req.TemplateId,
		Kind:        req.Body.Kind,
		Name:        req.Body.Name,
		Description: req.Body.Description,
		DesignData:  design,
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpdateCertificationTemplate200JSONResponse{
		Code: 200, Data: template, Timestamp: time.Now(), Module: module,
	}, nil
}
