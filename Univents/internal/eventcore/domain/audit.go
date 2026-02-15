package domain

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type AuditBuilder struct {
	eventID   uuid.UUID
	actorType ActorType
	actorID   *uuid.UUID
	audit     Audit // current building audit
	entries   []Audit
}

// Start initializes builder with base fields
func Start(eventID uuid.UUID, actorType ActorType, actorID *uuid.UUID) *AuditBuilder {
	return &AuditBuilder{
		eventID:   eventID,
		actorType: actorType,
		actorID:   actorID,
		entries:   make([]Audit, 0),
	}
}

// Action sets the action for current audit
func (b *AuditBuilder) Action(a AuditAction) *AuditBuilder {
	b.audit = Audit{
		EventID:   b.eventID,
		ActorType: b.actorType,
		ActorID:   b.actorID,
		Action:    a,
	}
	return b
}

// StatusChange sets from/to status for current audit
func (b *AuditBuilder) StatusChange(from, to Status) *AuditBuilder {
	b.audit = Audit{
		EventID:    b.eventID,
		ActorType:  b.actorType,
		ActorID:    b.actorID,
		FromStatus: &from,
		ToStatus:   &to,
	}
	// Auto-infer action from status transition
	b.audit.Action = AuditAction(inferStatusAction(string(from), string(to)))
	return b
}

// Metadata sets metadata for current audit
func (b *AuditBuilder) Metadata(v any) *AuditBuilder {
	blob, _ := json.Marshal(v)
	raw := json.RawMessage(blob)
	b.audit.Metadata = &raw
	return b
}

// Commit finalizes current and flushes all to storage
func (b *AuditBuilder) Commit(ctx context.Context, appendAudit func(ctx context.Context, audit Audit) (*Audit, error)) error {
	for _, a := range b.entries {
		if _, err := appendAudit(ctx, a); err != nil {
			return err
		}
	}
	return nil
}

// push moves current audit to entries if valid
func (b *AuditBuilder) push() {
	if b.audit.Action != "" || b.audit.FromStatus != nil {
		b.entries = append(b.entries, b.audit)
		b.audit = Audit{}
	}
}

func inferStatusAction(from, to string) string {
	return from + "_to_" + to
}
