package queries

import (
	"univents/ports"
)

type Queries struct {
	certs ports.CertificationRepo
}

func NewQueries(certs ports.CertificationRepo) *Queries {
	return &Queries{certs: certs}
}
