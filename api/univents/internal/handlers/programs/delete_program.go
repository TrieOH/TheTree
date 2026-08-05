package programs

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) DeleteProgram(ctx context.Context, req openapi.DeleteProgramRequestObject) (openapi.DeleteProgramResponseObject, error) {
	program, err := h.ops.DeleteProgram(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteProgram200JSONResponse{
		Code: 200, Data: program, Timestamp: time.Now(), Module: module,
	}, nil
}
