package signatures

import (
	"lib/email"
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	events     ports.EventRepo
	editions   ports.EditionRepo
	signatures ports.SignatureRepo
	requests   ports.SignatureRequestRepo
	email      *email.Client
	hmacSecret string
	authz      *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	signatures ports.SignatureRepo,
	requests ports.SignatureRequestRepo,
	email *email.Client,
	hmacSecret string,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events:     events,
		editions:   editions,
		signatures: signatures,
		requests:   requests,
		email:      email,
		hmacSecret: hmacSecret,
		authz:      authz,
	}
}
