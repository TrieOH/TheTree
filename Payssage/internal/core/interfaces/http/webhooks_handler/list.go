package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// ListWebhookEndpoints godoc
// @Summary List webhook endpoints
// @Description Lists all registered webhook endpoints for the given workspace
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {array} dto.WebhookEndpointListResponse "Endpoints retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/webhooks [get]
func (h *Handler) ListWebhookEndpoints(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	endpoints, err := h.queries.ListWebhookEndpoints(r.Context(), workspaceName)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	out := make([]dto.WebhookEndpointListResponse, 0, len(endpoints))
	for _, e := range endpoints {
		out = append(out, dto.WebhookEndpointListResponse{
			ID:          e.ID,
			WorkspaceID: e.WorkspaceID,
			URL:         e.URL,
			CreatedAt:   e.CreatedAt,
		})
	}

	resp.OK().WithData(out).Send(w)
}
