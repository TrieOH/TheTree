package intents

import (
	"encoding/json"
	"net/http"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/errx"
	"payssage/internal/shared/validation"

	_ "payssage/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
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

// CancelIntent godoc
// @Summary Cancel a payment intent
// @Description Cancels a pending payment intent
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param intent_id path string true "Intent ID"
// @Success 200 {object} contracts.Intent "Intent canceled successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /intents/{intent_id}/cancel [post]
func (h *Handler) CancelIntent(w http.ResponseWriter, r *http.Request) {
	intentID, rs := validation.GetUUID(r, "intent_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	intent, err := h.commands.CancelIntent(r.Context(), intentID)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			resp.NotFound("intent not found").Send(w)
			return
		}
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(intent).Send(w)
}

type CancelPixRequest struct {
	Provider           string    `json:"provider"`
	SellerCredentialID uuid.UUID `json:"seller_credential_id"`
}

// CancelPix godoc
// @Summary Cancel a payment intent and its pix
// @Description Cancels a pending payment intent and its related pix QR code
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param intent_id path string true "Intent ID"
// @Success 200 {object} contracts.Intent "Pix canceled successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /intents/{intent_id}/cancel-pix [post]
func (h *Handler) CancelPix(w http.ResponseWriter, r *http.Request) {
	intentID, rs := validation.GetUUID(r, "intent_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req CancelPixRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := CancelPixInput{
		Provider:           req.Provider,
		IntentID:           intentID,
		SellerCredentialID: req.SellerCredentialID,
	}

	intent, err := h.commands.CancelPix(r.Context(), in)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			resp.NotFound("intent not found").Send(w)
			return
		}
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Pix canceled successfully").WithData(intent).Send(w)
}

// GetByID godoc
// @Summary Get a payment intent
// @Description Retrieves a payment intent by ID
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param intent_id path string true "Intent ID"
// @Success 200 {object} contracts.Intent "Intent retrieved successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /intents/{intent_id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	intentID, rs := validation.GetUUID(r, "intent_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	intent, err := h.queries.GetByID(r.Context(), intentID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(intent).Send(w)
}

type CreateIntentRequest struct {
	Amount               int64           `json:"amount"`
	Currency             string          `json:"currency"`
	Provider             string          `json:"provider"`
	Metadata             json.RawMessage `json:"metadata"`
	PaymentMethodID      string          `json:"payment_method_id"`
	Installments         int             `json:"installments"`
	CardToken            string          `json:"card_token"`
	PaymentMethodType    string          `json:"payment_method_type"`
	SellerCredentialID   uuid.UUID       `json:"seller_credential_id"`
	PayerEmail           string          `json:"payer_email"`
	IdentificationNumber string          `json:"identification_number"`
	IdentificationType   string          `json:"identification_type"`
}

// InitiateCheckout godoc
// @Summary Create a payment intent
// @Description Creates a new payment intent for the authenticated workspace
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param request body CreateIntentRequest true "Intent details"
// @Success 201 {object} contracts.Intent "Intent created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /intents [post]
func (h *Handler) InitiateCheckout(w http.ResponseWriter, r *http.Request) {
	var req CreateIntentRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := CreateIntentInput{
		Amount:               req.Amount,
		Currency:             req.Currency,
		Provider:             req.Provider,
		Metadata:             req.Metadata,
		PaymentMethodID:      req.PaymentMethodID,
		Installments:         req.Installments,
		CardToken:            req.CardToken,
		PaymentMethodType:    req.PaymentMethodType,
		SellerCredentialID:   req.SellerCredentialID,
		PayerEmail:           req.PayerEmail,
		IdentificationNumber: req.IdentificationNumber,
		IdentificationType:   req.IdentificationType,
	}

	intent, err := h.commands.InitiateCheckout(r.Context(), in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(intent).Send(w)
}

// List godoc
// @Summary List payment intents
// @Description Lists all payment intents for the authenticated user. Accessible via API key or user session.
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string false "X-API-Key: tp_xxxxxxxx"
// @Param Cookie header string false "Cookie: access_token=xxx"
// @Security APIKey
// @Security Cookie
// @Success 200 {array} contracts.Intent "Intents retrieved successfully"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /intents [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	intents, err := h.queries.List(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	out := make([]contracts.Intent, 0, len(intents))
	for _, i := range intents {
		out = append(out, i)
	}

	resp.OK().WithData(out).Send(w)
}

// ListByWorkspace godoc
// @Summary List payment intents by workspace
// @Description Lists all payment intents for the authenticated workspace. Accessible via API key or user session.
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string false "X-API-Key: tp_xxxxxxxx"
// @Param Cookie header string false "Cookie: access_token=xxx"
// @Security APIKey
// @Security Cookie
// @Success 200 {array} contracts.Intent "Intents retrieved successfully"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /workspaces/{name}/intents [get]
func (h *Handler) ListByWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	intents, err := h.queries.ListByWorkspace(r.Context(), workspaceName)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	out := make([]contracts.Intent, 0, len(intents))
	for _, i := range intents {
		out = append(out, i)
	}

	resp.OK().WithData(out).Send(w)
}

type PayIntentRequest struct {
	SellerCredentialID uuid.UUID `json:"seller_credential_id"`
}

// Charge godoc
// @Summary Pay a payment intent
// @Description Charges the payment provider for a pending intent using the provided card token
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param intent_id path string true "Intent ID"
// @Param request body PayIntentRequest true "Payment details"
// @Success 200 {object} contracts.Intent "Intent charged successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /intents/{intent_id}/charge [post]
func (h *Handler) Charge(w http.ResponseWriter, r *http.Request) {
	intentID, rs := validation.GetUUID(r, "intent_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req PayIntentRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	intent, err := h.commands.Charge(r.Context(), intentID, req.SellerCredentialID)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			resp.NotFound("intent not found").Send(w)
			return
		}
		if errx.IsKind(err, "invalid") {
			resp.BadRequest(err.Error()).Send(w)
			return
		}
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(intent).Send(w)
}
