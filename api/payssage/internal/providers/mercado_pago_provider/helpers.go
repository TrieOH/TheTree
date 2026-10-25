package mercado_pago_provider

import (
	"encoding/json"
	"fmt"
	"payssage/internal/providers/mercado_pago"
	"payssage/models"
	"sort"
	"strings"
	"unicode/utf8"
)

func formatAmount(centavos int64) string {
	return fmt.Sprintf("%d.%02d", centavos/100, centavos%100)
}

// paymentDescription builds the MP top-level `description` from the item
// titles — the field MP shows on the buyer's payment receipt/email
// (`additional_info.items[].title` only feeds the checkout screen; without
// a description MP renders "Produto sem nome" on the receipt). Items are
// grouped by title with their total quantity ("2x Ticket Legal"). Falls
// back to a generic label when there are no usable titles so the receipt
// never shows a nameless product. Bounded length: MP truncates long
// descriptions.
func paymentDescription(ai *mercado_pago.AdditionalInfo) string {
	const fallback = "Payssage purchase"
	if ai == nil || len(ai.Items) == 0 {
		return fallback
	}

	counts := make(map[string]int)
	for _, item := range ai.Items {
		title := strings.TrimSpace(item.Title)
		if title == "" {
			continue
		}
		qty := max(item.Quantity, 1)
		counts[title] += qty
	}
	if len(counts) == 0 {
		return fallback
	}

	titles := make([]string, 0, len(counts))
	for title, count := range counts {
		if count > 1 {
			titles = append(titles, fmt.Sprintf("%dx %s", count, title))
		} else {
			titles = append(titles, title)
		}
	}
	sort.Strings(titles) // map iteration is unordered — stable output for receipts/tests

	desc := strings.Join(titles, ", ")
	return truncateDescription(desc, 250)
}

// truncateDescription caps a description at limit bytes without splitting a
// UTF-8 character or cutting a name in half: it backs off to the last ", "
// separator inside the limit (dropping whole trailing items), and only when
// there is no separator — a single over-long item — cuts at a rune boundary.
func truncateDescription(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	cut := s[:limit]
	for len(cut) > 0 && !utf8.ValidString(cut) {
		cut = cut[:len(cut)-1] // back off to a rune boundary (don't split ç, ã, …)
	}
	if i := strings.LastIndex(cut, ", "); i > 0 {
		return cut[:i] // whole-item boundary — drop the tail items, keep them intact
	}
	return cut
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
	case "partially_refunded":
		// A partial refund (e.g. initiated from the MP panel) still means the
		// order is being reversed — treat it as refunded rather than falling
		// through to pending (which would re-trigger pending-race logic).
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
