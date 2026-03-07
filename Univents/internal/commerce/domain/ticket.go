package domain

import (
	"encoding/json"
	"fmt"
	"time"
	"univents/internal/shared/errx"
	"univents/internal/shared/validation"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type Ticket struct {
	ID          uuid.UUID `json:"id"`
	ScopeID     uuid.UUID `json:"scope_id"`
	EditionID   uuid.UUID `json:"edition_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`

	CreatedBy uuid.UUID  `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateTicketSpec struct {
	EditionID   uuid.UUID `json:"edition_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}

func NewTicket(creatorID uuid.UUID, spec CreateTicketSpec) (*Ticket, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("ticket").SetMessage("error generating uuid").SetCause(err)
	}

	t := &Ticket{
		ID:          id,
		EditionID:   spec.EditionID,
		Name:        spec.Name,
		Description: spec.Description,
		CreatedBy:   creatorID,
	}

	if err := t.validate(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Ticket) AddScope(scopeID uuid.UUID) {
	t.ScopeID = scopeID
}

func (t *Ticket) validate() error {
	return validation.Run(
		validation.RequireUUID("ticket", "edition_id", t.EditionID),
		validation.RequireUUID("ticket", "created_by", t.CreatedBy),
		validation.RequireString("ticket", "name", t.Name),
	)
}

type PermissionType string

const (
	PermissionTypeActivity   PermissionType = "activity"
	PermissionTypeProduct    PermissionType = "product"
	PermissionTypeCheckpoint PermissionType = "checkpoint"
)

type TicketPermission struct {
	ID             uuid.UUID      `json:"id"`
	TicketID       uuid.UUID      `json:"ticket_id"`
	PermissionType PermissionType `json:"permission_type"`
	ActivityID     *uuid.UUID     `json:"activity_id"`
	ProductID      *uuid.UUID     `json:"product_id"`
	CheckpointID   *uuid.UUID     `json:"checkpoint_id"`
	CreatedAt      time.Time      `json:"created_at"`
}

type CreateTicketPermissionSpec struct {
	TicketID       uuid.UUID      `json:"ticket_id"`
	PermissionType PermissionType `json:"permission_type"`
	ActivityID     *uuid.UUID     `json:"activity_id"`
	ProductID      *uuid.UUID     `json:"product_id"`
	CheckpointID   *uuid.UUID     `json:"checkpoint_id"`
}

func NewTicketPermission(creatorID uuid.UUID, spec CreateTicketPermissionSpec) (*TicketPermission, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("ticket").SetMessage("error generating uuid").SetCause(err)
	}

	tp := &TicketPermission{
		ID:             id,
		TicketID:       spec.TicketID,
		PermissionType: spec.PermissionType,
		ActivityID:     spec.ActivityID,
		ProductID:      spec.ProductID,
		CheckpointID:   spec.CheckpointID,
	}

	if err := tp.validate(); err != nil {
		return nil, err
	}

	return tp, nil
}

func (tp *TicketPermission) validate() error {
	return validation.Run(
		validation.RequireUUID("ticket_permission", "ticket_id", tp.TicketID),
		validation.RequireString("ticket_permission", "permission_type", string(tp.PermissionType)),
		validation.Assert("ticket_permission",
			tp.PermissionType == PermissionTypeActivity ||
				tp.PermissionType == PermissionTypeProduct ||
				tp.PermissionType == PermissionTypeCheckpoint,
			"invalid permission type",
		),
		validation.AssertIf("ticket_permission",
			func() bool { return tp.PermissionType == PermissionTypeActivity },
			func() bool {
				return tp.ActivityID != nil && *tp.ActivityID != uuid.Nil && tp.ProductID == nil && tp.CheckpointID == nil
			},
			"activity permission must have activity_id only",
		),
		validation.AssertIf("ticket_permission",
			func() bool { return tp.PermissionType == PermissionTypeProduct },
			func() bool {
				return tp.ProductID != nil && *tp.ProductID != uuid.Nil && tp.ActivityID == nil && tp.CheckpointID == nil
			},
			"product permission must have product_id only",
		),
		validation.AssertIf("ticket_permission",
			func() bool { return tp.PermissionType == PermissionTypeCheckpoint },
			func() bool {
				return tp.CheckpointID != nil && *tp.CheckpointID != uuid.Nil && tp.ActivityID == nil && tp.ProductID == nil
			},
			"checkpoint permission must have checkpoint_id only",
		),
	)
}

const (
	TypeGrantTicketPermissions = "ticket:grant_permissions"
	MaxGrantRetries            = 5
)

type TicketGrant struct {
	TicketID uuid.UUID `json:"ticket_id"`
	UserID   uuid.UUID `json:"user_id"` // assigned_to_user_id if set, else buyer
}

type GrantTicketPermissionsPayload struct {
	Grants []TicketGrant `json:"grants"`
}

func NewGrantTicketPermissionsTask(grants []TicketGrant, paymentID string) (*asynq.Task, error) {
	payload, err := json.Marshal(GrantTicketPermissionsPayload{Grants: grants})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal grant permissions payload: %w", err)
	}
	return asynq.NewTask(TypeGrantTicketPermissions, payload,
		asynq.TaskID(fmt.Sprintf("%s:%s", paymentID, TypeReservationExpired)),
		asynq.MaxRetry(MaxGrantRetries),
	), nil
}
