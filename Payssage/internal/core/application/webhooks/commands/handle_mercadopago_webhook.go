package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/errx"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/spf13/viper"
)

func (uc *CommandService) HandleMercadoPagoWebhook(ctx context.Context, mpPaymentID string, eventID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.HandleMercadoPagoWebhook")
	defer span.End()

	log.Printf("[mp-webhook] looking up intent for provider_payment_id=%s", mpPaymentID)

	intent, err := uc.intents.GetByProviderPaymentID(ctx, mpPaymentID)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			log.Printf("[mp-webhook] no intent found for mp payment %s", mpPaymentID)
			return nil
		}
		return err
	}

	log.Printf("[mp-webhook] found intent=%s status=%s", intent.ID, intent.Status)

	cfg, err := config.New(viper.GetString("MP_ACCESS_TOKEN"))
	if err != nil {
		return err
	}
	client := payment.NewClient(cfg)
	resource, err := client.Get(ctx, parseInt(mpPaymentID))
	if err != nil {
		log.Printf("[mp-webhook] failed to fetch mp payment %s: %v", mpPaymentID, err)
		return err
	}

	log.Printf("[mp-webhook] mp payment %s status=%s", mpPaymentID, resource.Status)

	event := mapMPStatusToEvent(resource.Status)
	if event == "" {
		log.Printf("[mp-webhook] ignoring mp payment status=%s payment=%s", resource.Status, mpPaymentID)
		return nil
	}

	log.Printf("[mp-webhook] dispatching event=%s for intent=%s", event, intent.ID)
	return uc.HandleProviderWebhook(ctx, eventID, "mercadopago", intent.ID.String(), event)
}

func mapMPStatusToEvent(status string) string {
	switch status {
	case "approved":
		return domain.EventPaymentSucceeded
	case "rejected", "cancelled":
		return domain.EventPaymentFailed
	default:
		return ""
	}
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
