package app

import (
	"log"
	"net/http"
	"net/http/pprof"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/products"
	"univents/internal/features/ticket_types"

	fh "github.com/MintzyG/fun/handlers"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func (app *Univents) CreateRouter(middlewares middlewares, handlers handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.realIP)
	r.Use(middlewares.requestID)
	r.Use(middlewares.logger)
	r.Use(middlewares.metrics)
	r.Use(middlewares.recover)
	r.Use(middlewares.timeout)
	r.Use(middlewares.bodySize)
	r.Use(middlewares.ratelimit)
	r.Use(middlewares.cors)

	r.Handle("/metrics", promhttp.Handler())

	//r.With(middlewares.jwt).Get("/ws/token", deps.Security.WSAuth)
	events.RegisterRoutes(r, handlers.Events, middlewares.jwt)
	editions.RegisterRoutes(r, handlers.Editions, middlewares.jwt)
	ticket_types.RegisterRoutes(r, handlers.TicketTypes, middlewares.jwt)
	products.RegisterRoutes(r, handlers.Products, middlewares.jwt)
	//activities.RegisterRoutes(r, handlers.Activities, middlewares.jwt)
	//signatures.RegisterRoutes(r, handlers.signatures, middlewares.jwt)
	//certifications.RegisterRoutes(r, handlers.certs, middlewares.jwt)
	//checkpoints.Routes(r, handlers.Checkpoints, middlewares.jwt)
	//products.Routes(r, handlers.Products, middlewares.jwt)
	//purchases.Routes(r, handlers.Purchases, middlewares.jwt)

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

func servePprof(port string) {
	pmux := http.NewServeMux()
	pmux.HandleFunc("/debug/pprof/", pprof.Index)
	pmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	log.Printf("univents pprof listening on :%s", port)
	if err := http.ListenAndServe(":"+port, pmux); err != nil {
		log.Fatalf("univents pprof server error: %v", err)
	}
}
