package river

import (
	"lib/errx"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func NewClient(dbPool *pgxpool.Pool, workers *river.Workers, queues map[string]river.QueueConfig, periodicJobs []*river.PeriodicJob) *river.Client[pgx.Tx] {
	if queues == nil {
		queues = map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		}
	}
	client, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
		Queues:       queues,
		Workers:      workers,
		PeriodicJobs: periodicJobs,
	})
	if err != nil {
		errx.Exit(err, "error creating river client")
	}
	return client
}
