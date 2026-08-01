// Package programs implements the StrictServerInterface methods for the
// programs feature.
package programs

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/internal/services"
	"univents/models"
)

const module = "Univents"

type Handlers struct {
	ops *services.Programs
}

func New(ops *services.Programs) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListEditionPrograms(ctx context.Context, req openapi.ListEditionProgramsRequestObject) (openapi.ListEditionProgramsResponseObject, error) {
	programs, err := h.ops.ListProgramsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionPrograms200JSONResponse{
		Code: 200, Data: &programs, Timestamp: time.Now(), Module: module,
	}, nil
}

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

func (h *Handlers) GetProgram(ctx context.Context, req openapi.GetProgramRequestObject) (openapi.GetProgramResponseObject, error) {
	program, err := h.ops.GetProgramByID(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProgram200JSONResponse{
		Code: 200, Data: program, Timestamp: time.Now(), Module: module,
	}, nil
}

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

func (h *Handlers) DeleteProgram(ctx context.Context, req openapi.DeleteProgramRequestObject) (openapi.DeleteProgramResponseObject, error) {
	program, err := h.ops.DeleteProgram(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteProgram200JSONResponse{
		Code: 200, Data: program, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListProgramOccurrences(ctx context.Context, req openapi.ListProgramOccurrencesRequestObject) (openapi.ListProgramOccurrencesResponseObject, error) {
	occurrences, err := h.ops.ListOccurrencesByProgram(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProgramOccurrences200JSONResponse{
		Code: 200, Data: &occurrences, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateProgramOccurrence(ctx context.Context, req openapi.CreateProgramOccurrenceRequestObject) (openapi.CreateProgramOccurrenceResponseObject, error) {
	occurrence, err := h.ops.CreateOccurrence(ctx, models.CreateProgramOccurrenceInput{
		ProgramID:   req.ProgramId,
		StartsAt:    req.Body.StartsAt,
		EndsAt:      req.Body.EndsAt,
		MaxCapacity: req.Body.MaxCapacity,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateProgramOccurrence201JSONResponse{
		Code: 201, Data: occurrence, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListEditionOccurrences(ctx context.Context, req openapi.ListEditionOccurrencesRequestObject) (openapi.ListEditionOccurrencesResponseObject, error) {
	occurrences, err := h.ops.ListOccurrencesByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionOccurrences200JSONResponse{
		Code: 200, Data: &occurrences, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetOccurrence(ctx context.Context, req openapi.GetOccurrenceRequestObject) (openapi.GetOccurrenceResponseObject, error) {
	occurrence, err := h.ops.GetOccurrenceByID(ctx, req.OccurrenceId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOccurrence200JSONResponse{
		Code: 200, Data: occurrence, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PatchOccurrence(ctx context.Context, req openapi.PatchOccurrenceRequestObject) (openapi.PatchOccurrenceResponseObject, error) {
	occurrence, err := h.ops.PatchOccurrence(ctx, models.PatchProgramOccurrenceInput{
		OccurrenceID: req.OccurrenceId,
		StartsAt:     req.Body.StartsAt,
		EndsAt:       req.Body.EndsAt,
		MaxCapacity:  req.Body.MaxCapacity,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PatchOccurrence200JSONResponse{
		Code: 200, Data: occurrence, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) DeleteOccurrence(ctx context.Context, req openapi.DeleteOccurrenceRequestObject) (openapi.DeleteOccurrenceResponseObject, error) {
	occurrence, err := h.ops.DeleteOccurrence(ctx, req.OccurrenceId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteOccurrence200JSONResponse{
		Code: 200, Data: occurrence, Timestamp: time.Now(), Module: module,
	}, nil
}
