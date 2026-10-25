package service

import (
	"univents/internal/authz"
	"univents/ports"
)

type Deps struct {
	Events   ports.EventRepo
	Editions ports.EditionRepo
	Badges   ports.BadgeEditionOps
	Authz    *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
