// Package handlers implements the generated StrictServerInterface by
// delegating to one subpackage per feature. Each subpackage owns the
// methods of its feature; this package only wires them together.
//
// Auth, validation, and error mapping run in the strict middleware stack
// (see internal/app/auth_dispatch.go); the handlers here are pure
// domain logic + fun envelope construction.
package handlers

import (
	"context"

	"IdentityX/internal/handlers/actors"
	"IdentityX/internal/handlers/api_keys"
	"IdentityX/internal/handlers/authn"
	"IdentityX/internal/handlers/capabilities"
	"IdentityX/internal/handlers/organizations"
	"IdentityX/internal/handlers/profile_schemas"
	"IdentityX/internal/handlers/profiles"
	"IdentityX/internal/handlers/projects"
	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
)

// Server implements openapi.StrictServerInterface.
type Server struct {
	actors         *actors.Handlers
	apiKeys        *api_keys.Handlers
	authn          *authn.Handlers
	capabilities   *capabilities.Handlers
	organizations  *organizations.Handlers
	profileSchemas *profile_schemas.Handlers
	profiles       *profiles.Handlers
	projects       *projects.Handlers
}

// NewServer wires the per-feature handlers from the services aggregate.
func NewServer(ops *services.Operations) *Server {
	return &Server{
		actors:         actors.New(ops.Actors),
		apiKeys:        api_keys.New(ops.APIKeys),
		authn:          authn.New(ops.Authn),
		capabilities:   capabilities.New(ops.Capabilities),
		organizations:  organizations.New(ops.Organizations),
		profileSchemas: profile_schemas.New(ops.ProfileSchemas),
		profiles:       profiles.New(ops.Profiles),
		projects:       projects.New(ops.Projects),
	}
}

// ── StrictServerInterface ────────────────────────────────────────────────

func (s *Server) ListActors(ctx context.Context, req openapi.ListActorsRequestObject) (openapi.ListActorsResponseObject, error) {
	return s.actors.ListActors(ctx, req)
}

func (s *Server) CreateActor(ctx context.Context, req openapi.CreateActorRequestObject) (openapi.CreateActorResponseObject, error) {
	return s.actors.CreateActor(ctx, req)
}

func (s *Server) GetActor(ctx context.Context, req openapi.GetActorRequestObject) (openapi.GetActorResponseObject, error) {
	return s.actors.GetActor(ctx, req)
}

func (s *Server) GetActorByEmail(ctx context.Context, req openapi.GetActorByEmailRequestObject) (openapi.GetActorByEmailResponseObject, error) {
	return s.actors.GetActorByEmail(ctx, req)
}

func (s *Server) CreateAPIKey(ctx context.Context, req openapi.CreateAPIKeyRequestObject) (openapi.CreateAPIKeyResponseObject, error) {
	return s.apiKeys.CreateAPIKey(ctx, req)
}

func (s *Server) GetSetup(ctx context.Context, req openapi.GetSetupRequestObject) (openapi.GetSetupResponseObject, error) {
	return s.authn.GetSetup(ctx, req)
}

func (s *Server) PostSetup(ctx context.Context, req openapi.PostSetupRequestObject) (openapi.PostSetupResponseObject, error) {
	return s.authn.PostSetup(ctx, req)
}

func (s *Server) PostRegister(ctx context.Context, req openapi.PostRegisterRequestObject) (openapi.PostRegisterResponseObject, error) {
	return s.authn.PostRegister(ctx, req)
}

func (s *Server) PostLogin(ctx context.Context, req openapi.PostLoginRequestObject) (openapi.PostLoginResponseObject, error) {
	return s.authn.PostLogin(ctx, req)
}

func (s *Server) PostLogout(ctx context.Context, req openapi.PostLogoutRequestObject) (openapi.PostLogoutResponseObject, error) {
	return s.authn.PostLogout(ctx, req)
}

func (s *Server) PostRefresh(ctx context.Context, req openapi.PostRefreshRequestObject) (openapi.PostRefreshResponseObject, error) {
	return s.authn.PostRefresh(ctx, req)
}

func (s *Server) GetOAuthConnect(ctx context.Context, req openapi.GetOAuthConnectRequestObject) (openapi.GetOAuthConnectResponseObject, error) {
	return s.authn.GetOAuthConnect(ctx, req)
}

