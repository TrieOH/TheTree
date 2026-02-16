package inbounds

import (
	"GoAuth/internal/domain/scopes"
	"time"

	"github.com/google/uuid"
)

type CreateScopeInput struct {
	Name       string
	ExternalID *string
	ProjectID  uuid.UUID
	ParentID   *uuid.UUID // nil = default to project root
}

type GetScopeInput struct {
	ProjectID uuid.UUID
	ScopeID   uuid.UUID
}

type ScopeOutput struct {
	ID         uuid.UUID
	Type       scopes.ScopeType
	ParentID   *uuid.UUID
	ProjectID  *uuid.UUID
	Name       *string
	ExternalID *string
	CreatedAt  time.Time
}

func ScopeToScopeOutput(scope *scopes.Scope) *ScopeOutput {
	if scope == nil {
		return nil
	}

	return &ScopeOutput{
		ID:         scope.ID,
		Type:       scope.Type,
		ParentID:   scope.ParentID,
		ProjectID:  scope.ProjectID,
		Name:       scope.Name,
		ExternalID: scope.ExternalID,
		CreatedAt:  scope.CreatedAt,
	}
}

func ScopeSliceToScopeSliceOutput(scopes []scopes.Scope) []ScopeOutput {
	if scopes == nil {
		return nil
	}

	result := make([]ScopeOutput, 0, len(scopes))
	for _, scope := range scopes {
		result = append(result, ScopeOutput{
			ID:         scope.ID,
			Type:       scope.Type,
			ParentID:   scope.ParentID,
			ProjectID:  scope.ProjectID,
			Name:       scope.Name,
			ExternalID: scope.ExternalID,
			CreatedAt:  scope.CreatedAt,
		})
	}
	return result
}
