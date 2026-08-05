package certifications

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListCertificationTemplateLinks(ctx context.Context, req openapi.ListCertificationTemplateLinksRequestObject) (openapi.ListCertificationTemplateLinksResponseObject, error) {
	links, err := h.ops.ListCertTemplateLinks(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.ListCertificationTemplateLinks200JSONResponse{
		Code: 200, Data: &links, Timestamp: time.Now(), Module: module,
	}, nil
}
