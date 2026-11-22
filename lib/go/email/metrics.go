package email

import (
	"context"
	"errors"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Instruments for Grafana. Names carry their Prometheus suffixes and no unit,
// like lib/httpserver's, so OTLP ingestion stores them as is.
var (
	instrumentsOnce sync.Once
	sentCounter     metric.Int64Counter
	quotaUsedGauge  metric.Int64Gauge
)

func instruments() {
	instrumentsOnce.Do(func() {
		meter := otel.Meter("lib/email")
		sentCounter, _ = meter.Int64Counter("email_recipients_total",
			metric.WithDescription("Recipients per send attempt, by result: sent, failed or limited."))
		quotaUsedGauge, _ = meter.Int64Gauge("email_monthly_quota_used",
			metric.WithDescription("Recipients this service has mailed this calendar month, and its limit."))
	})
}

func recordSend(ctx context.Context, recipients int, err error) {
	instruments()
	if sentCounter == nil {
		return
	}
	result := "sent"
	switch {
	case errors.Is(err, ErrMonthlyLimit):
		result = "limited"
	case err != nil:
		result = "failed"
	}
	sentCounter.Add(ctx, int64(recipients), metric.WithAttributes(attribute.String("result", result)))
}

// RecordQuotaUsed reports this month's total after a reservation, so
// Grafana can show how close each service is to its SMTP_MONTHLY_LIMIT.
func RecordQuotaUsed(ctx context.Context, used, limit int) {
	instruments()
	if quotaUsedGauge == nil {
		return
	}
	quotaUsedGauge.Record(ctx, int64(used), metric.WithAttributes(attribute.Int("limit", limit)))
}
