package initialization

import (
	apiKeyCommands "TriePayments/internal/core/application/api_keys/commands"
	apiKeyQueries "TriePayments/internal/core/application/api_keys/queries"
	intentCommands "TriePayments/internal/core/application/intents/commands"
	intentQueries "TriePayments/internal/core/application/intents/queries"
	oauthCommands "TriePayments/internal/core/application/oauth/commands"
	async "TriePayments/internal/core/application/webhooks/asynq"
	webhooksCommands "TriePayments/internal/core/application/webhooks/commands"
	webhooksQueries "TriePayments/internal/core/application/webhooks/queries"
	workspaceCommands "TriePayments/internal/core/application/workspaces/commands"
	workspaceQueries "TriePayments/internal/core/application/workspaces/queries"
	"TriePayments/internal/core/domain"
	"TriePayments/internal/core/infrastructure"
	"TriePayments/internal/core/infrastructure/providers"
	apiKeysHandler "TriePayments/internal/core/interfaces/http/api_keys_handler"
	intents "TriePayments/internal/core/interfaces/http/intent_handler"
	"TriePayments/internal/core/interfaces/http/oauth_handler"
	webhooks "TriePayments/internal/core/interfaces/http/webhooks_handler"
	workspaces "TriePayments/internal/core/interfaces/http/workspaces_handler"
	"TriePayments/internal/interfaces/http/middleware"
	"TriePayments/internal/interfaces/http/router"
	"TriePayments/internal/interfaces/http/system"
	"TriePayments/internal/plataform/database"
	"TriePayments/internal/plataform/database/sqlc"
	"TriePayments/internal/plataform/telemetry"
	"TriePayments/internal/worker"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/go-co-op/gocron/v2"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type TriePayments struct {
	Port      string
	DB        *pgxpool.Pool
	Redis     *redis.Client
	Scheduler gocron.Scheduler
	GaClient  *goauth.Client

	Deps *router.HTTPDeps
}

func TriePaymentsSetup() *TriePayments {
	var app TriePayments

	LoadEnv(&app)
	SetupGoAuth(&app)
	SetupFUN()
	if viper.GetString("ENV") != "test" {
		SetupDB(&app, "./internal/plataform/database/migrations")
	} else {
		log.Println("WE'RE TESTING")
		SetupDB(&app, "../internal/plataform/database/migrations")
	}
	app.Redis = SetupRedis(15 * time.Second)
	SetupCron(app.DB, &app)

	return &app
}

func TriePaymentsStart(app *TriePayments, skipMux bool) {
	ctx := context.Background()

	defer app.DB.Close()
	defer app.Redis.Close()

	shutdown, err := telemetry.InitTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func(ctx context.Context) {
		err := shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(ctx)

	defer func() {
		err := app.Scheduler.StopJobs()
		if err != nil {
			log.Printf("Error stopping jobs: %v", err)
		}
		err = app.Scheduler.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	mpProvider, err := providers.NewMercadoPagoProvider(
		viper.GetString("MP_CLIENT_ID"),
		viper.GetString("MP_ACCESS_TOKEN"),
		viper.GetString("MP_REDIRECT_URI"), // https://triepayments.com/oauth/mercadopago/callback
	)
	if err != nil {
		log.Fatalf("Error creating mercado pago provider: %s", err.Error())
	}

	providerMap := map[string]domain.OAuthProvider{
		"mercadopago": mpProvider,
	}

	q := sqlc.New(app.DB)
	txRunner := database.NewPGXTxRunner(app.DB)
	tracer := otel.Tracer(string(telemetry.TriePaymentsTracer))
	logs := telemetry.Log()
	//ws := sockets.New()

	// Init Repos
	intentRepo := infrastructure.NewIntentsRepo(q, logs, tracer)
	workspaceRepo := infrastructure.NewWorkspaceRepo(q, logs, tracer)
	apiKeysRepo := infrastructure.NewApiKeyRepo(q, logs, tracer)
	endpointsRepo := infrastructure.NewWebhookEndpointRepo(q, logs, tracer)
	deliveriesRepo := infrastructure.NewWebhookDeliveryRepo(q, logs, tracer)
	oauthStatesRepo := infrastructure.NewOAuthStatesRepo(q, logs, tracer)
	providerCredentialsRepo := infrastructure.NewProviderCredentialsRepo(q, logs, tracer)

	authMW := middleware.NewAuthMiddleware(app.GaClient, apiKeysRepo, workspaceRepo, tracer)

	// init Async Handlers
	webhooksAsyncHandler := async.New(deliveriesRepo, app.GaClient, tracer, txRunner)
	server, asynqClient, scheduler, inspector, err := worker.InitAsynq(worker.Deps{
		WebhookAsynq: webhooksAsyncHandler,
	})
	defer func() {
		if err = inspector.Close(); err != nil {
			telemetry.Log().Error("error closing the asynq inspector", zap.Error(err))
		}
		scheduler.Shutdown()
		server.Shutdown()
		if err = asynqClient.Close(); err != nil {
			telemetry.Log().Error("error closing the asynq client", zap.Error(err))
		}
	}()

	// Init Commands and Queries
	intentC := intentCommands.New(intentRepo, workspaceRepo, app.GaClient, txRunner, tracer)
	intentQ := intentQueries.New(intentRepo, workspaceRepo, app.GaClient, txRunner, tracer)
	workspaceC := workspaceCommands.New(workspaceRepo, app.GaClient, txRunner, tracer)
	workspaceQ := workspaceQueries.New(workspaceRepo, app.GaClient, txRunner, tracer)
	apiKeyC := apiKeyCommands.New(apiKeysRepo, workspaceRepo, app.GaClient, txRunner, tracer)
	apiKeyQ := apiKeyQueries.New(apiKeysRepo, workspaceRepo, app.GaClient, txRunner, tracer)
	webhooksC := webhooksCommands.New(endpointsRepo, deliveriesRepo, workspaceRepo, intentRepo, asynqClient, app.GaClient, txRunner, tracer)
	webhooksQ := webhooksQueries.New(endpointsRepo, workspaceRepo, app.GaClient, txRunner, tracer)
	oauthC := oauthCommands.New(intentRepo, workspaceRepo, oauthStatesRepo, providerCredentialsRepo, providerMap, app.GaClient, txRunner, tracer)

	// Init Handlers
	systemHandler := system.NewSystemHandler()
	intentHandler := intents.NewIntentsHandler(intentC, intentQ)
	workspaceHandler := workspaces.NewWorkspacesHandler(workspaceC, workspaceQ)
	apiKeyHandler := apiKeysHandler.NewApiKeysHandler(apiKeyC, apiKeyQ)
	webhooksHandler := webhooks.NewWebhooksHandler(webhooksC, webhooksQ)
	oauthHandler := oauth_handler.NewOAuthHandler(oauthC)

	asynqmonHandler := asynqmon.New(asynqmon.Options{
		RootPath: "/admin/asynq",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	})

	deps := &router.HTTPDeps{
		SystemHandler:     systemHandler,
		IntentsHandler:    intentHandler,
		WorkspacesHandler: workspaceHandler,
		WebhooksHandler:   webhooksHandler,
		ApiKeysHandler:    apiKeyHandler,
		OauthHandler:      oauthHandler,
		AuthMiddleware:    authMW,
		AsynqmonHandler:   asynqmonHandler,
	}

	app.Deps = deps

	if !skipMux {
		mux := router.CreateRouter(deps)

		log.Printf("TriePayments listening on :%s", app.Port)
		log.Fatal(http.ListenAndServe(":"+app.Port, mux))
	}
}
