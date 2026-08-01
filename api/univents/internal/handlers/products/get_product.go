package products

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetProduct(ctx context.Context, req openapi.GetProductRequestObject) (openapi.GetProductResponseObject, error) {
	product, err := h.ops.GetProductByID(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProduct200JSONResponse{
		Code: 200, Data: product, Timestamp: time.Now(), Module: module,
	}, nil
}
