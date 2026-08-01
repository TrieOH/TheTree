package products

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) CreateProductVariant(ctx context.Context, req openapi.CreateProductVariantRequestObject) (openapi.CreateProductVariantResponseObject, error) {
	variant, err := h.ops.CreateVariant(ctx, models.CreateProductVariantInput{
		ProductID:   req.ProductId,
		VendorCode:  req.Body.VendorCode,
		Name:        req.Body.Name,
		Description: req.Body.Description,
		Price:       req.Body.Price,
		Stock:       req.Body.Stock,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateProductVariant201JSONResponse{
		Code: 201, Data: variant, Timestamp: time.Now(), Module: module,
	}, nil
}
