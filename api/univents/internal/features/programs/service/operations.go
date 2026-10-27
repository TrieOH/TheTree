package service

import (
	"context"
	"encoding/json"

	"lib/database"
	"lib/telemetry"
	"univents/internal/authz"
	"univents/internal/notify"
	"univents/ports"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Notifier is the LISTEN/NOTIFY publish surface for the store's stock
// deltas (register/de-register move occurrence availability). Satisfied by
// *database.Notifier (lib/go/database).
type Notifier interface {
	Notify(ctx context.Context, channel, payload string) error
}

// channelUniventsChanges is the single NOTIFY channel the store publishes
// on (contract in internal/notify).
const channelUniventsChanges = notify.Channel

type Deps struct {
	Events         ports.EventRepo
	Editions       ports.EditionRepo
	Programs       ports.ProgramRepo
	Occurrences    ports.ProgramOccurrenceRepo
	Registrations  ports.RegistrationRepo
	TicketTypes    ports.TicketTypeRepo
	Participations ports.ProgramParticipationRepo
	Authz          *authz.Service
	Notifier       Notifier
	Tx             database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// notifyStock publishes the occurrence's stock delta (D10 contract): item
// ids only — the SSE relay re-queries availability from the DB. A missed
// notification is a stale snapshot, never data loss.
func (o *Operations) notifyStock(ctx context.Context, editionID, occurrenceID uuid.UUID) {
	raw, err := json.Marshal(notify.StockNotification{
		Kind:      notify.KindStock,
		EditionID: editionID,
		ItemIDs:   []uuid.UUID{occurrenceID},
	})
	if err != nil {
		telemetry.Log().Error("programs: marshal notifier payload",
			zap.String("channel", channelUniventsChanges),
			zap.Error(err))
		return
	}
	err = o.Notifier.Notify(ctx, channelUniventsChanges, string(raw))
	if err != nil {
		telemetry.Log().Error("programs: publish to notifier",
			zap.String("channel", channelUniventsChanges),
			zap.String("payload", string(raw)),
			zap.Error(err))
	}
}
