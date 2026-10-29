package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionSignatures(ctx context.Context, req openapi.ListEditionSignaturesRequestObject) (openapi.ListEditionSignaturesResponseObject, error) {
	signatures, err := h.ops.ListByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionSignatures200JSONResponse{Items: signatures}, nil
}
