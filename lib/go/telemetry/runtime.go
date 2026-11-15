package telemetry

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"os"
	"runtime/debug"
	"runtime/metrics"
	"strconv"
	"strings"
	"syscall"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// processStart stands in for the process start time; package init runs
// within milliseconds of it.
var processStart = time.Now()

// gcQuantiles are the go_gc_duration_seconds quantiles, in the order
// debug.ReadGCStats fills PauseQuantiles.
var gcQuantiles = []string{"0", "0.25", "0.5", "0.75", "1"}

// runtimeInstruments observes the Go runtime and process metrics the
// dashboards read, under the names the Prometheus Go and process collectors
// used. Names carry their Prometheus suffixes and no unit, so OTLP
// ingestion stores them as is. The process metrics read /proc and are
// skipped where it's missing.
func runtimeInstruments(meter metric.Meter) error {
	var errs []error
	gauge := func(name, desc string) metric.Float64ObservableGauge {
		g, err := meter.Float64ObservableGauge(name, metric.WithDescription(desc))
		errs = append(errs, err)
		return g
	}
	counter := func(name, desc string) metric.Float64ObservableCounter {
		c, err := meter.Float64ObservableCounter(name, metric.WithDescription(desc))
		errs = append(errs, err)
		return c
	}

	goroutines := gauge("go_goroutines", "Number of goroutines that currently exist.")
	heapInuse := gauge("go_memstats_heap_inuse_bytes", "Number of heap bytes that are in use.")
	memLimit := gauge("go_gc_gomemlimit_bytes", "Go runtime memory limit (GOMEMLIMIT).")
	gcPause := gauge("go_gc_duration_seconds", "Quantiles of recent GC stop-the-world pauses.")
	startTime := gauge("process_start_time_seconds", "Start time of the process since unix epoch in seconds.")
	rss := gauge("process_resident_memory_bytes", "Resident memory size in bytes.")
	cpu := counter("process_cpu_seconds_total", "Total user and system CPU time spent in seconds.")
	netRx := counter("process_network_receive_bytes_total", "Number of bytes received by the process over the network.")
	netTx := counter("process_network_transmit_bytes_total", "Number of bytes sent by the process over the network.")
	err := errors.Join(errs...)
	if err != nil {
		return err
	}

	samples := []metrics.Sample{
		{Name: "/sched/goroutines:goroutines"},
		{Name: "/memory/classes/heap/objects:bytes"},
		{Name: "/memory/classes/heap/unused:bytes"},
		{Name: "/gc/gomemlimit:bytes"},
	}
	quantileAttrs := make([]metric.ObserveOption, len(gcQuantiles))
	for i, q := range gcQuantiles {
		quantileAttrs[i] = metric.WithAttributes(attribute.String("quantile", q))
	}

	_, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		metrics.Read(samples)
		o.ObserveFloat64(goroutines, float64(samples[0].Value.Uint64()))
		o.ObserveFloat64(heapInuse, float64(samples[1].Value.Uint64()+samples[2].Value.Uint64()))
		o.ObserveFloat64(memLimit, float64(samples[3].Value.Uint64()))

		stats := debug.GCStats{PauseQuantiles: make([]time.Duration, len(gcQuantiles))}
		debug.ReadGCStats(&stats)
		for i, p := range stats.PauseQuantiles {
			o.ObserveFloat64(gcPause, p.Seconds(), quantileAttrs[i])
		}

		o.ObserveFloat64(startTime, float64(processStart.UnixNano())/1e9)

		var ru syscall.Rusage
		if syscall.Getrusage(syscall.RUSAGE_SELF, &ru) == nil {
			o.ObserveFloat64(cpu, timevalSeconds(ru.Utime)+timevalSeconds(ru.Stime))
		}
		if pages, ok := residentPages(); ok {
			o.ObserveFloat64(rss, float64(pages*int64(os.Getpagesize())))
		}
		if in, out, ok := ipOctets(); ok {
			o.ObserveFloat64(netRx, in)
			o.ObserveFloat64(netTx, out)
		}
		return nil
	}, goroutines, heapInuse, memLimit, gcPause, startTime, rss, cpu, netRx, netTx)
	return err
}

func timevalSeconds(tv syscall.Timeval) float64 {
	return float64(tv.Sec) + float64(tv.Usec)/1e6
}

// residentPages is the second field of /proc/self/statm.
func residentPages() (int64, bool) {
	b, err := os.ReadFile("/proc/self/statm")
	if err != nil {
		return 0, false
	}
	fields := strings.Fields(string(b))
	if len(fields) < 2 {
		return 0, false
	}
	n, err := strconv.ParseInt(fields[1], 10, 64)
	return n, err == nil
}

// ipOctets reads InOctets and OutOctets from the IpExt lines of
// /proc/self/net/netstat: a header line of names, then one of values.
func ipOctets() (in, out float64, ok bool) {
	b, err := os.ReadFile("/proc/self/net/netstat")
	if err != nil {
		return 0, 0, false
	}
	var names []string
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) == 0 || fields[0] != "IpExt:" {
			continue
		}
		if names == nil {
			names = fields
			continue
		}
		var gotIn, gotOut bool
		for i := 1; i < len(fields) && i < len(names); i++ {
			v, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				continue
			}
			switch names[i] {
			case "InOctets":
				in, gotIn = v, true
			case "OutOctets":
				out, gotOut = v, true
			}
		}
		return in, out, gotIn && gotOut
	}
	return 0, 0, false
}
