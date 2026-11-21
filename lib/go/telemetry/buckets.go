package telemetry

import (
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// DurationBuckets are the histogram boundaries for every duration in
// seconds: 0.25ms doubling 16 times (~8s).
var DurationBuckets = func() []float64 {
	b := make([]float64, 16)
	for i := range b {
		b[i] = 0.00025 * float64(int(1)<<i)
	}
	return b
}()

// secondsHistograms gives every histogram in seconds DurationBuckets. The
// SDK's default boundaries (0, 5, 10 … 10000) are sized for milliseconds, so
// a seconds histogram that sets none of its own (otelpgx's
// db.client.operation.duration) put every observation in the first bucket
// and read 4.75s at any percentile.
var secondsHistograms = sdkmetric.NewView(
	sdkmetric.Instrument{Kind: sdkmetric.InstrumentKindHistogram, Unit: "s"},
	sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{Boundaries: DurationBuckets}},
)
