// Package cachetags names what the edge-cached reads hold (ADR-0020): one
// tag per entity, one per collection. A read tags its response with what
// it returns; a write purges the entity it changed and the collections it
// sits in. The cached payloads embed nothing from other entities (no
// stock, no counts), so a write never reaches past its own collections.
package cachetags

import "github.com/google/uuid"

// Events is the public event list.
const Events = "events"

// Event is one event, however it is addressed (by id or by slug).
func Event(id uuid.UUID) string { return "event:" + id.String() }

// EventEditions is an event's public edition list.
func EventEditions(eventID uuid.UUID) string { return "event:" + eventID.String() + ":editions" }

// EditionProducts is an edition's product list.
func EditionProducts(editionID uuid.UUID) string { return edition(editionID, "products") }

// EditionTicketTypes is an edition's ticket type list.
func EditionTicketTypes(editionID uuid.UUID) string { return edition(editionID, "ticket-types") }

// EditionPrograms is an edition's program list.
func EditionPrograms(editionID uuid.UUID) string { return edition(editionID, "programs") }

// EditionOccurrences is an edition's occurrence list.
func EditionOccurrences(editionID uuid.UUID) string { return edition(editionID, "occurrences") }

func edition(id uuid.UUID, collection string) string {
	return "edition:" + id.String() + ":" + collection
}
