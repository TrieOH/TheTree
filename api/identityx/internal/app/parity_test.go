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

	spec "IdentityX"
	"IdentityX/internal/openapi"
	"lib/authz"

	"github.com/go-chi/chi/v5"
)

// stubStrict implements openapi.StrictServerInterface for route walking;
// registration only stores handlers, they are never served.
type stubStrict struct{}

// errStub is returned by the walk-only stub; it is never served.
var errStub = errors.New("parity stub")

func (stubStrict) GetJWKS(_ context.Context, _ openapi.GetJWKSRequestObject) (openapi.GetJWKSResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetPlatformProfile(_ context.Context, _ openapi.GetPlatformProfileRequestObject) (openapi.GetPlatformProfileResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpsertPlatformProfile(_ context.Context, _ openapi.UpsertPlatformProfileRequestObject) (openapi.UpsertPlatformProfileResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetIntrospect(_ context.Context, _ openapi.GetIntrospectRequestObject) (openapi.GetIntrospectResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PostLogin(_ context.Context, _ openapi.PostLoginRequestObject) (openapi.PostLoginResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PostLogout(_ context.Context, _ openapi.PostLogoutRequestObject) (openapi.PostLogoutResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PostRefresh(_ context.Context, _ openapi.PostRefreshRequestObject) (openapi.PostRefreshResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PostRegister(_ context.Context, _ openapi.PostRegisterRequestObject) (openapi.PostRegisterResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetSetup(_ context.Context, _ openapi.GetSetupRequestObject) (openapi.GetSetupResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PostSetup(_ context.Context, _ openapi.PostSetupRequestObject) (openapi.PostSetupResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetOAuthCallback(_ context.Context, _ openapi.GetOAuthCallbackRequestObject) (openapi.GetOAuthCallbackResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetOAuthConnect(_ context.Context, _ openapi.GetOAuthConnectRequestObject) (openapi.GetOAuthConnectResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetOAuthProviders(_ context.Context, _ openapi.GetOAuthProvidersRequestObject) (openapi.GetOAuthProvidersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListProjectOAuthProviders(_ context.Context, _ openapi.ListProjectOAuthProvidersRequestObject) (openapi.ListProjectOAuthProvidersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateProjectOAuthProvider(_ context.Context, _ openapi.CreateProjectOAuthProviderRequestObject) (openapi.CreateProjectOAuthProviderResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpdateOAuthProvider(_ context.Context, _ openapi.UpdateOAuthProviderRequestObject) (openapi.UpdateOAuthProviderResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DeleteOAuthProvider(_ context.Context, _ openapi.DeleteOAuthProviderRequestObject) (openapi.DeleteOAuthProviderResponseObject, error) {
	return nil, errStub
}
func (stubStrict) DisableOAuthProvider(_ context.Context, _ openapi.DisableOAuthProviderRequestObject) (openapi.DisableOAuthProviderResponseObject, error) {
	return nil, errStub
}
func (stubStrict) EnableOAuthProvider(_ context.Context, _ openapi.EnableOAuthProviderRequestObject) (openapi.EnableOAuthProviderResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizations(_ context.Context, _ openapi.ListOrganizationsRequestObject) (openapi.ListOrganizationsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateOrganization(_ context.Context, _ openapi.CreateOrganizationRequestObject) (openapi.CreateOrganizationResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizationProjectActors(_ context.Context, _ openapi.ListOrganizationProjectActorsRequestObject) (openapi.ListOrganizationProjectActorsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateOrganizationProjectActor(_ context.Context, _ openapi.CreateOrganizationProjectActorRequestObject) (openapi.CreateOrganizationProjectActorResponseObject, error) {
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
func (stubStrict) ListOrganizationProjects(_ context.Context, _ openapi.ListOrganizationProjectsRequestObject) (openapi.ListOrganizationProjectsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateOrganizationProject(_ context.Context, _ openapi.CreateOrganizationProjectRequestObject) (openapi.CreateOrganizationProjectResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizationProjectMembers(_ context.Context, _ openapi.ListOrganizationProjectMembersRequestObject) (openapi.ListOrganizationProjectMembersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) AddOrganizationProjectMember(_ context.Context, _ openapi.AddOrganizationProjectMemberRequestObject) (openapi.AddOrganizationProjectMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RemoveOrganizationProjectMember(_ context.Context, _ openapi.RemoveOrganizationProjectMemberRequestObject) (openapi.RemoveOrganizationProjectMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetOrganizationProjectActor(_ context.Context, _ openapi.GetOrganizationProjectActorRequestObject) (openapi.GetOrganizationProjectActorResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListProjects(_ context.Context, _ openapi.ListProjectsRequestObject) (openapi.ListProjectsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateProject(_ context.Context, _ openapi.CreateProjectRequestObject) (openapi.CreateProjectResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListProjectMembers(_ context.Context, _ openapi.ListProjectMembersRequestObject) (openapi.ListProjectMembersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) AddProjectMember(_ context.Context, _ openapi.AddProjectMemberRequestObject) (openapi.AddProjectMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RemoveProjectMember(_ context.Context, _ openapi.RemoveProjectMemberRequestObject) (openapi.RemoveProjectMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListActors(_ context.Context, _ openapi.ListActorsRequestObject) (openapi.ListActorsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateActor(_ context.Context, _ openapi.CreateActorRequestObject) (openapi.CreateActorResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetActor(_ context.Context, _ openapi.GetActorRequestObject) (openapi.GetActorResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetActorByEmail(_ context.Context, _ openapi.GetActorByEmailRequestObject) (openapi.GetActorByEmailResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateAPIKey(_ context.Context, _ openapi.CreateAPIKeyRequestObject) (openapi.CreateAPIKeyResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListCapabilities(_ context.Context, _ openapi.ListCapabilitiesRequestObject) (openapi.ListCapabilitiesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateCapability(_ context.Context, _ openapi.CreateCapabilityRequestObject) (openapi.CreateCapabilityResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetPlatformProfileSchema(_ context.Context, _ openapi.GetPlatformProfileSchemaRequestObject) (openapi.GetPlatformProfileSchemaResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpsertPlatformProfileSchema(_ context.Context, _ openapi.UpsertPlatformProfileSchemaRequestObject) (openapi.UpsertPlatformProfileSchemaResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetProjectProfileSchema(_ context.Context, _ openapi.GetProjectProfileSchemaRequestObject) (openapi.GetProjectProfileSchemaResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpsertProjectProfileSchema(_ context.Context, _ openapi.UpsertProjectProfileSchemaRequestObject) (openapi.UpsertProjectProfileSchemaResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetProjectProfile(_ context.Context, _ openapi.GetProjectProfileRequestObject) (openapi.GetProjectProfileResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpsertProjectProfile(_ context.Context, _ openapi.UpsertProjectProfileRequestObject) (openapi.UpsertProjectProfileResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOutdatedPlatformProfiles(_ context.Context, _ openapi.ListOutdatedPlatformProfilesRequestObject) (openapi.ListOutdatedPlatformProfilesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOutdatedProjectProfiles(_ context.Context, _ openapi.ListOutdatedProjectProfilesRequestObject) (openapi.ListOutdatedProjectProfilesResponseObject, error) {
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
// the spec's security blocks, runs the middlewares the spec declares: the
// setup guard (except the two setup routes), the scheme combination's
// middleware, and the platform-client scope guard for client-only
// operations. Chains are keyed by generated-form operationID; the
// client-only list is validated against the spec inside the resolver, so
// this test no longer re-derives expectations from the same list it feeds.
func TestAuthMatrixMatchesSpec(t *testing.T) {
	mw := middlewares{
		jwtAuth:    labelMW("JWT"),
		apiKeyAuth: labelMW("APIKey"),
		anyAuth:    labelMW("AnyAuth"),
		clientOnly: labelMW("ClientOnly"),
	}
	resolver, err := authResolver(mw, labelMW("setupGuard"))
	if err != nil {
		t.Fatalf("authResolver: %v", err)
	}
	chains := resolver.Chains()

	ops, err := authz.SpecOperations(spec.OpenAPISpec)
	if err != nil {
		t.Fatalf("SpecOperations: %v", err)
	}
	clientOnly := make(map[string]bool, len(clientOnlyOps))
	for _, op := range clientOnlyOps {
		clientOnly[op] = true
	}

	var mismatches []string
	for _, op := range ops {
		var want []string
		if op.OperationID != "getSetup" && op.OperationID != "postSetup" {
			want = append(want, "setupGuard")
		}
		if len(op.Schemes) > 0 {
			switch strings.Join(op.Schemes, "+") {
			case "bearerAuth":
				want = append(want, "JWT")
			case "apiKeyAuth+bearerAuth":
				want = append(want, "AnyAuth")
			default:
				mismatches = append(mismatches, op.OperationID+": unexpected scheme combination "+strings.Join(op.Schemes, "+"))
				continue
			}
		}
		if clientOnly[op.OperationID] {
			want = append(want, "ClientOnly")
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
