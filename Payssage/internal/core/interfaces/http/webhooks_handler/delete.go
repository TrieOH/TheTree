package workspaces_handler

import (
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// DeleteWebhookEndpoint godoc
// @Summary Delete a webhook endpoint
// @Description Deletes a registered webhook endpoint
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param endpoint_id path string true "Endpoint ID"
// @Success 200 {object} object "Endpoint deleted successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/webhooks/{endpoint_id} [delete]
func (h *Handler) DeleteWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	endpointID, rs := validation.GetUUID(r, "endpoint_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	if err := h.commands.DeleteWebhookEndpoint(r.Context(), workspaceName, endpointID); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("endpoint deleted").Send(w)
}
