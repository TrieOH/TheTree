package contracts

import (
	"time"

	"univents/internal/shared/errx"
	"univents/internal/shared/validation"

	"github.com/google/uuid"
)

type CheckpointType string

const (
	CheckpointTypeEntry   CheckpointType = "entry"
	CheckpointTypeZone    CheckpointType = "zone"
	CheckpointTypeAmenity CheckpointType = "amenity"
	CheckpointTypeSession CheckpointType = "session"
	CheckpointTypeExit    CheckpointType = "exit"
)

type CheckpointAccess string

const (
	CheckpointAccessOpen      CheckpointAccess = "open"
	CheckpointAccessTicket    CheckpointAccess = "ticket"
	CheckpointAccessStaffOnly CheckpointAccess = "staff_only"
)

type Checkpoint struct {
	ID         uuid.UUID        `json:"id"`
	ScopeID    uuid.UUID        `json:"scope_id"`
	EditionID  uuid.UUID        `json:"edition_id"`
	Name       string           `json:"name"`
	Type       CheckpointType   `json:"type"`
	AccessMode CheckpointAccess `json:"access_mode"`
	StartsAt   *time.Time       `json:"starts_at"`
	EndsAt     *time.Time       `json:"ends_at"`
	CreatedBy  uuid.UUID        `json:"created_by"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  *time.Time       `json:"deleted_at"`
}

type CreateCheckpointSpec struct {
	EditionID  uuid.UUID        `json:"edition_id"`
	StartsAt   *time.Time       `json:"starts_at"`
	EndsAt     *time.Time       `json:"ends_at"`
	Name       string           `json:"name"`
	Type       CheckpointType   `json:"type"`
	AccessMode CheckpointAccess `json:"access_mode"`
}

func NewCheckpoint(creatorID uuid.UUID, spec CreateCheckpointSpec, edition *Edition) (*Checkpoint, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("checkpoint").SetMessage("error generating uuid").SetCause(err)
	}

	c := &Checkpoint{
		ID:         id,
		EditionID:  spec.EditionID,
		Name:       spec.Name,
		StartsAt:   spec.StartsAt,
		EndsAt:     spec.EndsAt,
		Type:       spec.Type,
		AccessMode: spec.AccessMode,
		CreatedBy:  creatorID,
	}

	if err := c.validate(edition); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Checkpoint) validate(edition *Edition) error {
	return validation.Run(
		validation.RequireUUID("checkpoint", "edition_id", c.EditionID),
		validation.RequireUUID("checkpoint", "created_by", c.CreatedBy),
		validation.RequireString("checkpoint", "name", c.Name),
		validation.AssertIf("checkpoint",
			func() bool { return c.StartsAt != nil && c.EndsAt != nil },
			func() bool { return c.StartsAt.Before(*c.EndsAt) },
			"checkpoint start must be before end",
		),
		validation.AssertIf("checkpoint",
			func() bool { return c.StartsAt != nil },
			func() bool { return !c.StartsAt.Before(edition.StartsAt) && c.StartsAt.Before(edition.EndsAt) },
			"checkpoint start must be within edition duration",
		),
		validation.AssertIf("checkpoint",
			func() bool { return c.EndsAt != nil },
			func() bool { return c.EndsAt.After(edition.StartsAt) && !c.EndsAt.After(edition.EndsAt) },
			"checkpoint end must be within edition duration",
		),
	)
}

func (e *Checkpoint) AddScope(scopeID uuid.UUID) {
	e.ScopeID = scopeID
}
