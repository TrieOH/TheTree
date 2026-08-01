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
func (o *Operations) ArchiveForm(ctx context.Context, namespaceID, formID uuid.UUID) (*models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.ArchiveForm")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	form, err := o.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusClosed {
		return nil, fun.ErrBadRequest("cannot archive a form not on closed")
	}

	err = authz.Service.CheckNamespace(ctx, ident.Sub.ID, namespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.forms.Close(ctx, form.ID)
}
