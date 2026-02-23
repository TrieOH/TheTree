package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ActorType string

const (
	ActorTypeUnknown     ActorType = "unknown"
	ActorTypeOwner       ActorType = "owner"
	ActorTypeAdmin       ActorType = "admin"
	ActorTypeStaff       ActorType = "staff"
	ActorTypeParticipant ActorType = "participant"
	ActorTypePresenter   ActorType = "presenter"
	ActorTypeSystem      ActorType = "system"
)

type AuditActionState string

const (
	ActionStateSucceeded AuditActionState = "succeeded"
	ActionStateFailed    AuditActionState = "failed"
	ActionStatePending   AuditActionState = "pending"
	ActionStateUnset     AuditActionState = "unset"
)

type Audit struct {
	ID         uuid.UUID        `json:"id"`
	ResourceID uuid.UUID        `json:"resource_id"`
	ActorType  ActorType        `json:"actor_type"`
	ActorID    *uuid.UUID       `json:"actor_id"`
	Action     string           `json:"action"`
	State      AuditActionState `json:"state"`
	FromStatus *string          `json:"from_status"`
	ToStatus   *string          `json:"to_status"`
	Metadata   *json.RawMessage `json:"metadata"`
	CreatedAt  time.Time        `json:"created_at"`
}

type AuditBuilder struct {
	resourceID uuid.UUID
	actorType  ActorType
	actorID    *uuid.UUID
	audit      Audit // current building audit
	entries    []Audit
	metadata   map[string]any
}

// StartAudit initializes builder with base fields
func StartAudit(resourceID uuid.UUID, actorType ActorType, actorID *uuid.UUID) *AuditBuilder {
	return &AuditBuilder{
		resourceID: resourceID,
		actorType:  actorType,
		actorID:    actorID,
		entries:    make([]Audit, 0),
	}
}

func (b *AuditBuilder) GetActor() ActorType {
	return b.audit.ActorType
}

func (b *AuditBuilder) Continue(resourceID uuid.UUID, actorType ActorType, actorID *uuid.UUID) *AuditBuilder {
	b.resourceID = resourceID
	b.actorType = actorType
	b.actorID = actorID
	return b
}

// Action sets the action for current audit
func (b *AuditBuilder) Action(a string) *AuditBuilder {
	b.audit.Action = a
	return b
}

func (b *AuditBuilder) Actor(a ActorType) *AuditBuilder {
	b.actorType = a
	b.audit.ActorType = a
	return b
}

// State sets the action state for current audit
func (b *AuditBuilder) State(s AuditActionState) *AuditBuilder {
	b.audit.State = s
	return b
}

// StatusChange sets from/to status for current audit
func (b *AuditBuilder) StatusChange(from, to string) *AuditBuilder {
	b.audit.FromStatus = &from
	b.audit.ToStatus = &to
	return b
}

func (b *AuditBuilder) StatusFrom(from string) *AuditBuilder {
	b.audit.FromStatus = &from
	return b
}

func (b *AuditBuilder) StatusTo(to string) *AuditBuilder {
	b.audit.ToStatus = &to
	return b
}

func (b *AuditBuilder) AddMetadata(key string, value any) *AuditBuilder {
	if b.metadata == nil {
		b.metadata = make(map[string]any)
	}
	b.metadata[key] = value
	return b
}

func (b *AuditBuilder) GetAudits() []Audit {
	return b.entries
}

// Emit moves current audit to entries
func (b *AuditBuilder) Emit() {
	if b.audit.Action == "" {
		return
	}

	if b.resourceID == uuid.Nil {
		return
	}

	// Inject base fields
	b.audit.ResourceID = b.resourceID
	b.audit.ActorID = b.actorID
	b.audit.CreatedAt = time.Now()

	if b.audit.State == "" {
		b.audit.State = ActionStateUnset
	}

	if b.audit.ActorType == "" {
		b.audit.ActorType = b.actorType
	}

	if b.metadata != nil {
		if blob, err := json.Marshal(b.metadata); err == nil {
			raw := json.RawMessage(blob)
			b.audit.Metadata = &raw
		}
	}

	b.entries = append(b.entries, b.audit)

	// reset
	b.metadata = nil
	b.audit = Audit{State: ActionStateUnset}
}
