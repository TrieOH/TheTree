package programs

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionPrograms(ctx context.Context, req openapi.ListEditionProgramsRequestObject) (openapi.ListEditionProgramsResponseObject, error) {
	programs, err := h.ops.ListProgramsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionPrograms200JSONResponse{
		Code: 200, Data: &programs, Timestamp: time.Now(), Module: module,
	}, nil
}
