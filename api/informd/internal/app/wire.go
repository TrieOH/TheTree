package app

import (
	"Informd/internal/features/answers"
	"Informd/internal/features/fields"
	"Informd/internal/features/forms"
	"Informd/internal/features/namespaces"
	"Informd/internal/features/responders"
	"Informd/internal/features/responses"
	"Informd/internal/features/steps"
	"Informd/internal/sqlc"
	"Informd/ports"
	"lib/database"
	"lib/errx"
	"lib/telemetry"
	"lib/xslices"
	"net/http"
	"strings"
	"time"

	fm "github.com/MintzyG/fun/middlewares"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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

func (app *Informd) tracer() trace.Tracer {
	return otel.Tracer(app.cfg.AppName)
}

func (app *Informd) txRunner() database.TxRunner {
	return database.NewPGXTxRunner(app.db)
}

func (app *Informd) initRepos(q *sqlc.Queries) repos {
	return repos{
		namespaces: namespaces.NewRepo(q, app.tracer()),
		forms:      forms.NewRepo(q, app.tracer()),
		steps:      steps.NewRepo(q, app.tracer()),
		fields:     fields.NewRepos(q, app.tracer()),
		answers:    answers.NewRepo(q, app.tracer()),
		responders: responders.NewRepo(q, app.tracer()),
		responses:  responses.NewRepo(q, app.tracer()),
	}
}

func (app *Informd) initQueries(r repos) queries {
	return queries{
		namespaces: namespaces.NewQueries(r.namespaces, r.forms, r.steps, r.fields, r.answers, r.responses, r.responders, app.txRunner(), app.tracer()),
		forms:      forms.NewQueries(r.forms, r.steps, r.fields, r.answers, r.responses, r.responders, r.namespaces, app.txRunner(), app.tracer()),
		steps:      steps.NewQueries(r.forms, r.steps, r.namespaces, app.txRunner(), app.tracer()),
		fields:     fields.NewQueries(r.forms, r.steps, r.fields, r.namespaces, app.txRunner(), app.tracer()),
	}
}

func (app *Informd) initCommands(r repos) commands {
	return commands{
		namespaces: namespaces.NewCommands(r.namespaces, r.forms, app.txRunner(), app.tracer()),
		forms:      forms.NewCommands(r.forms, r.steps, r.namespaces, app.txRunner(), app.tracer()),
		steps:      steps.NewCommands(r.forms, r.steps, r.namespaces, app.txRunner(), app.tracer()),
		fields:     fields.NewCommands(r.forms, r.steps, r.fields, r.namespaces, app.txRunner(), app.tracer()),
		responses:  responses.NewCommands(r.responders, r.responses, r.answers, r.forms, app.txRunner(), app.tracer()),
	}
}

func (app *Informd) initMiddlewares() middlewares {
	var mw middlewares
	authMW := app.setupAuthMiddlewares()
	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	mw.bodySize = fm.MaxBodySize(1 << 20)
	mw.requestID = fm.RequestID(fm.RequestIDConfig{Header: "X-Request-ID"})
	mw.logger = fm.Logs(fm.Config{Logger: telemetry.Log(), SkipPrefixes: []string{"/metrics", "/health"}, RequestIDHeader: "X-Request-ID"})
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
	mw.recover = fm.Recover(telemetry.Log())
	mw.timeout = fm.Timeout(60 * time.Second)
	mw.ratelimit = fm.RateLimit(fm.RateLimitConfig{RPS: 400, Burst: 20,
		KeyExtractor: func(r *http.Request) string { return r.RemoteAddr },
	})
	return mw
}

func (app *Informd) initHandlers(c commands, q queries) handlers {
	return handlers{
		namespaces: namespaces.NewHandler(c.namespaces, q.namespaces),
		forms:      forms.NewHandlers(c.forms, q.forms),
		steps:      steps.NewHandlers(c.steps, q.steps),
		fields:     fields.NewHandlers(c.fields, q.fields),
		responses:  responses.NewHandlers(c.responses),
	}
}
