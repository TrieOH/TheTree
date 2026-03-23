package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/errx"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (uc *CommandService) HandleMercadoPagoWebhook(ctx context.Context, mpOrderID string, eventID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.HandleMercadoPagoWebhook")
	defer span.End()

	log.Printf("[mp-webhook] looking up intent for order_id=%s", mpOrderID)

	intent, err := uc.intents.GetByMPOrderID(ctx, mpOrderID)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			log.Printf("[mp-webhook] no intent found for mp order %s", mpOrderID)
			return nil
		}
		return err
	}

	log.Printf("[mp-webhook] found intent=%s status=%s", intent.ID, intent.Status)

	if intent.SellerCredentialID == nil {
		return fmt.Errorf("intent %s has no seller credential", intent.ID)
	}

	cred, err := uc.credentials.GetByID(ctx, *intent.SellerCredentialID)
	if err != nil {
		return fmt.Errorf("failed to fetch seller credential: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.mercadopago.com/v1/payments/"+mpOrderID,
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cred.Credentials.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[mp-webhook] failed to fetch mp order %s: %v", mpOrderID, err)
		return err
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("mercadopago fetch order error %d: %s", resp.StatusCode, string(rawBody))
	}

	var mpOrder struct {
		ID           int64  `json:"id"`
		Status       string `json:"status"`
		StatusDetail string `json:"status_detail"`
	}
	if err := json.Unmarshal(rawBody, &mpOrder); err != nil {
		return err
	}

	log.Printf("[mp-webhook] mp order %s status=%s status_detail=%s", mpOrderID, mpOrder.Status, mpOrder.StatusDetail)

	event := mapMPOrderStatusToEvent(mpOrder.Status, mpOrder.StatusDetail)
	if event == "" {
		log.Printf("[mp-webhook] ignoring mp order status=%s detail=%s order=%s", mpOrder.Status, mpOrder.StatusDetail, mpOrderID)
		return nil
	}

	log.Printf("[mp-webhook] dispatching event=%s for intent=%s", event, intent.ID)
	return uc.HandleProviderWebhook(ctx, eventID, "mercadopago", intent.ID.String(), event)
}

func mapMPOrderStatusToEvent(status, statusDetail string) string {
	switch status {
	case "approved":
		return domain.EventPaymentSucceeded
	case "pending", "in_process", "authorized":
		return ""
	case "rejected", "cancelled", "refunded", "charged_back":
		return domain.EventPaymentFailed
	default:
		return ""
	}
}
