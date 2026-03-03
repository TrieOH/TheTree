package testing

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"univents/initialization"
	"univents/internal/core/application/edition/async"
	editionCommands "univents/internal/core/application/edition/commands"
	editionQueries "univents/internal/core/application/edition/queries"
	"univents/internal/core/application/event/commands"
	"univents/internal/core/application/event/queries"
	"univents/internal/core/infrastructure"
	eventhttp "univents/internal/core/interfaces/http"
	editionhttp "univents/internal/core/interfaces/http/editions"
	"univents/internal/interfaces/http/middleware"
	"univents/internal/interfaces/http/router"
	"univents/internal/interfaces/http/system"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
	"univents/internal/plataform/telemetry"
	"univents/internal/worker"

	"github.com/gavv/httpexpect/v2"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// ============================================================================
// TEST FRAMEWORK - Core Infrastructure
// ============================================================================

// TestSuite manages the entire test environment
type TestSuite struct {
	Server *httptest.Server
	App    *initialization.UniventsApp
	t      *testing.T
}

func NewTestSuite(t *testing.T) *TestSuite {
	suite := &TestSuite{t: t}
	suite.setup()

	t.Cleanup(func() {
		suite.teardown()
	})

	return suite
}

func (s *TestSuite) setup() {
	s.App = initialization.UniventsSetup()

	ctx := context.Background()

	defer s.App.DB.Close()
	defer s.App.Redis.Close()

	shutdown, err := telemetry.InitTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func(ctx context.Context) {
		err := shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(ctx)

	defer func() {
		err := s.App.Scheduler.StopJobs()
		if err != nil {
			log.Printf("Error stopping jobs: %v", err)
		}
		err = s.App.Scheduler.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	q := sqlc.New(s.App.DB)
	txRunner := database.NewPGXTxRunner(s.App.DB)
	tracer := otel.Tracer(string(telemetry.UniventsTracer))
	logs := telemetry.Log()

	authMW := middleware.NewAuthMiddleware(s.App.GaClient, tracer)

	eventRepo := infrastructure.NewEventRepo(q, logs, tracer)
	editionRepo := infrastructure.NewEditionRepo(q, logs, tracer)

	workerHandlers := async.New(editionRepo, s.App.GaClient, tracer, txRunner)
	server, asynqClient, scheduler, err := worker.InitAsynq(worker.Deps{
		Handlers: workerHandlers,
	})
	defer func() {
		scheduler.Shutdown()
		server.Shutdown()
		if err = asynqClient.Close(); err != nil {
			telemetry.Log().Error("error closing the asynq client", zap.Error(err))
		}
	}()

	eventCommands := commands.New(eventRepo, s.App.GaClient, tracer, txRunner)
	eventQueries := queries.New(eventRepo, s.App.GaClient, tracer, txRunner)
	editionC := editionCommands.New(eventRepo, editionRepo, asynqClient, s.App.GaClient, tracer, txRunner)
	editionQ := editionQueries.New(eventRepo, editionRepo, s.App.GaClient, tracer, txRunner)

	eventHandler := eventhttp.NewEventsHandler(eventCommands, eventQueries)
	editionHandler := editionhttp.NewEditionsHandler(editionC, editionQ)

	systemHandler := system.NewUniventsHandler()

	asynqmonHandler := asynqmon.New(asynqmon.Options{
		RootPath: "/admin/asynq",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	})

	deps := &router.HTTPDeps{
		EventsHandler:   eventHandler,
		EditionsHandler: editionHandler,
		SystemHandler:   systemHandler,
		AuthMiddleware:  authMW,
		AsynqmonHandler: asynqmonHandler,
	}

	r := createTestRouter(deps)
	s.Server = httptest.NewServer(r)
}

func (s *TestSuite) teardown() {
	if s.Server != nil {
		s.Server.Close()
	}
	if s.App.DB != nil {
		s.App.DB.Close()
	}
	if s.App.Redis != nil {
		s.App.Redis.Close()
	}
}

// NewClient creates a new API client for testing
func (s *TestSuite) NewClient(t *testing.T) *Client {
	return &Client{
		expect: httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  s.Server.URL,
			Reporter: httpexpect.NewAssertReporter(t),
		}),
		t:       t,
		baseURL: s.Server.URL,
	}
}

func createTestRouter(deps *router.HTTPDeps) http.Handler {
	return router.CreateTestRouter(deps)
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
