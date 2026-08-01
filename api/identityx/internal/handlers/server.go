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

	"IdentityX/internal/handler"
	"IdentityX/internal/handlers/actors"
	"IdentityX/internal/handlers/api_keys"
	"IdentityX/internal/handlers/authn"
	"IdentityX/internal/handlers/capabilities"
	"IdentityX/internal/handlers/organizations"
	"IdentityX/internal/handlers/profile_schemas"
	"IdentityX/internal/handlers/profiles"
	"IdentityX/internal/handlers/projects"
	"IdentityX/internal/services"
)

// Server implements handler.StrictServerInterface.
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

func (s *Server) ListActors(ctx context.Context, req handler.ListActorsRequestObject) (handler.ListActorsResponseObject, error) {
	return s.actors.ListActors(ctx, req)
}

func (s *Server) CreateActor(ctx context.Context, req handler.CreateActorRequestObject) (handler.CreateActorResponseObject, error) {
	return s.actors.CreateActor(ctx, req)
}

func (s *Server) GetActor(ctx context.Context, req handler.GetActorRequestObject) (handler.GetActorResponseObject, error) {
	return s.actors.GetActor(ctx, req)
}

func (s *Server) GetActorByEmail(ctx context.Context, req handler.GetActorByEmailRequestObject) (handler.GetActorByEmailResponseObject, error) {
	return s.actors.GetActorByEmail(ctx, req)
}

func (s *Server) CreateAPIKey(ctx context.Context, req handler.CreateAPIKeyRequestObject) (handler.CreateAPIKeyResponseObject, error) {
	return s.apiKeys.CreateAPIKey(ctx, req)
}

func (s *Server) GetSetup(ctx context.Context, req handler.GetSetupRequestObject) (handler.GetSetupResponseObject, error) {
	return s.authn.GetSetup(ctx, req)
}

func (s *Server) PostSetup(ctx context.Context, req handler.PostSetupRequestObject) (handler.PostSetupResponseObject, error) {
	return s.authn.PostSetup(ctx, req)
}

func (s *Server) PostRegister(ctx context.Context, req handler.PostRegisterRequestObject) (handler.PostRegisterResponseObject, error) {
	return s.authn.PostRegister(ctx, req)
}

func (s *Server) PostLogin(ctx context.Context, req handler.PostLoginRequestObject) (handler.PostLoginResponseObject, error) {
	return s.authn.PostLogin(ctx, req)
}

func (s *Server) PostLogout(ctx context.Context, req handler.PostLogoutRequestObject) (handler.PostLogoutResponseObject, error) {
	return s.authn.PostLogout(ctx, req)
}

func (s *Server) PostRefresh(ctx context.Context, req handler.PostRefreshRequestObject) (handler.PostRefreshResponseObject, error) {
	return s.authn.PostRefresh(ctx, req)
}

func (s *Server) GetOAuthConnect(ctx context.Context, req handler.GetOAuthConnectRequestObject) (handler.GetOAuthConnectResponseObject, error) {
	return s.authn.GetOAuthConnect(ctx, req)
}

func (s *Server) GetOAuthCallback(ctx context.Context, req handler.GetOAuthCallbackRequestObject) (handler.GetOAuthCallbackResponseObject, error) {
	return s.authn.GetOAuthCallback(ctx, req)
}

func (s *Server) GetJWKS(ctx context.Context, req handler.GetJWKSRequestObject) (handler.GetJWKSResponseObject, error) {
	return s.authn.GetJWKS(ctx, req)
}

func (s *Server) GetIntrospect(ctx context.Context, req handler.GetIntrospectRequestObject) (handler.GetIntrospectResponseObject, error) {
	return s.authn.GetIntrospect(ctx, req)
}

func (s *Server) ListOrganizations(ctx context.Context, req handler.ListOrganizationsRequestObject) (handler.ListOrganizationsResponseObject, error) {
	return s.organizations.ListOrganizations(ctx, req)
}

func (s *Server) CreateOrganization(ctx context.Context, req handler.CreateOrganizationRequestObject) (handler.CreateOrganizationResponseObject, error) {
	return s.organizations.CreateOrganization(ctx, req)
}

func (s *Server) ListOrganizationMembers(ctx context.Context, req handler.ListOrganizationMembersRequestObject) (handler.ListOrganizationMembersResponseObject, error) {
	return s.organizations.ListOrganizationMembers(ctx, req)
}

func (s *Server) AddOrganizationMember(ctx context.Context, req handler.AddOrganizationMemberRequestObject) (handler.AddOrganizationMemberResponseObject, error) {
	return s.organizations.AddOrganizationMember(ctx, req)
}

func (s *Server) RemoveOrganizationMember(ctx context.Context, req handler.RemoveOrganizationMemberRequestObject) (handler.RemoveOrganizationMemberResponseObject, error) {
	return s.organizations.RemoveOrganizationMember(ctx, req)
}

func (s *Server) ListOrganizationProjects(ctx context.Context, req handler.ListOrganizationProjectsRequestObject) (handler.ListOrganizationProjectsResponseObject, error) {
	return s.organizations.ListOrganizationProjects(ctx, req)
}

func (s *Server) CreateOrganizationProject(ctx context.Context, req handler.CreateOrganizationProjectRequestObject) (handler.CreateOrganizationProjectResponseObject, error) {
	return s.organizations.CreateOrganizationProject(ctx, req)
}

