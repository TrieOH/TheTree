package app

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"

	"lib/database"
	libriver "lib/river"
	"lib/telemetry"
	"payssage/internal/database/sqlc"
	"payssage/internal/features/api_keys"
	"payssage/internal/features/intents"
	"payssage/internal/features/oauth"
	"payssage/internal/features/webhooks"
	"payssage/internal/features/workspaces"
	"payssage/internal/jobs"
	"payssage/internal/platform/providers"
	"payssage/ports"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

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

type middlewares struct {
	authMW *AuthMiddleware
}

func (app *Payssage) run() {
	ctx := context.Background()
	var rt runtime
	rt.repoQueries = sqlc.New(app.db)
	rt.logger = telemetry.Log()
	rt.txRunner = database.NewPGXTxRunner(app.db, rt.logger)
	rt.tracer = otel.Tracer("Payssage")
	rt.repos = app.startRepos(rt)
	rt.middlewares = app.startMiddlewares(rt)

	app.river = libriver.NewClient(app.db, libriver.NewWorkers(
		libriver.Register[jobs.DeliverWebhookArgs](jobs.NewDeliverWebhookWorker(rt.repos.deliveriesRepo)),
	), nil, nil)
	if err := app.river.Start(ctx); err != nil {
		telemetry.Log().Fatal("failed to start river client", zap.Error(err))
	}
	defer app.river.Stop(ctx)

	rt.paymentProviders = app.startPaymentProviders()
	rt.commands = app.startCommands(rt, rt.repos)
	rt.queries = app.startQueries(rt, rt.repos)
	rt.handlers = app.startHandlers(rt)
	mux := CreateRouter(rt.handlers)
	if pp := viper.GetString("PROFILE_PORT"); pp != "" {
		go func() {
			pmux := http.NewServeMux()
			pmux.HandleFunc("/debug/pprof/", pprof.Index)
			pmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			pmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			pmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			pmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
			log.Printf("payssage pprof listening on :%s", pp)
			log.Println(http.ListenAndServe(":"+pp, pmux))
		}()
	}
	port := viper.GetString("port")
	log.Printf("payssage listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func (app *Payssage) startHandlers(rt runtime) *HTTPDeps {
	var handlers HTTPDeps
	handlers.AuthMiddleware = rt.middlewares.authMW
	handlers.IntentsHandler = intents.NewHandler(rt.commands.intents, rt.queries.intents)
	handlers.WorkspacesHandler = workspaces.NewHandler(rt.commands.workspaces, rt.queries.workspaces)
	handlers.ApiKeysHandler = api_keys.NewHandler(rt.commands.apiKeys, rt.queries.apiKeys)
	handlers.WebhooksHandler = webhooks.NewHandler(rt.commands.webhooks, rt.queries.webhooks)
	handlers.OauthHandler = oauth.NewHandler(rt.commands.oauth, rt.queries.oauth)
	return &handlers
}

func (app *Payssage) startCommands(rt runtime, r repos) commands {
	var cmd commands
	cmd.webhooks = webhooks.NewCommandService(r.endpointsRepo, r.deliveriesRepo, r.eventsRepo, r.workspaceRepo, r.intentRepo, r.providerCredentialsRepo, app.river, app.sdb, rt.txRunner, rt.tracer)
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
	mw.authMW = NewAuthMiddleware(app.ga, rt.repos.apiKeysRepo, rt.repos.workspaceRepo, rt.tracer)
	return mw
}
