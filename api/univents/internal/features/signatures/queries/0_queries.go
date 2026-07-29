package queries

import (
	"univents/ports"
)

type Queries struct {
	editions   ports.EditionRepo
	signatures ports.SignatureRepo
	requests   ports.SignatureRequestRepo
}

func NewQueries(
	editions ports.EditionRepo,
	signatures ports.SignatureRepo,
	requests ports.SignatureRequestRepo,
) *Queries {
	return &Queries{
		editions:   editions,
		signatures: signatures,
		requests:   requests,
	}
}
