package workspaces_handler

import (
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// GetByID godoc
// @Summary Get a payment intent
// @Description Retrieves a payment intent by ID
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param intent_id path string true "Intent ID"
// @Success 200 {object} domain.Intent "Intent retrieved successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
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
