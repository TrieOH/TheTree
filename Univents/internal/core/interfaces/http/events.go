package http

import (
	"net/http"
	"univents/internal/core/application/event/commands"
	"univents/internal/core/application/event/queries"
	"univents/internal/core/domain"
	"univents/internal/core/interfaces/http/dto"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type EventsHandler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewEventsHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *EventsHandler {
	return &EventsHandler{
		commands: commands,
		queries:  queries,
	}
}

func (handler *EventsHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := domain.CreateEventSpec{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Acronym:        req.Acronym,
		Slug:           req.Slug,
		Tagline:        req.Tagline,
		Description:    req.Description,
		IsSeries:       req.IsSeries,
		LogoUrl:        req.LogoUrl,
		ContactEmail:   req.ContactEmail,
	}

	ctx := r.Context()
	out, err := handler.commands.CreateEvent(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}

func (handler *EventsHandler) PatchEvent(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.PatchEventRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := domain.PatchEventSpec{
		ID:           eventID,
		Name:         req.Name,
		Acronym:      req.Acronym,
		Slug:         req.Slug,
		Tagline:      req.Tagline,
		Description:  req.Description,
		IsSeries:     req.IsSeries,
		LogoUrl:      req.LogoUrl,
		BannerUrl:    req.BannerUrl,
		HasGallery:   req.HasGallery,
		ContactEmail: req.ContactEmail,
		SocialLinks:  req.SocialLinks,
	}

	ctx := r.Context()
	out, warns, err := handler.commands.PatchEvent(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).AddTrace(warns).Send(w)
}

func (handler *EventsHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := handler.queries.ListEvents(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
}

func (handler *EventsHandler) ListOwnEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := handler.queries.ListOwnEvents(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
}

func (handler *EventsHandler) PublishEvent(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.PublishEvent(ctx, eventID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}

func (handler *EventsHandler) ListEventAudits(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.ListEventAudits(ctx, eventID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
}
