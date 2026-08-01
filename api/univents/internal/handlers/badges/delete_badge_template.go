package badges

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) DeleteBadgeTemplate(ctx context.Context, req openapi.DeleteBadgeTemplateRequestObject) (openapi.DeleteBadgeTemplateResponseObject, error) {
	err := h.ops.DeleteTemplate(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteBadgeTemplate204Response{}, nil
}
