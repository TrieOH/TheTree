package app

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/features/actors"
	"IdentityX/internal/features/authn"
	"IdentityX/internal/features/blacklist"
	"IdentityX/internal/features/crypto_keys"
	"IdentityX/internal/features/organizations"
	"IdentityX/internal/features/platform_roles"
	"IdentityX/ports"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"log"
	"net/http"
	"strings"

	"github.com/MintzyG/fun/middlewares"
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
}

type queries struct {
	authn *authn.Queries
	orgs  *organizations.Queries
}

type commands struct {
	authn *authn.Commands
	orgs  *organizations.Commands
}

type mws struct {
	logger            func(http.Handler) http.Handler
	cors              func(http.Handler) http.Handler
	jwtAuth           func(http.Handler) http.Handler
	apiKeyAuth        func(http.Handler) http.Handler
	anyAuth           func(http.Handler) http.Handler
	clientOnly        func(http.Handler) http.Handler
	projectClientOnly func(http.Handler) http.Handler
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
	return r
}

func (app *IdentityX) startQueries(rt runtime, r repos) queries {
	var cmd queries
	cmd.authn = authn.NewQueries(ports.AuthnDeps{
		CryptoKeys: r.cryptoKeys,
		Logger:     rt.logger,
		Tracer:     rt.tracer,
		Tx:         rt.tx,
	})
	cmd.orgs = organizations.NewQueries(ports.OrganizationDeps{
		Actors: r.actors,
		Orgs:   r.orgs,
		Logger: rt.logger,
		Tracer: rt.tracer,
		Tx:     rt.tx,
	})
	return cmd
}

func (app *IdentityX) startCommands(rt runtime, r repos) commands {
	var cmd commands
	cmd.authn = authn.NewCommands(ports.AuthnDeps{
		Actors:             r.actors,
		PlatformRoles:      r.platformRoles,
		CryptoKeys:         r.cryptoKeys,
		Blacklist:          r.blacklist,
		ExternalIdentities: r.externalIdentities,
		Logger:             rt.logger,
		Tracer:             rt.tracer,
		Tx:                 rt.tx,
	})
	cmd.orgs = organizations.NewCommands(ports.OrganizationDeps{
		Actors: r.actors,
		Orgs:   r.orgs,
		Logger: rt.logger,
		Tracer: rt.tracer,
		Tx:     rt.tx,
	})
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
		Authn:             authn.NewHandlers(rt.commands.authn, rt.queries.authn),
		Orgs:              organizations.NewHandlers(rt.commands.orgs, rt.queries.orgs),
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
	//collectors, err := middlewares.NewCollectors(prometheus.DefaultRegisterer)
	//if err != nil {
	//	errx.Exit(err, "Failed to create collectors")
	//}
	//mw.metrics = middlewares.Metrics(collectors, middlewares.MetricsConfig{SkipPrefixes: []string{"/metrics", "/health"}})
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
