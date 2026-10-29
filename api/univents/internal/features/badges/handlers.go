package badges

import (
	"context"
	"encoding/json"

	"univents/errs"
	"univents/internal/features/badges/service"
	"univents/internal/openapi"
	"univents/models"

	"github.com/google/uuid"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) CreateBadgeTemplate(ctx context.Context, req openapi.CreateBadgeTemplateRequestObject) (openapi.CreateBadgeTemplateResponseObject, error) {
	design, err := json.Marshal(req.Body.DesignData)
	if err != nil {
		return nil, errs.ErrInvalidDesignData
	}
	template, err := h.ops.CreateTemplate(ctx, models.CreateBadgeTemplateInput{
		EditionID:    req.EditionId,
		TicketTypeID: req.Body.TicketTypeId,
		Origin:       req.Body.Origin,
		Name:         req.Body.Name,
		DesignData:   design,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateBadgeTemplate201JSONResponse(*template), nil
}

func (h *Handlers) DeleteBadgeTemplate(ctx context.Context, req openapi.DeleteBadgeTemplateRequestObject) (openapi.DeleteBadgeTemplateResponseObject, error) {
	err := h.ops.DeleteTemplate(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteBadgeTemplate204Response{}, nil
}

func (h *Handlers) GetBadgeTemplate(ctx context.Context, req openapi.GetBadgeTemplateRequestObject) (openapi.GetBadgeTemplateResponseObject, error) {
	template, err := h.ops.GetTemplate(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.GetBadgeTemplate200JSONResponse(*template), nil
}

func (h *Handlers) GetEditionBadgesPrint(ctx context.Context, req openapi.GetEditionBadgesPrintRequestObject) (openapi.GetEditionBadgesPrintResponseObject, error) {
	var emissionIDs []uuid.UUID
	if req.Params.EmissionIds != nil {
		emissionIDs = *req.Params.EmissionIds
	}

	items, err := h.ops.PrintByEdition(ctx, req.EditionId, emissionIDs)
	if err != nil {
		return nil, err
	}
	return openapi.GetEditionBadgesPrint200JSONResponse{Items: items}, nil
}

func (h *Handlers) ListBadgeTemplates(ctx context.Context, req openapi.ListBadgeTemplatesRequestObject) (openapi.ListBadgeTemplatesResponseObject, error) {
	templates, err := h.ops.ListTemplates(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListBadgeTemplates200JSONResponse{Items: templates}, nil
}

func (h *Handlers) ListEditionBadgeEmissions(ctx context.Context, req openapi.ListEditionBadgeEmissionsRequestObject) (openapi.ListEditionBadgeEmissionsResponseObject, error) {
	items, err := h.ops.ListByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionBadgeEmissions200JSONResponse{Items: items}, nil
}

func (h *Handlers) ListUserBadges(ctx context.Context, req openapi.ListUserBadgesRequestObject) (openapi.ListUserBadgesResponseObject, error) {
	groups, err := h.ops.ListByUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return openapi.ListUserBadges200JSONResponse(*groups), nil
}

func (h *Handlers) UpdateBadgeTemplate(ctx context.Context, req openapi.UpdateBadgeTemplateRequestObject) (openapi.UpdateBadgeTemplateResponseObject, error) {
	input := models.UpdateBadgeTemplateInput{TemplateID: req.TemplateId}

	if req.Body.Name != nil {
		name := *req.Body.Name
		input.Name = &name
	}
	if req.Body.DesignData != nil {
		design, err := json.Marshal(req.Body.DesignData)
		if err != nil {
			return nil, err
		}
		input.DesignData = design
	}

	template, err := h.ops.UpdateTemplate(ctx, input)
	if err != nil {
		return nil, err
	}
	return openapi.UpdateBadgeTemplate200JSONResponse(*template), nil
}
