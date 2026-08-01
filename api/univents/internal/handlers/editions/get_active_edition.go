package editions

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetActiveEdition(ctx context.Context, req openapi.GetActiveEditionRequestObject) (openapi.GetActiveEditionResponseObject, error) {
	edition, err := h.ops.GetActive(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActiveEdition200JSONResponse{
		Code: 200, Data: edition, Timestamp: time.Now(), Module: module,
	}, nil
}
