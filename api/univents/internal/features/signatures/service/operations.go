package service

import (
	"lib/email"
	"univents/internal/authz"
	"univents/ports"
)

type Deps struct {
	Events     ports.EventRepo
	Editions   ports.EditionRepo
	Signatures ports.SignatureRepo
	Requests   ports.SignatureRequestRepo
	Email      *email.Client
	HmacSecret string
	Authz      *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
