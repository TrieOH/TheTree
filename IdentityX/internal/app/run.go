package app

import (
	"IdentityX/internal/features/account"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/auth"
	"IdentityX/internal/features/projects"
	"IdentityX/internal/features/security"
	"IdentityX/internal/features/sessions"
	"IdentityX/internal/interfaces/http/middleware"
	"IdentityX/internal/interfaces/http/router"
	"IdentityX/internal/interfaces/http/system"
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/platform/email"
	"IdentityX/internal/platform/memory/redis"
	"IdentityX/internal/platform/telemetry"
	"IdentityX/internal/shared/ports"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type runtime struct {
	middlewares middlewares
	handlers    router.Handlers
	commands    commands
	queries     queries
	repos       repos
	repoQueries *sqlc.Queries
	txRunner    database.TxRunner
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
	security *security.CommandService
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

type middlewares struct {
	authMW *middleware.AuthMiddleware
}

func (app *IdentityX) run() {
	var rt runtime
	rt.repoQueries = sqlc.New(app.db)
	rt.txRunner = database.NewPGTxRunner(app.db)
	rt.tracer = otel.Tracer(string(telemetry.IdentityXTracer))
	rt.logger = telemetry.Log()
	rt.repos = app.startRepos(rt)
	rt.renderer, rt.mailer = email.NewMailPair(rt.logger, rt.tracer)
	rt.commands = app.startCommands(rt, rt.repos)
	rt.queries = app.startQueries(rt, rt.repos)
	rt.middlewares = app.startMiddlewares(rt)
	rt.handlers = app.startHandlers(rt)
	mux := router.CreateRouter(rt.handlers)
	port := viper.GetString("PORT")
	log.Printf("IdentityX listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func (app *IdentityX) startHandlers(rt runtime) router.Handlers {
	var h router.Handlers
	h.System = system.NewHandler()
	h.Users = auth.NewHandler(*rt.commands.auth, *rt.queries.auth)
	h.Accounts = account.NewHandler(*rt.commands.accounts)
	h.Projects = projects.NewHandler(*rt.commands.projects, *rt.queries.projects)
	h.Sessions = sessions.NewHandler(*rt.commands.sessions, *rt.queries.sessions)
	h.ApiKeys = api_keys.NewHandler(*rt.commands.apiKeys)
	h.AuthMW = *rt.middlewares.authMW
	return h
}

func (app *IdentityX) startCommands(rt runtime, r repos) commands {
	var cmd commands
	cmd.apiKeys = api_keys.NewCommandService(r.apiKeys, r.projects, rt.logger, rt.tracer, rt.txRunner)
	cmd.projects = projects.NewCommandService(r.projects, r.keys, rt.logger, rt.tracer, rt.txRunner)
	cmd.security = security.NewCommandService(r.sessions, r.projects, r.keys, r.apiKeys, rt.logger, rt.tracer, rt.txRunner)
	cmd.sessions = sessions.NewCommandService(r.sessions, r.keys, rt.logger, rt.tracer, rt.txRunner)
	cmd.auth = auth.NewCommandService(r.users, r.sessions, r.projects, r.keys, rt.renderer, rt.mailer, rt.logger, rt.tracer, rt.txRunner)
	cmd.accounts = account.NewCommandService(r.users, r.accounts, r.sessions, r.keys, r.tokenReuseList, rt.renderer, rt.mailer, rt.logger, rt.tracer, rt.txRunner)
	return cmd
}

func (app *IdentityX) startQueries(rt runtime, r repos) queries {
	var q queries
	q.auth = auth.NewQueryService(r.keys, rt.logger, rt.tracer, rt.txRunner)
	q.projects = projects.NewQueryService(r.projects, r.users, rt.logger, rt.tracer, rt.txRunner)
	q.sessions = sessions.NewQueryService(r.sessions, rt.logger, rt.tracer, rt.txRunner)
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

func (app *IdentityX) startMiddlewares(rt runtime) middlewares {
	var mw middlewares
	mw.authMW = middleware.NewAuthMiddleware(*rt.commands.security, rt.tracer, redis.NewCache(app.redis), viper.GetString("ISSUER"))
	return mw
}
