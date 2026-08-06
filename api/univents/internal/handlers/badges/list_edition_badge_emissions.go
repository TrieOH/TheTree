package badges

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionBadgeEmissions(ctx context.Context, req openapi.ListEditionBadgeEmissionsRequestObject) (openapi.ListEditionBadgeEmissionsResponseObject, error) {
	items, err := h.ops.ListByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionBadgeEmissions200JSONResponse{
		Code: 200, Data: &items, Timestamp: time.Now(), Module: module,
	}, nil
}
