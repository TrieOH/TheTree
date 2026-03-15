package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// ListWebhookEvents godoc
// @Summary List webhook events
// @Description Lists all inbound provider webhook events for the given workspace
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {array} dto.WebhookEventResponse "Events retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/webhook-events [get]
func (h *Handler) ListWebhookEvents(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	events, err := h.queries.ListWebhookEvents(r.Context(), workspaceName)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	out := make([]dto.WebhookEventResponse, 0, len(events))
	for _, e := range events {
		out = append(out, dto.WebhookEventResponse{
			ID:          e.ID,
			Provider:    e.Provider,
			EventType:   e.EventType,
			ExternalID:  e.ExternalID,
			WorkspaceID: e.WorkspaceID,
			IntentID:    e.IntentID,
			ReceivedAt:  e.ReceivedAt,
		})
	}

	resp.OK().WithData(out).Send(w)
}
