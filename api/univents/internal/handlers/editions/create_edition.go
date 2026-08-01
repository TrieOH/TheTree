package editions

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) CreateEdition(ctx context.Context, req openapi.CreateEditionRequestObject) (openapi.CreateEditionResponseObject, error) {
	edition, err := h.ops.Create(ctx, models.CreateEditionInput{
		EventID:  req.EventId,
		Name:     req.Body.Name,
		Slug:     req.Body.Slug,
		StartsAt: req.Body.StartsAt,
		EndsAt:   req.Body.EndsAt,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateEdition201JSONResponse{
		Code: 201, Data: edition, Timestamp: time.Now(), Module: module,
	}, nil
}
