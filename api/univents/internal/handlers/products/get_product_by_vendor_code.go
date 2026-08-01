package products

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetProductByVendorCode(ctx context.Context, req openapi.GetProductByVendorCodeRequestObject) (openapi.GetProductByVendorCodeResponseObject, error) {
	product, err := h.ops.GetProductByVendorCode(ctx, req.EditionId, req.VendorCode)
	if err != nil {
		return nil, err
	}
	return openapi.GetProductByVendorCode200JSONResponse{
		Code: 200, Data: product, Timestamp: time.Now(), Module: module,
	}, nil
}
