package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionPrograms(ctx context.Context, req openapi.ListEditionProgramsRequestObject) (openapi.ListEditionProgramsResponseObject, error) {
	programs, err := h.ops.ListProgramsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionPrograms200JSONResponse{Items: programs}, nil
}
