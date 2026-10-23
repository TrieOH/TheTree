package jobs

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	libjobs "lib/jobs"
)

type ping struct{}

func (ping) Kind() string { return "ping" }

func TestInProcessRefusesBeforeServe(t *testing.T) {
	c, err := newCarrier("")
	if err != nil {
		t.Fatal(err)
	}
	err = c.Enqueue(context.Background(), ping{})
	if !errors.Is(err, ErrNotServing) {
		t.Fatalf("err = %v, want ErrNotServing", err)
	}
}

func TestInProcessRunsAfterServe(t *testing.T) {
	c, err := newCarrier("")
	if err != nil {
		t.Fatal(err)
	}
	ran := make(chan struct{}, 1)
	reg := libjobs.NewRegistry()
	libjobs.Register(reg, func(context.Context, ping) error {
		ran <- struct{}{}
		return nil
	})
	c.Serve(reg)
	defer func() { _ = c.close(context.Background()) }()

	if !c.Served() || c.Handler() != nil {
		t.Fatalf("served = %v, handler = %v; want served with no handler in-process", c.Served(), c.Handler())
	}
	err = c.Enqueue(context.Background(), ping{})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-ran:
	case <-time.After(2 * time.Second):
		t.Fatal("job did not run")
	}
}

func TestHTTPEnqueuesBeforeServeAndHandlesAfter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	c, err := newCarrier(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Enqueue(context.Background(), ping{})
	if err != nil {
		t.Fatalf("enqueue before Serve over HTTP: %v", err)
	}
	if c.Handler() != nil {
		t.Fatal("handler set before Serve")
	}
	c.Serve(libjobs.NewRegistry())
	if c.Handler() == nil {
		t.Fatal("no handler after Serve")
	}
}

func TestServeTwicePanics(t *testing.T) {
	c, err := newCarrier("http://example.invalid")
	if err != nil {
		t.Fatal(err)
	}
	c.Serve(libjobs.NewRegistry())
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic on the second Serve")
		}
	}()
	c.Serve(libjobs.NewRegistry())
}
