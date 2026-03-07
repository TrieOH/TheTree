package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// ListAPIKeys godoc
// @Summary List API keys
// @Description Lists all API keys for the given workspace (raw keys are never returned)
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {array} dto.APIKeyResponse "API keys retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/keys [get]
func (h *Handler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	keys, err := h.queries.List(r.Context(), workspaceName)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	out := make([]dto.APIKeyResponse, 0, len(keys))
	for _, k := range keys {
		out = append(out, dto.APIKeyResponse{
			ID:        k.ID,
			Name:      k.Name,
			Prefix:    k.KeyPrefix,
			CreatedAt: k.CreatedAt,
			RevokedAt: k.RevokedAt,
		})
	}

	resp.OK().WithData(out).Send(w)
}
