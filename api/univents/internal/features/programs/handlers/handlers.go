package handlers

import (
	"context"

	idx "sdk/identityx"
	"univents/internal/features/programs/service"
	"univents/internal/openapi"
	"univents/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) CheckInOccurrence(ctx context.Context, req openapi.CheckInOccurrenceRequestObject) (openapi.CheckInOccurrenceResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	part, err := h.ops.CheckIn(ctx, req.OccurrenceId, req.Body.AttendeeId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.CheckInOccurrence200JSONResponse(*part), nil
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
	return openapi.CreateProgram201JSONResponse(*program), nil
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
	return openapi.CreateProgramOccurrence201JSONResponse(*occurrence), nil
}

func (h *Handlers) DeleteOccurrence(ctx context.Context, req openapi.DeleteOccurrenceRequestObject) (openapi.DeleteOccurrenceResponseObject, error) {
	occurrence, err := h.ops.DeleteOccurrence(ctx, req.OccurrenceId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteOccurrence200JSONResponse(*occurrence), nil
}

func (h *Handlers) DeleteProgram(ctx context.Context, req openapi.DeleteProgramRequestObject) (openapi.DeleteProgramResponseObject, error) {
	program, err := h.ops.DeleteProgram(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteProgram200JSONResponse(*program), nil
}

func (h *Handlers) DeregisterOccurrence(ctx context.Context, req openapi.DeregisterOccurrenceRequestObject) (openapi.DeregisterOccurrenceResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	part, err := h.ops.Deregister(ctx, req.OccurrenceId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.DeregisterOccurrence200JSONResponse(*part), nil
}

func (h *Handlers) GetOccurrence(ctx context.Context, req openapi.GetOccurrenceRequestObject) (openapi.GetOccurrenceResponseObject, error) {
	occurrence, err := h.ops.GetOccurrenceByID(ctx, req.OccurrenceId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOccurrence200JSONResponse(*occurrence), nil
}

func (h *Handlers) GetProgram(ctx context.Context, req openapi.GetProgramRequestObject) (openapi.GetProgramResponseObject, error) {
	program, err := h.ops.GetProgramByID(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProgram200JSONResponse(*program), nil
}

func (h *Handlers) ListEditionOccurrences(ctx context.Context, req openapi.ListEditionOccurrencesRequestObject) (openapi.ListEditionOccurrencesResponseObject, error) {
	occurrences, err := h.ops.ListOccurrencesByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionOccurrences200JSONResponse{Items: occurrences}, nil
}

func (h *Handlers) ListEditionPrograms(ctx context.Context, req openapi.ListEditionProgramsRequestObject) (openapi.ListEditionProgramsResponseObject, error) {
	programs, err := h.ops.ListProgramsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionPrograms200JSONResponse{Items: programs}, nil
}

func (h *Handlers) ListMyParticipations(ctx context.Context, req openapi.ListMyParticipationsRequestObject) (openapi.ListMyParticipationsResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	participations, err := h.ops.MyParticipations(ctx, req.EditionId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.ListMyParticipations200JSONResponse{Items: participations}, nil
}

func (h *Handlers) ListOccurrenceParticipants(ctx context.Context, req openapi.ListOccurrenceParticipantsRequestObject) (openapi.ListOccurrenceParticipantsResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	participants, err := h.ops.Participants(ctx, req.OccurrenceId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.ListOccurrenceParticipants200JSONResponse{Items: participants}, nil
}

func (h *Handlers) ListProgramOccurrences(ctx context.Context, req openapi.ListProgramOccurrencesRequestObject) (openapi.ListProgramOccurrencesResponseObject, error) {
	occurrences, err := h.ops.ListOccurrencesByProgram(ctx, req.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProgramOccurrences200JSONResponse{Items: occurrences}, nil
}

func (h *Handlers) MarkParticipationAttended(ctx context.Context, req openapi.MarkParticipationAttendedRequestObject) (openapi.MarkParticipationAttendedResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	part, err := h.ops.MarkAttended(ctx, req.ParticipationId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.MarkParticipationAttended200JSONResponse(*part), nil
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
	return openapi.PatchOccurrence200JSONResponse(*occurrence), nil
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
	return openapi.PatchProgram200JSONResponse(*program), nil
}

func (h *Handlers) RegisterOccurrence(ctx context.Context, req openapi.RegisterOccurrenceRequestObject) (openapi.RegisterOccurrenceResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	part, err := h.ops.Register(ctx, req.OccurrenceId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.RegisterOccurrence201JSONResponse(*part), nil
}
