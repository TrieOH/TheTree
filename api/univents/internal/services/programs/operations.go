package programs

import (
	"context"
	"encoding/json"

	"lib/database"
	"lib/telemetry"
	"univents/internal/authz"
	"univents/internal/services/notify"
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
// on (contract in internal/services/notify).
const channelUniventsChanges = notify.Channel

type Operations struct {
	events         ports.EventRepo
	editions       ports.EditionRepo
	programs       ports.ProgramRepo
	occurrences    ports.ProgramOccurrenceRepo
	registrations  ports.RegistrationRepo
	ticketTypes    ports.TicketTypeRepo
	participations ports.ProgramParticipationRepo
	authz          *authz.Service
	notifier       Notifier
	tx             database.TxRunner
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	programs ports.ProgramRepo,
	occurrences ports.ProgramOccurrenceRepo,
	registrations ports.RegistrationRepo,
	ticketTypes ports.TicketTypeRepo,
	participations ports.ProgramParticipationRepo,
	authz *authz.Service,
	notifier Notifier,
	tx database.TxRunner,
) *Operations {
	return &Operations{
		events:         events,
		editions:       editions,
		programs:       programs,
		occurrences:    occurrences,
		registrations:  registrations,
		ticketTypes:    ticketTypes,
		participations: participations,
		authz:          authz,
		notifier:       notifier,
		tx:             tx,
	}
}

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
	err = o.notifier.Notify(ctx, channelUniventsChanges, string(raw))
	if err != nil {
		telemetry.Log().Error("programs: publish to notifier",
			zap.String("channel", channelUniventsChanges),
			zap.String("payload", string(raw)),
			zap.Error(err))
	}
}
