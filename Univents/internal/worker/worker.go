package worker

import (
	"context"
	"log"
	"univents/internal/core/application/edition/async"
	"univents/internal/core/domain"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

type Deps struct {
	Handlers *async.AsynqHandlers
}

func InitAsynq(deps Deps) (*asynq.Server, *asynq.Client, *asynq.Scheduler, error) {
	redisOpt := asynq.RedisClientOpt{Addr: viper.GetString("REDIS_ADDR"), Password: viper.GetString("REDIS_PASSWORD"), DB: viper.GetInt("REDIS_DB")}

	// Client for enqueueing tasks
	client := asynq.NewClient(redisOpt)

	// Server for processing tasks
	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Printf("task failed: type=%s, err=%v", task.Type(), err)
		}),
		RetryDelayFunc: asynq.DefaultRetryDelayFunc,
	})

	// Scheduler for periodic tasks (if needed later)
	scheduler := asynq.NewScheduler(redisOpt, nil)

	// Setup handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(domain.AsynqEditionOpen, deps.Handlers.HandleOpenEdition)
	mux.HandleFunc(domain.AsynqEditionStart, deps.Handlers.HandleStartEdition)
	mux.HandleFunc(domain.AsynqEditionEnd, deps.Handlers.HandleEndEdition)

	// Run server in background
	go func() {
		if err := server.Run(mux); err != nil {
			log.Fatalf("asynq server error: %v", err)
		}
	}()

	// Run scheduler in background (if using periodic tasks)
	go func() {
		if err := scheduler.Run(); err != nil {
			log.Fatalf("asynq scheduler error: %v", err)
		}
	}()

	return server, client, scheduler, nil
}
