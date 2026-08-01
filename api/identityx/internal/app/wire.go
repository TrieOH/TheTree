package app

import (
	"IdentityX/internal/authz"
	"IdentityX/internal/features/actors"
	actorsHandlers "IdentityX/internal/features/actors/handlers"
	actorsRepos "IdentityX/internal/features/actors/repos"
	apikeys "IdentityX/internal/features/api_keys"
	apikeysHandlers "IdentityX/internal/features/api_keys/handlers"
	apikeysRepos "IdentityX/internal/features/api_keys/repos"
	"IdentityX/internal/features/authn"
	authnHandlers "IdentityX/internal/features/authn/handlers"
	authnRepos "IdentityX/internal/features/authn/repos"
	blacklistRepos "IdentityX/internal/features/blacklist/repos"
	"IdentityX/internal/features/capabilities"
	capabilitiesHandlers "IdentityX/internal/features/capabilities/handlers"
	capabilitiesRepos "IdentityX/internal/features/capabilities/repos"
	cryptoKeysRepos "IdentityX/internal/features/crypto_keys/repos"
	"IdentityX/internal/features/organizations"
	orgsHandlers "IdentityX/internal/features/organizations/handlers"
	orgsRepos "IdentityX/internal/features/organizations/repos"
	platformRolesRepos "IdentityX/internal/features/platform_roles/repos"
	"IdentityX/internal/features/profile_schemas"
	profileSchemasHandlers "IdentityX/internal/features/profile_schemas/handlers"
	profileSchemasRepos "IdentityX/internal/features/profile_schemas/repos"
	"IdentityX/internal/features/profiles"
	profilesHandlers "IdentityX/internal/features/profiles/handlers"
	profilesRepos "IdentityX/internal/features/profiles/repos"
	"IdentityX/internal/features/projects"
	projectsHandlers "IdentityX/internal/features/projects/handlers"
	projectsRepos "IdentityX/internal/features/projects/repos"
	"IdentityX/internal/sqlc"
	"IdentityX/ports"
	"net/http"
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

type operations struct {
	authn          *authn.Operations
	orgs           *organizations.Operations
	projects       *projects.Operations
	actors         *actors.Operations
	apiKeys        *apikeys.Operations
	capabilities   *capabilities.Operations
	profiles       *profiles.Operations
	profileSchemas *profile_schemas.Operations
}

type middlewares struct {
	jwtAuth           func(http.Handler) http.Handler
	apiKeyAuth        func(http.Handler) http.Handler
	anyAuth           func(http.Handler) http.Handler
	clientOnly        func(http.Handler) http.Handler
	projectClientOnly func(http.Handler) http.Handler
}

type handlers struct {
	Actors         *actorsHandlers.Handlers
	APIKeys        *apikeysHandlers.Handlers
	Authn          *authnHandlers.Handlers
	Orgs           *orgsHandlers.Handlers
	Projects       *projectsHandlers.Handlers
	Capabilities   *capabilitiesHandlers.Handlers
	Profiles       *profilesHandlers.Handlers
	ProfileSchemas *profileSchemasHandlers.Handlers
}

// ── Init functions ────────────────────────────────────────────────────────

func (app *IdentityX) initRepos(q *sqlc.Queries) *repos {
	return &repos{
		actors:             actorsRepos.NewRepo(q),
		apiKeys:            apikeysRepos.NewRepo(q),
		capabilities:       capabilitiesRepos.NewRepo(q),
		platformRoles:      platformRolesRepos.NewRepo(q),
		cryptoKeys:         cryptoKeysRepos.NewRepo(q),
		blacklist:          blacklistRepos.NewRepo(q),
		externalIdentities: authnRepos.NewRepo(q),
		orgs:               orgsRepos.NewRepo(q),
		projects:           projectsRepos.NewRepo(q),
		profileSchemas:     profileSchemasRepos.NewSchemaRepo(q),
		profiles:           profilesRepos.NewProfileRepo(q),
	}
}

func (app *IdentityX) initOperations(r *repos) operations {
	authzSvc := authz.New(r.orgs, r.projects)
	return operations{
		authn:          authn.NewOperations(r.actors, r.projects, r.platformRoles, r.cryptoKeys, r.blacklist, r.externalIdentities),
		orgs:           organizations.NewOperations(r.projects, r.actors, r.orgs, authzSvc),
		projects:       projects.NewOperations(r.cryptoKeys, r.projects, r.actors, authzSvc),
		actors:         actors.NewOperations(r.actors, r.projects, authzSvc),
		apiKeys:        apikeys.NewOperations([]byte(app.cfg.HmacSecret), r.actors, r.apiKeys, r.capabilities, r.projects, authzSvc),
		capabilities:   capabilities.NewOperations(r.actors, r.capabilities, r.projects, authzSvc),
		profiles:       profiles.NewOperations(r.profiles, r.profileSchemas, r.projects, authzSvc),
		profileSchemas: profile_schemas.NewOperations(r.profileSchemas, r.projects, authzSvc),
	}
}

func (app *IdentityX) initMiddlewares(r *repos) middlewares {
	var mw middlewares
	authMW := app.SetupAuthMiddlewares(r.cryptoKeys, r.apiKeys, r.actors, r.capabilities)
	mw.jwtAuth = authMW.JWT()
	mw.apiKeyAuth = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	mw.clientOnly = ClientOnly()
	mw.projectClientOnly = ProjectClientOnly()
	return mw
}

func (app *IdentityX) initHandlers(ops operations) handlers {
	return handlers{
		Actors:         actorsHandlers.NewHandlers(ops.actors),
		APIKeys:        apikeysHandlers.NewHandlers(ops.apiKeys),
		Authn:          authnHandlers.NewHandlers(ops.authn),
		Orgs:           orgsHandlers.NewHandlers(ops.orgs),
		Projects:       projectsHandlers.NewHandlers(ops.projects),
		Capabilities:   capabilitiesHandlers.NewHandlers(ops.capabilities),
		Profiles:       profilesHandlers.New(ops.profiles),
		ProfileSchemas: profileSchemasHandlers.New(ops.profileSchemas),
	}
}
