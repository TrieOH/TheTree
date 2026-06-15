package river

import (
	"github.com/riverqueue/river"
)

type WorkerRegistrar func(*river.Workers)

func NewWorkers(registrars ...WorkerRegistrar) *river.Workers {
	workers := river.NewWorkers()
	for _, r := range registrars {
		r(workers)
	}
	return workers
}

func Register[T river.JobArgs](worker river.Worker[T]) WorkerRegistrar {
	return func(workers *river.Workers) {
		river.AddWorker(workers, worker)
	}
}
