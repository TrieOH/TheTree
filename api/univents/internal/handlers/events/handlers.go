// Package events implements the StrictServerInterface methods for the
// events feature.
package events

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/internal/services"
	"univents/models"
)

const module = "Univents"

type Handlers struct {
	ops *services.Events
}

func New(ops *services.Events) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListPublicEvents(ctx context.Context, _ openapi.ListPublicEventsRequestObject) (openapi.ListPublicEventsResponseObject, error) {
	events, err := h.ops.ListPublic(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListPublicEvents200JSONResponse{
		Code: 200, Data: &events, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetEventBySlug(ctx context.Context, req openapi.GetEventBySlugRequestObject) (openapi.GetEventBySlugResponseObject, error) {
	event, err := h.ops.GetBySlug(ctx, req.EventSlug)
	if err != nil {
		return nil, err
	}
	return openapi.GetEventBySlug200JSONResponse{
		Code: 200, Data: event, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateEvent(ctx context.Context, req openapi.CreateEventRequestObject) (openapi.CreateEventResponseObject, error) {
	event, err := h.ops.Create(ctx, models.CreateEventInput{
		FullName:     req.Body.FullName,
		Acronym:      req.Body.Acronym,
		Slug:         req.Body.Slug,
		Description:  req.Body.Description,
		ContactEmail: req.Body.ContactEmail,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateEvent201JSONResponse{
		Code: 201, Data: event, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListOwnedEvents(ctx context.Context, _ openapi.ListOwnedEventsRequestObject) (openapi.ListOwnedEventsResponseObject, error) {
	events, err := h.ops.ListOwned(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListOwnedEvents200JSONResponse{
		Code: 200, Data: &events, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListJoinedEvents(ctx context.Context, _ openapi.ListJoinedEventsRequestObject) (openapi.ListJoinedEventsResponseObject, error) {
	events, err := h.ops.ListJoined(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListJoinedEvents200JSONResponse{
		Code: 200, Data: &events, Timestamp: time.Now(), Module: module,
	}, nil
}

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

func (h *Handlers) PublishEvent(ctx context.Context, req openapi.PublishEventRequestObject) (openapi.PublishEventResponseObject, error) {
	err := h.ops.Publish(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.PublishEvent204Response{}, nil
}

func (h *Handlers) DiscontinueEvent(ctx context.Context, req openapi.DiscontinueEventRequestObject) (openapi.DiscontinueEventResponseObject, error) {
	err := h.ops.Discontinue(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.DiscontinueEvent204Response{}, nil
}

func (h *Handlers) ListEventMembers(ctx context.Context, req openapi.ListEventMembersRequestObject) (openapi.ListEventMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEventMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) AddEventMember(ctx context.Context, req openapi.AddEventMemberRequestObject) (openapi.AddEventMemberResponseObject, error) {
	member, err := h.ops.AddMember(ctx, req.EventId, models.AddEventMemberInput{
		Email: req.Body.Email,
		Role:  req.Body.Role,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddEventMember201JSONResponse{
		Code: 201, Data: member, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) RemoveEventMember(ctx context.Context, req openapi.RemoveEventMemberRequestObject) (openapi.RemoveEventMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, req.EventId, models.RemoveMemberInput{
		Email: req.Body.Email,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveEventMember204Response{}, nil
}
