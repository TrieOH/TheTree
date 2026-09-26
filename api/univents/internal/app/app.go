package app

import (
	"context"
	"net/http"

	"lib/database"
	"lib/email"
	"lib/httpserver"
	"lib/objectstorage"
	libriver "lib/river"
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
// river workers and jobs — happens in the start hook. Storage never
// crosses Boot's interface: the pool is created here and closed in the
// returned shutdown.
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

		pool, err := database.SetupDB(cfg.ToDBConfig(), constraintMessages())
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

		// River must exist before the operations: the webhook receiver
		// cancels the expiry job on approve via the client (best-effort;
		// split 7 checkout schedules the job).
		riverClient, riverUI, err := app.initRiver(ctx, repos, notifier, tx)
		if err != nil {
			return nil, nil, err
		}

		ops := app.initOperations(repos, notifier, riverClient, tx)
		primitives := app.initMiddlewares()
		handlers := app.initHandlers(ops)

		mux := app.CreateRouter(primitives, handlers, riverUI)
		return mux, func(ctx context.Context) error {
			libriver.LogStop(ctx, riverClient)
			database.CloseDB(pool)
			return nil
		}, nil
	}

	return httpserver.Boot(httpserver.Config{
		AppName:            cfg.AppName,
		Port:               cfg.Port,
		ProfilePort:        cfg.ProfilePort,
		CorsAllowedOrigins: cfg.AllowedOrigins,
		CorsAllowedHeaders: cfg.AllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
	}, start)
}
