package app

import (
	"Informd/internal/database/sqlc"
	"Informd/internal/features/answers"
	"Informd/internal/features/fields"
	"Informd/internal/features/forms"
	"Informd/internal/features/namespaces"
	"Informd/internal/features/responders"
	"Informd/internal/features/responses"
	"Informd/internal/features/steps"
	"Informd/ports"
	"lib/database"
	"lib/errx"
	"lib/xslices"
	"net/http"
	"strings"
	"time"

	fm "github.com/MintzyG/fun/middlewares"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ── Wire types ────────────────────────────────────────────────────────────

type repos struct {
	namespaces ports.NamespaceRepo
	forms      ports.FormsRepo
	steps      ports.StepRepo
	fields     ports.FieldsRepo
	answers    ports.AnswerRepo
	responders ports.ResponderRepo
	responses  ports.ResponseRepo
}

type queries struct {
	namespaces *namespaces.Queries
	forms      *forms.Queries
	steps      *steps.Queries
	fields     *fields.Queries
}

type commands struct {
	namespaces *namespaces.Commands
	forms      *forms.Commands
	steps      *steps.Commands
	fields     *fields.Commands
	responses  *responses.Commands
}

type handlers struct {
	namespaces *namespaces.Handlers
	forms      *forms.Handlers
	steps      *steps.Handlers
	fields     *fields.Handlers
	responses  *responses.Handlers
}

type middlewares struct {
	logger    func(http.Handler) http.Handler
	requestID func(http.Handler) http.Handler
	bodySize  func(http.Handler) http.Handler
	metrics   func(http.Handler) http.Handler
	cors      func(http.Handler) http.Handler
	realIP    func(http.Handler) http.Handler
	recover   func(http.Handler) http.Handler
	timeout   func(http.Handler) http.Handler
	ratelimit func(http.Handler) http.Handler
	jwt       func(http.Handler) http.Handler
	apiKey    func(http.Handler) http.Handler
	anyAuth   func(http.Handler) http.Handler
}

// ── Init functions ────────────────────────────────────────────────────────

func initRepos(q *sqlc.Queries, logger *zap.Logger, tracer trace.Tracer) repos {
	return repos{
		namespaces: namespaces.NewRepo(q, logger, tracer),
		forms:      forms.NewRepo(q, logger, tracer),
		steps:      steps.NewRepo(q, logger, tracer),
		fields:     fields.NewRepos(q, logger, tracer),
		answers:    answers.NewRepo(q, logger, tracer),
		responders: responders.NewRepo(q, logger, tracer),
		responses:  responses.NewRepo(q, logger, tracer),
	}
}

func initQueries(r repos, logger *zap.Logger, tracer trace.Tracer, tx database.TxRunner) queries {
	return queries{
		namespaces: namespaces.NewQueries(r.namespaces, r.forms, r.steps, r.fields, r.answers, r.responses, r.responders, logger, tx, tracer),
		forms:      forms.NewQueries(r.forms, r.steps, r.fields, r.answers, r.responses, r.responders, r.namespaces, logger, tx, tracer),
		steps:      steps.NewQueries(r.forms, r.steps, r.namespaces, logger, tx, tracer),
		fields:     fields.NewQueries(r.forms, r.steps, r.fields, r.namespaces, logger, tx, tracer),
	}
}

func initCommands(r repos, logger *zap.Logger, tracer trace.Tracer, tx database.TxRunner) commands {
	return commands{
		namespaces: namespaces.NewCommands(r.namespaces, r.forms, logger, tx, tracer),
		forms:      forms.NewCommands(r.forms, r.steps, r.namespaces, logger, tx, tracer),
		steps:      steps.NewCommands(r.forms, r.steps, r.namespaces, logger, tx, tracer),
		fields:     fields.NewCommands(r.forms, r.steps, r.fields, r.namespaces, logger, tx, tracer),
		responses:  responses.NewCommands(r.responders, r.responses, r.answers, r.forms, logger, tx, tracer),
	}
}

func initHandlers(c commands, q queries) handlers {
	return handlers{
		namespaces: namespaces.NewHandler(c.namespaces, q.namespaces),
		forms:      forms.NewHandlers(c.forms, q.forms),
		steps:      steps.NewHandlers(c.steps, q.steps),
		fields:     fields.NewHandlers(c.fields, q.fields),
		responses:  responses.NewHandlers(c.responses),
	}
}

func initMiddlewares(logger *zap.Logger) middlewares {
	var mw middlewares
	authMW := setupAuthMiddlewares()
	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	mw.bodySize = fm.MaxBodySize(1 << 20)
	mw.requestID = fm.RequestID(fm.RequestIDConfig{Header: "X-Request-ID"})
	mw.logger = fm.Logs(fm.Config{Logger: logger, SkipPrefixes: []string{"/metrics", "/health"}, RequestIDHeader: "X-Request-ID"})
	collectors, err := fm.NewCollectors(prometheus.DefaultRegisterer)
	if err != nil {
		errx.Exit(err, "Failed to create collectors")
	}
	mw.metrics = fm.Metrics(collectors, fm.MetricsConfig{SkipPrefixes: []string{"/metrics", "/health"}})
	mw.cors = fm.CORS(fm.CORSConfig{
		AllowedOrigins:   xslices.Clean(strings.Split(app.cfg.CorsAllowedOrigins, ",")),
		AllowCredentials: true,
	})
	mw.realIP = fm.RealIP()
	mw.recover = fm.Recover(logger)
	mw.timeout = fm.Timeout(60 * time.Second)
	mw.ratelimit = fm.RateLimit(fm.RateLimitConfig{RPS: 400, Burst: 20,
		KeyExtractor: func(r *http.Request) string { return r.RemoteAddr },
	})
	return mw
}
