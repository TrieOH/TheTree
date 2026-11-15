package telemetry

import (
	"os"
	"testing"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// TestRuntimeInstrumentsObserve — every metric the dashboards read is
// observed under its Prometheus name with a plausible value.
func TestRuntimeInstrumentsObserve(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	err := runtimeInstruments(sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader)).Meter("test"))
	if err != nil {
		t.Fatal(err)
	}

	var rm metricdata.ResourceMetrics
	err = reader.Collect(t.Context(), &rm)
	if err != nil {
		t.Fatal(err)
	}
	points := map[string][]float64{}
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch d := m.Data.(type) {
			case metricdata.Gauge[float64]:
				for _, p := range d.DataPoints {
					points[m.Name] = append(points[m.Name], p.Value)
				}
			case metricdata.Sum[float64]:
				if !d.IsMonotonic {
					t.Errorf("%s is not monotonic", m.Name)
				}
				for _, p := range d.DataPoints {
					points[m.Name] = append(points[m.Name], p.Value)
				}
			}
		}
	}

	positive := []string{"go_goroutines", "go_memstats_heap_inuse_bytes", "go_gc_gomemlimit_bytes", "process_start_time_seconds", "process_cpu_seconds_total"}
	_, err = os.Stat("/proc/self/statm")
	if err == nil {
		positive = append(positive, "process_resident_memory_bytes")
	}
	for _, name := range positive {
		if v := points[name]; len(v) != 1 || v[0] <= 0 {
			t.Errorf("%s = %v, want one positive value", name, v)
		}
	}
	if got := len(points["go_gc_duration_seconds"]); got != len(gcQuantiles) {
		t.Errorf("go_gc_duration_seconds has %d quantiles, want %d", got, len(gcQuantiles))
	}
}
