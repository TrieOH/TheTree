// Package steps is the steps feature: its repo, handlers and operations.
// Build wires it.
package steps

import (
	"Informd/errs"
	"Informd/internal/authz"
	"Informd/models"
	"Informd/ports"
	"context"
	"lib/xslices"
	"uuid"
)

type Deps struct {
	Forms      ports.FormsRepo
	Steps      ports.StepRepo
	Namespaces ports.NamespaceRepo
	Authz      *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) BulkEdit(ctx context.Context, formID uuid.UUID, payload []models.UpdateFormStepInput) error {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	for _, p := range payload {
		if p.FormID != formID {
			return errs.ErrStepsAcrossForms
		}
	}

	steps := xslices.MapSlice(payload, models.UpdateFormStepInputToStep)
	return o.Steps.BulkEdit(ctx, steps)
}

func (o *Operations) Create(ctx context.Context, payload models.CreateFormStepInput) (*models.Step, error) {
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

func (o *Operations) List(ctx context.Context, formID uuid.UUID) ([]models.Step, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Steps.List(ctx, formID)
}
