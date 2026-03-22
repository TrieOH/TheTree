package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/errx"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// Charge godoc
// @Summary Pay a payment intent
// @Description Charges the payment provider for a pending intent using the provided card token
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param intent_id path string true "Intent ID"
// @Param request body dto.PayIntentRequest true "Payment details"
// @Success 200 {object} domain.Intent "Intent charged successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /intents/{intent_id}/charge [post]
func (h *Handler) Charge(w http.ResponseWriter, r *http.Request) {
	intentID, rs := validation.GetUUID(r, "intent_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.PayIntentRequest
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
