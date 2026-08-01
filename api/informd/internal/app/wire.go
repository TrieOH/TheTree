package app

import (
	"Informd/internal/authz"
	"Informd/internal/features/answers"
	"Informd/internal/features/fields"
	"Informd/internal/features/forms"
	"Informd/internal/features/namespaces"
	"Informd/internal/features/responders"
	"Informd/internal/features/responses"
	"Informd/internal/features/steps"
	"Informd/internal/sqlc"
	"Informd/ports"
	"net/http"
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
	jwt     func(http.Handler) http.Handler
	apiKey  func(http.Handler) http.Handler
	anyAuth func(http.Handler) http.Handler
}

// ── Init functions ────────────────────────────────────────────────────────

func (app *Informd) initRepos(q *sqlc.Queries) repos {
	r := repos{
		namespaces: namespaces.NewRepo(q),
		forms:      forms.NewRepo(q),
		steps:      steps.NewRepo(q),
		fields:     fields.NewRepos(q),
		answers:    answers.NewRepo(q),
		responders: responders.NewRepo(q),
		responses:  responses.NewRepo(q),
	}
	authz.Service = authz.New(r.forms, r.namespaces)
	return r
}

func (app *Informd) initQueries(r repos) queries {
	return queries{
		namespaces: namespaces.NewQueries(r.namespaces, r.forms, r.steps, r.fields, r.answers, r.responses, r.responders),
		forms:      forms.NewQueries(r.forms, r.steps, r.fields, r.answers, r.responses, r.responders, r.namespaces),
		steps:      steps.NewQueries(r.forms, r.steps, r.namespaces),
		fields:     fields.NewQueries(r.forms, r.steps, r.fields, r.namespaces),
	}
}

func (app *Informd) initCommands(r repos) commands {
	return commands{
		namespaces: namespaces.NewCommands(r.namespaces, r.forms),
		forms:      forms.NewCommands(r.forms, r.steps, r.namespaces),
		steps:      steps.NewCommands(r.forms, r.steps, r.namespaces),
		fields:     fields.NewCommands(r.forms, r.steps, r.fields, r.namespaces),
		responses:  responses.NewCommands(r.responders, r.responses, r.answers, r.forms),
	}
}

func (app *Informd) initMiddlewares() middlewares {
	var mw middlewares
	authMW := app.setupAuthMiddlewares()
	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
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
