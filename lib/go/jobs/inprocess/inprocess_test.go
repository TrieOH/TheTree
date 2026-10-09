package inprocess_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"lib/jobs"
	"lib/jobs/inprocess"
)

type ping struct{}

func (ping) Kind() string { return "ping" }

func TestRetriesUntilSuccess(t *testing.T) {
	reg := jobs.NewRegistry()
	var calls atomic.Int32
	done := make(chan struct{})
	jobs.Register(reg, func(context.Context, ping) error {
		if calls.Add(1) < 3 {
			return errors.New("transient")
		}
		close(done)
		return nil
	})
	r := inprocess.New(reg, inprocess.Options{BaseBackoff: time.Millisecond})
	r.Start()
	defer func() { _ = r.Stop(context.Background()) }()

	err := r.Enqueue(context.Background(), ping{})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("job did not succeed; calls = %d", calls.Load())
	}
}

func TestPeriodicRunOnStart(t *testing.T) {
	reg := jobs.NewRegistry()
	ran := make(chan struct{}, 1)
	jobs.Register(reg, func(context.Context, ping) error {
		select {
		case ran <- struct{}{}:
		default:
		}
		return nil
	})
	r := inprocess.New(reg, inprocess.Options{})
	r.Start(jobs.Periodic{Job: ping{}, Every: time.Hour, RunOnStart: true})
	defer func() { _ = r.Stop(context.Background()) }()

	select {
	case <-ran:
	case <-time.After(2 * time.Second):
		t.Fatal("periodic job did not run on start")
	}
}

func TestStopDrainsAndRefuses(t *testing.T) {
	reg := jobs.NewRegistry()
	var calls atomic.Int32
	jobs.Register(reg, func(context.Context, ping) error {
		time.Sleep(10 * time.Millisecond)
		calls.Add(1)
		return nil
	})
	r := inprocess.New(reg, inprocess.Options{Workers: 1})
	r.Start()
	for range 3 {
		err := r.Enqueue(context.Background(), ping{})
		if err != nil {
			t.Fatal(err)
		}
	}
	err := r.Stop(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 3 {
		t.Fatalf("drained %d jobs, want 3", calls.Load())
	}
	err = r.Enqueue(context.Background(), ping{})
	if !errors.Is(err, inprocess.ErrStopped) {
		t.Fatalf("enqueue after stop = %v, want ErrStopped", err)
	}
}
