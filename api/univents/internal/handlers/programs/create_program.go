package programs

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) CreateProgram(ctx context.Context, req openapi.CreateProgramRequestObject) (openapi.CreateProgramResponseObject, error) {
	staffOnly := false
	if req.Body.StaffOnly != nil {
		staffOnly = *req.Body.StaffOnly
	}
	program, err := h.ops.CreateProgram(ctx, models.CreateProgramInput{
		EditionID:      req.EditionId,
		Kind:           req.Body.Kind,
		Name:           req.Body.Name,
		Description:    req.Body.Description,
		MinAccessLevel: req.Body.MinAccessLevel,
		StaffOnly:      staffOnly,
		Price:          req.Body.Price,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateProgram201JSONResponse{
		Code: 201, Data: program, Timestamp: time.Now(), Module: module,
	}, nil
}
