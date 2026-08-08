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

	"lib/authz"
	spec "univents"
	"univents/internal/openapi"

	"github.com/go-chi/chi/v5"
)

// stubStrict implements openapi.StrictServerInterface for route walking;
// registration only stores handlers, they are never served.
type stubStrict struct{}

// errStub is returned by the walk-only stub; it is never served.
var errStub = errors.New("parity stub")

func (stubStrict) DeleteBadgeTemplate(_ context.Context, _ openapi.DeleteBadgeTemplateRequestObject) (openapi.DeleteBadgeTemplateResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetBadgeTemplate(_ context.Context, _ openapi.GetBadgeTemplateRequestObject) (openapi.GetBadgeTemplateResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListMyCertifications(_ context.Context, _ openapi.ListMyCertificationsRequestObject) (openapi.ListMyCertificationsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DeleteCertificationTemplate(_ context.Context, _ openapi.DeleteCertificationTemplateRequestObject) (openapi.DeleteCertificationTemplateResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetCertificationTemplate(_ context.Context, _ openapi.GetCertificationTemplateRequestObject) (openapi.GetCertificationTemplateResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpdateCertificationTemplate(_ context.Context, _ openapi.UpdateCertificationTemplateRequestObject) (openapi.UpdateCertificationTemplateResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UnlinkCertificationTemplate(_ context.Context, _ openapi.UnlinkCertificationTemplateRequestObject) (openapi.UnlinkCertificationTemplateResponseObject, error) {
	return nil, errStub
}
func (stubStrict) LinkCertificationTemplate(_ context.Context, _ openapi.LinkCertificationTemplateRequestObject) (openapi.LinkCertificationTemplateResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListCertificationTemplateLinks(_ context.Context, _ openapi.ListCertificationTemplateLinksRequestObject) (openapi.ListCertificationTemplateLinksResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetCertification(_ context.Context, _ openapi.GetCertificationRequestObject) (openapi.GetCertificationResponseObject, error) {
	return nil, errStub
}
func (stubStrict) InvalidateCertification(_ context.Context, _ openapi.InvalidateCertificationRequestObject) (openapi.InvalidateCertificationResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListBadgeTemplates(_ context.Context, _ openapi.ListBadgeTemplatesRequestObject) (openapi.ListBadgeTemplatesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateBadgeTemplate(_ context.Context, _ openapi.CreateBadgeTemplateRequestObject) (openapi.CreateBadgeTemplateResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpdateBadgeTemplate(_ context.Context, _ openapi.UpdateBadgeTemplateRequestObject) (openapi.UpdateBadgeTemplateResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListUserBadges(_ context.Context, _ openapi.ListUserBadgesRequestObject) (openapi.ListUserBadgesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListEditionBadgeEmissions(_ context.Context, _ openapi.ListEditionBadgeEmissionsRequestObject) (openapi.ListEditionBadgeEmissionsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetEditionBadgesPrint(_ context.Context, _ openapi.GetEditionBadgesPrintRequestObject) (openapi.GetEditionBadgesPrintResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListEditionCertifications(_ context.Context, _ openapi.ListEditionCertificationsRequestObject) (openapi.ListEditionCertificationsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListCertificationEmissionErrors(_ context.Context, _ openapi.ListCertificationEmissionErrorsRequestObject) (openapi.ListCertificationEmissionErrorsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListCertificationTemplates(_ context.Context, _ openapi.ListCertificationTemplatesRequestObject) (openapi.ListCertificationTemplatesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateCertificationTemplate(_ context.Context, _ openapi.CreateCertificationTemplateRequestObject) (openapi.CreateCertificationTemplateResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListEditionOccurrences(_ context.Context, _ openapi.ListEditionOccurrencesRequestObject) (openapi.ListEditionOccurrencesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListEditionProducts(_ context.Context, _ openapi.ListEditionProductsRequestObject) (openapi.ListEditionProductsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateInitialProduct(_ context.Context, _ openapi.CreateInitialProductRequestObject) (openapi.CreateInitialProductResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetProductByVendorCode(_ context.Context, _ openapi.GetProductByVendorCodeRequestObject) (openapi.GetProductByVendorCodeResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListEditionPrograms(_ context.Context, _ openapi.ListEditionProgramsRequestObject) (openapi.ListEditionProgramsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateProgram(_ context.Context, _ openapi.CreateProgramRequestObject) (openapi.CreateProgramResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListEditionSignatureRequests(_ context.Context, _ openapi.ListEditionSignatureRequestsRequestObject) (openapi.ListEditionSignatureRequestsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateSignatureRequest(_ context.Context, _ openapi.CreateSignatureRequestRequestObject) (openapi.CreateSignatureRequestResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListEditionSignatures(_ context.Context, _ openapi.ListEditionSignaturesRequestObject) (openapi.ListEditionSignaturesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateSignature(_ context.Context, _ openapi.CreateSignatureRequestObject) (openapi.CreateSignatureResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListTicketTypes(_ context.Context, _ openapi.ListTicketTypesRequestObject) (openapi.ListTicketTypesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateTicketType(_ context.Context, _ openapi.CreateTicketTypeRequestObject) (openapi.CreateTicketTypeResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetVariantByVendorCode(_ context.Context, _ openapi.GetVariantByVendorCodeRequestObject) (openapi.GetVariantByVendorCodeResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListPublicEvents(_ context.Context, _ openapi.ListPublicEventsRequestObject) (openapi.ListPublicEventsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateEvent(_ context.Context, _ openapi.CreateEventRequestObject) (openapi.CreateEventResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListJoinedEvents(_ context.Context, _ openapi.ListJoinedEventsRequestObject) (openapi.ListJoinedEventsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOwnedEvents(_ context.Context, _ openapi.ListOwnedEventsRequestObject) (openapi.ListOwnedEventsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PatchEvent(_ context.Context, _ openapi.PatchEventRequestObject) (openapi.PatchEventResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DiscontinueEvent(_ context.Context, _ openapi.DiscontinueEventRequestObject) (openapi.DiscontinueEventResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CompleteEventPayments(_ context.Context, _ openapi.CompleteEventPaymentsRequestObject) (openapi.CompleteEventPaymentsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ConnectEventPayments(_ context.Context, _ openapi.ConnectEventPaymentsRequestObject) (openapi.ConnectEventPaymentsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DisconnectEventPayments(_ context.Context, _ openapi.DisconnectEventPaymentsRequestObject) (openapi.DisconnectEventPaymentsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListPublicEditions(_ context.Context, _ openapi.ListPublicEditionsRequestObject) (openapi.ListPublicEditionsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateEdition(_ context.Context, _ openapi.CreateEditionRequestObject) (openapi.CreateEditionResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetActiveEdition(_ context.Context, _ openapi.GetActiveEditionRequestObject) (openapi.GetActiveEditionResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListDraftEditions(_ context.Context, _ openapi.ListDraftEditionsRequestObject) (openapi.ListDraftEditionsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListPastEditions(_ context.Context, _ openapi.ListPastEditionsRequestObject) (openapi.ListPastEditionsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListUpcomingEditions(_ context.Context, _ openapi.ListUpcomingEditionsRequestObject) (openapi.ListUpcomingEditionsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PatchEdition(_ context.Context, _ openapi.PatchEditionRequestObject) (openapi.PatchEditionResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PublishEdition(_ context.Context, _ openapi.PublishEditionRequestObject) (openapi.PublishEditionResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListEventMembers(_ context.Context, _ openapi.ListEventMembersRequestObject) (openapi.ListEventMembersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) AddEventMember(_ context.Context, _ openapi.AddEventMemberRequestObject) (openapi.AddEventMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RemoveEventMember(_ context.Context, _ openapi.RemoveEventMemberRequestObject) (openapi.RemoveEventMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PublishEvent(_ context.Context, _ openapi.PublishEventRequestObject) (openapi.PublishEventResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetEventBySlug(_ context.Context, _ openapi.GetEventBySlugRequestObject) (openapi.GetEventBySlugResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetEditionBySlug(_ context.Context, _ openapi.GetEditionBySlugRequestObject) (openapi.GetEditionBySlugResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DeleteOccurrence(_ context.Context, _ openapi.DeleteOccurrenceRequestObject) (openapi.DeleteOccurrenceResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetOccurrence(_ context.Context, _ openapi.GetOccurrenceRequestObject) (openapi.GetOccurrenceResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PatchOccurrence(_ context.Context, _ openapi.PatchOccurrenceRequestObject) (openapi.PatchOccurrenceResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DeleteProduct(_ context.Context, _ openapi.DeleteProductRequestObject) (openapi.DeleteProductResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetProduct(_ context.Context, _ openapi.GetProductRequestObject) (openapi.GetProductResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PatchProduct(_ context.Context, _ openapi.PatchProductRequestObject) (openapi.PatchProductResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListProductVariants(_ context.Context, _ openapi.ListProductVariantsRequestObject) (openapi.ListProductVariantsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateProductVariant(_ context.Context, _ openapi.CreateProductVariantRequestObject) (openapi.CreateProductVariantResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DeleteProgram(_ context.Context, _ openapi.DeleteProgramRequestObject) (openapi.DeleteProgramResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetProgram(_ context.Context, _ openapi.GetProgramRequestObject) (openapi.GetProgramResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PatchProgram(_ context.Context, _ openapi.PatchProgramRequestObject) (openapi.PatchProgramResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListProgramOccurrences(_ context.Context, _ openapi.ListProgramOccurrencesRequestObject) (openapi.ListProgramOccurrencesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateProgramOccurrence(_ context.Context, _ openapi.CreateProgramOccurrenceRequestObject) (openapi.CreateProgramOccurrenceResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DenySignatureRequest(_ context.Context, _ openapi.DenySignatureRequestRequestObject) (openapi.DenySignatureRequestResponseObject, error) {
	return nil, errStub
}
func (stubStrict) FulfillSignatureRequest(_ context.Context, _ openapi.FulfillSignatureRequestRequestObject) (openapi.FulfillSignatureRequestResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetSignatureRequest(_ context.Context, _ openapi.GetSignatureRequestRequestObject) (openapi.GetSignatureRequestResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CancelSignatureRequest(_ context.Context, _ openapi.CancelSignatureRequestRequestObject) (openapi.CancelSignatureRequestResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RevokeSignature(_ context.Context, _ openapi.RevokeSignatureRequestObject) (openapi.RevokeSignatureResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DeleteSignature(_ context.Context, _ openapi.DeleteSignatureRequestObject) (openapi.DeleteSignatureResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetSignature(_ context.Context, _ openapi.GetSignatureRequestObject) (openapi.GetSignatureResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetTicketType(_ context.Context, _ openapi.GetTicketTypeRequestObject) (openapi.GetTicketTypeResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PatchTicketType(_ context.Context, _ openapi.PatchTicketTypeRequestObject) (openapi.PatchTicketTypeResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DeleteProductVariant(_ context.Context, _ openapi.DeleteProductVariantRequestObject) (openapi.DeleteProductVariantResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PatchProductVariant(_ context.Context, _ openapi.PatchProductVariantRequestObject) (openapi.PatchProductVariantResponseObject, error) {
	return nil, errStub
}
func (stubStrict) VerifyCertification(_ context.Context, _ openapi.VerifyCertificationRequestObject) (openapi.VerifyCertificationResponseObject, error) {
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
	var next http.Handler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	for i := range slices.Backward(chain) {
		next = chain[i](next)
	}
	next.ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))
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
