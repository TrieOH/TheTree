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

func mwJWT(next http.Handler) http.Handler { return next }

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
	"DELETE /badges/{template_id}":                              "deleteBadgeTemplate",
	"DELETE /certifications/templates/{template_id}":            "deleteCertificationTemplate",
	"DELETE /certifications/templates/{template_id}/link":       "unlinkCertificationTemplate",
	"DELETE /events/{event_id}/members/{user_id}":               "removeEventMember",
	"DELETE /occurrences/{occurrence_id}":                       "deleteOccurrence",
	"DELETE /products/{product_id}":                             "deleteProduct",
	"DELETE /programs/{program_id}":                             "deleteProgram",
	"DELETE /signatures/{signature_id}":                         "deleteSignature",
	"DELETE /variants/{variant_id}":                             "deleteProductVariant",
	"GET /badges/{template_id}":                                 "getBadgeTemplate",
	"GET /certifications":                                       "listMyCertifications",
	"GET /certifications/templates/{template_id}":               "getCertificationTemplate",
	"GET /certifications/templates/{template_id}/links":         "listCertificationTemplateLinks",
	"GET /certifications/{cert_id}":                             "getCertification",
	"GET /docs/openapi.yml":                                     "getOpenAPISpec",
	"GET /editions/{edition_id}/badges":                         "listBadgeTemplates",
	"GET /editions/{edition_id}/certifications":                 "listEditionCertifications",
	"GET /editions/{edition_id}/certifications/emission-errors": "listCertificationEmissionErrors",
	"GET /editions/{edition_id}/certifications/templates":       "listCertificationTemplates",
	"GET /editions/{edition_id}/occurrences":                    "listEditionOccurrences",
	"GET /editions/{edition_id}/products":                       "listEditionProducts",
	"GET /editions/{edition_id}/products/{vendor_code}:by-code": "getProductByVendorCode",
	"GET /editions/{edition_id}/programs":                       "listEditionPrograms",
	"GET /editions/{edition_id}/signature-requests":             "listEditionSignatureRequests",
	"GET /editions/{edition_id}/signatures":                     "listEditionSignatures",
	"GET /editions/{edition_id}/ticket-types":                   "listTicketTypes",
	"GET /editions/{edition_id}/variants/{vendor_code}:by-code": "getVariantByVendorCode",
	"GET /events":                                                      "listPublicEvents",
	"GET /events/joined":                                               "listJoinedEvents",
	"GET /events/owned":                                                "listOwnedEvents",
	"GET /events/{event_id}/editions":                                  "listPublicEditions",
	"GET /events/{event_id}/editions/active":                           "getActiveEdition",
	"GET /events/{event_id}/editions/draft":                            "listDraftEditions",
	"GET /events/{event_id}/editions/past":                             "listPastEditions",
	"GET /events/{event_id}/editions/upcoming":                         "listUpcomingEditions",
	"GET /events/{event_id}/members":                                   "listEventMembers",
	"GET /events/{event_slug}:by-slug":                                 "getEventBySlug",
	"GET /events/{event_slug}:by-slug/editions/{edition_slug}:by-slug": "getEditionBySlug",
	"GET /health":                                                      "getHealth",
	"GET /occurrences/{occurrence_id}":                                 "getOccurrence",
	"GET /products/{product_id}":                                       "getProduct",
	"GET /products/{product_id}/variants":                              "listProductVariants",
	"GET /programs/{program_id}":                                       "getProgram",
	"GET /programs/{program_id}/occurrences":                           "listProgramOccurrences",
	"GET /signature-requests/{request_id}":                             "getSignatureRequest",
	"GET /signatures/{signature_id}":                                   "getSignature",
	"GET /ticket-types/{ticket_type_id}":                               "getTicketType",
	"GET /verify/{hash}":                                               "verifyCertification",
	"PATCH /events/{event_id}":                                         "patchEvent",
	"PATCH /events/{event_id}/editions/{edition_id}":                   "patchEdition",
	"PATCH /occurrences/{occurrence_id}":                               "patchOccurrence",
	"PATCH /products/{product_id}":                                     "patchProduct",
	"PATCH /programs/{program_id}":                                     "patchProgram",
	"PATCH /ticket-types/{ticket_type_id}":                             "patchTicketType",
	"PATCH /variants/{variant_id}":                                     "patchProductVariant",
	"POST /certifications/templates/{template_id}/link":                "linkCertificationTemplate",
	"POST /certifications/{cert_id}/invalidate":                        "invalidateCertification",
	"POST /editions/{edition_id}/badges":                               "createBadgeTemplate",
	"POST /editions/{edition_id}/certifications/templates":             "createCertificationTemplate",
	"POST /editions/{edition_id}/products":                             "createInitialProduct",
	"POST /editions/{edition_id}/programs":                             "createProgram",
	"POST /editions/{edition_id}/signature-requests":                   "createSignatureRequest",
	"POST /editions/{edition_id}/signatures":                           "createSignature",
	"POST /editions/{edition_id}/ticket-types":                         "createTicketType",
	"POST /events":                                                     "createEvent",
	"POST /events/{event_id}/discontinue":                              "discontinueEvent",
	"POST /events/{event_id}/editions":                                 "createEdition",
	"POST /events/{event_id}/editions/{edition_id}/publish":            "publishEdition",
	"POST /events/{event_id}/members":                                  "addEventMember",
	"POST /events/{event_id}/publish":                                  "publishEvent",
	"POST /products/{product_id}/variants":                             "createProductVariant",
	"POST /programs/{program_id}/occurrences":                          "createProgramOccurrence",
	"POST /signature-requests/deny":                                    "denySignatureRequest",
	"POST /signature-requests/fulfill":                                 "fulfillSignatureRequest",
	"POST /signature-requests/{request_id}/cancel":                     "cancelSignatureRequest",
	"POST /signatures/revoke":                                          "revokeSignature",
	"PUT /certifications/templates/{template_id}":                      "updateCertificationTemplate",
}