func (s *Server) GetOAuthCallback(ctx context.Context, req openapi.GetOAuthCallbackRequestObject) (openapi.GetOAuthCallbackResponseObject, error) {
	return s.authn.GetOAuthCallback(ctx, req)
}

func (s *Server) GetJWKS(ctx context.Context, req openapi.GetJWKSRequestObject) (openapi.GetJWKSResponseObject, error) {
	return s.authn.GetJWKS(ctx, req)
}

func (s *Server) GetIntrospect(ctx context.Context, req openapi.GetIntrospectRequestObject) (openapi.GetIntrospectResponseObject, error) {
	return s.authn.GetIntrospect(ctx, req)
}

func (s *Server) ListOrganizations(ctx context.Context, req openapi.ListOrganizationsRequestObject) (openapi.ListOrganizationsResponseObject, error) {
	return s.organizations.ListOrganizations(ctx, req)
}

func (s *Server) CreateOrganization(ctx context.Context, req openapi.CreateOrganizationRequestObject) (openapi.CreateOrganizationResponseObject, error) {
	return s.organizations.CreateOrganization(ctx, req)
}

func (s *Server) ListOrganizationMembers(ctx context.Context, req openapi.ListOrganizationMembersRequestObject) (openapi.ListOrganizationMembersResponseObject, error) {
	return s.organizations.ListOrganizationMembers(ctx, req)
}

func (s *Server) AddOrganizationMember(ctx context.Context, req openapi.AddOrganizationMemberRequestObject) (openapi.AddOrganizationMemberResponseObject, error) {
	return s.organizations.AddOrganizationMember(ctx, req)
}

func (s *Server) RemoveOrganizationMember(ctx context.Context, req openapi.RemoveOrganizationMemberRequestObject) (openapi.RemoveOrganizationMemberResponseObject, error) {
	return s.organizations.RemoveOrganizationMember(ctx, req)
}

func (s *Server) ListOrganizationProjects(ctx context.Context, req openapi.ListOrganizationProjectsRequestObject) (openapi.ListOrganizationProjectsResponseObject, error) {
	return s.organizations.ListOrganizationProjects(ctx, req)
}

func (s *Server) CreateOrganizationProject(ctx context.Context, req openapi.CreateOrganizationProjectRequestObject) (openapi.CreateOrganizationProjectResponseObject, error) {
	return s.organizations.CreateOrganizationProject(ctx, req)
}

func (s *Server) ListOrganizationProjectActors(ctx context.Context, req openapi.ListOrganizationProjectActorsRequestObject) (openapi.ListOrganizationProjectActorsResponseObject, error) {
	return s.organizations.ListOrganizationProjectActors(ctx, req)
}

func (s *Server) CreateOrganizationProjectActor(ctx context.Context, req openapi.CreateOrganizationProjectActorRequestObject) (openapi.CreateOrganizationProjectActorResponseObject, error) {
	return s.organizations.CreateOrganizationProjectActor(ctx, req)
}

func (s *Server) ListOrganizationProjectMembers(ctx context.Context, req openapi.ListOrganizationProjectMembersRequestObject) (openapi.ListOrganizationProjectMembersResponseObject, error) {
	return s.organizations.ListOrganizationProjectMembers(ctx, req)
}

func (s *Server) AddOrganizationProjectMember(ctx context.Context, req openapi.AddOrganizationProjectMemberRequestObject) (openapi.AddOrganizationProjectMemberResponseObject, error) {
	return s.organizations.AddOrganizationProjectMember(ctx, req)
}

func (s *Server) RemoveOrganizationProjectMember(ctx context.Context, req openapi.RemoveOrganizationProjectMemberRequestObject) (openapi.RemoveOrganizationProjectMemberResponseObject, error) {
	return s.organizations.RemoveOrganizationProjectMember(ctx, req)
}

func (s *Server) GetOrganizationProjectActor(ctx context.Context, req openapi.GetOrganizationProjectActorRequestObject) (openapi.GetOrganizationProjectActorResponseObject, error) {
	return s.organizations.GetOrganizationProjectActor(ctx, req)
}

