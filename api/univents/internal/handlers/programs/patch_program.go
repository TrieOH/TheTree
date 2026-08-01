package programs

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) PatchProgram(ctx context.Context, req openapi.PatchProgramRequestObject) (openapi.PatchProgramResponseObject, error) {
	staffOnly := false
	if req.Body.StaffOnly != nil {
		staffOnly = *req.Body.StaffOnly
	}
	program, err := h.ops.PatchProgram(ctx, models.PatchProgramInput{
		ProgramID:      req.ProgramId,
		Kind:           req.Body.Kind,
		Name:           req.Body.Name,
		Description:    req.Body.Description,
		MinAccessLevel: req.Body.MinAccessLevel,
		StaffOnly:      staffOnly,
		Price:          req.Body.Price,
		BannerURL:      req.Body.BannerUrl,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PatchProgram200JSONResponse{
		Code: 200, Data: program, Timestamp: time.Now(), Module: module,
	}, nil
}
