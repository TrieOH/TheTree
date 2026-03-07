package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// RegisterWebhookEndpoint godoc
// @Summary Register a webhook endpoint
// @Description Registers a URL to receive normalized payment events for the given workspace
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param request body dto.RegisterWebhookEndpointRequest true "Endpoint details"
// @Success 201 {object} dto.WebhookEndpointResponse "Endpoint registered successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/webhooks [post]
func (h *Handler) RegisterWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	var req dto.RegisterWebhookEndpointRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	endpoint, err := h.commands.RegisterWebhookEndpoint(r.Context(), workspaceName, req.URL)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(dto.WebhookEndpointResponse{
		ID:          endpoint.ID,
		WorkspaceID: endpoint.WorkspaceID,
		URL:         endpoint.URL,
		Secret:      endpoint.Secret, // only time secret is returned
		CreatedAt:   endpoint.CreatedAt,
	}).Send(w)
}
