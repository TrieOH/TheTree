package dto

import (
	"TriePayments/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

type WorkspaceResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Sandbox   bool      `json:"sandbox"`
	CreatedAt time.Time `json:"created_at"`
}

func MapWorkspaceResponse(w *domain.Workspace) WorkspaceResponse {
	return WorkspaceResponse{
		ID:        w.ID,
		Name:      w.Name,
		Sandbox:   w.Sandbox,
		CreatedAt: w.CreatedAt,
	}
}
