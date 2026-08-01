package certifications

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) DeleteCertificationTemplate(ctx context.Context, req openapi.DeleteCertificationTemplateRequestObject) (openapi.DeleteCertificationTemplateResponseObject, error) {
	err := h.ops.DeleteTemplate(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteCertificationTemplate204Response{}, nil
}
