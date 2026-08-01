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

	"IdentityX/internal/handler"

	"github.com/go-chi/chi/v5"
)

// stubStrict implements handler.StrictServerInterface for route walking;
// registration only stores handlers, they are never served.
type stubStrict struct{}

// errStub is returned by the walk-only stub; it is never served.
var errStub = errors.New("parity stub")

func (stubStrict) GetJWKS(_ context.Context, _ handler.GetJWKSRequestObject) (handler.GetJWKSResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetPlatformProfile(_ context.Context, _ handler.GetPlatformProfileRequestObject) (handler.GetPlatformProfileResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpsertPlatformProfile(_ context.Context, _ handler.UpsertPlatformProfileRequestObject) (handler.UpsertPlatformProfileResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetIntrospect(_ context.Context, _ handler.GetIntrospectRequestObject) (handler.GetIntrospectResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PostLogin(_ context.Context, _ handler.PostLoginRequestObject) (handler.PostLoginResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PostLogout(_ context.Context, _ handler.PostLogoutRequestObject) (handler.PostLogoutResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PostRefresh(_ context.Context, _ handler.PostRefreshRequestObject) (handler.PostRefreshResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PostRegister(_ context.Context, _ handler.PostRegisterRequestObject) (handler.PostRegisterResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetSetup(_ context.Context, _ handler.GetSetupRequestObject) (handler.GetSetupResponseObject, error) {
	return nil, errStub
}
func (stubStrict) PostSetup(_ context.Context, _ handler.PostSetupRequestObject) (handler.PostSetupResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetOAuthCallback(_ context.Context, _ handler.GetOAuthCallbackRequestObject) (handler.GetOAuthCallbackResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetOAuthConnect(_ context.Context, _ handler.GetOAuthConnectRequestObject) (handler.GetOAuthConnectResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizations(_ context.Context, _ handler.ListOrganizationsRequestObject) (handler.ListOrganizationsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateOrganization(_ context.Context, _ handler.CreateOrganizationRequestObject) (handler.CreateOrganizationResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizationProjectActors(_ context.Context, _ handler.ListOrganizationProjectActorsRequestObject) (handler.ListOrganizationProjectActorsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateOrganizationProjectActor(_ context.Context, _ handler.CreateOrganizationProjectActorRequestObject) (handler.CreateOrganizationProjectActorResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RemoveOrganizationMember(_ context.Context, _ handler.RemoveOrganizationMemberRequestObject) (handler.RemoveOrganizationMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizationMembers(_ context.Context, _ handler.ListOrganizationMembersRequestObject) (handler.ListOrganizationMembersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) AddOrganizationMember(_ context.Context, _ handler.AddOrganizationMemberRequestObject) (handler.AddOrganizationMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizationProjects(_ context.Context, _ handler.ListOrganizationProjectsRequestObject) (handler.ListOrganizationProjectsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateOrganizationProject(_ context.Context, _ handler.CreateOrganizationProjectRequestObject) (handler.CreateOrganizationProjectResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetOrganizationProjectActor(_ context.Context, _ handler.GetOrganizationProjectActorRequestObject) (handler.GetOrganizationProjectActorResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RemoveOrganizationProjectMember(_ context.Context, _ handler.RemoveOrganizationProjectMemberRequestObject) (handler.RemoveOrganizationProjectMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListOrganizationProjectMembers(_ context.Context, _ handler.ListOrganizationProjectMembersRequestObject) (handler.ListOrganizationProjectMembersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) AddOrganizationProjectMember(_ context.Context, _ handler.AddOrganizationProjectMemberRequestObject) (handler.AddOrganizationProjectMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetPlatformProfileSchema(_ context.Context, _ handler.GetPlatformProfileSchemaRequestObject) (handler.GetPlatformProfileSchemaResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpsertPlatformProfileSchema(_ context.Context, _ handler.UpsertPlatformProfileSchemaRequestObject) (handler.UpsertPlatformProfileSchemaResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListProjects(_ context.Context, _ handler.ListProjectsRequestObject) (handler.ListProjectsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateProject(_ context.Context, _ handler.CreateProjectRequestObject) (handler.CreateProjectResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListActors(_ context.Context, _ handler.ListActorsRequestObject) (handler.ListActorsResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateActor(_ context.Context, _ handler.CreateActorRequestObject) (handler.CreateActorResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetActorByEmail(_ context.Context, _ handler.GetActorByEmailRequestObject) (handler.GetActorByEmailResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetActor(_ context.Context, _ handler.GetActorRequestObject) (handler.GetActorResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetProjectProfile(_ context.Context, _ handler.GetProjectProfileRequestObject) (handler.GetProjectProfileResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpsertProjectProfile(_ context.Context, _ handler.UpsertProjectProfileRequestObject) (handler.UpsertProjectProfileResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateAPIKey(_ context.Context, _ handler.CreateAPIKeyRequestObject) (handler.CreateAPIKeyResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListCapabilities(_ context.Context, _ handler.ListCapabilitiesRequestObject) (handler.ListCapabilitiesResponseObject, error) {
	return nil, errStub
}
func (stubStrict) CreateCapability(_ context.Context, _ handler.CreateCapabilityRequestObject) (handler.CreateCapabilityResponseObject, error) {
	return nil, errStub
}
func (stubStrict) RemoveProjectMember(_ context.Context, _ handler.RemoveProjectMemberRequestObject) (handler.RemoveProjectMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) ListProjectMembers(_ context.Context, _ handler.ListProjectMembersRequestObject) (handler.ListProjectMembersResponseObject, error) {
	return nil, errStub
}
func (stubStrict) AddProjectMember(_ context.Context, _ handler.AddProjectMemberRequestObject) (handler.AddProjectMemberResponseObject, error) {
	return nil, errStub
}
func (stubStrict) GetProjectProfileSchema(_ context.Context, _ handler.GetProjectProfileSchemaRequestObject) (handler.GetProjectProfileSchemaResponseObject, error) {
	return nil, errStub
}
func (stubStrict) UpsertProjectProfileSchema(_ context.Context, _ handler.UpsertProjectProfileSchemaRequestObject) (handler.UpsertProjectProfileSchemaResponseObject, error) {
	return nil, errStub
}

func mwJWT(next http.Handler) http.Handler        { return next }
func mwAnyAuth(next http.Handler) http.Handler    { return next }
func mwClientOnly(next http.Handler) http.Handler { return next }

func mwName(mw func(http.Handler) http.Handler) string {
	fn := runtime.FuncForPC(reflect.ValueOf(mw).Pointer())
	name := fn.Name()
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	return strings.TrimPrefix(name, "mw")
}

var routeOperation = map[string]string{
	"GET /health":                                  "getHealth",
	"GET /docs/openapi.yml":                        "getOpenAPISpec",
	"GET /auth/setup":                              "getSetup",
	"POST /auth/setup":                             "postSetup",
	"POST /auth/register":                          "postRegister",
	"POST /auth/login":                             "postLogin",
	"POST /auth/logout":                            "postLogout",
	"POST /auth/refresh":                           "postRefresh",
	"GET /auth/{provider}/connect":                 "getOAuthConnect",
	"GET /auth/{provider}/callback":                "getOAuthCallback",
	"GET /.well-known/jwks.json":                   "getJWKS",
	"GET /auth/introspect":                         "getIntrospect",
	"GET /projects/{project_id}/actors/{actor_id}": "getActor",
	"GET /projects/{project_id}/actors/{actor_email}:by_email":                     "getActorByEmail",
	"POST /projects/{project_id}/actors":                                           "createActor",
	"GET /projects/{project_id}/actors":                                            "listActors",
	"POST /projects/{project_id}/api_keys":                                         "createAPIKey",
	"GET /projects/{project_id}/capabilities":                                      "listCapabilities",
	"POST /projects/{project_id}/capabilities":                                     "createCapability",
	"GET /organizations":                                                           "listOrganizations",
	"POST /organizations":                                                          "createOrganization",
	"GET /organizations/{organization_id}/members":                                 "listOrganizationMembers",
	"POST /organizations/{organization_id}/members":                                "addOrganizationMember",
	"DELETE /organizations/{organization_id}/members":                              "removeOrganizationMember",
	"GET /organizations/{organization_id}/projects":                                "listOrganizationProjects",
	"POST /organizations/{organization_id}/projects":                               "createOrganizationProject",
	"POST /organizations/{org_id}/projects/{project_id}/actors":                    "createOrganizationProjectActor",
	"GET /organizations/{org_id}/projects/{project_id}/actors":                     "listOrganizationProjectActors",
	"GET /organizations/{organization_id}/projects/{project_id}/members":           "listOrganizationProjectMembers",
	"POST /organizations/{organization_id}/projects/{project_id}/members":          "addOrganizationProjectMember",
	"DELETE /organizations/{organization_id}/projects/{project_id}/members":        "removeOrganizationProjectMember",
	"GET /organizations/{organization_id}/projects/{project_id}/actors/{actor_id}": "getOrganizationProjectActor",
	"GET /projects":                                        "listProjects",
	"POST /projects":                                       "createProject",
	"GET /projects/{project_id}/members":                   "listProjectMembers",
	"POST /projects/{project_id}/members":                  "addProjectMember",
	"DELETE /projects/{project_id}/members":                "removeProjectMember",
	"GET /actors/{actor_id}/profile":                       "getPlatformProfile",
	"PUT /actors/{actor_id}/profile":                       "upsertPlatformProfile",
	"GET /projects/{project_id}/actors/{actor_id}/profile": "getProjectProfile",
	"PUT /projects/{project_id}/actors/{actor_id}/profile": "upsertProjectProfile",
	"GET /projects/{project_id}/profile-schema":            "getProjectProfileSchema",
	"PUT /projects/{project_id}/profile-schema":            "upsertProjectProfileSchema",
	"GET /profile-schema":                                  "getPlatformProfileSchema",
	"PUT /profile-schema":                                  "upsertPlatformProfileSchema",
}

//nolint:gosec // auth-chain labels ("JWT"), not credentials
var expectedOps = map[string]string{
	"getSetup":                        "public",
	"postSetup":                       "public",
	"postRegister":                    "public",
	"postLogin":                       "public",
	"postRefresh":                     "public",
	"getOAuthConnect":                 "public",
	"getOAuthCallback":                "public",
	"getJWKS":                         "public",
	"postLogout":                      "JWT",
	"getIntrospect":                   "AnyAuth",
	"getActor":                        "AnyAuth+ClientOnly",
	"getActorByEmail":                 "AnyAuth+ClientOnly",
	"createActor":                     "AnyAuth+ClientOnly",
	"listActors":                      "AnyAuth+ClientOnly",
	"createAPIKey":                    "AnyAuth+ClientOnly",
	"listCapabilities":                "AnyAuth",
	"createCapability":                "JWT+ClientOnly",
	"listOrganizations":               "JWT+ClientOnly",
	"createOrganization":              "JWT+ClientOnly",
	"listOrganizationMembers":         "JWT+ClientOnly",
	"addOrganizationMember":           "JWT+ClientOnly",
	"removeOrganizationMember":        "JWT+ClientOnly",
	"listOrganizationProjects":        "JWT+ClientOnly",
	"createOrganizationProject":       "JWT+ClientOnly",
	"createOrganizationProjectActor":  "JWT+ClientOnly",
	"listOrganizationProjectActors":   "JWT+ClientOnly",
	"listOrganizationProjectMembers":  "JWT+ClientOnly",
	"addOrganizationProjectMember":    "JWT+ClientOnly",
	"removeOrganizationProjectMember": "JWT+ClientOnly",
	"getOrganizationProjectActor":     "JWT+ClientOnly",
	"listProjects":                    "AnyAuth+ClientOnly",
	"createProject":                   "AnyAuth+ClientOnly",
	"listProjectMembers":              "AnyAuth+ClientOnly",
	"addProjectMember":                "AnyAuth+ClientOnly",
	"removeProjectMember":             "AnyAuth+ClientOnly",
	"getPlatformProfile":              "JWT+ClientOnly",
	"upsertPlatformProfile":           "JWT+ClientOnly",
	"getProjectProfile":               "JWT+ClientOnly",
	"upsertProjectProfile":            "JWT+ClientOnly",
	"getProjectProfileSchema":         "JWT+ClientOnly",
	"upsertProjectProfileSchema":      "JWT+ClientOnly",
	"getPlatformProfileSchema":        "JWT+ClientOnly",
	"upsertPlatformProfileSchema":     "JWT+ClientOnly",
}

// operation runs with no auth middleware, otherwise the chain names joined
// with "+". Mirrors the pre-swap parity matrix.

func TestRouterParity(t *testing.T) {
	r := chi.NewRouter()
	// harness-owned routes; mirror their registration
	r.Get("/health", http.NotFoundHandler().ServeHTTP)
	r.Get("/docs/openapi.yml", http.NotFoundHandler().ServeHTTP)
	handler.HandlerWithOptions(handler.NewStrictHandler(stubStrict{}, nil), handler.ChiServerOptions{
		BaseRouter: r,
	})

	got := map[string]string{}
	_ = chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		got[normalizeRoute(method+" "+route)] = "public"
		return nil
	})

	// route coverage: walked == routeOperation keys
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

	// auth matrix: every operation must have a dispatch chain matching expectedOps
	chains := authChains(middlewares{
		jwtAuth:    mwJWT,
		anyAuth:    mwAnyAuth,
		clientOnly: mwClientOnly,
	})
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
