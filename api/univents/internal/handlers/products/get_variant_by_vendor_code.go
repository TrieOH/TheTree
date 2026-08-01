package products

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetVariantByVendorCode(ctx context.Context, req openapi.GetVariantByVendorCodeRequestObject) (openapi.GetVariantByVendorCodeResponseObject, error) {
	variant, err := h.ops.GetVariantByVendorCode(ctx, req.EditionId, req.VendorCode)
	if err != nil {
		return nil, err
	}
	return openapi.GetVariantByVendorCode200JSONResponse{
		Code: 200, Data: variant, Timestamp: time.Now(), Module: module,
	}, nil
}
