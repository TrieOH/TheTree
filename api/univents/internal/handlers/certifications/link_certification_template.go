package certifications

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) LinkCertificationTemplate(ctx context.Context, req openapi.LinkCertificationTemplateRequestObject) (openapi.LinkCertificationTemplateResponseObject, error) {
	err := h.ops.LinkCertTemplate(ctx, req.TemplateId, req.Body.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.LinkCertificationTemplate201Response{}, nil
}
