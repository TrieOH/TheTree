package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// ListWebhookDeliveries godoc
// @Summary List webhook deliveries
// @Description Lists all webhook deliveries for the given endpoint
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param endpoint_id path string true "Endpoint ID"
// @Success 200 {array} dto.WebhookDeliveryResponse "Deliveries retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/webhooks/{endpoint_id}/deliveries [get]
func (h *Handler) ListWebhookDeliveries(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	endpointID, rs := validation.GetUUID(r, "endpoint_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	deliveries, err := h.queries.ListWebhookDeliveries(r.Context(), workspaceName, endpointID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	out := make([]dto.WebhookDeliveryResponse, 0, len(deliveries))
	for _, d := range deliveries {
		out = append(out, dto.WebhookDeliveryResponse{
			ID:              d.ID,
			EndpointID:      d.EndpointID,
			IntentID:        d.IntentID,
			Event:           d.Event,
			Status:          string(d.Status),
			Attempts:        d.Attempts,
			LastAttemptedAt: d.LastAttemptedAt,
			CreatedAt:       d.CreatedAt,
		})
	}

	resp.OK().WithData(out).Send(w)
}
