// Package services aggregates every feature's operations layer. Import this
// package instead of the per-feature subpackages:
//
//	ops := services.NewOperations(r, deps...)
package services

import (
	"univents/internal/authz"
	"univents/internal/repos"
	"univents/internal/services/badges"
	"univents/internal/services/certifications"
	"univents/internal/services/checkouts"
	"univents/internal/services/editions"
	"univents/internal/services/events"
	"univents/internal/services/payments"
	"univents/internal/services/products"
	"univents/internal/services/programs"
	"univents/internal/services/purchases"
	"univents/internal/services/signatures"
	"univents/internal/services/store"
	"univents/internal/services/ticket_types"
	"univents/internal/services/webhooks"
	"univents/internal/services/ws"

	"lib/database"
	"lib/email"
	"lib/objectstorage"
	idx "sdk/identityx"
	payssage "sdk/payssage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
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
	Payments    = payments.Operations
	Webhooks    = webhooks.Operations
	Checkouts   = checkouts.Operations
	Purchases   = purchases.Operations
	WS          = ws.Operations
	Store       = store.Operations
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
	NewPayments    = payments.NewOperations
	NewWebhooks    = webhooks.NewOperations
	NewCheckouts   = checkouts.NewOperations
	NewPurchases   = purchases.NewOperations
	NewWS          = ws.NewOperations
	NewStore       = store.NewOperations
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
	Payments    *Payments
	Webhooks    *Webhooks
	Checkouts   *Checkouts
	Purchases   *Purchases
	WS          *WS
	Store       *Store
}

// NewOperations wires every feature's operations from the shared repos and
// the app's external dependencies (object storage, IdentityX client, email,
// HMAC secret). Authorization arrives by injection through the same seam —
// no service-locator globals.
func NewOperations(
	r *repos.Repos,
	authzSvc *authz.Service,
	objStorage *objectstorage.Client,
	idxClient *idx.Client,
	emailClient *email.Client,
	hmacSecret string,
	payssageClient *payssage.Client,
	platformWalletID uuid.UUID,
	notifier *database.Notifier,
	riverClient *river.Client[pgx.Tx],
	tx database.TxRunner,
	webhookSecret string,
) *Operations {
	badgesOps := NewBadges(r.Badges, r.Badges, r.Registrations, r.Editions, r.Events, checkouts.NewSDKActorResolver(idxClient.Actors), emailClient, authzSvc)
	wsOps := NewWS(r.WsTokens, r.Purchases, payssageClient, notifier)
	return &Operations{
		Events:      NewEvents(r.Events, objStorage, idxClient, authzSvc, badgesOps),
		Editions:    NewEditions(r.Events, r.Editions, authzSvc, badgesOps),
		TicketTypes: NewTicketTypes(r.Events, r.Editions, r.TicketTypes, authzSvc),
		Products:    NewProducts(r.Events, r.Editions, r.Products, authzSvc),
		Programs:    NewPrograms(r.Events, r.Editions, r.Programs, r.Occurrences, r.Registrations, r.TicketTypes, r.Programs, authzSvc, notifier, tx),
		Badges:      badgesOps,
		Signatures:  NewSignatures(r.Events, r.Editions, r.Signatures, r.SignatureRequests, emailClient, hmacSecret, authzSvc),
		Certs:       NewCerts(r.Events, r.Editions, r.Certs, r.Programs, emailClient, riverClient, authzSvc),
		Payments:    NewPayments(r.Events, payssageClient, authzSvc, platformWalletID),
		Webhooks:    NewWebhooks(r.Purchases, r.Registrations, r.Products, r.Programs, badgesOps, notifier, riverClient, tx, webhookSecret),
		Checkouts:   NewCheckouts(r.Purchases, r.Editions, r.Events, r.TicketTypes, r.Products, r.Programs, r.Occurrences, r.Registrations, r.Products, r.Programs, badgesOps, notifier, riverClient, tx, payssageClient, payssageClient, platformWalletID, wsOps, checkouts.NewSDKActorResolver(idxClient.Actors), authzSvc),
		Purchases:   NewPurchases(r.Purchases),
		WS:          wsOps,
		Store:       NewStore(r.Purchases, r.Editions, notifier),
	}
}
