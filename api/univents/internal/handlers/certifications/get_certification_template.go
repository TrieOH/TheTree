package certifications

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetCertificationTemplate(ctx context.Context, req openapi.GetCertificationTemplateRequestObject) (openapi.GetCertificationTemplateResponseObject, error) {
	template, err := h.ops.GetTemplateByID(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.GetCertificationTemplate200JSONResponse{
		Code: 200, Data: template, Timestamp: time.Now(), Module: module,
	}, nil
}
