// Package services aggregates every feature's operations layer. Import this
// package instead of the per-feature subpackages:
//
//	ops := services.NewOperations(r, authzSvc, hmacSecret)
//	actor, err := ops.Actors.GetByID(ctx, id)
package services

import (
	"IdentityX/internal/repos"
	"IdentityX/internal/services/actors"
	"IdentityX/internal/services/api_keys"
	"IdentityX/internal/services/authn"
	"IdentityX/internal/services/capabilities"
	"IdentityX/internal/services/organizations"
	"IdentityX/internal/services/profile_schemas"
	"IdentityX/internal/services/profiles"
	"IdentityX/internal/services/projects"

	"IdentityX/internal/authz"
)

// Type and constructor aliases for each feature's operations package.
type (
	Actors         = actors.Operations
	APIKeys        = api_keys.Operations
	Authn          = authn.Operations
	Capabilities   = capabilities.Operations
	Organizations  = organizations.Operations
	ProfileSchemas = profile_schemas.Operations
	Profiles       = profiles.Operations
	Projects       = projects.Operations
)

var (
	NewActors         = actors.NewOperations
	NewAPIKeys        = api_keys.NewOperations
	NewAuthn          = authn.NewOperations
	NewCapabilities   = capabilities.NewOperations
	NewOrganizations  = organizations.NewOperations
	NewProfileSchemas = profile_schemas.NewOperations
	NewProfiles       = profiles.NewOperations
	NewProjects       = projects.NewOperations
)

// Operations is the aggregate of every feature's operations, constructed
// once at startup and consumed by the HTTP handlers.
type Operations struct {
	Actors         *Actors
	APIKeys        *APIKeys
	Authn          *Authn
	Capabilities   *Capabilities
	Organizations  *Organizations
	ProfileSchemas *ProfileSchemas
	Profiles       *Profiles
	Projects       *Projects
}

// NewOperations wires every feature's operations from the shared repos.
// hmacSecret is the API-key signing secret (app config), passed through to
// the api_keys service.
func NewOperations(r *repos.Repos, authzSvc *authz.Service, hmacSecret string) *Operations {
	return &Operations{
		Actors:         NewActors(r.Actors, r.Projects, authzSvc),
		APIKeys:        NewAPIKeys([]byte(hmacSecret), r.Actors, r.APIKeys, r.Capabilities, r.Projects, authzSvc),
		Authn:          NewAuthn(r.Actors, r.Projects, r.PlatformRoles, r.CryptoKeys, r.Blacklist, r.ExternalIdentities),
		Capabilities:   NewCapabilities(r.Actors, r.Capabilities, r.Projects, authzSvc),
		Organizations:  NewOrganizations(r.Projects, r.Actors, r.Organizations, authzSvc),
		ProfileSchemas: NewProfileSchemas(r.ProfileSchemas, r.Projects, authzSvc),
		Profiles:       NewProfiles(r.Profiles, r.ProfileSchemas, r.Projects, authzSvc),
		Projects:       NewProjects(r.CryptoKeys, r.Projects, r.Actors, authzSvc),
	}
}
