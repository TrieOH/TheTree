package service

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
)

type Deps struct {
	Actors       ports.ActorRepo
	Capabilities ports.CapabilityRepo
	Projects     ports.ProjectRepo
	Authz        *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
