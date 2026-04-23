package app

import (
	"log"
	"net/http"
	"payssage/internal/features/api_keys"
	"payssage/internal/features/intents"
	"payssage/internal/features/oauth"
	"payssage/internal/features/webhooks"
	"payssage/internal/features/workspaces"
	"payssage/internal/interfaces/http/middleware"
	"payssage/internal/interfaces/http/router"
	"payssage/internal/interfaces/http/system"
	"payssage/internal/platform/database"
	"payssage/internal/platform/database/sqlc"
	"payssage/internal/platform/providers"
	"payssage/internal/platform/queue"
	"payssage/internal/platform/telemetry"
	"payssage/internal/shared/ports"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type runtime struct {
	middlewares      middlewares
	handlers         *router.HTTPDeps
	commands         commands
	queries          queries
	repos            repos
	repoQueries      *sqlc.Queries
	txRunner         database.TxRunner
	tracer           trace.Tracer
	logger           *zap.Logger
	asynq            asynqDeps
	paymentProviders paymentProviders
}

type paymentProviders struct {
	oauth    map[string]ports.OAuthProvider
	payments map[string]ports.PaymentAbstractionLayer
}

type commands struct {
	webhooks   *webhooks.CommandService
	intents    *intents.CommandService
	workspaces *workspaces.CommandService
	apiKeys    *api_keys.CommandService
	oauth      *oauth.CommandService
}

type queries struct {
	webhooks   *webhooks.QueryService
	intents    *intents.QueryService
	workspaces *workspaces.QueryService
	apiKeys    *api_keys.QueryService
	oauth      *oauth.QueryService
}

type repos struct {
	intentRepo              ports.IntentRepository
	workspaceRepo           ports.WorkspaceRepo
	apiKeysRepo             ports.ApiKeysRepo
	endpointsRepo           ports.WebhookEndpointRepo
	deliveriesRepo          ports.WebhookDeliveryRepo
	eventsRepo              ports.WebhookEventRepo
	oauthStatesRepo         ports.OAuthStateRepo
	providerCredentialsRepo ports.ProviderCredentialRepo
	marketplaceRepo         ports.MarketplaceConfigRepo
}

type middlewares struct {
	authMW *middleware.AuthMiddleware
}

type asynqDeps struct {
	client    *asynq.Client
	inspector *asynq.Inspector
	scheduler *asynq.Scheduler
	server    *asynq.Server
}

func (app *Payssage) run() {
	var rt runtime
	rt.repoQueries = sqlc.New(app.db)
	rt.txRunner = database.NewPGXTxRunner(app.db)
	rt.tracer = otel.Tracer(string(telemetry.PayssageTracer))
	rt.logger = telemetry.Log()
	rt.repos = app.startRepos(rt)
	rt.middlewares = app.startMiddlewares(rt)
	rt.asynq = app.startAsynq()
	defer app.stopAsynq(rt.asynq)
	rt.paymentProviders = app.startPaymentProviders()
	rt.commands = app.startCommands(rt, rt.repos)
	rt.queries = app.startQueries(rt, rt.repos)
	rt.handlers = app.startHandlers(rt)
	mux := router.CreateRouter(rt.handlers)
	port := viper.GetString("port")
	log.Printf("payssage listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func (app *Payssage) startHandlers(rt runtime) *router.HTTPDeps {
	var handlers router.HTTPDeps
	handlers.AsynqmonHandler = asynqmon.New(asynqmon.Options{
		RootPath: "/admin/asynq",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	})
	handlers.AuthMiddleware = rt.middlewares.authMW
	handlers.SystemHandler = system.NewHandler()
	handlers.IntentsHandler = intents.NewHandler(rt.commands.intents, rt.queries.intents)
	handlers.WorkspacesHandler = workspaces.NewHandler(rt.commands.workspaces, rt.queries.workspaces)
	handlers.ApiKeysHandler = api_keys.NewHandler(rt.commands.apiKeys, rt.queries.apiKeys)
	handlers.WebhooksHandler = webhooks.NewHandler(rt.commands.webhooks, rt.queries.webhooks)
	handlers.OauthHandler = oauth.NewHandler(rt.commands.oauth, rt.queries.oauth)
	return &handlers
}

func (app *Payssage) startCommands(rt runtime, r repos) commands {
	var cmd commands
	cmd.webhooks = webhooks.NewCommandService(r.endpointsRepo, r.deliveriesRepo, r.eventsRepo, r.workspaceRepo, r.intentRepo, r.providerCredentialsRepo, rt.asynq.client, app.sdb, rt.txRunner, rt.tracer)
	cmd.intents = intents.NewCommandService(r.intentRepo, r.workspaceRepo, r.providerCredentialsRepo, r.marketplaceRepo, cmd.webhooks, rt.paymentProviders.oauth, rt.paymentProviders.payments, rt.txRunner, rt.tracer)
	cmd.workspaces = workspaces.NewCommandService(r.workspaceRepo, app.sdb, rt.txRunner, rt.tracer)
	cmd.apiKeys = api_keys.NewCommandService(r.apiKeysRepo, r.workspaceRepo, app.sdb, rt.txRunner, rt.tracer)
	cmd.oauth = oauth.NewCommandService(r.intentRepo, r.workspaceRepo, r.oauthStatesRepo, r.providerCredentialsRepo, r.marketplaceRepo, rt.paymentProviders.oauth, rt.txRunner, rt.tracer)

	return cmd
}

func (app *Payssage) startQueries(rt runtime, r repos) queries {
	var q queries
	q.webhooks = webhooks.NewQueryService(r.endpointsRepo, r.deliveriesRepo, r.eventsRepo, r.workspaceRepo, app.sdb, rt.txRunner, rt.tracer)
	q.intents = intents.NewQueryService(r.intentRepo, r.workspaceRepo, rt.txRunner, rt.tracer)
	q.workspaces = workspaces.NewQueryService(r.workspaceRepo, rt.txRunner, rt.tracer)
	q.apiKeys = api_keys.NewQueryService(r.apiKeysRepo, r.workspaceRepo, rt.txRunner, rt.tracer)
	q.oauth = oauth.NewQueryService(r.workspaceRepo, r.marketplaceRepo, rt.txRunner, rt.tracer)
	return q
}

func (app *Payssage) startPaymentProviders() paymentProviders {
	var pp paymentProviders
	mpProvider, err := providers.NewMercadoPagoProvider(
		viper.GetString("MP_CLIENT_ID"),
		viper.GetString("MP_ACCESS_TOKEN"),
		viper.GetString("MP_CLIENT_SECRET"),
		viper.GetString("MP_REDIRECT_URI"), // https://triepayments.com/oauth/mercadopago/callback
		viper.GetString("MP_WEBHOOK_SECRET"),
	)
	if err != nil {
		log.Fatalf("Error creating mercado pago provider: %s", err.Error())
	}

	pp.oauth = map[string]ports.OAuthProvider{
		"mercadopago": mpProvider,
	}

	pp.payments = map[string]ports.PaymentAbstractionLayer{
		"mercadopago": mpProvider,
	}

	return pp
}

func (app *Payssage) startRepos(rt runtime) repos {
	var r repos
	r.intentRepo = intents.NewIntentsRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.workspaceRepo = workspaces.NewWorkspaceRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.apiKeysRepo = api_keys.NewApiKeyRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.endpointsRepo = webhooks.NewWebhookEndpointRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.deliveriesRepo = webhooks.NewWebhookDeliveryRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.eventsRepo = webhooks.NewWebhookEventRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.oauthStatesRepo = oauth.NewOAuthStatesRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.providerCredentialsRepo = oauth.NewProviderCredentialsRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.marketplaceRepo = oauth.NewMarketplaceConfigRepo(rt.repoQueries, rt.logger, rt.tracer)
	return r
}

func (app *Payssage) startMiddlewares(rt runtime) middlewares {
	var mw middlewares
	mw.authMW = middleware.NewAuthMiddleware(app.ga, rt.repos.apiKeysRepo, rt.repos.workspaceRepo, rt.tracer)
	return mw
}

func (app *Payssage) startAsynq() asynqDeps {
	var err error
	var deps asynqDeps
	deps.server, deps.client, deps.scheduler, deps.inspector, err = queue.InitAsynq(queue.Deps{})
	if err != nil {
		telemetry.Log().Fatal("failed to init Asynq", zap.Error(err))
	}
	return deps
}

func (app *Payssage) stopAsynq(deps asynqDeps) {
	if err := deps.inspector.Close(); err != nil {
		telemetry.Log().Error("error closing the asynq inspector", zap.Error(err))
	}
	deps.scheduler.Shutdown()
	deps.server.Shutdown()
	if err := deps.client.Close(); err != nil {
		telemetry.Log().Error("error closing the asynq client", zap.Error(err))
	}
}
