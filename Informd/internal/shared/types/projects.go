package types

import (
	"TrieForms/internal/shared/validation"
	"time"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	ScopeID   uuid.UUID `json:"scope_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProject(ownerID uuid.UUID, name string) (*Project, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fun.NewErrorf("error generating uuid for project: %s", err.Error()).Internal()
	}

	w := &Project{
		ID:        id,
		OwnerID:   ownerID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := w.validate(); err != nil {
		return nil, err
	}

	return w, nil
}

func (p *Project) validate() error {
	return validation.Run(
		validation.RequireUUID("project", "owner_id", p.OwnerID),
		validation.RequireString("project", "name", p.Name),
	)
}

func (p *Project) AddScope(scopeID uuid.UUID) {
	p.ScopeID = scopeID
}
