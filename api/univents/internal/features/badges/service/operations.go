package service

import (
	"lib/email"
	"univents/internal/authz"
	"univents/ports"
)

type Deps struct {
	Repo          ports.BadgeTemplateRepo
	Emissions     ports.BadgeEmissionRepo
	Registrations ports.RegistrationRepo
	Editions      ports.EditionRepo
	Events        ports.EventRepo
	Actors        ports.ActorResolver
	Email         *email.Client
	Authz         *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
