package jobs_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"lib/jobs"
)

type flaky struct {
	Err string `json:"err"`
}

func (flaky) Kind() string { return "flaky" }

type flakyGaveUp struct {
	Err string `json:"err"`
}

func (flakyGaveUp) Kind() string { return "flaky.gave_up" }

var policy = jobs.Policy{MaxAttempts: 3, Backoff: jobs.Exponential(time.Second, 3*time.Second)}

// handleRegistry fails a flaky job with its Err ("permanent" marks it so),
// and records each give-up's Err.
func handleRegistry(gaveUp *[]string, giveUpErr error) *jobs.Registry {
	reg := jobs.NewRegistry(jobs.WithPolicy(policy))
	jobs.Register(reg, func(_ context.Context, j flaky) error {
		switch j.Err {
		case "":
			return nil
		case "permanent":
			return jobs.Permanent(errors.New(j.Err))
		default:
			return errors.New(j.Err)
		}
	}, jobs.GiveUpWith(func(j flaky) jobs.Job { return flakyGaveUp(j) }))
	jobs.Register(reg, func(_ context.Context, j flakyGaveUp) error {
		*gaveUp = append(*gaveUp, j.Err)
		return giveUpErr
	})
	return reg
}

func delivery(t *testing.T, job jobs.Job, attempt int, dead bool) jobs.Delivery {
	t.Helper()
	env, err := jobs.Encode(context.Background(), job)
	if err != nil {
		t.Fatal(err)
	}
	return jobs.Delivery{Envelope: env, Attempt: attempt, Dead: dead}
}

func TestHandleOutcomes(t *testing.T) {
	cases := []struct {
		name      string
		job       jobs.Job
		attempt   int
		dead      bool
		giveUpErr error
		want      jobs.Result
		gaveUp    []string
	}{
		{name: "success", job: flaky{}, attempt: 1, want: jobs.Result{Outcome: jobs.Done}},
		{name: "failure backs off", job: flaky{Err: "down"}, attempt: 1, want: jobs.Result{Outcome: jobs.Retry, After: time.Second}},
		{name: "backoff doubles", job: flaky{Err: "down"}, attempt: 2, want: jobs.Result{Outcome: jobs.Retry, After: 2 * time.Second}},
		{name: "last attempt gives up", job: flaky{Err: "down"}, attempt: 3, want: jobs.Result{Outcome: jobs.Drop}, gaveUp: []string{"down"}},
		{name: "permanent drops without giving up", job: flaky{Err: "permanent"}, attempt: 1, want: jobs.Result{Outcome: jobs.Drop}},
		{name: "dead gives up without running", job: flaky{}, attempt: 1, dead: true, want: jobs.Result{Outcome: jobs.Drop}, gaveUp: []string{""}},
		{name: "failed give-up retries", job: flaky{Err: "down"}, attempt: 3, giveUpErr: errors.New("db"), want: jobs.Result{Outcome: jobs.Retry, After: 3 * time.Second}, gaveUp: []string{"down"}},
		{name: "unknown kind drops", job: greet{}, attempt: 1, want: jobs.Result{Outcome: jobs.Drop}},
		{name: "attempt 0 counts as the first", job: flaky{Err: "down"}, attempt: 0, want: jobs.Result{Outcome: jobs.Retry, After: time.Second}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gaveUp []string
			reg := handleRegistry(&gaveUp, c.giveUpErr)
			got := reg.Handle(context.Background(), delivery(t, c.job, c.attempt, c.dead))
			if got != c.want {
				t.Errorf("Handle = %+v, want %+v", got, c.want)
			}
			if len(gaveUp) != len(c.gaveUp) || (len(gaveUp) > 0 && gaveUp[0] != c.gaveUp[0]) {
				t.Errorf("gave up = %q, want %q", gaveUp, c.gaveUp)
			}
		})
	}
}

func TestHandleDropsUndecodableArgs(t *testing.T) {
	var gaveUp []string
	reg := handleRegistry(&gaveUp, nil)
	got := reg.Handle(context.Background(), jobs.Delivery{Kind: "flaky", Args: []byte(`"nope"`), Attempt: 1})
	if got.Outcome != jobs.Drop {
		t.Fatalf("Handle = %+v, want a drop", got)
	}
}

func TestExponentialCaps(t *testing.T) {
	b := jobs.Exponential(30*time.Second, time.Hour)
	for attempt, want := range map[int]time.Duration{1: 30 * time.Second, 2: time.Minute, 7: 32 * time.Minute, 8: time.Hour, 50: time.Hour} {
		if got := b(attempt); got != want {
			t.Errorf("attempt %d: %v, want %v", attempt, got, want)
		}
	}
}
