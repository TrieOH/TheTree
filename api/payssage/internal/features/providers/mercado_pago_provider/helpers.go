package mercado_pago_provider

import (
	"encoding/json"
	"fmt"
	"payssage/internal/providers/mercado_pago"
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

func wrapMPError(err error) error {
	return fmt.Errorf("mercadopago: %w", err)
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
