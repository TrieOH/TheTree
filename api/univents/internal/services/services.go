// Package services aggregates every feature's operations layer. Import this
// package instead of the per-feature subpackages:
//
//	ops := services.NewOperations(r, deps...)
package services

import (
	"univents/internal/repos"
	"univents/internal/services/badges"
	"univents/internal/services/certifications"
	"univents/internal/services/editions"
	"univents/internal/services/events"
	"univents/internal/services/products"
	"univents/internal/services/programs"
	"univents/internal/services/signatures"
	"univents/internal/services/ticket_types"

	"lib/email"
	"lib/objectstorage"
	idx "sdk/identityx"
)

// Type and constructor aliases for each feature's operations package.
type (
	Events      = events.Operations
	Editions    = editions.Operations
	TicketTypes = ticket_types.Operations
	Products    = products.Operations
	Programs    = programs.Operations
	Badges      = badges.Operations
	Signatures  = signatures.Operations
	Certs       = certifications.Operations
)

var (
	NewEvents      = events.NewOperations
	NewEditions    = editions.NewOperations
	NewTicketTypes = ticket_types.NewOperations
	NewProducts    = products.NewOperations
	NewPrograms    = programs.NewOperations
	NewBadges      = badges.NewOperations
	NewSignatures  = signatures.NewOperations
	NewCerts       = certifications.NewOperations
)

// Operations is the aggregate of every feature's operations, constructed
// once at startup and consumed by the HTTP handlers.
type Operations struct {
	Events      *Events
	Editions    *Editions
	TicketTypes *TicketTypes
	Products    *Products
	Programs    *Programs
	Badges      *Badges
	Signatures  *Signatures
	Certs       *Certs
}

// NewOperations wires every feature's operations from the shared repos and
// the app's external dependencies (object storage, IdentityX client, email,
// HMAC secret).
func NewOperations(r *repos.Repos, objStorage *objectstorage.Client, idxClient *idx.Client, emailClient *email.Client, hmacSecret string) *Operations {
	return &Operations{
		Events:      NewEvents(r.Events, objStorage, idxClient),
		Editions:    NewEditions(r.Events, r.Editions),
		TicketTypes: NewTicketTypes(r.Events, r.Editions, r.TicketTypes),
		Products:    NewProducts(r.Events, r.Editions, r.Products),
		Programs:    NewPrograms(r.Events, r.Editions, r.Programs, r.Occurrences),
		Badges:      NewBadges(r.Badges),
		Signatures:  NewSignatures(r.Events, r.Editions, r.Signatures, r.SignatureRequests, emailClient, hmacSecret),
		Certs:       NewCerts(r.Events, r.Editions, r.Certs, r.Programs, emailClient),
	}
}
