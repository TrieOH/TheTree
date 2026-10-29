package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) FulfillSignatureRequest(ctx context.Context, req openapi.FulfillSignatureRequestRequestObject) (openapi.FulfillSignatureRequestResponseObject, error) {
	signature, err := h.ops.FulfillRequest(ctx, req.Params.Token, req.Body.ImageUrl)
	if err != nil {
		return nil, err
	}
	return openapi.FulfillSignatureRequest201JSONResponse(*signature), nil
}
