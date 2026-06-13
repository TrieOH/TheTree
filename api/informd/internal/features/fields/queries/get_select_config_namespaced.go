package queries

import (
	"context"

	"Informd/models"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Queries) GetSelectConfigNamespaced(ctx context.Context, formID, namespaceID, fieldID uuid.UUID) (*models.FieldSelectConfig, error) {
	ctx, span := s.tracer.Start(ctx, "FieldService.GetSelectConfigNamespaced")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.namespaces.GetMember(ctx, sub.ID, namespaceID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}
	if fun.Is(err, fun.CodeNotFound) {
		_, err = s.forms.GetMember(ctx, sub.ID, formID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	return s.fields.GetSelectConfig(ctx, fieldID)
}
