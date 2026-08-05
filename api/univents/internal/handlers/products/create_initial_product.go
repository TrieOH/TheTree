package products

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) CreateInitialProduct(ctx context.Context, req openapi.CreateInitialProductRequestObject) (openapi.CreateInitialProductResponseObject, error) {
	requiresRegistration := false
	if req.Body.RequiresRegistration != nil {
		requiresRegistration = *req.Body.RequiresRegistration
	}
	product, err := h.ops.CreateInitial(ctx, models.CreateInitialProductInput{
		EditionID:            req.EditionId,
		RequiresRegistration: requiresRegistration,
		VendorCode:           req.Body.VendorCode,
		VariantVendorCode:    req.Body.VariantVendorCode,
		Name:                 req.Body.Name,
		Description:          req.Body.Description,
		Price:                req.Body.Price,
		Stock:                req.Body.Stock,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateInitialProduct201JSONResponse{
		Code: 201, Data: product, Timestamp: time.Now(), Module: module,
	}, nil
}
