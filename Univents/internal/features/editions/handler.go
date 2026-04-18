package editions

import (
	"net/http"
	"time"
	"univents/internal/shared/contracts"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
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

// Create godoc
// @Summary Create a new edition
// @Description Creates a new edition for an event.
// @Tags editions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param request body CreateEditionRequest true "Edition creation request"
// @Success 201 {object} object "Edition created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions [post]
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req CreateEditionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
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
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}

// List godoc
// @Summary List editions for an event
// @Description Retrieves a list of editions for a specific event.
// @Tags editions
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} object "Editions retrieved successfully"
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions [get]
func (handler *Handler) List(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.ListEditions(ctx, eventID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}

// ListAdmin godoc
// @Summary List editions for an event including draft ones
// @Description Retrieves a list of editions for a specific event and its draft editions if you have editions:read.
// @Tags editions
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} object "Editions retrieved successfully"
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/admin [get]
func (handler *Handler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.ListEditionsAdmin(ctx, eventID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}

// Announce godoc
// @Summary Announce an edition
// @Description Announces an edition making it publicly available.
// @Tags editions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Success 200 {object} object "Edition announced successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/announce [post]
func (handler *Handler) Announce(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.Announce(ctx, editionID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}

// ConnectPaymentAccount godoc
// @Summary Connect a payment account to an edition
// @Description Connects a payment provider credential to an edition.
// @Tags editions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param credential_id query string true "Trie Payments Credential ID (UUID)"
// @Param provider query string true "Payment provider name"
// @Success 200 {object} object "Payment account connected successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/payments/connect [post]
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
		resp.BadRequest("missing credential_id").Send(w)
		return
	}

	credID, err := uuid.Parse(triePaymentsCredentialID)
	if err != nil {
		resp.BadRequest("invalid credential_id: " + err.Error()).Send(w)
		return
	}

	provider := r.URL.Query().Get("provider")
	if provider == "" {
		resp.BadRequest("missing provider").Send(w)
		return
	}

	publicKey := r.URL.Query().Get("public_key")
	if publicKey == "" {
		resp.BadRequest("missing public_key").Send(w)
		return
	}

	ctx := r.Context()
	err = handler.commands.ConnectPayments(ctx, credID, editionID, provider, publicKey)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Payment account connected successfully").Send(w)
}

// DisconnectPaymentAccount godoc
// @Summary Disconnect payment account from an edition
// @Description Removes the connected payment provider credential from an edition.
// @Tags editions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Success 200 {object} object "Payment account disconnected successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/payments/disconnect [post]
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
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}
