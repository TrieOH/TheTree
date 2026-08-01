package products

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) PatchProduct(ctx context.Context, req openapi.PatchProductRequestObject) (openapi.PatchProductResponseObject, error) {
	product, err := h.ops.PatchProduct(ctx, models.PatchProductInput{
		ProductID:            req.ProductId,
		VendorCode:           req.Body.VendorCode,
		RequiresRegistration: req.Body.RequiresRegistration,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PatchProduct200JSONResponse{
		Code: 200, Data: product, Timestamp: time.Now(), Module: module,
	}, nil
}
