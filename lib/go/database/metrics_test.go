package database

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestPoolCollectorDescribesSevenMetrics(t *testing.T) {
	c := poolCollector{}
	ch := make(chan *prometheus.Desc, 16)
	c.Describe(ch)
	close(ch)

	n := 0
	for range ch {
		n++
	}
	if n != 7 {
		t.Fatalf("pool collector described %d metrics, want 7", n)
	}
}

func TestPoolCollectorNilPoolNoPanic(_ *testing.T) {
	poolMu.Lock()
	poolRef = nil
	poolMu.Unlock()

	c := poolCollector{}
	ch := make(chan prometheus.Metric)
	go func() {
		defer close(ch)
		c.Collect(ch) // must not panic when no pool is registered
	}()
	for range ch { //nolint:revive // drain; expect zero metrics
	}
}
