package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionProducts(ctx context.Context, req openapi.ListEditionProductsRequestObject) (openapi.ListEditionProductsResponseObject, error) {
	products, err := h.ops.ListProductsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionProducts200JSONResponse{Items: products}, nil
}
