// Package repos aggregates every feature's repository layer. Import this
// package instead of the per-feature subpackages:
//
//	r := repos.New(app.db)
package repos

import (
	"univents/internal/sqlc"

	"univents/internal/repos/badges"
	"univents/internal/repos/certifications"
	"univents/internal/repos/editions"
	"univents/internal/repos/events"
	"univents/internal/repos/products"
	"univents/internal/repos/programs"
	"univents/internal/repos/signatures"
	"univents/internal/repos/ticket_types"
)

// Type and constructor aliases for each feature's repo package.
type (
	Events      = events.Repo
	Editions    = editions.Repo
	TicketTypes = ticket_types.Repo
	Products    = products.Repo
	Programs    = programs.Repo
	Badges      = badges.Repo
	Signatures  = signatures.Repo
	Certs       = certifications.Repo
)

var (
	NewEvents      = events.NewRepo
	NewEditions    = editions.NewRepo
	NewTicketTypes = ticket_types.NewRepo
	NewProducts    = products.NewRepo
	NewPrograms    = programs.NewRepo
	NewBadges      = badges.NewRepo
	NewSignatures  = signatures.NewRepo
	NewCerts       = certifications.NewRepo
)

// Repos is the aggregate of every feature repo, constructed once at startup.
// Occurrences and SignatureRequests are served by the programs and
// signatures repos respectively.
type Repos struct {
	Events            *Events
	Editions          *Editions
	TicketTypes       *TicketTypes
	Products          *Products
	Programs          *Programs
	Occurrences       *Programs
	Badges            *Badges
	Signatures        *Signatures
	SignatureRequests *Signatures
	Certs             *Certs
}

// New constructs every feature repo from the shared query handle.
func New(q *sqlc.Queries) *Repos {
	return &Repos{
		Events:            NewEvents(q),
		Editions:          NewEditions(q),
		TicketTypes:       NewTicketTypes(q),
		Products:          NewProducts(q),
		Programs:          NewPrograms(q),
		Occurrences:       NewPrograms(q),
		Badges:            NewBadges(q),
		Signatures:        NewSignatures(q),
		SignatureRequests: NewSignatures(q),
		Certs:             NewCerts(q),
	}
}
