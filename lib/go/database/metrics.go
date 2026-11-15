package database

import (
	"context"
	"errors"
	"sync"

	"lib/errx"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// pgxpool pool statistics, read on each collection and pushed over OTLP
// with the HTTP metrics. Names carry their Prometheus suffixes and no unit,
// so OTLP ingestion stores them as is: dashboards query these names.
var (
	poolMu   sync.RWMutex
	poolRef  *pgxpool.Pool
	poolOnce sync.Once
)

// poolInstruments creates the pool instruments on meter and the callback
// that observes them.
func poolInstruments(meter metric.Meter) error {
	var errs []error
	gauge := func(name, desc string) metric.Int64ObservableGauge {
		g, err := meter.Int64ObservableGauge(name, metric.WithDescription(desc))
		errs = append(errs, err)
		return g
	}
	counter := func(name, desc string) metric.Int64ObservableCounter {
		c, err := meter.Int64ObservableCounter(name, metric.WithDescription(desc))
		errs = append(errs, err)
		return c
	}

	acquired := gauge("pgx_pool_acquired_conns", "Connections currently acquired by clients.")
	idle := gauge("pgx_pool_idle_conns", "Connections currently idle in the pool.")
	total := gauge("pgx_pool_total_conns", "Total connections currently in the pool.")
	maxConns := gauge("pgx_pool_max_conns", "Maximum number of connections the pool can hold.")
	constructing := gauge("pgx_pool_constructing_conns", "Connections currently being constructed.")
	emptyAcq := counter("pgx_pool_empty_acquires_total", "Total acquires that had to wait because the pool was empty.")
	canceledAcq := counter("pgx_pool_canceled_acquires_total", "Total acquires canceled while waiting for a connection.")
	err := errors.Join(errs...)
	if err != nil {
		return err
	}

	_, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		poolMu.RLock()
		p := poolRef
		poolMu.RUnlock()
		if p == nil {
			return nil
		}

		st := p.Stat()
		o.ObserveInt64(acquired, int64(st.AcquiredConns()))
		o.ObserveInt64(idle, int64(st.IdleConns()))
		o.ObserveInt64(total, int64(st.TotalConns()))
		o.ObserveInt64(maxConns, int64(st.MaxConns()))
		o.ObserveInt64(constructing, int64(st.ConstructingConns()))
		o.ObserveInt64(emptyAcq, st.EmptyAcquireCount())
		o.ObserveInt64(canceledAcq, st.CanceledAcquireCount())
		return nil
	}, acquired, idle, total, maxConns, constructing, emptyAcq, canceledAcq)
	return err
}

// RegisterPoolMetrics points the shared instruments at pool, creating them
// on the global meter the first time. Call once per process.
func RegisterPoolMetrics(pool *pgxpool.Pool) {
	poolMu.Lock()
	poolRef = pool
	poolMu.Unlock()

	poolOnce.Do(func() {
		err := poolInstruments(otel.Meter("lib/database"))
		if err != nil {
			errx.Exit(err, "error registering pool metrics")
		}
	})
}
