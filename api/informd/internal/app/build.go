package app

import (
	"context"

	"Informd/internal/authz"
	"Informd/internal/features/answers"
	"Informd/internal/features/fields"
	"Informd/internal/features/forms"
	"Informd/internal/features/namespaces"
	"Informd/internal/features/responders"
	"Informd/internal/features/responses"
	"Informd/internal/features/steps"
	"Informd/internal/sqlc"
	"lib/backend"

	"github.com/go-chi/chi/v5"
)

// Build wires Informd on the kit (ADR-0014): repos over the postgres
// module, each feature's operations over the repos (ADR-0015), auth from
// the identityx module, and the strict server.
func Build(_ context.Context, d *backend.Deps, r chi.Router) error {
	q := sqlc.New(d.DB)
	namespacesRepo := namespaces.NewRepo(q)
	formsRepo := forms.NewRepo(q)
	stepsRepo := steps.NewRepo(q)
	fieldsRepo := fields.NewRepo(q)
	answersRepo := answers.NewRepo(q)
	respondersRepo := responders.NewRepo(q)
	responsesRepo := responses.NewRepo(q)
	az := authz.New(formsRepo, namespacesRepo)

	srv := &server{
		namespaceHandlers: namespaces.NewHandlers(namespaces.New(namespaces.Deps{
			Namespaces: namespacesRepo, Forms: formsRepo, Steps: stepsRepo, Fields: fieldsRepo,
			Answers: answersRepo, Responses: responsesRepo, Responders: respondersRepo, Authz: az, Tx: d.Tx,
		})),
		formHandlers: forms.NewHandlers(forms.New(forms.Deps{
			Forms: formsRepo, Steps: stepsRepo, Namespaces: namespacesRepo, Fields: fieldsRepo,
			Answers: answersRepo, Responses: responsesRepo, Responders: respondersRepo, Authz: az, Tx: d.Tx,
		})),
		stepHandlers: steps.NewHandlers(steps.New(steps.Deps{
			Forms: formsRepo, Steps: stepsRepo, Namespaces: namespacesRepo, Authz: az,
		})),
		fieldHandlers: fields.NewHandlers(fields.New(fields.Deps{
			Forms: formsRepo, Steps: stepsRepo, Fields: fieldsRepo, Namespaces: namespacesRepo, Authz: az, Tx: d.Tx,
		})),
		responseHandlers: responses.NewHandlers(responses.New(responses.Deps{
			Responders: respondersRepo, Responses: responsesRepo, Answers: answersRepo, Forms: formsRepo, Tx: d.Tx,
		})),
	}

	chains, err := d.Chains(d.IDX.Auth)
	if err != nil {
		return err
	}
	mountStrict(r, srv, chains)
	return nil
}
