package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) GetProduct(ctx context.Context, req openapi.GetProductRequestObject) (openapi.GetProductResponseObject, error) {
	product, err := h.ops.GetProductByID(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProduct200JSONResponse(*product), nil
}
