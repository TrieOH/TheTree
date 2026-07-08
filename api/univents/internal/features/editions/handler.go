package editions

import (
	"net/http"
	"time"

	"univents/contracts"
	"univents/internal/shared/validation"

	fun "github.com/MintzyG/fun"
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
	r.Route("/events/{event_id}/editions", func(r chi.Router) {
		r.Get("/", h.List)
		r.With(jwt).Post("/", h.Create)
		r.With(jwt).Get("/admin", h.ListAdmin)
		r.With(jwt).Route("/{edition_id}", func(r chi.Router) {
			r.Post("/announce", h.Announce)
			r.Post("/payments/connect", h.ConnectPaymentAccount)
			r.Post("/payments/disconnect", h.DisconnectPaymentAccount)
		})
	})
}

type CreateEditionRequest struct {
	Type                 contracts.EditionType `json:"type"`
	EditionName          string                `json:"edition_name" validate:"required,min=3,max=256"`
	Tagline              *string               `json:"tagline" validate:"omitempty,max=512"`
	Description          *string               `json:"description" validate:"omitempty,max=8000"`
	RegistrationOpensAt  *time.Time            `json:"registration_opens_at"`
	RegistrationClosesAt *time.Time            `json:"registration_closes_at"`
	StartsAt             time.Time             `json:"starts_at"`
	EndsAt               time.Time             `json:"ends_at"`
	Timezone             string                `json:"timezone"`
	LocationName         string                `json:"location_name"`
	LocationAddress      string                `json:"location_address"`
	LogoUrl              *string               `json:"logo_url" validate:"omitempty,url"`
	BannerUrl            *string               `json:"banner_url" validate:"omitempty,url"`
	ContactEmail         *string               `json:"contact_email" validate:"omitempty,email"`
	ContactPhone         *string               `json:"contact_phone"`
	OrganizerName        *string               `json:"organizer_name"`
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req CreateEditionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	in := contracts.CreateEditionSpec{
		EventID:              eventID,
		Type:                 req.Type,
		EditionName:          req.EditionName,
		Tagline:              req.Tagline,
		Description:          req.Description,
		RegistrationOpensAt:  req.RegistrationOpensAt,
		RegistrationClosesAt: req.RegistrationClosesAt,
		StartsAt:             req.StartsAt,
		EndsAt:               req.EndsAt,
		Timezone:             req.Timezone,
		LocationName:         req.LocationName,
		LocationAddress:      req.LocationAddress,
		LogoUrl:              req.LogoUrl,
		BannerUrl:            req.BannerUrl,
		ContactEmail:         req.ContactEmail,
		ContactPhone:         req.ContactPhone,
		OrganizerName:        req.OrganizerName,
	}

	ctx := r.Context()
	out, err := handler.commands.Create(ctx, in)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
}

func (handler *Handler) List(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.ListEditions(ctx, eventID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
}

func (handler *Handler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.ListEditionsAdmin(ctx, eventID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
}

func (handler *Handler) Announce(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.Announce(ctx, editionID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().Send(w)
}

func (handler *Handler) ConnectPaymentAccount(w http.ResponseWriter, r *http.Request) {
	_, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	triePaymentsCredentialID := r.URL.Query().Get("credential_id")
	if triePaymentsCredentialID == "" {
		fun.BadRequest("missing credential_id").Send(w)
		return
	}

	credID, err := uuid.Parse(triePaymentsCredentialID)
	if err != nil {
		fun.BadRequest("invalid credential_id: " + err.Error()).Send(w)
		return
	}

	provider := r.URL.Query().Get("provider")
	if provider == "" {
		fun.BadRequest("missing provider").Send(w)
		return
	}

	publicKey := r.URL.Query().Get("public_key")
	if publicKey == "" {
		fun.BadRequest("missing public_key").Send(w)
		return
	}

	ctx := r.Context()
	err = handler.commands.ConnectPayments(ctx, credID, editionID, provider, publicKey)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Payment account connected successfully").Send(w)
}

func (handler *Handler) DisconnectPaymentAccount(w http.ResponseWriter, r *http.Request) {
	_, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.DisconnectPayments(ctx, editionID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().Send(w)
}
