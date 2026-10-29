package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionSignatureRequests(ctx context.Context, req openapi.ListEditionSignatureRequestsRequestObject) (openapi.ListEditionSignatureRequestsResponseObject, error) {
	requests, err := h.ops.ListRequestsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionSignatureRequests200JSONResponse{Items: requests}, nil
}
