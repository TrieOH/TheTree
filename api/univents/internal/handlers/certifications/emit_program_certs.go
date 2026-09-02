package certifications

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) EmitProgramCertifications(ctx context.Context, req openapi.EmitProgramCertificationsRequestObject) (openapi.EmitProgramCertificationsResponseObject, error) {
	err := h.ops.EmitCertsForProgram(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.EmitProgramCertifications200JSONResponse{
		Code:      200,
		Data:      &openapi.ProgramCertEmissionResult{ProgramId: req.ProgramId, Queued: true},
		Timestamp: time.Now(),
		Module:    module,
	}, nil
}