// expectedOps is the auth matrix keyed by operationId: "public" when the
// operation runs with no auth middleware, otherwise the chain names joined
// with "+". Harness-owned routes (getHealth, getOpenAPISpec) are excluded —
// they never run through the dispatch.
var expectedOps = map[string]string{
	"addEventMember":                  "JWT",
	"cancelSignatureRequest":          "JWT",
	"createBadgeTemplate":             "JWT",
	"createCertificationTemplate":     "JWT",
	"createEdition":                   "JWT",
	"createEvent":                     "JWT",
	"createInitialProduct":            "JWT",
	"createProductVariant":            "JWT",
	"createProgram":                   "JWT",
	"createProgramOccurrence":         "JWT",
	"createSignature":                 "JWT",
	"createSignatureRequest":          "JWT",
	"createTicketType":                "JWT",
	"deleteBadgeTemplate":             "JWT",
	"deleteCertificationTemplate":     "JWT",
	"deleteOccurrence":                "JWT",
	"deleteProduct":                   "JWT",
	"deleteProductVariant":            "JWT",
	"deleteProgram":                   "JWT",
	"deleteSignature":                 "JWT",
	"discontinueEvent":                "JWT",
	"getBadgeTemplate":                "JWT",
	"getCertification":                "JWT",
	"invalidateCertification":         "JWT",
	"linkCertificationTemplate":       "JWT",
	"listBadgeTemplates":              "JWT",
	"listCertificationEmissionErrors": "JWT",
	"listDraftEditions":               "JWT",
	"listEditionCertifications":       "JWT",
	"listEventMembers":                "JWT",
	"listJoinedEvents":                "JWT",
	"listMyCertifications":            "JWT",
	"listOwnedEvents":                 "JWT",
	"patchEdition":                    "JWT",
	"patchEvent":                      "JWT",
	"patchOccurrence":                 "JWT",
	"patchProduct":                    "JWT",
	"patchProductVariant":             "JWT",
	"patchProgram":                    "JWT",
	"patchTicketType":                 "JWT",
	"publishEdition":                  "JWT",
	"publishEvent":                    "JWT",
	"removeEventMember":               "JWT",
	"unlinkCertificationTemplate":     "JWT",
	"updateCertificationTemplate":     "JWT",
	"denySignatureRequest":            "public",
	"fulfillSignatureRequest":         "public",
	"getActiveEdition":                "public",
	"getCertificationTemplate":        "public",
	"getEditionBySlug":                "public",
	"getEventBySlug":                  "public",
	"getOccurrence":                   "public",
	"getProduct":                      "public",
	"getProductByVendorCode":          "public",
	"getProgram":                      "public",
	"getSignature":                    "public",
	"getSignatureRequest":             "public",
	"getTicketType":                   "public",
	"getVariantByVendorCode":          "public",
	"listCertificationTemplateLinks":  "public",
	"listCertificationTemplates":      "public",
	"listEditionOccurrences":          "public",
	"listEditionProducts":             "public",
	"listEditionPrograms":             "public",
	"listEditionSignatureRequests":    "public",
	"listEditionSignatures":           "public",
	"listPastEditions":                "public",
	"listProductVariants":             "public",
	"listProgramOccurrences":          "public",
	"listPublicEditions":              "public",
	"listPublicEvents":                "public",
	"listTicketTypes":                 "public",
	"listUpcomingEditions":            "public",
	"revokeSignature":                 "public",
	"verifyCertification":             "public",
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

	chains := authChains(middlewares{jwt: mwJWT})
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
