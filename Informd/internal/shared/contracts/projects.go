package contracts

import (
	"time"

	"github.com/MintzyG/FastUtilitiesNet"
	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id" validate:"required"`
	Name      string    `json:"name"     validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProject(ownerID uuid.UUID, name string) (*Project, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fun.Errf("error generating uuid for project: %s", err.Error()).Internal()
	}

	p := &Project{
		ID:        id,
		OwnerID:   ownerID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err = validate.Struct(p); err != nil {
		return nil, err
	}

	return p, nil
}
