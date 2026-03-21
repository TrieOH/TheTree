package workspaces_handler

import (
	"TriePayments/internal/core/domain"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// List godoc
// @Summary List payment intents
// @Description Lists all payment intents for the authenticated user. Accessible via API key or user session.
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string false "X-API-Key: tp_xxxxxxxx"
// @Param Cookie header string false "Cookie: access_token=xxx"
// @Security APIKey
// @Security Cookie
// @Success 200 {array} domain.Intent "Intents retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /intents [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	intents, err := h.queries.List(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	out := make([]domain.Intent, 0, len(intents))
	for _, i := range intents {
		out = append(out, i)
	}

	resp.OK().WithData(out).Send(w)
}

// ListByWorkspace godoc
// @Summary List payment intents by workspace
// @Description Lists all payment intents for the authenticated workspace. Accessible via API key or user session.
// @Tags intents
// @Accept json
// @Produce json
// @Param X-API-Key header string false "X-API-Key: tp_xxxxxxxx"
// @Param Cookie header string false "Cookie: access_token=xxx"
// @Security APIKey
// @Security Cookie
// @Success 200 {array} domain.Intent "Intents retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/intents [get]
func (h *Handler) ListByWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	intents, err := h.queries.ListByWorkspace(r.Context(), workspaceName)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	out := make([]domain.Intent, 0, len(intents))
	for _, i := range intents {
		out = append(out, i)
	}

	resp.OK().WithData(out).Send(w)
}
