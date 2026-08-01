package namespaces

import (
	"context"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — CheckForm already anchors via the form's namespace.
func (o *Operations) OpenForm(ctx context.Context, namespaceID, formID uuid.UUID) (*models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.OpenForm")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	form, err := o.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusDraft {
		return nil, fun.ErrBadRequest("cannot open a form not on draft")
	}

	err = authz.Service.CheckNamespace(ctx, ident.Sub.ID, namespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.forms.Open(ctx, form.ID)
}
