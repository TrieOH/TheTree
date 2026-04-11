package app

import (
	"TrieForms/internal/features/forms"
	"TrieForms/internal/features/keys"
	"TrieForms/internal/features/projects"
	"TrieForms/internal/interfaces/http/middleware"
	"TrieForms/internal/interfaces/http/router"
	"TrieForms/internal/interfaces/http/system"
	"TrieForms/internal/platform/database"
	"TrieForms/internal/platform/database/sqlc"
	"TrieForms/internal/platform/queue"
	"TrieForms/internal/platform/telemetry"
	"TrieForms/internal/shared/ports"
	"log"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type runtime struct {
	middlewares middlewares
	handlers    *router.HTTPDeps
	commands    commands
	queries     queries
	repos       repos
	repoQueries *sqlc.Queries
	txRunner    database.TxRunner
	tracer      trace.Tracer
	logger      *zap.Logger
	asynq       asynqDeps
}

type commands struct {
	projects *projects.CommandService
	apiKeys  *keys.CommandService
	forms    *forms.CommandService
}

type queries struct {
	projects *projects.QueryService
	apiKeys  *keys.QueryService
	forms    *forms.QueryService
}

type repos struct {
	projects ports.ProjectsRepo
	apiKeys  ports.ApiKeysRepo
	forms    ports.FormsRepo
}

type middlewares struct {
	authMW *middleware.AuthMiddleware
}

type asynqDeps struct {
	client    *asynq.Client
	inspector *asynq.Inspector
	scheduler *asynq.Scheduler
	server    *asynq.Server
}

func (app *TrieForms) run() runtime {
	var rt runtime
	rt.repoQueries = sqlc.New(app.db)
	rt.txRunner = database.NewPGXTxRunner(app.db)
	rt.tracer = otel.Tracer(string(telemetry.TrieFormsTracer))
	rt.logger = telemetry.Log()
	rt.repos = app.startRepos(rt)
	rt.middlewares = app.startMiddlewares(rt)
	rt.asynq = app.startAsynq()
	defer app.stopAsynq(rt.asynq)
	rt.commands = app.startCommands(rt)
	rt.queries = app.startQueries(rt)
	rt.handlers = app.startHandlers(rt)
	mux := router.CreateRouter(rt.handlers)
	port := viper.GetString("port")
	log.Printf("TrieForms listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
	return rt
}

func (app *TrieForms) startHandlers(rt runtime) *router.HTTPDeps {
	var handlers router.HTTPDeps
	handlers.AsynqmonHandler = asynqmon.New(asynqmon.Options{
		RootPath: "/admin/asynq",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	})
	handlers.SystemHandler = system.NewSystemHandler(app.gaClient)
	handlers.ProjectsHandler = projects.NewProjectHandler(rt.commands.projects, rt.queries.projects)
	handlers.ApiKeysHandler = keys.NewApiKeysHandler(rt.commands.apiKeys, rt.queries.apiKeys)
	handlers.FormsHandler = forms.NewFormsHandler(rt.commands.forms, rt.queries.forms)
	return &handlers
}

func (app *TrieForms) startCommands(rt runtime) commands {
	var cmd commands
	cmd.projects = projects.NewProjectCommandService(rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	cmd.apiKeys = keys.NewApiKeyCommandService(rt.repos.apiKeys, rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	cmd.forms = forms.NewFormCommandService(rt.repos.forms, rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	return cmd
}

func (app *TrieForms) startQueries(rt runtime) queries {
	var q queries
	q.projects = projects.NewProjectQueryService(rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	q.apiKeys = keys.NewApiKeyQueryService(rt.repos.apiKeys, rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	q.forms = forms.NewFormQueryService(rt.repos.forms, rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	return q
}

func (app *TrieForms) startRepos(rt runtime) repos {
	var r repos
	r.projects = projects.NewProjectRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.apiKeys = keys.NewApiKeyRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.forms = forms.NewFormRepo(rt.repoQueries, rt.logger, rt.tracer)
	return r
}

func (app *TrieForms) startMiddlewares(rt runtime) middlewares {
	var mw middlewares
	mw.authMW = middleware.NewAuthMiddleware(app.gaClient, rt.repos.apiKeys, rt.repos.projects, rt.tracer)
	return mw
}

func (app *TrieForms) startAsynq() asynqDeps {
	var err error
	var deps asynqDeps
	deps.server, deps.client, deps.scheduler, deps.inspector, err = queue.InitAsynq(queue.Deps{})
	if err != nil {
		telemetry.Log().Fatal("failed to init Asynq", zap.Error(err))
	}
	return deps
}

func (app *TrieForms) stopAsynq(deps asynqDeps) {
	if err := deps.inspector.Close(); err != nil {
		telemetry.Log().Error("error closing the asynq inspector", zap.Error(err))
	}
	deps.scheduler.Shutdown()
	deps.server.Shutdown()
	if err := deps.client.Close(); err != nil {
		telemetry.Log().Error("error closing the asynq client", zap.Error(err))
	}
}
