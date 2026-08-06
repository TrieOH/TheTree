package badges

import (
	"context"
	"time"

	"univents/internal/openapi"

	"github.com/google/uuid"
)

func (h *Handlers) GetEditionBadgesPrint(ctx context.Context, req openapi.GetEditionBadgesPrintRequestObject) (openapi.GetEditionBadgesPrintResponseObject, error) {
	var emissionIDs []uuid.UUID
	if req.Params.EmissionIds != nil {
		emissionIDs = *req.Params.EmissionIds
	}

	items, err := h.ops.PrintByEdition(ctx, req.EditionId, emissionIDs)
	if err != nil {
		return nil, err
	}
	return openapi.GetEditionBadgesPrint200JSONResponse{
		Code: 200, Data: &items, Timestamp: time.Now(), Module: module,
	}, nil
}
