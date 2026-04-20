package app

import (
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/auth"
	"IdentityX/internal/features/keys"
	"IdentityX/internal/features/projects"
	"IdentityX/internal/features/sessions"
	"IdentityX/internal/features/tokens"
	"IdentityX/internal/features/users"
	"IdentityX/internal/interfaces/http/middleware"
	"IdentityX/internal/interfaces/http/router"
	"IdentityX/internal/interfaces/http/system"
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/platform/memory/imc"
	"IdentityX/internal/platform/memory/redis"
	"IdentityX/internal/platform/telemetry"
	"IdentityX/internal/shared/ports"
	"log"
	"net/http"
	"time"

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
	caches      caches
	renderer    ports.EmailRenderer
	mailer      ports.Mailer
}

type commands struct {
	users    *users.CommandService
	sessions *sessions.CommandService
	projects *projects.CommandService
	keys     *keys.CommandService
	tokens   *tokens.CommandService
	apiKeys  *api_keys.CommandService
	auth     *auth.CommandService
}

type queries struct {
}

type repos struct {
	users          ports.UserRepository
	sessions       ports.SessionRepository
	projects       ports.ProjectRepository
	keys           ports.KeysRepository
	tokenReuseList ports.TokenReuseListRepository
	apiKeys        ports.ApiKeyRepository
}

type middlewares struct {
	authMW *middleware.AuthMiddleware
}

type caches struct {
	publicCache  *imc.InMemoryCache
	privateCache *imc.InMemoryCache
	sharedCache  *redis.Cache
}

func (app *IdentityX) run() {
	var rt runtime
	rt.repoQueries = sqlc.New(app.db)
	rt.txRunner = database.NewPGTxRunner(app.db)
	rt.tracer = otel.Tracer(string(telemetry.IdentityXTracer))
	rt.logger = telemetry.Log()
	rt.repos = app.startRepos(rt)
	rt.caches = app.startCaches()
	rt.renderer, rt.mailer = NewBundle(rt)
	rt.commands = app.startCommands(rt, rt.repos)
	rt.middlewares = app.startMiddlewares(rt)
	rt.handlers = app.startHandlers(rt)
	//rt.queries = app.startQueries(rt)
	mux := router.CreateRouter(rt.handlers)
	port := viper.GetString("PORT")
	log.Printf("IdentityX listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func (app *IdentityX) startHandlers(rt runtime) router.Handlers {
	var h router.Handlers
	h.System = system.NewHandler()
	h.Users = users.NewHandler(*rt.commands.users, redis.NewCache(app.redis))
	h.Projects = projects.NewHandler(*rt.commands.projects)
	h.Sessions = sessions.NewHandler(*rt.commands.sessions, redis.NewCache(app.redis))
	h.ApiKeys = api_keys.NewHandler(*rt.commands.apiKeys)
	h.AuthMW = *rt.middlewares.authMW
	return h
}

func (app *IdentityX) startCommands(rt runtime, r repos) commands {
	var cmd commands
	cmd.apiKeys = api_keys.NewCommandService(r.apiKeys, r.projects, rt.logger, rt.tracer, rt.txRunner)
	cmd.keys = keys.NewCommandService(r.keys, rt.caches.privateCache, rt.caches.publicCache, rt.logger, rt.tracer, rt.txRunner)
	cmd.tokens = tokens.NewCommandService(*cmd.keys)
	cmd.projects = projects.NewCommandService(r.users, r.projects, r.keys, rt.logger, rt.tracer, rt.txRunner)
	cmd.auth = auth.NewCommandService(r.sessions, r.projects, *cmd.tokens, r.apiKeys, rt.logger, rt.tracer, rt.txRunner)
	cmd.sessions = sessions.NewCommandService(r.sessions, *cmd.tokens, rt.logger, rt.tracer, rt.txRunner)
	cmd.users = users.NewCommandService(r.users, r.sessions, r.projects, r.keys, r.tokenReuseList, redis.NewCache(app.redis), *cmd.keys, *cmd.tokens, rt.renderer, rt.mailer, rt.logger, rt.tracer, rt.txRunner)
	return cmd
}

/*
func (app *IdentityX) startQueries(rt runtime) queries {
	var q queries
	q.users = users.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	q.sessions = sessions.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	q.projects = projects.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	q.projectUsers = project_users.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	q.keys = keys.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	q.tokenReuseList = tokens.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	q.apiKeys = api_keys.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	return q
}
*/

func (app *IdentityX) startRepos(rt runtime) repos {
	var r repos
	r.users = users.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.sessions = sessions.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.projects = projects.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.keys = keys.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.tokenReuseList = tokens.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.apiKeys = api_keys.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	return r
}

func (app *IdentityX) startMiddlewares(rt runtime) middlewares {
	var mw middlewares
	mw.authMW = middleware.NewAuthMiddleware(*rt.commands.auth, rt.tracer, redis.NewCache(app.redis), viper.GetString("ISSUER"))
	return mw
}

func (app *IdentityX) startCaches() caches {
	var c caches
	cacheTTLStr := viper.GetString("KEYS_CACHE_TTL")
	cacheTTL, err := time.ParseDuration(cacheTTLStr)
	if err != nil {
		cacheTTL = time.Hour
	}

	c.privateCache = imc.NewInMemoryCache(100, cacheTTL)
	c.publicCache = imc.NewInMemoryCache(1000, cacheTTL)
	c.sharedCache = redis.NewCache(app.redis)
	return c
}
