package programs

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListProgramOccurrences(ctx context.Context, req openapi.ListProgramOccurrencesRequestObject) (openapi.ListProgramOccurrencesResponseObject, error) {
	occurrences, err := h.ops.ListOccurrencesByProgram(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProgramOccurrences200JSONResponse{
		Code: 200, Data: &occurrences, Timestamp: time.Now(), Module: module,
	}, nil
}
