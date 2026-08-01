// Package badges implements the StrictServerInterface methods for the
// badges feature.
package badges

import (
	"context"
	"encoding/json"
	"time"

	"univents/internal/openapi"
	"univents/internal/services"
	"univents/models"

	"github.com/MintzyG/fun"
)

const module = "Univents"

type Handlers struct {
	ops *services.Badges
}

func New(ops *services.Badges) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListBadgeTemplates(ctx context.Context, req openapi.ListBadgeTemplatesRequestObject) (openapi.ListBadgeTemplatesResponseObject, error) {
	templates, err := h.ops.ListTemplates(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListBadgeTemplates200JSONResponse{
		Code: 200, Data: &templates, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateBadgeTemplate(ctx context.Context, req openapi.CreateBadgeTemplateRequestObject) (openapi.CreateBadgeTemplateResponseObject, error) {
	design, err := json.Marshal(req.Body.DesignData)
	if err != nil {
		return nil, fun.ErrBadRequest("invalid design_data")
	}
	template, err := h.ops.CreateTemplate(ctx, models.CreateBadgeTemplateInput{
		EditionID:    req.EditionId,
		TicketTypeID: req.Body.TicketTypeId,
		Name:         req.Body.Name,
		DesignData:   design,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateBadgeTemplate201JSONResponse{
		Code: 201, Data: template, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetBadgeTemplate(ctx context.Context, req openapi.GetBadgeTemplateRequestObject) (openapi.GetBadgeTemplateResponseObject, error) {
	template, err := h.ops.GetTemplate(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.GetBadgeTemplate200JSONResponse{
		Code: 200, Data: template, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) DeleteBadgeTemplate(ctx context.Context, req openapi.DeleteBadgeTemplateRequestObject) (openapi.DeleteBadgeTemplateResponseObject, error) {
	err := h.ops.DeleteTemplate(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteBadgeTemplate204Response{}, nil
}
