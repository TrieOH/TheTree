package programs

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetProgram(ctx context.Context, req openapi.GetProgramRequestObject) (openapi.GetProgramResponseObject, error) {
	program, err := h.ops.GetProgramByID(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProgram200JSONResponse{
		Code: 200, Data: program, Timestamp: time.Now(), Module: module,
	}, nil
}
