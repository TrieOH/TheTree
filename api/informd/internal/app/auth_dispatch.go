package app

import (
	"context"
	"net/http"
	"slices"

	"Informd/internal/openapi"
)

// authChains maps every operationId to the middleware chain it must run
// through, mirroring the parity-test matrix. Namespace routes require a
// JWT; form/step/field routes accept JWT or API key (anyAuth); the
// answerable view and response submission are public.
func authChains(mw middlewares) map[string][]func(http.Handler) http.Handler {
	chains := map[string][]func(http.Handler) http.Handler{
		"addFormMember":                  {mw.anyAuth},
		"addFormMemberNamespaced":        {mw.jwt},
		"addNamespaceMember":             {mw.jwt},
		"archiveForm":                    {mw.anyAuth},
		"archiveFormNamespaced":          {mw.jwt},
		"bulkEditFields":                 {mw.anyAuth},
		"bulkEditFieldsNamespaced":       {mw.anyAuth},
		"bulkEditSteps":                  {mw.anyAuth},
		"bulkEditStepsNamespaced":        {mw.anyAuth},
		"closeForm":                      {mw.anyAuth},
		"closeFormNamespaced":            {mw.jwt},
		"createField":                    {mw.anyAuth},
		"createFieldNamespaced":          {mw.anyAuth},
		"createForm":                     {mw.anyAuth},
		"createNamespace":                {mw.jwt},
		"createNamespaceForm":            {mw.jwt},
		"createStep":                     {mw.anyAuth},
		"createStepNamespaced":           {mw.anyAuth},
		"deleteField":                    {mw.anyAuth},
		"deleteFieldNamespaced":          {mw.anyAuth},
		"editSelectConfig":               {mw.anyAuth},
		"editSelectConfigNamespaced":     {mw.anyAuth},
		"getAnswerableForm":              {},
		"getFormResponseCount":           {mw.anyAuth},
		"getFormResponseCountNamespaced": {mw.jwt},
		"getFullForm":                    {mw.anyAuth},
		"getFullFormNamespaced":          {mw.jwt},
		"getSelectConfig":                {mw.anyAuth},
		"getSelectConfigNamespaced":      {mw.anyAuth},
		"listFields":                     {mw.anyAuth},
		"listFieldsNamespaced":           {mw.anyAuth},
		"listFormMembers":                {mw.anyAuth},
		"listFormMembersNamespaced":      {mw.jwt},
		"listMyArchivedForms":            {mw.anyAuth},
		"listMyForms":                    {mw.anyAuth},
		"listNamespaceArchivedForms":     {mw.jwt},
		"listNamespaceForms":             {mw.jwt},
		"listNamespaceMembers":           {mw.jwt},
		"listNamespaces":                 {mw.jwt},
		"listSteps":                      {mw.anyAuth},
		"listStepsNamespaced":            {mw.anyAuth},
		"openForm":                       {mw.anyAuth},
		"openFormNamespaced":             {mw.jwt},
		"redraftForm":                    {mw.anyAuth},
		"redraftFormNamespaced":          {mw.jwt},
		"removeFormMember":               {mw.anyAuth},
		"removeFormMemberNamespaced":     {mw.jwt},
		"removeNamespaceMember":          {mw.jwt},
		"submitResponse":                 {},
	}
	return chains
}

// authDispatch is the strict-server middleware that resolves the auth
// chain for each operation (by operationId) and runs it around the
// handler. Public operations (empty chain) pass through untouched.
func authDispatch(mw middlewares) openapi.StrictMiddlewareFunc {
	chains := authChains(mw)
	return func(f openapi.StrictHandlerFunc, operationID string) openapi.StrictHandlerFunc {
		chain := chains[operationID]
		if len(chain) == 0 {
			return f
		}
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
			var resp any
			var ferr error
			var called bool
			var next http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				resp, ferr = f(ctx, w, r, request)
			})
			for i := range slices.Backward(chain) {
				next = chain[i](next)
			}
			next.ServeHTTP(w, r)
			if !called {
				// An auth middleware rejected the request and already
				// wrote the response; nothing to return.
				return nil, nil
			}
			return resp, ferr
		}
	}
}