func (s *Server) ListOrganizationProjectActors(ctx context.Context, req handler.ListOrganizationProjectActorsRequestObject) (handler.ListOrganizationProjectActorsResponseObject, error) {
	return s.organizations.ListOrganizationProjectActors(ctx, req)
}

func (s *Server) CreateOrganizationProjectActor(ctx context.Context, req handler.CreateOrganizationProjectActorRequestObject) (handler.CreateOrganizationProjectActorResponseObject, error) {
	return s.organizations.CreateOrganizationProjectActor(ctx, req)
}

func (s *Server) ListOrganizationProjectMembers(ctx context.Context, req handler.ListOrganizationProjectMembersRequestObject) (handler.ListOrganizationProjectMembersResponseObject, error) {
	return s.organizations.ListOrganizationProjectMembers(ctx, req)
}

func (s *Server) AddOrganizationProjectMember(ctx context.Context, req handler.AddOrganizationProjectMemberRequestObject) (handler.AddOrganizationProjectMemberResponseObject, error) {
	return s.organizations.AddOrganizationProjectMember(ctx, req)
}

func (s *Server) RemoveOrganizationProjectMember(ctx context.Context, req handler.RemoveOrganizationProjectMemberRequestObject) (handler.RemoveOrganizationProjectMemberResponseObject, error) {
	return s.organizations.RemoveOrganizationProjectMember(ctx, req)
}

func (s *Server) GetOrganizationProjectActor(ctx context.Context, req handler.GetOrganizationProjectActorRequestObject) (handler.GetOrganizationProjectActorResponseObject, error) {
	return s.organizations.GetOrganizationProjectActor(ctx, req)
}

func (s *Server) GetPlatformProfile(ctx context.Context, req handler.GetPlatformProfileRequestObject) (handler.GetPlatformProfileResponseObject, error) {
	return s.profiles.GetPlatformProfile(ctx, req)
}

func (s *Server) UpsertPlatformProfile(ctx context.Context, req handler.UpsertPlatformProfileRequestObject) (handler.UpsertPlatformProfileResponseObject, error) {
	return s.profiles.UpsertPlatformProfile(ctx, req)
}

func (s *Server) GetProjectProfile(ctx context.Context, req handler.GetProjectProfileRequestObject) (handler.GetProjectProfileResponseObject, error) {
	return s.profiles.GetProjectProfile(ctx, req)
}

func (s *Server) UpsertProjectProfile(ctx context.Context, req handler.UpsertProjectProfileRequestObject) (handler.UpsertProjectProfileResponseObject, error) {
	return s.profiles.UpsertProjectProfile(ctx, req)
}

func (s *Server) GetPlatformProfileSchema(ctx context.Context, req handler.GetPlatformProfileSchemaRequestObject) (handler.GetPlatformProfileSchemaResponseObject, error) {
	return s.profileSchemas.GetPlatformProfileSchema(ctx, req)
}

func (s *Server) UpsertPlatformProfileSchema(ctx context.Context, req handler.UpsertPlatformProfileSchemaRequestObject) (handler.UpsertPlatformProfileSchemaResponseObject, error) {
	return s.profileSchemas.UpsertPlatformProfileSchema(ctx, req)
}

func (s *Server) GetProjectProfileSchema(ctx context.Context, req handler.GetProjectProfileSchemaRequestObject) (handler.GetProjectProfileSchemaResponseObject, error) {
	return s.profileSchemas.GetProjectProfileSchema(ctx, req)
}

func (s *Server) UpsertProjectProfileSchema(ctx context.Context, req handler.UpsertProjectProfileSchemaRequestObject) (handler.UpsertProjectProfileSchemaResponseObject, error) {
	return s.profileSchemas.UpsertProjectProfileSchema(ctx, req)
}

func (s *Server) ListProjects(ctx context.Context, req handler.ListProjectsRequestObject) (handler.ListProjectsResponseObject, error) {
	return s.projects.ListProjects(ctx, req)
}

func (s *Server) CreateProject(ctx context.Context, req handler.CreateProjectRequestObject) (handler.CreateProjectResponseObject, error) {
	return s.projects.CreateProject(ctx, req)
}

func (s *Server) ListProjectMembers(ctx context.Context, req handler.ListProjectMembersRequestObject) (handler.ListProjectMembersResponseObject, error) {
	return s.projects.ListProjectMembers(ctx, req)
}

func (s *Server) AddProjectMember(ctx context.Context, req handler.AddProjectMemberRequestObject) (handler.AddProjectMemberResponseObject, error) {
	return s.projects.AddProjectMember(ctx, req)
}

func (s *Server) RemoveProjectMember(ctx context.Context, req handler.RemoveProjectMemberRequestObject) (handler.RemoveProjectMemberResponseObject, error) {
	return s.projects.RemoveProjectMember(ctx, req)
}

func (s *Server) ListCapabilities(ctx context.Context, req handler.ListCapabilitiesRequestObject) (handler.ListCapabilitiesResponseObject, error) {
	return s.capabilities.ListCapabilities(ctx, req)
}

func (s *Server) CreateCapability(ctx context.Context, req handler.CreateCapabilityRequestObject) (handler.CreateCapabilityResponseObject, error) {
	return s.capabilities.CreateCapability(ctx, req)
}
