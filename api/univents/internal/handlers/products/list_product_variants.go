package products

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListProductVariants(ctx context.Context, req openapi.ListProductVariantsRequestObject) (openapi.ListProductVariantsResponseObject, error) {
	variants, err := h.ops.ListVariantsByProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProductVariants200JSONResponse{
		Code: 200, Data: &variants, Timestamp: time.Now(), Module: module,
	}, nil
}
