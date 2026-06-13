package contracts

import (
	"time"

	"payssage/internal/shared/errx"
	"payssage/internal/shared/validation"

	"github.com/google/uuid"
)

type Workspace struct {
	ID        uuid.UUID `json:"id"`
	ScopeID   uuid.UUID `json:"scope_id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Sandbox   bool      `json:"sandbox"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewWorkspace(userID uuid.UUID, name string) (*Workspace, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("product").SetMessage("error generating uuid").SetCause(err)
	}

	w := &Workspace{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := w.validate(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Workspace) validate() error {
	return validation.Run(
		validation.RequireUUID("workspace", "user_id", w.UserID),
		validation.RequireString("workspace", "name", w.Name),
	)
}

func (w *Workspace) AddScope(scopeID uuid.UUID) {
	w.ScopeID = scopeID
}
