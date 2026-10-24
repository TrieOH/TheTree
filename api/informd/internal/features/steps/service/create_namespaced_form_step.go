package service

import (
	"context"

	"Informd/models"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) CreateNamespaced(ctx context.Context, payload models.CreateNamespacedFormStepInput) (*models.Step, error) {
	step, err := models.NewStep(payload.FormID, payload.Title, payload.Description, payload.PositionHint)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Form(ctx, payload.FormID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	created, err := o.Steps.Create(ctx, *step)
	if err != nil {
		return nil, err
	}

	return created, nil
}
