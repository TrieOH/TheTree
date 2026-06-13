package events

import (
	"encoding/json"
	"net/http"

	"univents/internal/shared/contracts"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewHandler(
	commands *CommandService,
	queries *QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func Routes(
	r *chi.Mux,
	h *Handler,
	jwt func(http.Handler) http.Handler,
) {
	r.Route("/events", func(r chi.Router) {
		r.Get("/", h.ListEvents)
		r.Use(jwt)
		r.Post("/", h.CreateEvent)
		r.Get("/own", h.ListOwnEvents)
		r.Route("/{event_id}", func(r chi.Router) {
			r.Patch("/", h.PatchEvent)
			r.Post("/publish", h.PublishEvent)
			r.Post("/gallery", h.AddGalleryImage)
			r.Delete("/gallery", h.RemoveGalleryImage)
			r.Put("/logo", h.SetLogo)
			r.Delete("/logo", h.UnsetLogo)
			r.Put("/banner", h.SetBanner)
			r.Delete("/banner", h.UnsetBanner)
		})
	})
}

type CreateEventRequest struct {
	OrganizationID *uuid.UUID `json:"organization_id"`
	Name           string     `json:"name" validate:"required,min=2"`
	Acronym        *string    `json:"acronym"`
	Slug           string     `json:"slug" validate:"required,min=2"`
	Tagline        *string    `json:"tagline"`
	Description    *string    `json:"description"`
	IsSeries       bool       `json:"is_series"`
	LogoUrl        *string    `json:"logo_url"`
	BannerUrl      *string    `json:"banner_url"`
	ContactEmail   *string    `json:"contact_email" validate:"required,email"`
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Creates a new event with the provided details.
// @Tags events
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body CreateEventRequest true "Event creation request"
// @Success 201 {object} object "Event created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events [post]
func (handler *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	in := contracts.CreateEventSpec{
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
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
}

type PatchEventRequest struct {
	Name         string           `json:"name" validate:"required,min=3,max=256"`
	Acronym      *string          `json:"acronym" validate:"omitempty,min=2,max=32"`
	Slug         string           `json:"slug" validate:"required,min=2,max=32"`
	Tagline      *string          `json:"tagline" validate:"omitempty,max=512"`
	Description  *string          `json:"description"`
	IsSeries     bool             `json:"is_series"`
	LogoUrl      *string          `json:"logo_url" validate:"omitempty,url"`
	BannerUrl    *string          `json:"banner_url" validate:"omitempty,url"`
	HasGallery   bool             `json:"has_gallery"`
	ContactEmail *string          `json:"contact_email" validate:"omitempty,email"`
	SocialLinks  *json.RawMessage `json:"social_links" validate:"omitempty,json"`
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
// @Param request body PatchEventRequest true "Event patch request"
// @Success 201 {object} object "Event patched successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id} [patch]
func (handler *Handler) PatchEvent(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req PatchEventRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	in := contracts.PatchEventSpec{
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
	out, _, err := handler.commands.PatchEvent(ctx, in)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
}

// ListEvents godoc
// @Summary List all public events
// @Description Retrieves a list of all public events available.
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {object} object "Events retrieved successfully"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events [get]
func (handler *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := handler.queries.ListEvents(ctx)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
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
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/own [get]
func (handler *Handler) ListOwnEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := handler.queries.ListOwnEvents(ctx)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
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
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/publish [post]
func (handler *Handler) PublishEvent(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.PublishEvent(ctx, eventID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().Send(w)
}

type ImageURLRequest struct {
	URL string `json:"url" validate:"required,url"`
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
// @Param request body ImageURLRequest true "Image URL"
// @Success 200 {object} contracts.Event "Image added to gallery"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/gallery [post]
func (handler *Handler) AddGalleryImage(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.AddGalleryImage(ctx, eventID, req.URL)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Image added to gallery").WithData(product).Send(w)
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
// @Param request body ImageURLRequest true "Image URL"
// @Success 200 {object} contracts.Event "Image removed from gallery"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/gallery [delete]
func (handler *Handler) RemoveGalleryImage(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.RemoveGalleryImage(ctx, eventID, req.URL)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Image removed from gallery").WithData(product).Send(w)
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
// @Param request body ImageURLRequest true "Image URL"
// @Success 200 {object} contracts.Event "Logo set"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/logo [put]
func (handler *Handler) SetLogo(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.SetLogo(ctx, eventID, req.URL)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Logo set").WithData(product).Send(w)
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
// @Success 200 {object} contracts.Event "Logo unset"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/logo [delete]
func (handler *Handler) UnsetLogo(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.UnsetLogo(ctx, eventID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Logo unset").WithData(product).Send(w)
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
// @Param request body ImageURLRequest true "Image URL"
// @Success 200 {object} contracts.Event "Banner set"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/banner [put]
func (handler *Handler) SetBanner(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.SetBanner(ctx, eventID, req.URL)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Banner set").WithData(product).Send(w)
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
// @Success 200 {object} contracts.Event "Banner unset"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/banner [delete]
func (handler *Handler) UnsetBanner(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.UnsetBanner(ctx, eventID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Logo unset").WithData(product).Send(w)
}
