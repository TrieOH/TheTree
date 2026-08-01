package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"Informd/internal/openapi"

	"github.com/go-chi/chi/v5"
)

// stubStrict implements openapi.StrictServerInterface for route walking;
// registration only stores handlers, they are never served.
type stubStrict struct{}

// errStub is returned by the walk-only stub; it is never served.
var errStub = errors.New("parity stub")

func (stubStrict) ListMyForms(_ context.Context, _ openapi.ListMyFormsRequestObject) (openapi.ListMyFormsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateForm(_ context.Context, _ openapi.CreateFormRequestObject) (openapi.CreateFormResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListMyArchivedForms(_ context.Context, _ openapi.ListMyArchivedFormsRequestObject) (openapi.ListMyArchivedFormsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ArchiveForm(_ context.Context, _ openapi.ArchiveFormRequestObject) (openapi.ArchiveFormResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetAnswerableForm(_ context.Context, _ openapi.GetAnswerableFormRequestObject) (openapi.GetAnswerableFormResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CloseForm(_ context.Context, _ openapi.CloseFormRequestObject) (openapi.CloseFormResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetFullForm(_ context.Context, _ openapi.GetFullFormRequestObject) (openapi.GetFullFormResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RemoveFormMember(_ context.Context, _ openapi.RemoveFormMemberRequestObject) (openapi.RemoveFormMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListFormMembers(_ context.Context, _ openapi.ListFormMembersRequestObject) (openapi.ListFormMembersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) AddFormMember(_ context.Context, _ openapi.AddFormMemberRequestObject) (openapi.AddFormMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) OpenForm(_ context.Context, _ openapi.OpenFormRequestObject) (openapi.OpenFormResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RedraftForm(_ context.Context, _ openapi.RedraftFormRequestObject) (openapi.RedraftFormResponseObject, error) {
	return nil, errStub
}
func (stubStrict) SubmitResponse(_ context.Context, _ openapi.SubmitResponseRequestObject) (openapi.SubmitResponseResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetFormResponseCount(_ context.Context, _ openapi.GetFormResponseCountRequestObject) (openapi.GetFormResponseCountResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListSteps(_ context.Context, _ openapi.ListStepsRequestObject) (openapi.ListStepsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateStep(_ context.Context, _ openapi.CreateStepRequestObject) (openapi.CreateStepResponseObject, error) {
	return nil, errStub
}
func (stubStrict) BulkEditSteps(_ context.Context, _ openapi.BulkEditStepsRequestObject) (openapi.BulkEditStepsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListFields(_ context.Context, _ openapi.ListFieldsRequestObject) (openapi.ListFieldsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateField(_ context.Context, _ openapi.CreateFieldRequestObject) (openapi.CreateFieldResponseObject, error) {
	return nil, errStub
}
func (stubStrict) BulkEditFields(_ context.Context, _ openapi.BulkEditFieldsRequestObject) (openapi.BulkEditFieldsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DeleteField(_ context.Context, _ openapi.DeleteFieldRequestObject) (openapi.DeleteFieldResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetSelectConfig(_ context.Context, _ openapi.GetSelectConfigRequestObject) (openapi.GetSelectConfigResponseObject, error) {
	return nil, errStub
}
func (stubStrict) EditSelectConfig(_ context.Context, _ openapi.EditSelectConfigRequestObject) (openapi.EditSelectConfigResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListNamespaces(_ context.Context, _ openapi.ListNamespacesRequestObject) (openapi.ListNamespacesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateNamespace(_ context.Context, _ openapi.CreateNamespaceRequestObject) (openapi.CreateNamespaceResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListNamespaceForms(_ context.Context, _ openapi.ListNamespaceFormsRequestObject) (openapi.ListNamespaceFormsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateNamespaceForm(_ context.Context, _ openapi.CreateNamespaceFormRequestObject) (openapi.CreateNamespaceFormResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListNamespaceArchivedForms(_ context.Context, _ openapi.ListNamespaceArchivedFormsRequestObject) (openapi.ListNamespaceArchivedFormsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ArchiveFormNamespaced(_ context.Context, _ openapi.ArchiveFormNamespacedRequestObject) (openapi.ArchiveFormNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CloseFormNamespaced(_ context.Context, _ openapi.CloseFormNamespacedRequestObject) (openapi.CloseFormNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetFullFormNamespaced(_ context.Context, _ openapi.GetFullFormNamespacedRequestObject) (openapi.GetFullFormNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RemoveFormMemberNamespaced(_ context.Context, _ openapi.RemoveFormMemberNamespacedRequestObject) (openapi.RemoveFormMemberNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListFormMembersNamespaced(_ context.Context, _ openapi.ListFormMembersNamespacedRequestObject) (openapi.ListFormMembersNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) AddFormMemberNamespaced(_ context.Context, _ openapi.AddFormMemberNamespacedRequestObject) (openapi.AddFormMemberNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) OpenFormNamespaced(_ context.Context, _ openapi.OpenFormNamespacedRequestObject) (openapi.OpenFormNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RedraftFormNamespaced(_ context.Context, _ openapi.RedraftFormNamespacedRequestObject) (openapi.RedraftFormNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetFormResponseCountNamespaced(_ context.Context, _ openapi.GetFormResponseCountNamespacedRequestObject) (openapi.GetFormResponseCountNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListStepsNamespaced(_ context.Context, _ openapi.ListStepsNamespacedRequestObject) (openapi.ListStepsNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateStepNamespaced(_ context.Context, _ openapi.CreateStepNamespacedRequestObject) (openapi.CreateStepNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) BulkEditStepsNamespaced(_ context.Context, _ openapi.BulkEditStepsNamespacedRequestObject) (openapi.BulkEditStepsNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListFieldsNamespaced(_ context.Context, _ openapi.ListFieldsNamespacedRequestObject) (openapi.ListFieldsNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateFieldNamespaced(_ context.Context, _ openapi.CreateFieldNamespacedRequestObject) (openapi.CreateFieldNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) BulkEditFieldsNamespaced(_ context.Context, _ openapi.BulkEditFieldsNamespacedRequestObject) (openapi.BulkEditFieldsNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DeleteFieldNamespaced(_ context.Context, _ openapi.DeleteFieldNamespacedRequestObject) (openapi.DeleteFieldNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetSelectConfigNamespaced(_ context.Context, _ openapi.GetSelectConfigNamespacedRequestObject) (openapi.GetSelectConfigNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) EditSelectConfigNamespaced(_ context.Context, _ openapi.EditSelectConfigNamespacedRequestObject) (openapi.EditSelectConfigNamespacedResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RemoveNamespaceMember(_ context.Context, _ openapi.RemoveNamespaceMemberRequestObject) (openapi.RemoveNamespaceMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListNamespaceMembers(_ context.Context, _ openapi.ListNamespaceMembersRequestObject) (openapi.ListNamespaceMembersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) AddNamespaceMember(_ context.Context, _ openapi.AddNamespaceMemberRequestObject) (openapi.AddNamespaceMemberResponseObject, error) {
	return nil, errStub
}

func mwJWT(next http.Handler) http.Handler     { return next }
func mwAnyAuth(next http.Handler) http.Handler { return next }

func mwName(mw func(http.Handler) http.Handler) string {
	fn := runtime.FuncForPC(reflect.ValueOf(mw).Pointer())
	name := fn.Name()
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	return strings.TrimPrefix(name, "mw")
}

// routeOperation maps every walked route to its spec operationId.
var routeOperation = map[string]string{
	"DELETE /forms/{form_id}/members":                                                     "removeFormMember",
	"DELETE /forms/{form_id}/steps/{step_id}/fields/{field_id}":                           "deleteField",
	"DELETE /namespaces/{namespace_id}/forms/{form_id}/members":                           "removeFormMemberNamespaced",
	"DELETE /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}": "deleteFieldNamespaced",
	"DELETE /namespaces/{namespace_id}/members":                                           "removeNamespaceMember",
	"GET /docs/openapi.yml":                                                               "getOpenAPISpec",
	"GET /forms":                                                                          "listMyForms",
	"GET /forms/archived":                                                                 "listMyArchivedForms",
	"GET /forms/{form_id}/asnwerable":                                                     "getAnswerableForm",
	"GET /forms/{form_id}/full":                                                           "getFullForm",
	"GET /forms/{form_id}/members":                                                        "listFormMembers",
	"GET /forms/{form_id}/responses/count":                                                "getFormResponseCount",
	"GET /forms/{form_id}/steps":                                                          "listSteps",
	"GET /forms/{form_id}/steps/{step_id}/fields":                                         "listFields",
	"GET /forms/{form_id}/steps/{step_id}/fields/{field_id}/select":                       "getSelectConfig",
	"GET /health":                                                                             "getHealth",
	"GET /namespaces":                                                                         "listNamespaces",
	"GET /namespaces/{namespace_id}/forms":                                                    "listNamespaceForms",
	"GET /namespaces/{namespace_id}/forms/archived":                                           "listNamespaceArchivedForms",
	"GET /namespaces/{namespace_id}/forms/{form_id}/full":                                     "getFullFormNamespaced",
	"GET /namespaces/{namespace_id}/forms/{form_id}/members":                                  "listFormMembersNamespaced",
	"GET /namespaces/{namespace_id}/forms/{form_id}/responses/count":                          "getFormResponseCountNamespaced",
	"GET /namespaces/{namespace_id}/forms/{form_id}/steps":                                    "listStepsNamespaced",
	"GET /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields":                   "listFieldsNamespaced",
	"GET /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}/select": "getSelectConfigNamespaced",
	"GET /namespaces/{namespace_id}/members":                                                  "listNamespaceMembers",
	"POST /forms":                                                                             "createForm",
	"POST /forms/{form_id}/archive":                                                           "archiveForm",
	"POST /forms/{form_id}/close":                                                             "closeForm",
	"POST /forms/{form_id}/members":                                                           "addFormMember",
	"POST /forms/{form_id}/open":                                                              "openForm",
	"POST /forms/{form_id}/redraft":                                                           "redraftForm",
	"POST /forms/{form_id}/responses":                                                         "submitResponse",
	"POST /forms/{form_id}/steps":                                                             "createStep",
	"POST /forms/{form_id}/steps/{step_id}/fields":                                            "createField",
	"POST /namespaces":                                                                        "createNamespace",
	"POST /namespaces/{namespace_id}/forms":                                                   "createNamespaceForm",
	"POST /namespaces/{namespace_id}/forms/{form_id}/archive":                                 "archiveFormNamespaced",
	"POST /namespaces/{namespace_id}/forms/{form_id}/close":                                   "closeFormNamespaced",
	"POST /namespaces/{namespace_id}/forms/{form_id}/members":                                 "addFormMemberNamespaced",
	"POST /namespaces/{namespace_id}/forms/{form_id}/open":                                    "openFormNamespaced",
	"POST /namespaces/{namespace_id}/forms/{form_id}/redraft":                                 "redraftFormNamespaced",
	"POST /namespaces/{namespace_id}/forms/{form_id}/steps":                                   "createStepNamespaced",
	"POST /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields":                  "createFieldNamespaced",
	"POST /namespaces/{namespace_id}/members":                                                 "addNamespaceMember",
	"PUT /forms/{form_id}/steps":                                                              "bulkEditSteps",
	"PUT /forms/{form_id}/steps/{step_id}/fields":                                             "bulkEditFields",
	"PUT /forms/{form_id}/steps/{step_id}/fields/{field_id}/select":                           "editSelectConfig",
	"PUT /namespaces/{namespace_id}/forms/{form_id}/steps":                                    "bulkEditStepsNamespaced",
	"PUT /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields":                   "bulkEditFieldsNamespaced",
	"PUT /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}/select": "editSelectConfigNamespaced",
}

// expectedOps is the auth matrix keyed by operationId: "public" when the
// operation runs with no auth middleware, otherwise the chain names joined
// with "+". Harness-owned routes (getHealth, getOpenAPISpec) are excluded —
// they never run through the dispatch.
var expectedOps = map[string]string{
	"addFormMemberNamespaced":        "JWT",
	"addNamespaceMember":             "JWT",
	"archiveFormNamespaced":          "JWT",
	"closeFormNamespaced":            "JWT",
	"createNamespace":                "JWT",
	"createNamespaceForm":            "JWT",
	"getFormResponseCountNamespaced": "JWT",
	"getFullFormNamespaced":          "JWT",
	"listFormMembersNamespaced":      "JWT",
	"listNamespaceArchivedForms":     "JWT",
	"listNamespaceForms":             "JWT",
	"listNamespaceMembers":           "JWT",
	"listNamespaces":                 "JWT",
	"openFormNamespaced":             "JWT",
	"redraftFormNamespaced":          "JWT",
	"removeFormMemberNamespaced":     "JWT",
	"removeNamespaceMember":          "JWT",
	"addFormMember":                  "AnyAuth",
	"archiveForm":                    "AnyAuth",
	"bulkEditFields":                 "AnyAuth",
	"bulkEditFieldsNamespaced":       "AnyAuth",
	"bulkEditSteps":                  "AnyAuth",
	"bulkEditStepsNamespaced":        "AnyAuth",
	"closeForm":                      "AnyAuth",
	"createField":                    "AnyAuth",
	"createFieldNamespaced":          "AnyAuth",
	"createForm":                     "AnyAuth",
	"createStep":                     "AnyAuth",
	"createStepNamespaced":           "AnyAuth",
	"deleteField":                    "AnyAuth",
	"deleteFieldNamespaced":          "AnyAuth",
	"editSelectConfig":               "AnyAuth",
	"editSelectConfigNamespaced":     "AnyAuth",
	"getFormResponseCount":           "AnyAuth",
	"getFullForm":                    "AnyAuth",
	"getSelectConfig":                "AnyAuth",
	"getSelectConfigNamespaced":      "AnyAuth",
	"listFields":                     "AnyAuth",
	"listFieldsNamespaced":           "AnyAuth",
	"listFormMembers":                "AnyAuth",
	"listMyArchivedForms":            "AnyAuth",
	"listMyForms":                    "AnyAuth",
	"listSteps":                      "AnyAuth",
	"listStepsNamespaced":            "AnyAuth",
	"openForm":                       "AnyAuth",
	"redraftForm":                    "AnyAuth",
	"removeFormMember":               "AnyAuth",
	"getAnswerableForm":              "public",
	"submitResponse":                 "public",
}

func TestRouterParity(t *testing.T) {
	r := chi.NewRouter()
	// harness-owned routes; mirror their registration
	r.Get("/health", http.NotFoundHandler().ServeHTTP)
	r.Get("/docs/openapi.yml", http.NotFoundHandler().ServeHTTP)
	openapi.HandlerWithOptions(openapi.NewStrictHandler(stubStrict{}, nil), openapi.ChiServerOptions{
		BaseRouter: r,
	})

	got := map[string]string{}
	_ = chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		got[normalizeRoute(method+" "+route)] = "public"
		return nil
	})

	var missing, extra []string
	for want := range routeOperation {
		if _, ok := got[want]; !ok {
			missing = append(missing, want)
		}
	}
	for gotRoute := range got {
		if _, ok := routeOperation[gotRoute]; !ok {
			extra = append(extra, gotRoute)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)
	if len(missing) > 0 || len(extra) > 0 {
		t.Fatalf("route parity mismatch\nroutes expected but not walked:\n%s\nroutes walked but not expected:\n%s",
			strings.Join(missing, "\n"), strings.Join(extra, "\n"))
	}

	chains := authChains(middlewares{jwt: mwJWT, anyAuth: mwAnyAuth})
	var authMismatch, missingChain []string
	for opID, want := range expectedOps {
		chain, ok := chains[opID]
		if !ok {
			missingChain = append(missingChain, opID)
			continue
		}
		names := make([]string, 0, len(chain))
		for _, mw := range chain {
			if n := mwName(mw); n != "" && !strings.HasPrefix(n, "func") {
				names = append(names, n)
			}
		}
		gotAuth := strings.Join(names, "+")
		if gotAuth == "" {
			gotAuth = "public"
		}
		if gotAuth != want {
			authMismatch = append(authMismatch, fmt.Sprintf("%s: want %s, got %s", opID, want, gotAuth))
		}
	}
	for opID := range chains {
		if opID == "getHealth" || opID == "getOpenAPISpec" {
			continue
		}
		if _, ok := expectedOps[opID]; !ok {
			authMismatch = append(authMismatch, "chain present but not expected: "+opID)
		}
	}
	sort.Strings(missingChain)
	sort.Strings(authMismatch)
	if len(missingChain) > 0 || len(authMismatch) > 0 {
		t.Fatalf("auth matrix mismatch\noperations without a chain:\n%s\nmismatches:\n%s",
			strings.Join(missingChain, "\n"), strings.Join(authMismatch, "\n"))
	}

	t.Logf("parity ok: %d routes, %d operations with matching auth chains", len(got), len(chains))
}

func normalizeRoute(r string) string {
	parts := strings.SplitN(r, " ", 2)
	if len(parts) != 2 {
		return r
	}
	path := strings.TrimSuffix(parts[1], "/")
	if path == "" {
		path = "/"
	}
	return parts[0] + " " + path
}
