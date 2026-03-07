package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/errx"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// CancelIntent godoc
// @Summary Cancel a payment intent
// @Description Cancels a pending payment intent
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param intent_id path string true "Intent ID"
// @Success 200 {object} dto.IntentResponse "Intent cancelled successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
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

	resp.OK().WithData(dto.MapIntentResponse(intent)).Send(w)
}
