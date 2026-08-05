package app

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"sort"
	"strings"
	"testing"

	spec "Informd"
	"Informd/internal/openapi"
	"lib/authz"

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

// labeled middleware stubs record their names when run.
var parityInvocations []string

func labelMW(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parityInvocations = append(parityInvocations, name)
			next.ServeHTTP(w, r)
		})
	}
}

// runChain executes a chain and returns the middleware names that ran.
func runChain(chain []func(http.Handler) http.Handler) []string {
	parityInvocations = nil
	var next http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	for i := len(chain) - 1; i >= 0; i-- {
		next = chain[i](next)
	}
	next.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	return parityInvocations
}

// TestRouterRoutesMatchSpec asserts the router serves exactly the spec's
// paths, and nothing else. The harness-owned routes are declared in the
// spec (getHealth, getOpenAPISpec) and registered by the harness.
func TestRouterRoutesMatchSpec(t *testing.T) {
	r := chi.NewRouter()
	// harness-owned routes (excluded from codegen); the harness registers
	// them, mirroring httpserver.NewRouter.
	r.Get("/health", http.NotFoundHandler().ServeHTTP)
	r.Get("/docs/openapi.yml", http.NotFoundHandler().ServeHTTP)
	openapi.HandlerWithOptions(openapi.NewStrictHandler(stubStrict{}, nil), openapi.ChiServerOptions{
		BaseRouter: r,
	})

	walked := map[string]bool{}
	_ = chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		walked[normalizeRoute(method+" "+route)] = true
		return nil
	})

	ops, err := authz.SpecOperations(spec.OpenAPISpec)
	if err != nil {
		t.Fatalf("SpecOperations: %v", err)
	}
	expected := make(map[string]bool, len(ops))
	for _, op := range ops {
		expected[normalizeRoute(op.Method+" "+op.Path)] = true
	}

	var missing, extra []string
	for want := range expected {
		if !walked[want] {
			missing = append(missing, want)
		}
	}
	for got := range walked {
		if !expected[got] {
			extra = append(extra, got)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)
	if len(missing) > 0 || len(extra) > 0 {
		t.Fatalf("route parity mismatch\nroutes in spec but not served:\n%s\nroutes served but not in spec:\n%s",
			strings.Join(missing, "\n"), strings.Join(extra, "\n"))
	}
	t.Logf("route parity ok: %d routes", len(walked))
}

// TestAuthMatrixMatchesSpec asserts every operation's chain, composed from
// the spec's security blocks, runs exactly the middlewares the spec
// declares: public operations get none, protected operations get the JWT
// middleware.
func TestAuthMatrixMatchesSpec(t *testing.T) {
	mw := middlewares{jwt: labelMW("JWT")}
	resolver, err := authResolver(mw)
	if err != nil {
		t.Fatalf("authResolver: %v", err)
	}
	chains := resolver.Chains()

	ops, err := authz.SpecOperations(spec.OpenAPISpec)
	if err != nil {
		t.Fatalf("SpecOperations: %v", err)
	}
	var mismatches []string
	for _, op := range ops {
		var want []string
		if len(op.Schemes) > 0 {
			switch strings.Join(op.Schemes, "+") {
			case "bearerAuth":
				want = append(want, "JWT")
			default:
				mismatches = append(mismatches, op.OperationID+": unexpected scheme combination "+strings.Join(op.Schemes, "+"))
				continue
			}
		}
		if got := runChain(chains[authz.GeneratedOperationID(op.OperationID)]); !slices.Equal(got, want) {
			mismatches = append(mismatches, op.OperationID+": want "+strings.Join(want, "+")+", got "+strings.Join(got, "+"))
		}
	}
	sort.Strings(mismatches)
	if len(mismatches) > 0 {
		t.Fatalf("auth matrix mismatch\n%s", strings.Join(mismatches, "\n"))
	}
	t.Logf("auth matrix ok: %d operations match the spec", len(ops))
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
