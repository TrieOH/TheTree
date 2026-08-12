package store

import (
	"context"
	"encoding/json"
	"slices"
	"sync"

	"lib/telemetry"

	"univents/internal/services/notify"
	"univents/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RelayManager owns the stock-delta fan-out: ONE notifier subscription for
// kind="stock" (registered at construction, never torn down) routing to one
// Relay per edition. Relays are created lazily on the first SSE conn and
// dropped when the last conn leaves — the subscription itself lives for the
// process, so there is nothing to unsubscribe (the notifier has no per-
// handler teardown).
type RelayManager struct {
	avail    AvailabilityReader
	notifier Notifier

	mu     sync.Mutex
	relays map[uuid.UUID]*Relay
}

func newRelayManager(avail AvailabilityReader, notifier Notifier) *RelayManager {
	m := &RelayManager{
		avail:    avail,
		notifier: notifier,
		relays:   make(map[uuid.UUID]*Relay),
	}
	notifier.Subscribe(notify.Channel, func(payload string) {
		// Off the listen loop: the relay re-queries the DB, and a missed or
		// reordered delta is harmless — every broadcast re-reads the current
		// availability, so the numbers are always the authoritative truth.
		go m.handleStock(payload)
	})
	return m
}

// relayFor returns the edition's shared relay, creating it on first use.
func (m *RelayManager) relayFor(editionID uuid.UUID) *Relay {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r, ok := m.relays[editionID]; ok {
		return r
	}
	r := &Relay{editionID: editionID, avail: m.avail, manager: m}
	m.relays[editionID] = r
	return r
}

func (m *RelayManager) handleStock(payload string) {
	kind, stock, _ := notify.Parse(payload)
	if kind != notify.KindStock || stock == nil || len(stock.ItemIDs) == 0 {
		return
	}

	m.mu.Lock()
	relay, ok := m.relays[stock.EditionID]
	m.mu.Unlock()
	if !ok {
		return // no conns for this edition — nothing to update
	}
	relay.broadcast(stock)
}

func (m *RelayManager) drop(relay *Relay) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r, ok := m.relays[relay.editionID]; ok && r == relay {
		delete(m.relays, relay.editionID)
	}
}

// Relay is one edition's shared delta fan-out: its conns all get the same
// broadcast, and the availability re-query happens once per NOTIFY
// regardless of conn count (D10).
type Relay struct {
	editionID uuid.UUID
	avail     AvailabilityReader
	manager   *RelayManager

	mu    sync.Mutex
	conns map[*sseStream]struct{}
}

// subscribe registers a stream; the returned handle unregisters on close
// and drops the relay from the manager when the last conn leaves.
func (r *Relay) subscribe(s *sseStream) *subscription {
	r.mu.Lock()
	if r.conns == nil {
		r.conns = make(map[*sseStream]struct{})
	}
	r.conns[s] = struct{}{}
	r.mu.Unlock()
	return &subscription{relay: r, stream: s}
}

type subscription struct {
	relay  *Relay
	stream *sseStream
}

func (s *subscription) close() {
	r := s.relay
	r.mu.Lock()
	delete(r.conns, s.stream)
	empty := len(r.conns) == 0
	r.mu.Unlock()
	if empty {
		r.manager.drop(r)
	}
}

// broadcast re-queries availability for the NOTIFY'd item_ids and writes one
// `event: stock` per item to every conn. The DB read is the authority: the
// publisher's item_ids are a hint, the numbers come from here (D10).
func (r *Relay) broadcast(stock *notify.StockNotification) {
	availability, err := r.avail.Availability(context.Background(), r.editionID)
	if err != nil {
		telemetry.Log().Error("store: relay availability read failed",
			zap.String("edition_id", r.editionID.String()),
			zap.Error(err))
		return
	}

	// Filter the full availability read down to the affected items. An id
	// that misses (wrong edition, deleted item) is skipped — its stock
	// simply does not update.
	affected := make(map[uuid.UUID]models.ItemAvailability, len(stock.ItemIDs))
	for _, a := range availability {
		if slices.Contains(stock.ItemIDs, a.ItemID) {
			affected[a.ItemID] = a
		}
	}
	if len(affected) == 0 {
		return
	}

	r.mu.Lock()
	conns := make([]*sseStream, 0, len(r.conns))
	for c := range r.conns {
		conns = append(conns, c)
	}
	r.mu.Unlock()

	for _, a := range affected {
		raw, err := json.Marshal(stockItem{ID: a.ItemID, ItemType: string(a.ItemType), Stock: available(a)})
		if err != nil {
			continue
		}
		for _, c := range conns {
			// Best-effort: a dead conn drops the write and is removed by
			// waitClosed — never block the relay on a stuck stream.
			_ = c.event("stock", raw)
		}
	}
}
