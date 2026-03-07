package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"context"
	"log"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// HandleProviderWebhook godoc
// @Summary Receive provider webhook
// @Description Receives a webhook from a payment provider, normalizes it and forwards to registered endpoints
// @Tags webhooks
// @Accept json
// @Produce json
// @Param provider path string true "Provider name (e.g. mock, stripe)"
// @Param request body dto.ProviderWebhookRequest true "Provider webhook payload"
// @Success 200 {object} object "Received"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /webhooks/{provider} [post]
func (h *Handler) HandleProviderWebhook(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	var req dto.ProviderWebhookRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("received").Send(w)

	go func() {
		if err := h.commands.HandleProviderWebhook(context.Background(), provider, req.IntentID, req.Event); err != nil {
			log.Printf("[webhook] provider=%s err=%v", provider, err)
		}
	}()
}
