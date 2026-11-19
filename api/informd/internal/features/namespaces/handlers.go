package namespaces

import (
	"context"

	"Informd/internal/openapi"
	"Informd/models"
)

type Handlers struct {
	ops *Operations
}

func NewHandlers(ops *Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) AddNamespaceMember(ctx context.Context, req openapi.AddNamespaceMemberRequestObject) (openapi.AddNamespaceMemberResponseObject, error) {
	err := h.ops.AddMember(ctx, models.AddNamespaceMemberInput{
		UserID:      req.Body.UserId,
		Role:        req.Body.Role,
		NamespaceID: req.NamespaceId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddNamespaceMember201Response{}, nil
}

func (h *Handlers) CreateNamespace(ctx context.Context, req openapi.CreateNamespaceRequestObject) (openapi.CreateNamespaceResponseObject, error) {
	namespace, err := h.ops.Create(ctx, req.Body.Name)
	if err != nil {
		return nil, err
	}
	return openapi.CreateNamespace201JSONResponse(*namespace), nil
}

func (h *Handlers) CreateNamespaceForm(ctx context.Context, req openapi.CreateNamespaceFormRequestObject) (openapi.CreateNamespaceFormResponseObject, error) {
	form, err := h.ops.CreateForm(ctx, req.Body.Title, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.CreateNamespaceForm201JSONResponse(*form), nil
}

func (h *Handlers) ListNamespaceArchivedForms(ctx context.Context, req openapi.ListNamespaceArchivedFormsRequestObject) (openapi.ListNamespaceArchivedFormsResponseObject, error) {
	forms, err := h.ops.ListArchivedForms(ctx, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaceArchivedForms200JSONResponse{Items: forms}, nil
}

func (h *Handlers) ListNamespaceForms(ctx context.Context, req openapi.ListNamespaceFormsRequestObject) (openapi.ListNamespaceFormsResponseObject, error) {
	forms, err := h.ops.ListForms(ctx, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaceForms200JSONResponse{Items: forms}, nil
}

func (h *Handlers) ListNamespaceMembers(ctx context.Context, req openapi.ListNamespaceMembersRequestObject) (openapi.ListNamespaceMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaceMembers200JSONResponse{Items: members}, nil
}

func (h *Handlers) ListNamespaces(ctx context.Context, _ openapi.ListNamespacesRequestObject) (openapi.ListNamespacesResponseObject, error) {
	namespaces, err := h.ops.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaces200JSONResponse{Items: namespaces}, nil
}

func (h *Handlers) RemoveNamespaceMember(ctx context.Context, req openapi.RemoveNamespaceMemberRequestObject) (openapi.RemoveNamespaceMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, models.RemoveNamespaceMemberInput{
		UserID:      req.Body.UserId,
		NamespaceID: req.NamespaceId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveNamespaceMember204Response{}, nil
}
