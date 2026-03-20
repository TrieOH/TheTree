package workspaces_handler

import (
	"TriePayments/internal/core/application/intents/commands"
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// CreateIntent godoc
// @Summary Create a payment intent
// @Description Creates a new payment intent for the authenticated workspace
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param request body dto.CreateIntentRequest true "Intent details"
// @Success 201 {object} dto.IntentResponse "Intent created successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /intents [post]
func (h *Handler) CreateIntent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateIntentRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := commands.CreateIntentInput{
		Amount:             req.Amount,
		Currency:           req.Currency,
		Provider:           req.Provider,
		Metadata:           req.Metadata,
		PaymentMethodID:    req.PaymentMethodID,
		Installments:       req.Installments,
		CardToken:          req.CardToken,
		PaymentMethodType:  req.PaymentMethodType,
		SellerCredentialID: req.SellerCredentialID,
		PayerEmail:         req.PayerEmail,
	}

	intent, err := h.commands.CreateIntent(r.Context(), in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(dto.MapIntentResponse(intent)).Send(w)
}
