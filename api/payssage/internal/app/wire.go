package app

import (
	"net/http"
	"payssage/internal/authz"
	"payssage/internal/features/providers"
	"payssage/internal/handlers"
	providers2 "payssage/internal/providers"
	"payssage/internal/repos"
	"payssage/internal/services"
	"payssage/internal/sqlc"
	"payssage/ports"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

// ── Wire types ────────────────────────────────────────────────────────────

type middlewares struct {
	jwtAuth    func(http.Handler) http.Handler
	apiKeyAuth func(http.Handler) http.Handler
	anyAuth    func(http.Handler) http.Handler
}

// ── Init functions ────────────────────────────────────────────────────────

func (app *Payssage) initRepos(q *sqlc.Queries) *repos.Repos {
	return repos.New(q)
}

func (app *Payssage) initProviders(r *repos.Repos) {
	mercadoPago := providers.NewMercadoPago(app.cfg.MercadoPagoConfig, r.Intents, r.Collectors, r.Sellers, r.Wallets, app.httpClient)

	providers2.PayssageProviders.OAuth = map[providers2.AvailableProviders]ports.OAuthProvider{
		providers2.MercadoPagoProvider: mercadoPago,
	}

	providers2.PayssageProviders.Payments = map[providers2.AvailableProviders]ports.PaymentAbstractionLayer{
		providers2.MercadoPagoProvider: mercadoPago,
	}

	providers2.PayssageProviders.Webhooks = map[providers2.AvailableProviders]ports.WebhookProvider{
		providers2.MercadoPagoProvider: mercadoPago,
	}
}

func (app *Payssage) initOperations(riverClient *river.Client[pgx.Tx], r *repos.Repos) *services.Operations {
	authzSvc := authz.New(r.Organizations, r.Wallets)
	return services.NewOperations(r, authzSvc, riverClient, app.idxClient)
}

func (app *Payssage) initMiddlewares() middlewares {
	var mw middlewares
	authMW := app.setupAuthMiddlewares()
	mw.jwtAuth = authMW.JWT()
	mw.apiKeyAuth = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	return mw
}

func (app *Payssage) initHandlers(ops *services.Operations) *handlers.Server {
	return handlers.NewServer(ops)
}
