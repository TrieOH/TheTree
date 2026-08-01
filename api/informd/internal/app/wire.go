package app

import (
	"Informd/internal/authz"
	answersRepos "Informd/internal/features/answers/repos"
	"Informd/internal/features/fields"
	fieldsHandlers "Informd/internal/features/fields/handlers"
	fieldsRepos "Informd/internal/features/fields/repos"
	"Informd/internal/features/forms"
	formsHandlers "Informd/internal/features/forms/handlers"
	formsRepos "Informd/internal/features/forms/repos"
	"Informd/internal/features/namespaces"
	namespacesHandlers "Informd/internal/features/namespaces/handlers"
	namespacesRepos "Informd/internal/features/namespaces/repos"
	respondersRepos "Informd/internal/features/responders/repos"
	"Informd/internal/features/responses"
	responsesHandlers "Informd/internal/features/responses/handlers"
	responsesRepos "Informd/internal/features/responses/repos"
	"Informd/internal/features/steps"
	stepsHandlers "Informd/internal/features/steps/handlers"
	stepsRepos "Informd/internal/features/steps/repos"
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

type operations struct {
	namespaces *namespaces.Operations
	forms      *forms.Operations
	steps      *steps.Operations
	fields     *fields.Operations
	responses  *responses.Operations
}

type handlers struct {
	namespaces *namespacesHandlers.Handlers
	forms      *formsHandlers.Handlers
	steps      *stepsHandlers.Handlers
	fields     *fieldsHandlers.Handlers
	responses  *responsesHandlers.Handlers
}

type middlewares struct {
	jwt     func(http.Handler) http.Handler
	apiKey  func(http.Handler) http.Handler
	anyAuth func(http.Handler) http.Handler
}

// ── Init functions ────────────────────────────────────────────────────────

func (app *Informd) initRepos(q *sqlc.Queries) repos {
	r := repos{
		namespaces: namespacesRepos.NewRepo(q),
		forms:      formsRepos.NewRepo(q),
		steps:      stepsRepos.NewRepo(q),
		fields:     fieldsRepos.NewRepo(q),
		answers:    answersRepos.NewRepo(q),
		responders: respondersRepos.NewRepo(q),
		responses:  responsesRepos.NewRepo(q),
	}
	authz.Service = authz.New(r.forms, r.namespaces)
	return r
}

func (app *Informd) initOperations(r repos) operations {
	return operations{
		namespaces: namespaces.NewOperations(r.namespaces, r.forms, r.steps, r.fields, r.answers, r.responses, r.responders),
		forms:      forms.NewOperations(r.forms, r.steps, r.namespaces, r.fields, r.answers, r.responses, r.responders),
		steps:      steps.NewOperations(r.forms, r.steps, r.namespaces),
		fields:     fields.NewOperations(r.forms, r.steps, r.fields, r.namespaces),
		responses:  responses.NewOperations(r.responders, r.responses, r.answers, r.forms),
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

func (app *Informd) initHandlers(ops operations) handlers {
	return handlers{
		namespaces: namespacesHandlers.NewHandler(ops.namespaces),
		forms:      formsHandlers.NewHandlers(ops.forms),
		steps:      stepsHandlers.NewHandlers(ops.steps),
		fields:     fieldsHandlers.NewHandlers(ops.fields),
		responses:  responsesHandlers.NewHandlers(ops.responses),
	}
}
