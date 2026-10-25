package service

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
)

type Deps struct {
	Profiles ports.ProfileRepo
	Schemas  ports.ProfileSchemaRepo
	Actors   ports.ActorRepo
	Authz    *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
