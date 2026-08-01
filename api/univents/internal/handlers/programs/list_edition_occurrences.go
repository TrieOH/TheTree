package programs

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionOccurrences(ctx context.Context, req openapi.ListEditionOccurrencesRequestObject) (openapi.ListEditionOccurrencesResponseObject, error) {
	occurrences, err := h.ops.ListOccurrencesByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionOccurrences200JSONResponse{
		Code: 200, Data: &occurrences, Timestamp: time.Now(), Module: module,
	}, nil
}
