package signatures

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionSignatureRequests(ctx context.Context, req openapi.ListEditionSignatureRequestsRequestObject) (openapi.ListEditionSignatureRequestsResponseObject, error) {
	requests, err := h.ops.ListRequestsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionSignatureRequests200JSONResponse{
		Code: 200, Data: &requests, Timestamp: time.Now(), Module: module,
	}, nil
}
