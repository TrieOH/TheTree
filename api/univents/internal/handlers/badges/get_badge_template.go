package badges

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetBadgeTemplate(ctx context.Context, req openapi.GetBadgeTemplateRequestObject) (openapi.GetBadgeTemplateResponseObject, error) {
	template, err := h.ops.GetTemplate(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.GetBadgeTemplate200JSONResponse{
		Code: 200, Data: template, Timestamp: time.Now(), Module: module,
	}, nil
}
