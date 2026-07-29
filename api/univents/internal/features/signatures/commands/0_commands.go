package commands

import (
	"lib/email"
	"univents/ports"
)

type Commands struct {
	events     ports.EventRepo
	editions   ports.EditionRepo
	signatures ports.SignatureRepo
	requests   ports.SignatureRequestRepo
	email      *email.Client
	hmacSecret string
}

func NewCommands(
	events ports.EventRepo,
	editions ports.EditionRepo,
	signatures ports.SignatureRepo,
	requests ports.SignatureRequestRepo,
	email *email.Client,
	hmacSecret string,
) *Commands {
	return &Commands{
		events:     events,
		editions:   editions,
		signatures: signatures,
		requests:   requests,
		email:      email,
		hmacSecret: hmacSecret,
	}
}
