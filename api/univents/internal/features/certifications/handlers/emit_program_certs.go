package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) EmitProgramCertifications(ctx context.Context, req openapi.EmitProgramCertificationsRequestObject) (openapi.EmitProgramCertificationsResponseObject, error) {
	err := h.ops.EmitCertsForProgram(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.EmitProgramCertifications200JSONResponse(openapi.ProgramCertEmissionResult{ProgramId: req.ProgramId, Queued: true}), nil
}
