package app

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/features/actors"
	apikeys "IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/authn"
	"IdentityX/internal/features/blacklist"
	"IdentityX/internal/features/capabilities"
	"IdentityX/internal/features/crypto_keys"
	"IdentityX/internal/features/organizations"
	"IdentityX/internal/features/platform_roles"
	"IdentityX/internal/features/profile_schemas"
	"IdentityX/internal/features/profiles"
	"IdentityX/internal/features/projects"
	"IdentityX/ports"
	"lib/database"
	"lib/errx"
	"lib/xslices"
	"net/http"
	"strings"

	mws "github.com/MintzyG/fun/middlewares"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ── Wire types ────────────────────────────────────────────────────────────

type repos struct {
	actors             ports.ActorRepo
	apiKeys            ports.APIKeysRepo
	capabilities       ports.CapabilityRepo
	platformRoles      ports.PlatformRolesRepo
	cryptoKeys         ports.CryptoKeysRepo
	blacklist          ports.BlacklistRepo
	externalIdentities ports.ExternalIdentitiesRepo
	orgs               ports.OrganizationRepo
	projects           ports.ProjectRepo
	profileSchemas     ports.ProfileSchemaRepo
	profiles           ports.ProfileRepo
}

type queries struct {
	authn          *authn.Queries
	orgs           *organizations.Queries
	projects       *projects.Queries
	actors         *actors.Queries
	capabilities   *capabilities.Queries
	profiles       *profiles.Queries
	profileSchemas *profile_schemas.Queries
}

type commands struct {
	authn          *authn.Commands
	actors         *actors.Commands
	apiKeys        *apikeys.Commands
	capabilities   *capabilities.Commands
	orgs           *organizations.Commands
	projects       *projects.Commands
	profiles       *profiles.Commands
	profileSchemas *profile_schemas.Commands
}

type middlewares struct {
	logger            func(http.Handler) http.Handler
	cors              func(http.Handler) http.Handler
	jwtAuth           func(http.Handler) http.Handler
	apiKeyAuth        func(http.Handler) http.Handler
	anyAuth           func(http.Handler) http.Handler
	clientOnly        func(http.Handler) http.Handler
	projectClientOnly func(http.Handler) http.Handler
	metrics           func(http.Handler) http.Handler
}

type handlers struct {
	Actors         *actors.Handlers
	APIKeys        *apikeys.Handlers
	Authn          *authn.Handlers
	Orgs           *organizations.Handlers
	Projects       *projects.Handlers
	Capabilities   *capabilities.Handlers
	Profiles       *profiles.Handlers
	ProfileSchemas *profile_schemas.Handlers
}

// ── Init functions ────────────────────────────────────────────────────────

func (app *IdentityX) initRepos(q *sqlc.Queries, logger *zap.Logger, tracer trace.Tracer) *repos {
	return &repos{
		actors:             actors.NewRepo(q, logger, tracer),
		apiKeys:            apikeys.NewRepo(q, logger, tracer),
		capabilities:       capabilities.NewRepos(q, logger, tracer),
		platformRoles:      platform_roles.NewRepo(q, logger, tracer),
		cryptoKeys:         crypto_keys.NewRepo(q, logger, tracer),
		blacklist:          blacklist.NewRepo(q, logger, tracer),
		externalIdentities: authn.NewRepo(q, logger, tracer),
		orgs:               organizations.NewRepo(q, logger, tracer),
		projects:           projects.NewRepos(q, logger, tracer),
		profileSchemas:     profile_schemas.NewRepo(q, logger, tracer),
		profiles:           profiles.NewRepo(q, logger, tracer),
	}
}

func (app *IdentityX) initQueries(r *repos, tx database.TxRunner, logger *zap.Logger, tracer trace.Tracer) queries {
	return queries{
		actors:         actors.NewQueries(r.projects, r.actors, logger, tracer, tx),
		authn:          authn.NewQueries(r.cryptoKeys, logger, tracer, tx),
		orgs:           organizations.NewQueries(r.projects, r.actors, r.orgs, logger, tracer, tx),
		projects:       projects.NewQueries(r.projects, logger, tracer, tx),
		capabilities:   capabilities.NewQueries(r.capabilities, r.projects, logger, tracer, tx),
		profiles:       profiles.NewQueries(r.profiles, r.projects, logger, tracer, tx),
		profileSchemas: profile_schemas.NewQueries(r.profileSchemas, r.projects, logger, tracer, tx),
	}
}

