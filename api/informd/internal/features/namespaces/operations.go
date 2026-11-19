// Package namespaces is the namespaces feature: its repo, handlers and
// operations. Build wires it.
package namespaces

import (
	"Informd/errs"
	"Informd/internal/authz"
	"Informd/models"
	"Informd/ports"
	"context"
	"errors"
	"lib/database"
	idx "sdk/identityx"
	"time"
	"uuid"
)

type Deps struct {
	Namespaces ports.NamespaceRepo
	Forms      ports.FormsRepo
	Steps      ports.StepRepo
	Fields     ports.FieldsRepo
	Answers    ports.AnswerRepo
	Responses  ports.ResponseRepo
	Responders ports.ResponderRepo
	Authz      *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded from
	// boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) CreateForm(ctx context.Context, title string, namespaceID uuid.UUID) (*models.Form, error) {
	namespace, err := o.Namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	ident, err := o.Authz.Namespace(ctx, namespace.ID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	form, err := models.NewForm(&namespaceID, namespace.OwnerID, ident.Sub.ID, title)
	if err != nil {
		return nil, err
	}

	created, err := o.Forms.Create(ctx, *form)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (o *Operations) ListArchivedForms(ctx context.Context, namespaceID uuid.UUID) (forms []models.Form, err error) {
	_, err = o.Authz.Namespace(ctx, namespaceID, models.NamespaceMemberRoleMember)
	if err != nil {
		return nil, err
	}

	forms, err = o.Forms.ListFromNamespaceArchived(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}

func (o *Operations) ListForms(ctx context.Context, namespaceID uuid.UUID) (forms []models.Form, err error) {
	_, err = o.Authz.Namespace(ctx, namespaceID, models.NamespaceMemberRoleMember)
	if err != nil {
		return nil, err
	}

	forms, err = o.Forms.ListFromNamespace(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}

func (o *Operations) ListMembers(ctx context.Context, namespaceID uuid.UUID) (members []models.NamespaceMember, err error) {
	_, err = o.Authz.Namespace(ctx, namespaceID, models.NamespaceMemberRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = o.Namespaces.ListMembers(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (o *Operations) ListNamespaces(ctx context.Context) (members []models.Namespace, err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownNamespaces, err := o.Namespaces.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedNamespaces, err := o.Namespaces.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownNamespaces, joinedNamespaces...), nil
}

func (o *Operations) AddMember(ctx context.Context, payload models.AddNamespaceMemberInput) error {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return errs.ErrAddSelfToNamespace
	}

	namespace, err := o.Namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if payload.UserID == namespace.OwnerID {
		return errs.ErrAddOwnerToNamespace
	}

	_, err = o.Authz.Namespace(ctx, payload.NamespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Namespaces.GetMember(ctx, payload.UserID, namespace.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrAlreadyNamespaceMember
	}

	newMember := models.NamespaceMember{
		UserID:      payload.UserID,
		NamespaceID: payload.NamespaceID,
		Role:        payload.Role,
		AddedAt:     time.Now(),
		AddedBy:     ident.Sub.ID,
	}

	err = o.Namespaces.AddMember(ctx, newMember)
	if err != nil {
		return err
	}
	return nil
}

func (o *Operations) Create(ctx context.Context, name string) (*models.Namespace, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	project, err := models.NewNamespace(ident.Sub.ID, name)
	if err != nil {
		return nil, err
	}

	var created *models.Namespace
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Namespaces.Create(ctx, *project)
		if err != nil {
			return err
		}

		return o.Namespaces.AddMember(ctx, models.NamespaceMember{
			UserID:      ident.Sub.ID,
			NamespaceID: created.ID,
			Role:        models.NamespaceMemberRoleOwner,
			AddedAt:     time.Now(),
			AddedBy:     ident.Sub.ID,
		})
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveNamespaceMemberInput) error {
	namespace, err := o.Namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Namespace(ctx, payload.NamespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Namespaces.GetMember(ctx, payload.UserID, namespace.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err != nil {
		return errs.ErrNotNamespaceMember
	}

	return o.Namespaces.RemoveMember(ctx, payload.UserID, payload.NamespaceID)
}
