package models

import (
	"time"

	"github.com/google/uuid"
)

type Namespace struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id" validate:"required"`
	Name      string    `json:"name"     validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewNamespace(ownerID uuid.UUID, name string) (*Namespace, error) {
	p := &Namespace{
		OwnerID: ownerID,
		Name:    name,
	}
	if err := validate.Struct(p); err != nil {
		return nil, err
	}
	return p, nil
}
