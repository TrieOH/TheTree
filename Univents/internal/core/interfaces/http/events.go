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

// CreateEvent godoc
// @Summary Create a new event
// @Description Creates a new event with the provided details.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body dto.CreateEventRequest true "Event creation request"
// @Success 201 {object} object "Event created successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events [post]
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

// PatchEvent godoc
// @Summary Patch an event
// @Description Updates an existing event with the provided details.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param request body dto.PatchEventRequest true "Event patch request"
// @Success 201 {object} object "Event patched successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id} [patch]
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

// ListEvents godoc
// @Summary List all public events
// @Description Retrieves a list of all public events available.
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {object} object "Events retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events [get]
func (handler *EventsHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := handler.queries.ListEvents(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
}

// ListOwnEvents godoc
// @Summary List events owned by the authenticated user
// @Description Retrieves a list of events owned by the authenticated user.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {object} object "Events retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/own [get]
func (handler *EventsHandler) ListOwnEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := handler.queries.ListOwnEvents(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
}

// PublishEvent godoc
// @Summary Publish an event
// @Description Publishes an event making it publicly available.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Success 200 {object} object "Event published successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/publish [post]
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

// ListEventAudits godoc
// @Summary List audit logs for an event
// @Description Retrieves the audit logs for a specific event.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Success 200 {object} object "Audit logs retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/audit [get]
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
