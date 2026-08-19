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
	spec "payssage"
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
func (stubStrict) RefundIntent(_ context.Context, _ openapi.RefundIntentRequestObject) (openapi.RefundIntentResponseObject, error) {
	return nil, errStub
}
func (stubStrict) TestmodeRefundIntent(_ context.Context, _ openapi.TestmodeRefundIntentRequestObject) (openapi.TestmodeRefundIntentResponseObject, error) {
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
	mw := middlewares{jwtAuth: labelMW("JWT"), apiKeyAuth: labelMW("APIKey"), anyAuth: labelMW("Any")}
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
			case "apiKeyAuth+bearerAuth":
				want = append(want, "Any")
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
