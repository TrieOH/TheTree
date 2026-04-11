package ports

import (
	"TrieForms/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type FormsRepo interface {
	Create(ctx context.Context, toCreate contracts.Form) (*contracts.Form, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]contracts.Form, error)
}
