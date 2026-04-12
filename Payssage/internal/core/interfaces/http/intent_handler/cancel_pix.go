package workspaces_handler

import (
	"TriePayments/internal/core/application/intents/commands"
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/errx"
	"TriePayments/internal/shared/validation"
	"net/http"

	_ "TriePayments/internal/core/domain"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// CancelPix godoc
// @Summary Cancel a payment intent and its pix
// @Description Cancels a pending payment intent and its related pix QR code
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param intent_id path string true "Intent ID"
// @Success 200 {object} domain.Intent "Pix canceled successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /intents/{intent_id}/cancel-pix [post]
func (h *Handler) CancelPix(w http.ResponseWriter, r *http.Request) {
	intentID, rs := validation.GetUUID(r, "intent_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.CancelPixRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := commands.CancelPixInput{
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
