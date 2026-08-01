package products

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) DeleteProductVariant(ctx context.Context, req openapi.DeleteProductVariantRequestObject) (openapi.DeleteProductVariantResponseObject, error) {
	err := h.ops.DeleteVariant(ctx, req.VariantId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteProductVariant204Response{}, nil
}