func (app *IdentityX) initCommands(r *repos, tx database.TxRunner, logger *zap.Logger, tracer trace.Tracer) commands {
	return commands{
		authn:          authn.NewCommands(r.actors, r.projects, r.platformRoles, r.cryptoKeys, r.blacklist, r.externalIdentities, logger, tracer, tx),
		actors:         actors.NewCommands(r.actors, r.projects, logger, tracer, tx),
		apiKeys:        apikeys.NewCommands([]byte(app.cfg.HmacSecret), r.actors, r.apiKeys, r.capabilities, r.projects, logger, tracer, tx),
		orgs:           organizations.NewCommands(r.projects, r.actors, r.orgs, logger, tracer, tx),
		projects:       projects.NewCommands(r.projects, r.actors, logger, tracer, tx),
		capabilities:   capabilities.NewCommands(r.actors, r.capabilities, r.projects, logger, tracer, tx),
		profiles:       profiles.NewCommands(r.profiles, r.profileSchemas, r.projects, logger, tracer, tx),
		profileSchemas: profile_schemas.NewCommands(r.profileSchemas, r.projects, logger, tracer, tx),
	}
}

func (app *IdentityX) initMiddlewares(r *repos, logger *zap.Logger, cfg Config) middlewares {
	var mw middlewares
	authMW := app.SetupAuthMiddlewares(r.cryptoKeys, r.apiKeys, r.actors, r.capabilities, logger)
	mw.jwtAuth = authMW.JWT()
	mw.apiKeyAuth = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	// mw.bodySize = mws.MaxBodySize(1 << 20)
	// mw.requestID = mws.RequestID(mws.RequestIDConfig{Header: "X-Request-ID"})
	mw.logger = mws.Logs(mws.Config{Logger: logger, SkipPrefixes: []string{"/metrics", "/health"}, RequestIDHeader: "X-Request-ID"})
	collectors, err := mws.NewCollectors(prometheus.DefaultRegisterer)
	if err != nil {
		errx.Exit(err, "Failed to create collectors")
	}
	mw.metrics = mws.Metrics(collectors, mws.MetricsConfig{SkipPrefixes: []string{"/metrics", "/health"}})
	mw.cors = mws.CORS(mws.CORSConfig{
		AllowedOrigins:   xslices.Clean(strings.Split(cfg.CorsAllowedOrigins, ",")),
		AllowedHeaders:   xslices.Clean(strings.Split(cfg.CorsAllowedHeaders, ",")),
		AllowCredentials: true,
	})
	// mw.realIP = mws.RealIP()
	// mw.recover = mws.Recover(logger)
	// mw.timeout = mws.Timeout(60 * time.Second)
	// mw.ratelimit = mws.RateLimit(mws.RateLimitConfig{RPS: 400, Burst: 20,
	//	 KeyExtractor: func(r *http.Request) string { return r.RemoteAddr },
	// })
	mw.clientOnly = ClientOnly()
	mw.projectClientOnly = ProjectClientOnly()
	return mw
}

func (app *IdentityX) initHandlers(q queries, c commands) handlers {
	return handlers{
		Actors:         actors.NewHandlers(q.actors, c.actors),
		APIKeys:        apikeys.NewHandlers(c.apiKeys),
		Authn:          authn.NewHandlers(c.authn, q.authn),
		Orgs:           organizations.NewHandlers(c.orgs, q.orgs),
		Projects:       projects.NewHandlers(c.projects, q.projects),
		Capabilities:   capabilities.NewHandlers(c.capabilities, q.capabilities),
		Profiles:       profiles.NewHandlers(q.profiles, c.profiles),
		ProfileSchemas: profile_schemas.NewHandlers(q.profileSchemas, c.profileSchemas),
	}
}
