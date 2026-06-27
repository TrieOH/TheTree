package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Queries) List(ctx context.Context, formID uuid.UUID) ([]models.Step, error) {
	ctx, span := s.tracer.Start(ctx, "StepService.List")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.forms.GetMember(ctx, ident.Sub.ID, formID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	return s.steps.List(ctx, formID)
}
