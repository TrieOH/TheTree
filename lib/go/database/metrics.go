package database

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

// pgxpool pool statistics exposed as Prometheus gauges on the process's
// /metrics endpoint. The scrape job label (e.g. "univents") comes from the
// VictoriaMetrics scrape config, matching the HTTP metrics.
var (
	poolMu   sync.RWMutex
	poolRef  *pgxpool.Pool
	poolOnce sync.Once

	descAcquired  = prometheus.NewDesc("pgx_pool_acquired_conns", "Connections currently acquired by clients.", nil, nil)
	descIdle      = prometheus.NewDesc("pgx_pool_idle_conns", "Connections currently idle in the pool.", nil, nil)
	descTotal     = prometheus.NewDesc("pgx_pool_total_conns", "Total connections currently in the pool.", nil, nil)
	descMax       = prometheus.NewDesc("pgx_pool_max_conns", "Maximum number of connections the pool can hold.", nil, nil)
	descConstruct = prometheus.NewDesc("pgx_pool_constructing_conns", "Connections currently being constructed.", nil, nil)
	descEmptyAcq  = prometheus.NewDesc("pgx_pool_empty_acquires_total", "Total acquires that had to wait because the pool was empty.", nil, nil)
	descCancelAcq = prometheus.NewDesc("pgx_pool_canceled_acquires_total", "Total acquires canceled while waiting for a connection.", nil, nil)
)

type poolCollector struct{}

func (poolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- descAcquired
	ch <- descIdle
	ch <- descTotal
	ch <- descMax
	ch <- descConstruct
	ch <- descEmptyAcq
	ch <- descCancelAcq
}

func (poolCollector) Collect(ch chan<- prometheus.Metric) {
	poolMu.RLock()
	p := poolRef
	poolMu.RUnlock()
	if p == nil {
		return
	}

	st := p.Stat()
	ch <- prometheus.MustNewConstMetric(descAcquired, prometheus.GaugeValue, float64(st.AcquiredConns()))
	ch <- prometheus.MustNewConstMetric(descIdle, prometheus.GaugeValue, float64(st.IdleConns()))
	ch <- prometheus.MustNewConstMetric(descTotal, prometheus.GaugeValue, float64(st.TotalConns()))
	ch <- prometheus.MustNewConstMetric(descMax, prometheus.GaugeValue, float64(st.MaxConns()))
	ch <- prometheus.MustNewConstMetric(descConstruct, prometheus.GaugeValue, float64(st.ConstructingConns()))
	ch <- prometheus.MustNewConstMetric(descEmptyAcq, prometheus.CounterValue, float64(st.EmptyAcquireCount()))
	ch <- prometheus.MustNewConstMetric(descCancelAcq, prometheus.CounterValue, float64(st.CanceledAcquireCount()))
}

// RegisterPoolMetrics points the shared collector at pool and registers it
// with the process's default Prometheus registry. Call once per process.
func RegisterPoolMetrics(pool *pgxpool.Pool) {
	poolMu.Lock()
	poolRef = pool
	poolMu.Unlock()

	poolOnce.Do(func() {
		prometheus.MustRegister(poolCollector{})
	})
}
