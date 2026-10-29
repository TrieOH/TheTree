package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) GetProductByVendorCode(ctx context.Context, req openapi.GetProductByVendorCodeRequestObject) (openapi.GetProductByVendorCodeResponseObject, error) {
	product, err := h.ops.GetProductByVendorCode(ctx, req.EditionId, req.VendorCode)
	if err != nil {
		return nil, err
	}
	return openapi.GetProductByVendorCode200JSONResponse(*product), nil
}
