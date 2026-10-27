package handlers

import (
	"encoding/json"

	"context"
	"time"

	"univents/errs"
	"univents/internal/openapi"
	"univents/models"
)

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
	return openapi.CreateBadgeTemplate201JSONResponse{
		Code: 201, Data: template, Timestamp: time.Now(), Module: module,
	}, nil
}
