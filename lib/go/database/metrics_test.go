package database

import (
	"testing"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// TestPoolInstrumentsNilPool — before a pool is registered the callback
// observes nothing and doesn't panic.
func TestPoolInstrumentsNilPool(t *testing.T) {
	poolMu.Lock()
	poolRef = nil
	poolMu.Unlock()

	reader := sdkmetric.NewManualReader()
	err := poolInstruments(sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader)).Meter("test"))
	if err != nil {
		t.Fatal(err)
	}

	var rm metricdata.ResourceMetrics
	err = reader.Collect(t.Context(), &rm)
	if err != nil {
		t.Fatal(err)
	}
	for _, sm := range rm.ScopeMetrics {
		if len(sm.Metrics) != 0 {
			t.Fatalf("observed %d metrics with no pool, want 0", len(sm.Metrics))
		}
	}
}
