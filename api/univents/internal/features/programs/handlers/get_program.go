package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) GetProgram(ctx context.Context, req openapi.GetProgramRequestObject) (openapi.GetProgramResponseObject, error) {
	program, err := h.ops.GetProgramByID(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProgram200JSONResponse(*program), nil
}
