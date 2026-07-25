package app

import (
	"IdentityX/internal/features/actors"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/capabilities"
	"net/http"

	"IdentityX/internal/features/authn"
	"IdentityX/internal/features/organizations"
	"IdentityX/internal/features/profile_schemas"
	"IdentityX/internal/features/profiles"
	"IdentityX/internal/features/projects"

	fh "github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func (app *IdentityX) CreateRouter(middlewares middlewares, handlers handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(otelchi.Middleware(app.cfg.AppName,
		otelchi.WithChiRoutes(r),
		otelchi.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/health" && r.URL.Path != "/metrics"
		}),
	))

	// r.Use(deps.RealIP)
	// r.Use(deps.RequestID)
	r.Use(middlewares.logger)
	r.Use(middlewares.metrics)
	// r.Use(deps.Recover)
	// r.Use(deps.Timeout)
	// r.Use(deps.BodySize)
	// r.Use(deps.RateLimit)
	r.Use(middlewares.cors)

	// endpoints := riverui.NewEndpoints(app.river, nil)
	//
	// handler, err := riverui.NewHandler(&riverui.HandlerOpts{
	//	 Endpoints: endpoints,
	//	 Logger:    slog.Default(),
	//	 Prefix:    "/riverui",
	// })
	// if err != nil {
	// 	 errx.Exit(err, "failed to create river handler")
	// }
	// err = handler.Start(context.Background())
	// if err != nil {
	//	 errx.Exit(err, "failed to start river handler")
	// }
	// r.Mount("/riverui", handler)

	r.Handle("/metrics", promhttp.Handler())

	actors.RegisterRoutes(r, handlers.Actors, middlewares.anyAuth, middlewares.clientOnly)
	api_keys.RegisterRoutes(r, handlers.APIKeys, middlewares.anyAuth, middlewares.clientOnly)
	authn.RegisterRoutes(r, handlers.Authn, middlewares.jwtAuth, middlewares.anyAuth)
	organizations.RegisterRoutes(r, handlers.Orgs, middlewares.jwtAuth, middlewares.clientOnly)
	projects.RegisterRoutes(r, handlers.Projects, middlewares.anyAuth, middlewares.clientOnly)
	capabilities.RegisterRoutes(r, handlers.Capabilities, middlewares.jwtAuth, middlewares.anyAuth, middlewares.clientOnly)
	profiles.RegisterRoutes(r, handlers.Profiles, middlewares.jwtAuth, middlewares.clientOnly)
	profile_schemas.RegisterRoutes(r, handlers.ProfileSchemas, middlewares.jwtAuth, middlewares.clientOnly)

	r.Get("/health", fh.Health(app.cfg.AppName).Handle)

	return otelhttp.NewHandler(r, "http.server",
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/health"
		}),
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/metrics"
		}),
	)
}
