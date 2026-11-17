// Package jobs holds the one declaration of Univents' background work: which
// job kinds exist. The handlers live with their feature
// (features/checkouts/jobs, features/certifications/jobs); the carrier is
// chosen by the app wiring, and nothing here knows it. Univents has no
// periodic jobs: the expiry and gift-email jobs are enqueued with a delay.
package jobs

import (
	"lib/database"
	"lib/email"
	libjobs "lib/jobs"
	badges "univents/internal/features/badges/jobs"
	certs "univents/internal/features/certifications/jobs"
	checkouts "univents/internal/features/checkouts/jobs"
	"univents/ports"
)

// Deps are what the handlers need.
type Deps struct {
	Certs            ports.CertificationRepo
	BadgeEmissions   ports.BadgeEmissionRepo
	Editions         ports.EditionRepo
	Events           ports.EventRepo
	Purchases        ports.PurchaseRepo
	Registrations    ports.RegistrationRepo
	ProductPurchases ports.ProductPurchaseRepo
	Participations   ports.ProgramParticipationRepo
	TicketTypes      ports.TicketTypeRepo
	Email            *email.Client
	Notifier         checkouts.Notifier
	Tx               database.TxRunner
	AppURL           string // the Univents front, for the links in emails
	HmacSecret       string // keys the certificates' public verification hashes
}

// Registry registers every Univents job kind.
func Registry(d Deps) *libjobs.Registry {
	reg := libjobs.NewRegistry()
	libjobs.Register(reg, certs.NewGrantCertsWorker(d.Certs, d.Editions, d.Events, d.Email, d.AppURL, d.HmacSecret).Work)
	libjobs.Register(reg, certs.NewGrantCertsForOccurrenceWorker(d.Certs, d.Editions, d.Events, d.Email, d.AppURL, d.HmacSecret).Work)
	libjobs.Register(reg, checkouts.NewExpirePurchaseWorker(d.Purchases, d.Registrations, d.ProductPurchases, d.Participations, d.Notifier, d.Tx).Work)
	libjobs.Register(reg, badges.NewSendBadgeEmailWorker(d.Registrations, d.BadgeEmissions, d.Editions, d.Events, d.Email, d.AppURL).Work)
	libjobs.Register(reg, checkouts.NewSendGiftEmailWorker(d.Registrations, d.Editions, d.Events, d.TicketTypes, d.Email, d.AppURL).Work)
	// TODO: schedule GrantCertsForEdition on edition end and GrantCertsForOccurrence on occurrence end
	return reg
}
