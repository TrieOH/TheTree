package workspaces_handler

import (
	"TriePayments/internal/core/infrastructure/providers"
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"context"
	"encoding/json"
	"log"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
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

	switch provider {
	case "mercadopago":
		secret := viper.GetString("MP_WEBHOOK_SECRET")
		if secret != "" && !providers.VerifyMercadoPagoSignature(r, secret) {
			resp.BadRequest("invalid signature").Send(w)
			return
		}

		var req dto.MercadoPagoWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.BadRequest("invalid payload").Send(w)
			return
		}

		log.Printf("[webhook] mercadopago received action=%s data.id=%s", req.Action, req.Data.ID)

		if req.Data.ID == "" {
			log.Printf("[webhook] mercadopago ignoring ping with no data.id")
			resp.OK("ignored").Send(w)
			return
		}

		resp.OK("received").Send(w)

		go func() {
			if err := h.commands.HandleMercadoPagoWebhook(context.Background(), req.Data.ID); err != nil {
				log.Printf("[webhook] mercadopago err=%v", err)
			}
		}()

	default:
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
}
