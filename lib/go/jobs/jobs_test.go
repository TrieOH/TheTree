package jobs_test

import (
	"context"
	"errors"
	"testing"

	"lib/jobs"
)

type greet struct {
	Name string `json:"name"`
}

func (greet) Kind() string { return "greet" }

func TestDispatchRoundTrip(t *testing.T) {
	reg := jobs.NewRegistry()
	var got string
	jobs.Register(reg, func(_ context.Context, j greet) error {
		got = j.Name
		return nil
	})

	env, err := jobs.Encode(greet{Name: "ada"})
	if err != nil {
		t.Fatal(err)
	}
	err = reg.Dispatch(context.Background(), env)
	if err != nil {
		t.Fatal(err)
	}
	if got != "ada" {
		t.Fatalf("handler got %q, want ada", got)
	}
}

func TestDispatchUnknownKind(t *testing.T) {
	err := jobs.NewRegistry().Dispatch(context.Background(), jobs.Envelope{Kind: "nope"})
	if !errors.Is(err, jobs.ErrUnknownKind) {
		t.Fatalf("err = %v, want ErrUnknownKind", err)
	}
}

func TestRegisterTwicePanics(t *testing.T) {
	reg := jobs.NewRegistry()
	jobs.Register(reg, func(context.Context, greet) error { return nil })
	defer func() {
		if recover() == nil {
			t.Fatal("second Register did not panic")
		}
	}()
	jobs.Register(reg, func(context.Context, greet) error { return nil })
}
