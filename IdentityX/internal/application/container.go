package application

import (
	"GoAuth/internal/adapters/email"
	"GoAuth/internal/adapters/memory/imc"
	"GoAuth/internal/adapters/memory/redis"
	"GoAuth/internal/adapters/persistence"
	"GoAuth/internal/application/apikey"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/application/authenticator"
	"GoAuth/internal/application/keys"
	"GoAuth/internal/application/permission"
	"GoAuth/internal/application/project"
	"GoAuth/internal/application/role"
	"GoAuth/internal/application/schema"
	"GoAuth/internal/application/schema_fields"
	"GoAuth/internal/application/schema_version"
	"GoAuth/internal/application/scope"
	"GoAuth/internal/application/session"
	"GoAuth/internal/application/subcontext"
	"GoAuth/internal/application/tokens"
	"GoAuth/internal/infrastructure"
	"GoAuth/internal/ports/inbounds"
	"time"

	"github.com/spf13/viper"
)

type Application struct {
	Auth           inbounds.AuthService
	Keys           inbounds.KeysService
	Project        inbounds.ProjectService
	Schema         inbounds.SchemaService
	SchemaVersions inbounds.SchemaVersionService
	SchemaFields   inbounds.SchemaFieldsService
	Session        inbounds.SessionService
	Authenticator  inbounds.RequestAuthenticator
	Permission     inbounds.PermissionService
	Role           inbounds.RoleService
	Scope          inbounds.ScopeService
	Verifier       inbounds.TokenVerifier
	ApiKey         inbounds.ApiKeyService
	SubContext     inbounds.SubContextService
}

func NewApplication(infra infrastructure.Infra) *Application {
	repos := persistence.NewRepositories(infra)

	cacheTTLStr := viper.GetString("KEYS_CACHE_TTL")
	cacheTTL, err := time.ParseDuration(cacheTTLStr)
	if err != nil {
		cacheTTL = time.Hour
	}

	privateCache := imc.NewInMemoryCache(100, cacheTTL)
	publicCache := imc.NewInMemoryCache(1000, cacheTTL)

	sharedCache := redis.NewRedisCache(infra.Redis)

	keyService := keys.New(repos.Keys, privateCache, publicCache)
	mailBundle := email.NewBundle(infra)
	tokensBundle := tokens.NewBundle(keyService)

	schemaService := schema.New(schema.Deps{
		Schemas:      repos.Schemas,
		Versions:     repos.SchemaVersions,
		Fields:       repos.SchemaFields,
		Projects:     repos.Projects,
		ProjectUsers: repos.ProjectUsers,
		Cache:        sharedCache,
	}, infra.Tx)

	authService := auth.New(auth.Deps{
		Users:          repos.Users,
		Sessions:       repos.Sessions,
		Schemas:        repos.Schemas,
		Versions:       repos.SchemaVersions,
		Fields:         repos.SchemaFields,
		Projects:       repos.Projects,
		ProjectUsers:   repos.ProjectUsers,
		Keys:           repos.Keys,
		TokenReuseList: repos.TokenReuseList,
		Cache:          sharedCache,
	}, infra, keyService, schemaService, tokensBundle, mailBundle)

	apiKeyService := apikey.New(apikey.Deps{
		ApiKey:  repos.ApiKey,
		Project: repos.Projects,
	})

	return &Application{
		Auth: authService,
		Keys: keyService,
		Project: project.New(
			repos.Projects,
			repos.ProjectUsers,
			repos.Scopes,
			repos.Keys,
			infra.Tx,
		),
		Schema: schemaService,
		SchemaVersions: schema_version.New(schema_version.Deps{
			Schemas:  repos.Schemas,
			Versions: repos.SchemaVersions,
			Fields:   repos.SchemaFields,
			Projects: repos.Projects,
			Cache:    sharedCache,
		}, infra.Tx),
		SchemaFields: schema_fields.New(schema_fields.Deps{
			Schemas:  repos.Schemas,
			Versions: repos.SchemaVersions,
			Fields:   repos.SchemaFields,
			Projects: repos.Projects,
		}, infra.Tx),
		Session: session.New(repos.Sessions, tokensBundle.Verifier, infra.Tx),
		Authenticator: authenticator.New(authenticator.Deps{
			Session:       repos.Sessions,
			TokenVerifier: tokensBundle.Verifier,
			ApiKey:        apiKeyService,
		}, infra.Tracer),
		Permission: permission.New(repos.Permissions, repos.Projects, repos.ProjectUsers, repos.Sessions, schemaService, infra.Tx),
		Role:       role.New(repos.Roles, repos.Permissions, repos.Projects, repos.ProjectUsers, repos.Sessions, infra.Tx),
		Scope:      scope.New(repos.Projects, repos.Scopes, infra.Tx),
		Verifier:   tokensBundle.Verifier,
		ApiKey:     apiKeyService,
		SubContext: subcontext.New(subcontext.Deps{
			Projects:     repos.Projects,
			ProjectUsers: repos.ProjectUsers,
		}, infra.Tx),
	}
}
