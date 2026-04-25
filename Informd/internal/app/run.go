package app

import (
	"Informd/internal/features/forms"
	"Informd/internal/features/keys"
	"Informd/internal/features/projects"
	"Informd/internal/platform/database"
	"Informd/internal/platform/database/sqlc"
	"Informd/internal/platform/queue"
	"Informd/internal/platform/telemetry"
	"Informd/internal/shared/authz"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/errx"
	"Informd/internal/shared/ports"
	"Informd/internal/shared/utils"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/MintzyG/FastUtilitiesNet/middlewares"
	idx "github.com/TrieOH/IdentityX-SDK-Go"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type runtime struct {
	middlewares mws
	handlers    *Deps
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

type mws struct {
	logger    func(http.Handler) http.Handler
	requestID func(http.Handler) http.Handler
	bodySize  func(http.Handler) http.Handler
	metrics   func(http.Handler) http.Handler
	cors      func(http.Handler) http.Handler
	realIP    func(http.Handler) http.Handler
	recover   func(http.Handler) http.Handler
	timeout   func(http.Handler) http.Handler
	ratelimit func(http.Handler) http.Handler
	jwt       func(http.Handler) http.Handler
	apiKey    func(http.Handler) http.Handler
	anyAuth   func(http.Handler) http.Handler
}

type asynqDeps struct {
	client    *asynq.Client
	inspector *asynq.Inspector
	scheduler *asynq.Scheduler
	server    *asynq.Server
}

func (app *Informd) run() runtime {
	var rt runtime
	rt.logger = telemetry.NewLogger(telemetry.LogConfig{
		Level:       "info",
		Development: false,
	})
	rt.repoQueries = sqlc.New(app.db)
	rt.txRunner = database.NewPGXTxRunner(app.db, rt.logger)
	rt.tracer = otel.Tracer(string(telemetry.InformdTracer))
	rt.repos = app.startRepos(rt)
	rt.middlewares = app.startMiddlewares(rt)
	rt.asynq = app.startAsynq()
	defer app.stopAsynq(rt.asynq)
	rt.commands = app.startCommands(rt)
	rt.queries = app.startQueries(rt)
	rt.handlers = app.startHandlers(rt)
	mux := CreateRouter(rt.handlers)
	log.Printf("Informd listening on :%s", app.Config.Port)
	log.Fatal(http.ListenAndServe(":"+app.Config.Port, mux))
	return rt
}

func (app *Informd) startHandlers(rt runtime) *Deps {
	var handlers Deps
	handlers.AsynqmonHandler = asynqmon.New(asynqmon.Options{
		RootPath: "/admin/asynq",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     app.Config.RedisAddr,
			Password: app.Config.RedisPassword,
			DB:       app.Config.RedisDB,
		},
	})
	handlers.ProjectsHandler = projects.NewProjectHandler(rt.commands.projects, rt.queries.projects)
	handlers.ApiKeysHandler = keys.NewApiKeysHandler(rt.commands.apiKeys, rt.queries.apiKeys)
	handlers.FormsHandler = forms.NewFormsHandler(rt.commands.forms, rt.queries.forms)

	handlers.BodySize = rt.middlewares.bodySize
	handlers.RequestID = rt.middlewares.requestID
	handlers.Logger = rt.middlewares.logger
	handlers.Metrics = rt.middlewares.metrics
	handlers.CORS = rt.middlewares.cors
	handlers.RealIP = rt.middlewares.realIP
	handlers.Recover = rt.middlewares.recover
	handlers.Timeout = rt.middlewares.timeout
	handlers.RateLimit = rt.middlewares.ratelimit
	handlers.Jwt = rt.middlewares.jwt
	handlers.ApiKey = rt.middlewares.apiKey
	handlers.AnyAuth = rt.middlewares.anyAuth
	return &handlers
}

func (app *Informd) startCommands(rt runtime) commands {
	var cmd commands
	cmd.projects = projects.NewProjectCommandService(rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	cmd.apiKeys = keys.NewApiKeyCommandService(rt.repos.apiKeys, rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	cmd.forms = forms.NewFormCommandService(rt.repos.forms, rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	return cmd
}

func (app *Informd) startQueries(rt runtime) queries {
	var q queries
	q.projects = projects.NewQueryService(rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	q.apiKeys = keys.NewApiKeyQueryService(rt.repos.apiKeys, rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	q.forms = forms.NewFormQueryService(rt.repos.forms, rt.repos.projects, app.sdbClient, rt.txRunner, rt.tracer)
	return q
}

func (app *Informd) startRepos(rt runtime) repos {
	var r repos
	r.projects = projects.NewProjectRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.apiKeys = keys.NewApiKeyRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.forms = forms.NewFormRepo(rt.repoQueries, rt.logger, rt.tracer)
	return r
}

func (app *Informd) startMiddlewares(rt runtime) mws {
	var mw mws

	keyFunc := func(ctx context.Context, tokenStr string) (*idx.AccessClaims, error) {
		return app.idxClient.Tokens.VerifyAccessToken(ctx, tokenStr)
	}

	jwtHook := func(ctx context.Context, claims *idx.AccessClaims) (context.Context, error) {
		return authz.WithSubject(ctx, &authz.UserSubject{
			ID:    claims.Sub.ID,
			Email: claims.Sub.Email,
		}), nil
	}

	apiKeyHook := func(ctx context.Context, rawKey string) (context.Context, error) {
		if len(rawKey) < 11 {
			return ctx, errors.New("invalid api key format")
		}
		prefix := rawKey[:11]

		candidates, err := rt.repos.apiKeys.GetByPrefix(ctx, prefix)
		if err != nil || len(candidates) == 0 {
			return ctx, errors.New("invalid api key")
		}

		var matched *contracts.APIKey
		for _, candidate := range candidates {
			if bcrypt.CompareHashAndPassword([]byte(candidate.KeyHash), []byte(rawKey)) == nil {
				matched = &candidate
				break
			}
		}
		if matched == nil {
			return ctx, errors.New("invalid api key")
		}

		project, err := rt.repos.projects.GetByID(ctx, matched.ProjectID)
		if err != nil {
			return ctx, errors.New("workspace not found")
		}

		return authz.WithProject(ctx, project), nil
	}

	authMW := middlewares.New[*idx.AccessClaims](keyFunc, jwtHook, apiKeyHook)
	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	mw.bodySize = middlewares.MaxBodySize(1 << 20)
	mw.requestID = middlewares.RequestID(middlewares.RequestIDConfig{Header: "X-Request-ID"})
	mw.logger = middlewares.Logs(middlewares.Config{Logger: rt.logger, SkipPrefixes: []string{"/admin/asynq"}, RequestIDHeader: "X-Request-ID"})
	collectors, err := middlewares.NewCollectors(prometheus.DefaultRegisterer)
	if err != nil {
		errx.Must(err, "Failed to create collectors")
	}
	mw.metrics = middlewares.Metrics(collectors, middlewares.MetricsConfig{SkipPrefixes: []string{"/metrics", "/health"}})
	mw.cors = middlewares.CORS(middlewares.CORSConfig{
		AllowedOrigins:   utils.SplitAndCleanCSV(app.Config.CorsAllowedOrigins),
		AllowCredentials: true,
	})
	mw.realIP = middlewares.RealIP()
	mw.recover = middlewares.Recover(rt.logger)
	mw.timeout = middlewares.Timeout(60 * time.Second)
	mw.ratelimit = middlewares.RateLimit(middlewares.RateLimitConfig{RPS: 400, Burst: 20,
		KeyExtractor: func(r *http.Request) string { return r.RemoteAddr },
	})
	return mw
}

func (app *Informd) startAsynq() asynqDeps {
	var err error
	var deps asynqDeps
	deps.server, deps.client, deps.scheduler, deps.inspector, err = queue.InitAsynq(queue.Deps{
		RedisAddr:     app.Config.RedisAddr,
		RedisPassword: app.Config.RedisPassword,
		RedisDB:       app.Config.RedisDB,
	})
	if err != nil {
		errx.Must(err, "failed to init Asynq")
	}
	return deps
}

func (app *Informd) stopAsynq(deps asynqDeps) {
	if err := deps.inspector.Close(); err != nil {
		errx.Must(err, "error closing the asynq inspector")
	}
	deps.scheduler.Shutdown()
	deps.server.Shutdown()
	if err := deps.client.Close(); err != nil {
		errx.Must(err, "error closing the asynq client")
	}
}
