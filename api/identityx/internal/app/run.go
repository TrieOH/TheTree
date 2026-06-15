package app

import (
	"lib/errx"
	"log"
	"net/http"
	"net/http/pprof"
	"strings"

	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/features/actors"
	"IdentityX/internal/features/authn"
	"IdentityX/internal/features/blacklist"
	"IdentityX/internal/features/crypto_keys"
	"IdentityX/internal/features/organizations"
	"IdentityX/internal/features/platform_roles"
	"IdentityX/internal/features/projects"
	"IdentityX/ports"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/MintzyG/fun/middlewares"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type runtime struct {
	sqlcQ  *sqlc.Queries
	tx     database.TxRunner
	tracer trace.Tracer
	logger *zap.Logger
	mws    mws

	repos    repos
	queries  queries
	commands commands
}

type repos struct {
	actors             ports.ActorRepo
	platformRoles      ports.PlatformRolesRepo
	cryptoKeys         ports.CryptoKeysRepo
	blacklist          ports.BlacklistRepo
	externalIdentities ports.ExternalIdentitiesRepo
	orgs               ports.OrganizationRepo
	projects           ports.ProjectRepo
}

type queries struct {
	authn    *authn.Queries
	orgs     *organizations.Queries
	projects *projects.Queries
}

type commands struct {
	authn    *authn.Commands
	orgs     *organizations.Commands
	projects *projects.Commands
}

type mws struct {
	logger            func(http.Handler) http.Handler
	cors              func(http.Handler) http.Handler
	jwtAuth           func(http.Handler) http.Handler
	apiKeyAuth        func(http.Handler) http.Handler
	anyAuth           func(http.Handler) http.Handler
	clientOnly        func(http.Handler) http.Handler
	projectClientOnly func(http.Handler) http.Handler
	metrics           func(http.Handler) http.Handler
}

func (app *IdentityX) run() {
	var rt runtime
	rt.sqlcQ = sqlc.New(app.db)
	rt.logger = telemetry.Log()
	rt.tx = database.NewPGXTxRunner(app.db, rt.logger)
	rt.tracer = otel.Tracer(app.cfg.AppName)
	rt.repos = app.startRepos(rt)
	rt.queries = app.startQueries(rt, rt.repos)
	rt.commands = app.startCommands(rt, rt.repos)
	rt.mws = app.startMiddlewares(rt)
	routerDeps := app.setupRouter(rt)
	mux := app.CreateRouter(routerDeps, app.cfg.DebugMode, app.cfg.DisableRateLimit)
	if app.cfg.ProfilePort != "" {
		go func() {
			pmux := http.NewServeMux()
			pmux.HandleFunc("/debug/pprof/", pprof.Index)
			pmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			pmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			pmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			pmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
			log.Printf("identityx pprof listening on :%s", app.cfg.ProfilePort)
			log.Println(http.ListenAndServe(":"+app.cfg.ProfilePort, pmux))
		}()
	}

	port := app.cfg.Port
	log.Printf("IdentityX listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func (app *IdentityX) startRepos(rt runtime) repos {
	var r repos
	r.actors = actors.NewRepo(rt.sqlcQ, rt.logger, rt.tracer)
	r.platformRoles = platform_roles.NewRepo(rt.sqlcQ, rt.logger, rt.tracer)
	r.cryptoKeys = crypto_keys.NewRepo(rt.sqlcQ, rt.logger, rt.tracer)
	r.blacklist = blacklist.NewRepo(rt.sqlcQ, rt.logger, rt.tracer)
	r.externalIdentities = authn.NewRepo(rt.sqlcQ, rt.logger, rt.tracer)
	r.orgs = organizations.NewRepo(rt.sqlcQ, rt.logger, rt.tracer)
	r.projects = projects.NewRepos(rt.sqlcQ, rt.logger, rt.tracer)
	return r
}

func (app *IdentityX) startQueries(rt runtime, r repos) queries {
	var cmd queries
	cmd.authn = authn.NewQueries(r.cryptoKeys, rt.logger, rt.tracer, rt.tx)
	cmd.orgs = organizations.NewQueries(r.projects, r.orgs, rt.logger, rt.tracer, rt.tx)
	cmd.projects = projects.NewQueries(r.projects, rt.logger, rt.tracer, rt.tx)
	return cmd
}

func (app *IdentityX) startCommands(rt runtime, r repos) commands {
	var cmd commands
	cmd.authn = authn.NewCommands(r.actors, r.projects, r.platformRoles, r.cryptoKeys, r.blacklist, r.externalIdentities, rt.logger, rt.tracer, rt.tx)
	cmd.orgs = organizations.NewCommands(r.projects, r.actors, r.orgs, rt.logger, rt.tracer, rt.tx)
	cmd.projects = projects.NewCommands(r.projects, r.actors, rt.logger, rt.tracer, rt.tx)
	return cmd
}

func (app *IdentityX) setupRouter(rt runtime) RouterDeps {
	return RouterDeps{
		AppName:           app.cfg.AppName,
		CORS:              rt.mws.cors,
		Logger:            rt.mws.logger,
		JwtAuth:           rt.mws.jwtAuth,
		ApiKeyAuth:        rt.mws.apiKeyAuth,
		AnyAuth:           rt.mws.anyAuth,
		ClientOnly:        rt.mws.clientOnly,
		ProjectClientOnly: rt.mws.projectClientOnly,
		Metrics:           rt.mws.metrics,
		Authn:             authn.NewHandlers(rt.commands.authn, rt.queries.authn),
		Orgs:              organizations.NewHandlers(rt.commands.orgs, rt.queries.orgs),
		Projects:          projects.NewHandlers(rt.commands.projects, rt.queries.projects),
	}
}

func (app *IdentityX) startMiddlewares(rt runtime) mws {
	var mw mws
	authMW := SetupAuthMiddlewares(rt)
	mw.jwtAuth = authMW.JWT()
	mw.apiKeyAuth = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	//mw.bodySize = middlewares.MaxBodySize(1 << 20)
	//mw.requestID = middlewares.RequestID(middlewares.RequestIDConfig{Header: "X-Request-ID"})
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
	//mw.realIP = middlewares.RealIP()
	//mw.recover = middlewares.Recover(rt.logger)
	//mw.timeout = middlewares.Timeout(60 * time.Second)
	//mw.ratelimit = middlewares.RateLimit(middlewares.RateLimitConfig{RPS: 400, Burst: 20,
	//	KeyExtractor: func(r *http.Request) string { return r.RemoteAddr },
	//})
	mw.clientOnly = ClientOnly()
	mw.projectClientOnly = ProjectClientOnly()
	return mw
}
