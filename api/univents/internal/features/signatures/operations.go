package signatures

import (
	"lib/email"
	"univents/ports"
)

type Operations struct {
	events     ports.EventRepo
	editions   ports.EditionRepo
	signatures ports.SignatureRepo
	requests   ports.SignatureRequestRepo
	email      *email.Client
	hmacSecret string
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	signatures ports.SignatureRepo,
	requests ports.SignatureRequestRepo,
	email *email.Client,
	hmacSecret string,
) *Operations {
	return &Operations{
		events:     events,
		editions:   editions,
		signatures: signatures,
		requests:   requests,
		email:      email,
		hmacSecret: hmacSecret,
	}
}
