package oauth_handler

import (
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// CompleteOAuth godoc
// @Summary OAuth callback from provider
// @Description Handles the provider callback, exchanges code for token, stores credential, redirects to final URL
// @Tags oauth
// @Param provider path string true "Provider name (e.g. mercadopago)"
// @Param code query string true "Authorization code from provider"
// @Param state query string true "State token"
// @Success 302 "Final redirect URL"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /oauth/{provider}/callback [get]
func (h *Handler) CompleteOAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		resp.BadRequest("code and state are required").Send(w)
		return
	}

	finalURL, err := h.commands.CompleteOAuth(r.Context(), provider, state, code)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(map[string]string{
		"url": finalURL,
	}).Send(w)
}
