package editions

import (
	"net/http"
	"univents/internal/core/application/edition/commands"
	"univents/internal/core/application/edition/queries"
	"univents/internal/core/domain"
	"univents/internal/core/interfaces/http/dto"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

type Handler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewEditionsHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
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
// @Param request body dto.CreateEditionRequest true "Edition creation request"
// @Success 201 {object} object "Edition created successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions [post]
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.CreateEditionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := domain.CreateEditionSpec{
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
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
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
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
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
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/announce [post]
func (handler *Handler) Announce(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
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
	err := handler.commands.Announce(ctx, eventID, editionID)
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
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
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

	ctx := r.Context()
	err = handler.commands.ConnectPayments(ctx, credID, editionID, provider)
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
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
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
