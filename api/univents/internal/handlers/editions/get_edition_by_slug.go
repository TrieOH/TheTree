package editions

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetEditionBySlug(ctx context.Context, req openapi.GetEditionBySlugRequestObject) (openapi.GetEditionBySlugResponseObject, error) {
	edition, err := h.ops.GetByEventAndEditionSlug(ctx, req.EventSlug, req.EditionSlug)
	if err != nil {
		return nil, err
	}
	return openapi.GetEditionBySlug200JSONResponse{
		Code: 200, Data: edition, Timestamp: time.Now(), Module: module,
	}, nil
}
