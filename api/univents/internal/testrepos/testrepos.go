// Package testrepos is for tests only: every feature's repo over one query
// handle, for the DB tests that seed and assert through several features.
// Build wires repos one by one (ADR-0015); nothing outside tests imports
// this.
package testrepos

import (
	"univents/internal/features/badges"
	"univents/internal/features/certifications"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/products"
	"univents/internal/features/programs"
	"univents/internal/features/purchases"
	"univents/internal/features/registrations"
	"univents/internal/features/signatures"
	"univents/internal/features/ticket_types"
	"univents/internal/features/ws_tokens"
	"univents/internal/sqlc"
)

type Repos struct {
	Events            *events.Repo
	Editions          *editions.Repo
	TicketTypes       *ticket_types.Repo
	Products          *products.Repo
	Programs          *programs.Repo
	Occurrences       *programs.Repo
	Badges            *badges.Repo
	Signatures        *signatures.Repo
	SignatureRequests *signatures.Repo
	Certs             *certifications.Repo
	Registrations     *registrations.Repo
	Purchases         *purchases.Repo
	WsTokens          *ws_tokens.Repo
}

func New(q *sqlc.Queries) *Repos {
	programsRepo := programs.NewRepo(q)
	signaturesRepo := signatures.NewRepo(q)
	return &Repos{
		Events:            events.NewRepo(q),
		Editions:          editions.NewRepo(q),
		TicketTypes:       ticket_types.NewRepo(q),
		Products:          products.NewRepo(q),
		Programs:          programsRepo,
		Occurrences:       programsRepo,
		Badges:            badges.NewRepo(q),
		Signatures:        signaturesRepo,
		SignatureRequests: signaturesRepo,
		Certs:             certifications.NewRepo(q),
		Registrations:     registrations.NewRepo(q),
		Purchases:         purchases.NewRepo(q),
		WsTokens:          ws_tokens.NewRepo(q),
	}
}
