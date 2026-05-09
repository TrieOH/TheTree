package app

import (
	commands2 "Informd/internal/features/forms/commands"
	"Informd/internal/features/forms/handler"
	queries2 "Informd/internal/features/forms/queries"
	"Informd/internal/features/forms/repo"
	"Informd/internal/features/keys"
	"Informd/internal/features/namespaces"
	"Informd/internal/platform/database"
	"Informd/internal/platform/database/sqlc"
	"Informd/internal/platform/telemetry"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/errx"
	"Informd/internal/shared/ports"
	"context"
	authz2 "lib/authz"
	database2 "lib/database"
	"lib/xslices"
	"log"
	"net/http"
	"strings"
	"time"

	"git.trieoh.com/TrieOH/IdentityX-SDK-Go"
	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/middlewares"
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
	perms       authz2.Checker
}

type commands struct {
	namespaces *namespaces.CommandService
	apiKeys    *keys.CommandService
	forms      *commands2.CommandService
}

type queries struct {
	namespaces *namespaces.QueryService
	apiKeys    *keys.QueryService
	forms      *queries2.QueryService
}

type repos struct {
	namespaces ports.NamespaceRepo
	apiKeys    ports.ApiKeysRepo
	forms      ports.FormsRepo
	steps      ports.StepRepo
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
	rt.perms = authz2.NewChecker(app.sdbClient)
	rt.repoQueries = sqlc.New(app.db)
	rt.txRunner = database2.NewPGXTxRunner(app.db, rt.logger)
	rt.tracer = otel.Tracer(string(telemetry.InformdTracer))
	rt.repos = app.startRepos(rt)
	rt.middlewares = app.startMiddlewares(rt)
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
	handlers.ProjectsHandler = namespaces.NewHandler(rt.commands.namespaces, rt.queries.namespaces)
	handlers.ApiKeysHandler = keys.NewHandler(rt.commands.apiKeys, rt.queries.apiKeys)
	handlers.FormsHandler = handler.NewHandler(rt.commands.forms, rt.queries.forms)

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
	handlers.AppName = app.Config.AppName
	return &handlers
}

func (app *Informd) startCommands(rt runtime) commands {
	var cmd commands
	cmd.namespaces = namespaces.NewCommands(rt.repos.namespaces, rt.perms, rt.txRunner, rt.tracer)
	cmd.apiKeys = keys.NewCommands(rt.repos.apiKeys, rt.repos.namespaces, rt.perms, rt.txRunner, rt.tracer)
	cmd.forms = commands2.NewCommands(rt.repos.forms, rt.repos.steps, rt.repos.namespaces, rt.perms, rt.txRunner, rt.tracer)
	return cmd
}

func (app *Informd) startQueries(rt runtime) queries {
	var q queries
	q.namespaces = namespaces.NewQueries(rt.repos.namespaces, app.sdbClient, rt.txRunner, rt.tracer)
	q.apiKeys = keys.NewQueries(rt.repos.apiKeys, app.sdbClient, rt.txRunner, rt.tracer)
	q.forms = queries2.NewQueries(rt.repos.forms, rt.repos.steps, rt.repos.namespaces, app.sdbClient, rt.txRunner, rt.tracer)
	return q
}

func (app *Informd) startRepos(rt runtime) repos {
	var r repos
	r.namespaces = namespaces.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.apiKeys = keys.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.forms = repo.NewFormRepo(rt.repoQueries, app.db, rt.logger, rt.tracer)
	r.steps = repo.NewStepRepo(rt.repoQueries, rt.logger, rt.tracer)
	return r
}

func (app *Informd) startMiddlewares(rt runtime) mws {
	var mw mws

	keyFunc := func(ctx context.Context, tokenStr string) (*idx.AccessClaims, error) {
		return app.idxClient.Tokens.VerifyAccessToken(ctx, tokenStr)
	}

	jwtHook := func(ctx context.Context, claims *idx.AccessClaims) (context.Context, error) {
		return authz2.WithSubject(ctx, &authz2.UserSubject{
			ID:    claims.Sub.ID,
			Email: claims.Sub.Email,
		}), nil
	}

	apiKeyHook := func(ctx context.Context, rawKey string) (context.Context, error) {
		if len(rawKey) < 11 {
			return ctx, fun.ErrUnauthorized("invalid api key format")
		}
		prefix := rawKey[:11]

		candidates, err := rt.repos.apiKeys.GetByPrefix(ctx, prefix)
		if err != nil || len(candidates) == 0 {
			return ctx, fun.ErrUnauthorized("invalid api key")
		}

		var matched *contracts.APIKey
		for _, candidate := range candidates {
			if bcrypt.CompareHashAndPassword([]byte(candidate.KeyHash), []byte(rawKey)) == nil {
				matched = &candidate
				break
			}
		}
		if matched == nil {
			return ctx, fun.ErrUnauthorized("invalid api key")
		}
		return authz2.WithSubject(ctx, &authz2.UserSubject{ID: matched.OwnerID}), nil
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
		AllowedOrigins:   xslices.Clean(strings.Split(app.Config.CorsAllowedOrigins, ",")),
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
