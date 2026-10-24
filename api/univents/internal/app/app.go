package app

import (
	"context"
	"net/http"

	"lib/database"
	"lib/email"
	"lib/httpserver"
	"lib/objectstorage"
	spec "univents"
	"univents/internal/config"

	idx "sdk/identityx"
	payssage "sdk/payssage"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Univents struct {
	db          *pgxpool.Pool
	idxClient   *idx.Client
	objStorage  *objectstorage.Client
	emailClient *email.Client
	payssage    *payssage.Client

	cfg config.Config
}

// Run boots Univents through the Harness: the process sequence (FUN
// runtime, tracer, serving, shutdown ordering) is httpserver.Boot's
// implementation; everything Univents varies on — its clients (IdentityX,
// Payssage, object storage, email), the platform-wallet fail-fast check,
// its Postgres adapter with the constraint messages, the store notifier,
// and the background-job carrier — happens in the start hook. Storage never
// crosses Boot's interface: the pool is created here and closed in the
// returned shutdown. The schema is never migrated here — `just
// univents-goose` owns it.
func Run() error {
	cfg := config.Load()
	app := &Univents{cfg: cfg}

	start := func(ctx context.Context) (http.Handler, func(context.Context) error, error) {
		idxClient, err := idx.Bootstrap(ctx, cfg.ToIdentityXConfig())
		if err != nil {
			return nil, nil, err
		}
		app.idxClient = idxClient

		objStorage, err := SetupObjectStorage(cfg)
		if err != nil {
			return nil, nil, err
		}
		app.objStorage = objStorage
		app.emailClient = email.NewClient(cfg.ToEmailConfig())
		app.payssage = SetupPayssage(cfg)

		// Fail fast on the platform wallet: wrong/missing PAYSSAGE_WALLET_ID
		// or an unreachable Payssage must stop the boot, not surface at
		// checkout (D6).
		err = VerifyPayssageWallet(ctx, app.payssage, cfg.PayssageWalletID)
		if err != nil {
			return nil, nil, err
		}

		pool, err := database.SetupDBWithoutMigrations(cfg.ToDBConfig(), constraintMessages())
		if err != nil {
			return nil, nil, err
		}
		app.db = pool

		// The notifier (lib/go/database) is the store's LISTEN/NOTIFY bridge:
		// the webhook receiver publishes on it (split 4); the SSE relay and
		// WS hub subscribe in split 6. Notify opens its own connection per
		// call, so nothing needs starting here.
		notifier := database.NewNotifier(cfg.ToDBConfig().DSN())

		tx := database.NewPGXTxRunner(pool)

		repos := app.initRepos()

		enqueuer, events, stopJobs, err := app.initJobs(repos, notifier, tx)
		if err != nil {
			database.CloseDB(pool)
			return nil, nil, err
		}

		ops := app.initOperations(repos, notifier, enqueuer, tx)
		primitives := app.initMiddlewares()
		handlers := app.initHandlers(ops)

		mux := app.CreateRouter(primitives, handlers, events)
		return mux, func(ctx context.Context) error {
			err := stopJobs(ctx)
			database.CloseDB(pool)
			return err
		}, nil
	}

	return httpserver.Boot(httpserver.Config{
		AppName:            cfg.AppName,
		Port:               cfg.Port,
		ProfilePort:        cfg.ProfilePort,
		CorsAllowedOrigins: cfg.AllowedOrigins,
		CorsAllowedHeaders: cfg.AllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		RateLimit:          cfg.ToRateLimit(),
	}, start)
}
