package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// A seconds histogram with no boundaries of its own gets DurationBuckets,
// so a 3ms query lands in a low bucket instead of the first 0–5s one.
func TestSecondsHistogramsGetDurationBuckets(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader), sdkmetric.WithView(secondsHistograms))
	h, err := mp.Meter("test").Float64Histogram("db.client.operation.duration", metric.WithUnit("s"))
	if err != nil {
		t.Fatal(err)
	}
	h.Record(context.Background(), 0.003)

	var rm metricdata.ResourceMetrics
	err = reader.Collect(context.Background(), &rm)
	if err != nil {
		t.Fatal(err)
	}
	hist, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
	if !ok {
		t.Fatalf("data = %T, want a float64 histogram", rm.ScopeMetrics[0].Metrics[0].Data)
	}
	bounds := hist.DataPoints[0].Bounds
	if len(bounds) != len(DurationBuckets) || bounds[0] != DurationBuckets[0] {
		t.Fatalf("bounds = %v, want DurationBuckets", bounds)
	}
}
