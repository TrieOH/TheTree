package events

import (
	"encoding/json"
	"net/http"

	"univents/contracts"
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
		r.With(jwt).Post("/", h.CreateEvent)
		r.With(jwt).Get("/own", h.ListOwnEvents)
		r.With(jwt).Route("/{event_id}", func(r chi.Router) {
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

func (handler *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := handler.queries.ListEvents(ctx)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
}

func (handler *Handler) ListOwnEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := handler.queries.ListOwnEvents(ctx)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
}

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
