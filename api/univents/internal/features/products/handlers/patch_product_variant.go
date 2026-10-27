package handlers

import (
	"encoding/json"

	"context"
	"time"

	"univents/errs"
	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) PatchProductVariant(ctx context.Context, req openapi.PatchProductVariantRequestObject) (openapi.PatchProductVariantResponseObject, error) {
	gallery, err := json.Marshal(req.Body.GalleryUrls)
	if err != nil {
		return nil, errs.ErrInvalidGalleryURLs
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
