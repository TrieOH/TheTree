package events

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) PatchEvent(ctx context.Context, req openapi.PatchEventRequestObject) (openapi.PatchEventResponseObject, error) {
	event, err := h.ops.Patch(ctx, models.PatchEventInput{
		EventID:      req.EventId,
		FullName:     req.Body.FullName,
		Acronym:      req.Body.Acronym,
		Slug:         req.Body.Slug,
		Description:  req.Body.Description,
		LogoURL:      req.Body.LogoUrl,
		BannerURL:    req.Body.BannerUrl,
		ContactEmail: req.Body.ContactEmail,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PatchEvent200JSONResponse{
		Code: 200, Data: event, Timestamp: time.Now(), Module: module,
	}, nil
}
