package app

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/features/account"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/auth"
	"IdentityX/internal/features/projects"
	"IdentityX/internal/features/security"
	"IdentityX/internal/features/sessions"
	"IdentityX/internal/platform/email"
	"IdentityX/internal/shared/feature_deps"
	"IdentityX/internal/shared/ports"
	"lib/database"
	"lib/errx"
	"lib/telemetry"
	"lib/xslices"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MintzyG/fun/middlewares"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type runtime struct {
	middlewares mws
	handlers    Handlers
	commands    commands
	queries     queries
	repos       repos
	repoQueries *sqlc.Queries
	tx          database.TxRunner
	tracer      trace.Tracer
	logger      *zap.Logger
	renderer    ports.EmailRenderer
	mailer      ports.Mailer
}

type commands struct {
	auth     *auth.CommandService
	accounts *account.CommandService
	sessions *sessions.CommandService
	projects *projects.CommandService
	apiKeys  *api_keys.CommandService
}

type queries struct {
	auth     *auth.QueryService
	projects *projects.QueryService
	sessions *sessions.QueryService
}

type repos struct {
	users          ports.UserRepository
	accounts       ports.AccountRepository
	sessions       ports.SessionRepository
	projects       ports.ProjectRepository
	keys           ports.KeysRepository
	tokenReuseList ports.TokenReuseListRepository
	apiKeys        ports.ApiKeyRepository
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

func (app *IdentityX) run() {
	var rt runtime
	rt.repoQueries = sqlc.New(app.db)
	rt.logger = telemetry.Log()
	rt.tx = database.NewPGXTxRunner(app.db, rt.logger)
	rt.tracer = otel.Tracer(app.cfg.AppName)
	rt.repos = app.startRepos(rt)
	rt.renderer, rt.mailer = email.NewMailPair(
		rt.logger,
		rt.tracer,
		app.cfg.AppUrl,
		app.cfg.SmtpHost,
		app.cfg.SmtpPort,
		app.cfg.SmtpUser,
		app.cfg.SmtpPass,
		app.cfg.SmtpFrom,
		app.cfg.SmtpTls,
		app.cfg.SmtpStartTls,
	)
	rt.commands = app.startCommands(rt, rt.repos)
	rt.queries = app.startQueries(rt, rt.repos)
	rt.middlewares = app.startMiddlewares(rt)
	rt.handlers = app.startHandlers(rt)
	mux := CreateRouter(rt.handlers, app.cfg.DebugMode, app.cfg.DisableRateLimit)
	port := app.cfg.Port
	log.Printf("IdentityX listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func (app *IdentityX) startHandlers(rt runtime) Handlers {
	var h Handlers
	h.Users = auth.NewHandler(*rt.commands.auth, *rt.queries.auth)
	h.Accounts = account.NewHandler(*rt.commands.accounts)
	h.Projects = projects.NewHandler(*rt.commands.projects, *rt.queries.projects)
	h.Sessions = sessions.NewHandler(*rt.commands.sessions, *rt.queries.sessions)
	h.ApiKeys = api_keys.NewHandler(*rt.commands.apiKeys)

	h.BodySize = rt.middlewares.bodySize
	h.RequestID = rt.middlewares.requestID
	h.Logger = rt.middlewares.logger
	h.Metrics = rt.middlewares.metrics
	h.CORS = rt.middlewares.cors
	h.RealIP = rt.middlewares.realIP
	h.Recover = rt.middlewares.recover
	h.Timeout = rt.middlewares.timeout
	h.RateLimit = rt.middlewares.ratelimit
	h.Jwt = rt.middlewares.jwt
	h.ApiKey = rt.middlewares.apiKey
	h.AnyAuth = rt.middlewares.anyAuth
	return h
}

func (app *IdentityX) startCommands(rt runtime, r repos) commands {
	var cmd commands
	cmd.apiKeys = api_keys.NewCommandService(feature_deps.ApiKeysCommandDeps{
		ApiKeys: r.apiKeys,
		Project: r.projects,
		Logger:  rt.logger,
		Tracer:  rt.tracer,
		Tx:      rt.tx,
	})
	cmd.projects = projects.NewCommandService(feature_deps.ProjectCommandDeps{
		KeyLifetime:   app.cfg.KeyLifetime,
		EncryptionKey: app.encryptionKey,
		Projects:      r.projects,
		Keys:          r.keys,
		Logger:        rt.logger,
		Tracer:        rt.tracer,
		Tx:            rt.tx,
	})
	cmd.sessions = sessions.NewCommandService(feature_deps.SessionCommandDeps{
		Sessions: r.sessions,
		Keys:     r.keys,
		Logger:   rt.logger,
		Tracer:   rt.tracer,
		Tx:       rt.tx,
	})
	cmd.auth = auth.NewCommandService(feature_deps.AuthCommandDeps{
		EncryptionKey: app.encryptionKey,
		Issuer:        app.cfg.Issuer,
		Users:         r.users,
		Sessions:      r.sessions,
		Projects:      r.projects,
		Keys:          r.keys,
		Renderer:      rt.renderer,
		Mailer:        rt.mailer,
		Logger:        rt.logger,
		Tracer:        rt.tracer,
		Tx:            rt.tx,
	})
	cmd.accounts = account.NewCommandService(feature_deps.AccountCommandDeps{
		EncryptionKey:  app.encryptionKey,
		Issuer:         app.cfg.Issuer,
		Users:          r.users,
		Accounts:       r.accounts,
		Sessions:       r.sessions,
		Keys:           r.keys,
		TokenReuseList: r.tokenReuseList,
		MailRenderer:   rt.renderer,
		MailSender:     rt.mailer,
		Logger:         rt.logger,
		Tracer:         rt.tracer,
		Tx:             rt.tx,
	})
	return cmd
}

func (app *IdentityX) startQueries(rt runtime, r repos) queries {
	var q queries
	q.auth = auth.NewQueryService(r.keys, rt.logger, rt.tracer, rt.tx)
	q.projects = projects.NewQueryService(r.projects, r.users, rt.logger, rt.tracer, rt.tx)
	q.sessions = sessions.NewQueryService(r.sessions, rt.logger, rt.tracer, rt.tx)
	return q
}

func (app *IdentityX) startRepos(rt runtime) repos {
	var r repos
	r.users = auth.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.accounts = account.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.sessions = sessions.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.projects = projects.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.keys = security.NewKeysRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.tokenReuseList = security.NewTokenReuseRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.apiKeys = api_keys.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	return r
}

func (app *IdentityX) startMiddlewares(rt runtime) mws {
	var mw mws
	authMW := SetupAuthMiddlewares(rt.repos.sessions, rt.repos.keys, rt.repos.apiKeys, rt.tracer, app.cfg.Issuer)
	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	mw.bodySize = middlewares.MaxBodySize(1 << 20)
	mw.requestID = middlewares.RequestID(middlewares.RequestIDConfig{Header: "X-Request-ID"})
	mw.logger = middlewares.Logs(middlewares.Config{Logger: rt.logger, SkipPrefixes: []string{"/metrics", "/health"}, RequestIDHeader: "X-Request-ID"})
	collectors, err := middlewares.NewCollectors(prometheus.DefaultRegisterer)
	if err != nil {
		errx.Exit(err, "Failed to create collectors")
	}
	mw.metrics = middlewares.Metrics(collectors, middlewares.MetricsConfig{SkipPrefixes: []string{"/metrics", "/health"}})
	mw.cors = middlewares.CORS(middlewares.CORSConfig{
		AllowedOrigins:   xslices.Clean(strings.Split(app.cfg.CorsAllowedOrigins, ",")),
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