func (s *Server) GetPlatformProfile(ctx context.Context, req openapi.GetPlatformProfileRequestObject) (openapi.GetPlatformProfileResponseObject, error) {
	return s.profiles.GetPlatformProfile(ctx, req)
}

func (s *Server) UpsertPlatformProfile(ctx context.Context, req openapi.UpsertPlatformProfileRequestObject) (openapi.UpsertPlatformProfileResponseObject, error) {
	return s.profiles.UpsertPlatformProfile(ctx, req)
}

func (s *Server) GetProjectProfile(ctx context.Context, req openapi.GetProjectProfileRequestObject) (openapi.GetProjectProfileResponseObject, error) {
	return s.profiles.GetProjectProfile(ctx, req)
}

func (s *Server) UpsertProjectProfile(ctx context.Context, req openapi.UpsertProjectProfileRequestObject) (openapi.UpsertProjectProfileResponseObject, error) {
	return s.profiles.UpsertProjectProfile(ctx, req)
}

func (s *Server) ListOutdatedPlatformProfiles(ctx context.Context, req openapi.ListOutdatedPlatformProfilesRequestObject) (openapi.ListOutdatedPlatformProfilesResponseObject, error) {
	return s.profiles.ListOutdatedPlatformProfiles(ctx, req)
}

func (s *Server) ListOutdatedProjectProfiles(ctx context.Context, req openapi.ListOutdatedProjectProfilesRequestObject) (openapi.ListOutdatedProjectProfilesResponseObject, error) {
	return s.profiles.ListOutdatedProjectProfiles(ctx, req)
}

func (s *Server) GetPlatformProfileSchema(ctx context.Context, req openapi.GetPlatformProfileSchemaRequestObject) (openapi.GetPlatformProfileSchemaResponseObject, error) {
	return s.profileSchemas.GetPlatformProfileSchema(ctx, req)
}

func (s *Server) UpsertPlatformProfileSchema(ctx context.Context, req openapi.UpsertPlatformProfileSchemaRequestObject) (openapi.UpsertPlatformProfileSchemaResponseObject, error) {
	return s.profileSchemas.UpsertPlatformProfileSchema(ctx, req)
}

func (s *Server) GetProjectProfileSchema(ctx context.Context, req openapi.GetProjectProfileSchemaRequestObject) (openapi.GetProjectProfileSchemaResponseObject, error) {
	return s.profileSchemas.GetProjectProfileSchema(ctx, req)
}

func (s *Server) UpsertProjectProfileSchema(ctx context.Context, req openapi.UpsertProjectProfileSchemaRequestObject) (openapi.UpsertProjectProfileSchemaResponseObject, error) {
	return s.profileSchemas.UpsertProjectProfileSchema(ctx, req)
}

func (s *Server) ListProjects(ctx context.Context, req openapi.ListProjectsRequestObject) (openapi.ListProjectsResponseObject, error) {
	return s.projects.ListProjects(ctx, req)
}

func (s *Server) CreateProject(ctx context.Context, req openapi.CreateProjectRequestObject) (openapi.CreateProjectResponseObject, error) {
	return s.projects.CreateProject(ctx, req)
}

func (s *Server) ListProjectMembers(ctx context.Context, req openapi.ListProjectMembersRequestObject) (openapi.ListProjectMembersResponseObject, error) {
	return s.projects.ListProjectMembers(ctx, req)
}

func (s *Server) AddProjectMember(ctx context.Context, req openapi.AddProjectMemberRequestObject) (openapi.AddProjectMemberResponseObject, error) {
	return s.projects.AddProjectMember(ctx, req)
}

func (s *Server) RemoveProjectMember(ctx context.Context, req openapi.RemoveProjectMemberRequestObject) (openapi.RemoveProjectMemberResponseObject, error) {
	return s.projects.RemoveProjectMember(ctx, req)
}

func (s *Server) ListCapabilities(ctx context.Context, req openapi.ListCapabilitiesRequestObject) (openapi.ListCapabilitiesResponseObject, error) {
	return s.capabilities.ListCapabilities(ctx, req)
}

func (s *Server) CreateCapability(ctx context.Context, req openapi.CreateCapabilityRequestObject) (openapi.CreateCapabilityResponseObject, error) {
	return s.capabilities.CreateCapability(ctx, req)
}
