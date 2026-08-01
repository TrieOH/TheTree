package certifications

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) UnlinkCertificationTemplate(ctx context.Context, req openapi.UnlinkCertificationTemplateRequestObject) (openapi.UnlinkCertificationTemplateResponseObject, error) {
	err := h.ops.UnlinkCertTemplate(ctx, req.TemplateId, req.Body.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.UnlinkCertificationTemplate204Response{}, nil
}
