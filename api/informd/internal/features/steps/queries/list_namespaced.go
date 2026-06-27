package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Queries) ListNamespaced(ctx context.Context, formID, namespaceID uuid.UUID) ([]models.Step, error) {
	ctx, span := s.tracer.Start(ctx, "StepService.ListNamespaced")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.namespaces.GetMember(ctx, ident.Sub.ID, namespaceID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}
	if fun.Is(err, fun.CodeNotFound) {
		_, err = s.forms.GetMember(ctx, ident.Sub.ID, formID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	return s.steps.List(ctx, formID)
}
