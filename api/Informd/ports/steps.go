package ports

import (
	"Informd/contracts"
	"context"

	"github.com/google/uuid"
)

type StepRepo interface {
	Create(ctx context.Context, toCreate contracts.Step) (*contracts.Step, error)
	List(ctx context.Context, formID uuid.UUID) ([]contracts.Step, error)
}
