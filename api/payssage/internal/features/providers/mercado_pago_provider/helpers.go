package mercado_pago_provider

import (
	"encoding/json"
	"fmt"
	"payssage/internal/providers/mercado_pago"
	"payssage/models"
	"strconv"
)

func formatAmount(centavos int64) string {
	return fmt.Sprintf("%d.%02d", centavos/100, centavos%100)
}

func parseAmount(s string) int64 {
	f, _ := strconv.ParseFloat(s, 64)
	return int64(f * 100)
}

func calcApplicationFee(amountCents int64, feeBps int) int64 {
	return (amountCents*int64(feeBps) + 5000) / 10000
}

func buildAdditionalInfo(ai *mercado_pago.AdditionalInfo, fallbackAmountCents int64) map[string]any {
	if ai == nil || len(ai.Items) == 0 {
		return map[string]any{
			"items": []map[string]any{
				{
					"title":      "Online Purchase",
					"quantity":   1,
					"unit_price": json.Number(formatAmount(fallbackAmountCents)),
				},
			},
		}
	}

	items := make([]map[string]any, len(ai.Items))
	for i, item := range ai.Items {
		priceCents := fallbackAmountCents
		if item.UnitPriceCents != nil {
			priceCents = *item.UnitPriceCents
		}
		items[i] = map[string]any{
			"title":      item.Title,
			"quantity":   item.Quantity,
			"unit_price": json.Number(formatAmount(priceCents)),
		}
	}

	return map[string]any{"items": items}
}

// NormalizeStatus maps a MercadoPago /v1/payments status into payssage's
// provider-agnostic models.IntentStatus.
// See: https://www.mercadopago.com/developers/en/docs/checkout-api/payment-management/status
func (p *Provider) NormalizeStatus(status string) models.IntentStatus {
	switch status {
	case "approved":
		return models.IntentStatusSucceeded
	case "pending", "in_process", "authorized":
		return models.IntentStatusPending
	case "rejected":
		return models.IntentStatusRejected
	case "cancelled":
		return models.IntentStatusCancelled
	case "refunded", "charged_back":
		return models.IntentStatusRefunded
	default:
		return models.IntentStatusPending
	}
}

// MapStatusDetail normalizes MercadoPago's status_detail codes into a
// small, provider-agnostic set for product/UI use, e.g. showing the payer
// a reason without leaking raw MP vocabulary. Returns nil when there is no
// meaningful detail to store (empty status_detail), matching the nullable
// status_detail column.
// See: https://www.mercadopago.com/developers/en/docs/checkout-api/response-handling/collection-results
func (p *Provider) MapStatusDetail(statusDetail string) *models.IntentStatusDetail {
	if statusDetail == "" {
		return nil
	}

	var detail models.IntentStatusDetail
	switch statusDetail {
	case "cc_rejected_insufficient_amount":
		detail = models.StatusDetailInsufficientFunds
	case "cc_rejected_high_risk", "cc_rejected_blacklist":
		detail = models.StatusDetailHighRisk
	case "cc_rejected_bad_filled_card_number", "cc_rejected_bad_filled_date",
		"cc_rejected_bad_filled_other", "cc_rejected_invalid_installments":
		detail = models.StatusDetailInvalidCard
	case "cc_rejected_card_disabled":
		detail = models.StatusDetailCardDisabled
	case "cc_rejected_card_error":
		detail = models.StatusDetailExpiredCard
	case "cc_rejected_bad_filled_security_code":
		detail = models.StatusDetailInvalidSecurityCode
	case "pending_review_manual", "pending_contingency":
		detail = models.StatusDetailPendingReview
	default:
		detail = models.StatusDetailOther
	}

	return &detail
}
