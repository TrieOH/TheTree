// Package products implements the StrictServerInterface methods for the
// products feature.
package products

import (
	"context"
	"encoding/json"
	"time"

	"univents/internal/openapi"
	"univents/internal/services"
	"univents/models"

	"github.com/MintzyG/fun"
)

const module = "Univents"

type Handlers struct {
	ops *services.Products
}

func New(ops *services.Products) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListEditionProducts(ctx context.Context, req openapi.ListEditionProductsRequestObject) (openapi.ListEditionProductsResponseObject, error) {
	products, err := h.ops.ListProductsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionProducts200JSONResponse{
		Code: 200, Data: &products, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetProductByVendorCode(ctx context.Context, req openapi.GetProductByVendorCodeRequestObject) (openapi.GetProductByVendorCodeResponseObject, error) {
	product, err := h.ops.GetProductByVendorCode(ctx, req.EditionId, req.VendorCode)
	if err != nil {
		return nil, err
	}
	return openapi.GetProductByVendorCode200JSONResponse{
		Code: 200, Data: product, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetVariantByVendorCode(ctx context.Context, req openapi.GetVariantByVendorCodeRequestObject) (openapi.GetVariantByVendorCodeResponseObject, error) {
	variant, err := h.ops.GetVariantByVendorCode(ctx, req.EditionId, req.VendorCode)
	if err != nil {
		return nil, err
	}
	return openapi.GetVariantByVendorCode200JSONResponse{
		Code: 200, Data: variant, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetProduct(ctx context.Context, req openapi.GetProductRequestObject) (openapi.GetProductResponseObject, error) {
	product, err := h.ops.GetProductByID(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProduct200JSONResponse{
		Code: 200, Data: product, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListProductVariants(ctx context.Context, req openapi.ListProductVariantsRequestObject) (openapi.ListProductVariantsResponseObject, error) {
	variants, err := h.ops.ListVariantsByProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProductVariants200JSONResponse{
		Code: 200, Data: &variants, Timestamp: time.Now(), Module: module,
	}, nil
}

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

func (h *Handlers) PatchProductVariant(ctx context.Context, req openapi.PatchProductVariantRequestObject) (openapi.PatchProductVariantResponseObject, error) {
	gallery, err := json.Marshal(req.Body.GalleryUrls)
	if err != nil {
		return nil, fun.ErrBadRequest("invalid gallery_urls")
	}
	variant, err := h.ops.PatchVariant(ctx, models.PatchProductVariantInput{
		VariantID:   req.VariantId,
		VendorCode:  req.Body.VendorCode,
		Name:        req.Body.Name,
		Description: req.Body.Description,
		Price:       req.Body.Price,
		Stock:       req.Body.Stock,
		GalleryURLs: gallery,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PatchProductVariant200JSONResponse{
		Code: 200, Data: variant, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) DeleteProduct(ctx context.Context, req openapi.DeleteProductRequestObject) (openapi.DeleteProductResponseObject, error) {
	err := h.ops.DeleteProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteProduct204Response{}, nil
}

func (h *Handlers) DeleteProductVariant(ctx context.Context, req openapi.DeleteProductVariantRequestObject) (openapi.DeleteProductVariantResponseObject, error) {
	err := h.ops.DeleteVariant(ctx, req.VariantId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteProductVariant204Response{}, nil
}
