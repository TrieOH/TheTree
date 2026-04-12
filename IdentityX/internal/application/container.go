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
	"GoAuth/internal/application/project"
	"GoAuth/internal/application/session"
	"GoAuth/internal/application/tokens"
	"GoAuth/internal/infrastructure"
	"GoAuth/internal/ports/inbounds"
	"time"

	"github.com/spf13/viper"
)

type Application struct {
	Auth          inbounds.AuthService
	Keys          inbounds.KeysService
	Project       inbounds.ProjectService
	Session       inbounds.SessionService
	Authenticator inbounds.RequestAuthenticator
	Verifier      inbounds.TokenVerifier
	ApiKey        inbounds.ApiKeyService
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

	authService := auth.New(auth.Deps{
		Users:          repos.Users,
		Sessions:       repos.Sessions,
		Projects:       repos.Projects,
		ProjectUsers:   repos.ProjectUsers,
		Keys:           repos.Keys,
		TokenReuseList: repos.TokenReuseList,
		Redis:          sharedCache,
	}, infra, keyService, tokensBundle, mailBundle)

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
			repos.Keys,
			infra.Tx,
		),
		Session: session.New(repos.Sessions, tokensBundle.Verifier, infra.Tx),
		Authenticator: authenticator.New(authenticator.Deps{
			Session:       repos.Sessions,
			TokenVerifier: tokensBundle.Verifier,
			ApiKey:        apiKeyService,
		}, infra.Tracer),
		Verifier: tokensBundle.Verifier,
		ApiKey:   apiKeyService,
	}
}
