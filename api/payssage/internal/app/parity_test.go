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

	"payssage/internal/openapi"

	"github.com/go-chi/chi/v5"
)

// stubStrict implements openapi.StrictServerInterface for route walking;
// registration only stores handlers, they are never served.
type stubStrict struct{}

// errStub is returned by the walk-only stub; it is never served.
var errStub = errors.New("parity stub")

func (stubStrict) ListCollectors(_ context.Context, _ openapi.ListCollectorsRequestObject) (openapi.ListCollectorsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetCollector(_ context.Context, _ openapi.GetCollectorRequestObject) (openapi.GetCollectorResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListIntentsByProfile(_ context.Context, _ openapi.ListIntentsByProfileRequestObject) (openapi.ListIntentsByProfileResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetIntent(_ context.Context, _ openapi.GetIntentRequestObject) (openapi.GetIntentResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CancelIntent(_ context.Context, _ openapi.CancelIntentRequestObject) (openapi.CancelIntentResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizations(_ context.Context, _ openapi.ListOrganizationsRequestObject) (openapi.ListOrganizationsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateOrganization(_ context.Context, _ openapi.CreateOrganizationRequestObject) (openapi.CreateOrganizationResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizationCollectors(_ context.Context, _ openapi.ListOrganizationCollectorsRequestObject) (openapi.ListOrganizationCollectorsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizationIntents(_ context.Context, _ openapi.ListOrganizationIntentsRequestObject) (openapi.ListOrganizationIntentsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetOrganizationMemberByEmail(_ context.Context, _ openapi.GetOrganizationMemberByEmailRequestObject) (openapi.GetOrganizationMemberByEmailResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetOrganizationMemberByID(_ context.Context, _ openapi.GetOrganizationMemberByIDRequestObject) (openapi.GetOrganizationMemberByIDResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RemoveOrganizationMember(_ context.Context, _ openapi.RemoveOrganizationMemberRequestObject) (openapi.RemoveOrganizationMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizationMembers(_ context.Context, _ openapi.ListOrganizationMembersRequestObject) (openapi.ListOrganizationMembersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) AddOrganizationMember(_ context.Context, _ openapi.AddOrganizationMemberRequestObject) (openapi.AddOrganizationMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizationWallets(_ context.Context, _ openapi.ListOrganizationWalletsRequestObject) (openapi.ListOrganizationWalletsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ProviderCallback(_ context.Context, _ openapi.ProviderCallbackRequestObject) (openapi.ProviderCallbackResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ConnectProvider(_ context.Context, _ openapi.ConnectProviderRequestObject) (openapi.ConnectProviderResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RevokeProvider(_ context.Context, _ openapi.RevokeProviderRequestObject) (openapi.RevokeProviderResponseObject, error) {
	return nil, errStub
}
func (stubStrict) HardCreateIntent(_ context.Context, _ openapi.HardCreateIntentRequestObject) (openapi.HardCreateIntentResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListWallets(_ context.Context, _ openapi.ListWalletsRequestObject) (openapi.ListWalletsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateWallet(_ context.Context, _ openapi.CreateWalletRequestObject) (openapi.CreateWalletResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetWallet(_ context.Context, _ openapi.GetWalletRequestObject) (openapi.GetWalletResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UnbindCollector(_ context.Context, _ openapi.UnbindCollectorRequestObject) (openapi.UnbindCollectorResponseObject, error) {
	return nil, errStub
}
func (stubStrict) BindCollector(_ context.Context, _ openapi.BindCollectorRequestObject) (openapi.BindCollectorResponseObject, error) {
	return nil, errStub
}
func (stubStrict) SetWalletFee(_ context.Context, _ openapi.SetWalletFeeRequestObject) (openapi.SetWalletFeeResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListWalletIntents(_ context.Context, _ openapi.ListWalletIntentsRequestObject) (openapi.ListWalletIntentsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) Checkout(_ context.Context, _ openapi.CheckoutRequestObject) (openapi.CheckoutResponseObject, error) {
	return nil, errStub
}
func (stubStrict) SetWalletSandbox(_ context.Context, _ openapi.SetWalletSandboxRequestObject) (openapi.SetWalletSandboxResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListWalletSellers(_ context.Context, _ openapi.ListWalletSellersRequestObject) (openapi.ListWalletSellersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListWebhookEndpoints(_ context.Context, _ openapi.ListWebhookEndpointsRequestObject) (openapi.ListWebhookEndpointsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateWebhookEndpoint(_ context.Context, _ openapi.CreateWebhookEndpointRequestObject) (openapi.CreateWebhookEndpointResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListWebhookEvents(_ context.Context, _ openapi.ListWebhookEventsRequestObject) (openapi.ListWebhookEventsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetWebhookDelivery(_ context.Context, _ openapi.GetWebhookDeliveryRequestObject) (openapi.GetWebhookDeliveryResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DeleteWebhookEndpoint(_ context.Context, _ openapi.DeleteWebhookEndpointRequestObject) (openapi.DeleteWebhookEndpointResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetWebhookEndpoint(_ context.Context, _ openapi.GetWebhookEndpointRequestObject) (openapi.GetWebhookEndpointResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListWebhookDeliveries(_ context.Context, _ openapi.ListWebhookDeliveriesRequestObject) (openapi.ListWebhookDeliveriesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetWebhookEvent(_ context.Context, _ openapi.GetWebhookEventRequestObject) (openapi.GetWebhookEventResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ReceiveWebhook(_ context.Context, _ openapi.ReceiveWebhookRequestObject) (openapi.ReceiveWebhookResponseObject, error) {
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
	"DELETE /organizations/{organization_id}/members": "removeOrganizationMember",
	"DELETE /wallets/{wallet_id}/collector":           "unbindCollector",
	"DELETE /webhooks/endpoints/{endpoint_id}":        "deleteWebhookEndpoint",
	"GET /collectors":                                 "listCollectors",
	"GET /collectors/{collector_id}":                  "getCollector",
	"GET /docs/openapi.yml":                           "getOpenAPISpec",
	"GET /health":                                     "getHealth",
	"GET /intents":                                    "listIntentsByProfile",
	"GET /intents/{intent_id}":                        "getIntent",
	"GET /organizations":                              "listOrganizations",
	"GET /organizations/{organization_id}/collectors": "listOrganizationCollectors",
	"GET /organizations/{organization_id}/intents":    "listOrganizationIntents",
	"GET /organizations/{organization_id}/member/{member_email}:by_email": "getOrganizationMemberByEmail",
	"GET /organizations/{organization_id}/member/{member_id}":             "getOrganizationMemberByID",
	"GET /organizations/{organization_id}/members":                        "listOrganizationMembers",
	"GET /organizations/{organization_id}/wallets":                        "listOrganizationWallets",
	"GET /providers/{provider}/callback":                                  "providerCallback",
	"GET /wallets":                                                        "listWallets",
	"GET /wallets/{wallet_id}":                                            "getWallet",
	"GET /wallets/{wallet_id}/intents":                                    "listWalletIntents",
	"GET /wallets/{wallet_id}/sellers":                                    "listWalletSellers",
	"GET /wallets/{wallet_id}/webhooks/endpoints":                         "listWebhookEndpoints",
	"GET /wallets/{wallet_id}/webhooks/events":                            "listWebhookEvents",
	"GET /webhooks/deliveries/{delivery_id}":                              "getWebhookDelivery",
	"GET /webhooks/endpoints/{endpoint_id}":                               "getWebhookEndpoint",
	"GET /webhooks/endpoints/{endpoint_id}/deliveries":                    "listWebhookDeliveries",
	"GET /webhooks/events/{event_id}":                                     "getWebhookEvent",
	"PATCH /wallets/{wallet_id}/fee":                                      "setWalletFee",
	"PATCH /wallets/{wallet_id}/sandbox":                                  "setWalletSandbox",
	"POST /intents/{intent_id}/cancel":                                    "cancelIntent",
	"POST /organizations":                                                 "createOrganization",
	"POST /organizations/{organization_id}/members":                       "addOrganizationMember",
	"POST /providers/{provider}/connect":                                  "connectProvider",
	"POST /providers/{provider}/revoke":                                   "revokeProvider",
	"POST /testmode/intents/create":                                       "hardCreateIntent",
	"POST /wallets":                                                       "createWallet",
	"POST /wallets/{wallet_id}/collector":                                 "bindCollector",
	"POST /wallets/{wallet_id}/intents":                                   "checkout",
	"POST /wallets/{wallet_id}/webhooks/endpoints":                        "createWebhookEndpoint",
	"POST /webhooks/{provider}":                                           "receiveWebhook",
}

// expectedOps is the auth matrix keyed by operationId: "public" when the
// operation runs with no auth middleware, otherwise the chain names joined
// with "+". Harness-owned routes (getHealth, getOpenAPISpec) are excluded —
// they never run through the dispatch.
var expectedOps = map[string]string{
	"addOrganizationMember":        "JWT",
	"bindCollector":                "JWT",
	"cancelIntent":                 "JWT",
	"checkout":                     "JWT",
	"connectProvider":              "JWT",
	"createOrganization":           "JWT",
	"createWallet":                 "JWT",
	"createWebhookEndpoint":        "JWT",
	"deleteWebhookEndpoint":        "JWT",
	"getCollector":                 "JWT",
	"getIntent":                    "JWT",
	"getOrganizationMemberByEmail": "JWT",
	"getOrganizationMemberByID":    "JWT",
	"getWallet":                    "JWT",
	"getWebhookDelivery":           "JWT",
	"getWebhookEndpoint":           "JWT",
	"getWebhookEvent":              "JWT",
	"hardCreateIntent":             "JWT",
	"listCollectors":               "JWT",
	"listIntentsByProfile":         "JWT",
	"listOrganizationCollectors":   "JWT",
	"listOrganizationIntents":      "JWT",
	"listOrganizationMembers":      "JWT",
	"listOrganizationWallets":      "JWT",
	"listOrganizations":            "JWT",
	"listWalletIntents":            "JWT",
	"listWalletSellers":            "JWT",
	"listWallets":                  "JWT",
	"listWebhookDeliveries":        "JWT",
	"listWebhookEndpoints":         "JWT",
	"listWebhookEvents":            "JWT",
	"providerCallback":             "public",
	"removeOrganizationMember":     "JWT",
	"revokeProvider":               "JWT",
	"setWalletFee":                 "JWT",
	"setWalletSandbox":             "JWT",
	"unbindCollector":              "JWT",
	"receiveWebhook":               "public",
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

	chains := authChains(middlewares{jwtAuth: mwJWT})
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
