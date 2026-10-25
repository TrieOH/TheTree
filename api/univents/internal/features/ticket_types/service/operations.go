package service

import (
	"univents/internal/authz"
	"univents/ports"
)

type Deps struct {
	Events      ports.EventRepo
	Editions    ports.EditionRepo
	TicketTypes ports.TicketTypeRepo
	Authz       *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
