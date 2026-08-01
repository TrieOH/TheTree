package products

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) DeleteProduct(ctx context.Context, req openapi.DeleteProductRequestObject) (openapi.DeleteProductResponseObject, error) {
	err := h.ops.DeleteProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteProduct204Response{}, nil
}
