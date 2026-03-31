package http

import (
	"net/http"
	"univents/internal/commerce/interfaces/http/dtos"
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

// AddGalleryImage godoc
// @Summary Add an image to the event gallery
// @Description Adds a MinIO URL to the event's gallery_urls array.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param request body dtos.ImageURLRequest true "Image URL"
// @Success 200 {object} domain.Event "Image added to gallery"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 403 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/gallery [post]
func (handler *EventsHandler) AddGalleryImage(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dtos.ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.AddGalleryImage(ctx, eventID, req.URL)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Image added to gallery").WithData(product).Send(w)
}

// RemoveGalleryImage godoc
// @Summary Remove an image from the event gallery
// @Description Removes a URL from the event's gallery_urls array and deletes the object from MinIO.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param request body dtos.ImageURLRequest true "Image URL"
// @Success 200 {object} domain.Event "Image removed from gallery"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 403 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/gallery [delete]
func (handler *EventsHandler) RemoveGalleryImage(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dtos.ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.RemoveGalleryImage(ctx, eventID, req.URL)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Image removed from gallery").WithData(product).Send(w)
}

// SetLogo godoc
// @Summary Set the event logo
// @Description Sets the event logo URL.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param request body dtos.ImageURLRequest true "Image URL"
// @Success 200 {object} domain.Event "Logo set"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 403 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/logo [put]
func (handler *EventsHandler) SetLogo(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dtos.ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.SetLogo(ctx, eventID, req.URL)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logo set").WithData(product).Send(w)
}

// UnsetLogo godoc
// @Summary Unset the event logo
// @Description Clears the event logo.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Success 200 {object} domain.Event "Logo unset"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 403 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/logo [delete]
func (handler *EventsHandler) UnsetLogo(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.UnsetLogo(ctx, eventID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logo unset").WithData(product).Send(w)
}

// SetBanner godoc
// @Summary Set the event banner
// @Description Sets the event banner URL.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param request body dtos.ImageURLRequest true "Image URL"
// @Success 200 {object} domain.Event "Banner set"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 403 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/banner [put]
func (handler *EventsHandler) SetBanner(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dtos.ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.SetBanner(ctx, eventID, req.URL)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Banner set").WithData(product).Send(w)
}

// UnsetBanner godoc
// @Summary Unset the event banner
// @Description Clears the event banner.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Success 200 {object} domain.Event "Banner unset"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 403 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/banner [delete]
func (handler *EventsHandler) UnsetBanner(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.UnsetBanner(ctx, eventID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logo unset").WithData(product).Send(w)
}
