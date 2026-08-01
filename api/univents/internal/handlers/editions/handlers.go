// Package editions implements the StrictServerInterface methods for the
// editions feature.
package editions

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/internal/services"
	"univents/models"
)

const module = "Univents"

type Handlers struct {
	ops *services.Editions
}

func New(ops *services.Editions) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) GetEditionBySlug(ctx context.Context, req openapi.GetEditionBySlugRequestObject) (openapi.GetEditionBySlugResponseObject, error) {
	edition, err := h.ops.GetByEventAndEditionSlug(ctx, req.EventSlug, req.EditionSlug)
	if err != nil {
		return nil, err
	}
	return openapi.GetEditionBySlug200JSONResponse{
		Code: 200, Data: edition, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListPublicEditions(ctx context.Context, req openapi.ListPublicEditionsRequestObject) (openapi.ListPublicEditionsResponseObject, error) {
	editions, err := h.ops.ListPublic(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListPublicEditions200JSONResponse{
		Code: 200, Data: &editions, Timestamp: time.Now(), Module: module,
	}, nil
}

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

func (h *Handlers) GetActiveEdition(ctx context.Context, req openapi.GetActiveEditionRequestObject) (openapi.GetActiveEditionResponseObject, error) {
	edition, err := h.ops.GetActive(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActiveEdition200JSONResponse{
		Code: 200, Data: edition, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListPastEditions(ctx context.Context, req openapi.ListPastEditionsRequestObject) (openapi.ListPastEditionsResponseObject, error) {
	editions, err := h.ops.GetPast(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListPastEditions200JSONResponse{
		Code: 200, Data: &editions, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListUpcomingEditions(ctx context.Context, req openapi.ListUpcomingEditionsRequestObject) (openapi.ListUpcomingEditionsResponseObject, error) {
	editions, err := h.ops.GetUpcoming(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListUpcomingEditions200JSONResponse{
		Code: 200, Data: &editions, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListDraftEditions(ctx context.Context, req openapi.ListDraftEditionsRequestObject) (openapi.ListDraftEditionsResponseObject, error) {
	editions, err := h.ops.ListDraft(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListDraftEditions200JSONResponse{
		Code: 200, Data: &editions, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PatchEdition(ctx context.Context, req openapi.PatchEditionRequestObject) (openapi.PatchEditionResponseObject, error) {
	edition, err := h.ops.Patch(ctx, models.PatchEditionInput{
		EditionID:           req.EditionId,
		Name:                req.Body.Name,
		Slug:                req.Body.Slug,
		Tagline:             req.Body.Tagline,
		Description:         req.Body.Description,
		RegistrationOpensAt: req.Body.RegistrationOpensAt,
		StartsAt:            req.Body.StartsAt,
		EndsAt:              req.Body.EndsAt,
		LocationName:        req.Body.LocationName,
		LocationAddress:     req.Body.LocationDescription,
		LogoURL:             req.Body.LogoUrl,
		BannerURL:           req.Body.BannerUrl,
		ContactEmail:        req.Body.ContactEmail,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PatchEdition200JSONResponse{
		Code: 200, Data: edition, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PublishEdition(ctx context.Context, req openapi.PublishEditionRequestObject) (openapi.PublishEditionResponseObject, error) {
	err := h.ops.Publish(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.PublishEdition204Response{}, nil
}
