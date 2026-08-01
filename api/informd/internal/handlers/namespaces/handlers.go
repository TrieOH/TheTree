// Package namespaces implements the StrictServerInterface methods for the
// namespaces feature, including the duplicated namespaced form routes
// (scheduled for removal; kept for parity).
package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/internal/services"
	"Informd/models"
)

const module = "Informd"

type Handlers struct {
	ops *services.Namespaces
}

func New(ops *services.Namespaces) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListNamespaces(ctx context.Context, _ openapi.ListNamespacesRequestObject) (openapi.ListNamespacesResponseObject, error) {
	namespaces, err := h.ops.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaces200JSONResponse{
		Code: 200, Data: &namespaces, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateNamespace(ctx context.Context, req openapi.CreateNamespaceRequestObject) (openapi.CreateNamespaceResponseObject, error) {
	namespace, err := h.ops.Create(ctx, req.Body.Name)
	if err != nil {
		return nil, err
	}
	return openapi.CreateNamespace201JSONResponse{
		Code: 201, Data: namespace, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListNamespaceMembers(ctx context.Context, req openapi.ListNamespaceMembersRequestObject) (openapi.ListNamespaceMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaceMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}

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

func (h *Handlers) RemoveNamespaceMember(ctx context.Context, req openapi.RemoveNamespaceMemberRequestObject) (openapi.RemoveNamespaceMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, models.RemoveNamespaceMemberInput{
		UserID:      req.Body.UserId,
		NamespaceID: req.NamespaceId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveNamespaceMember200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateNamespaceForm(ctx context.Context, req openapi.CreateNamespaceFormRequestObject) (openapi.CreateNamespaceFormResponseObject, error) {
	form, err := h.ops.CreateForm(ctx, req.Body.Title, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.CreateNamespaceForm201JSONResponse{
		Code: 201, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListNamespaceForms(ctx context.Context, req openapi.ListNamespaceFormsRequestObject) (openapi.ListNamespaceFormsResponseObject, error) {
	forms, err := h.ops.ListForms(ctx, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaceForms200JSONResponse{
		Code: 200, Data: &forms, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListNamespaceArchivedForms(ctx context.Context, req openapi.ListNamespaceArchivedFormsRequestObject) (openapi.ListNamespaceArchivedFormsResponseObject, error) {
	forms, err := h.ops.ListArchivedForms(ctx, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaceArchivedForms200JSONResponse{
		Code: 200, Data: &forms, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetFullFormNamespaced(ctx context.Context, req openapi.GetFullFormNamespacedRequestObject) (openapi.GetFullFormNamespacedResponseObject, error) {
	form, err := h.ops.GetFullForm(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFullFormNamespaced200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListFormMembersNamespaced(ctx context.Context, req openapi.ListFormMembersNamespacedRequestObject) (openapi.ListFormMembersNamespacedResponseObject, error) {
	members, err := h.ops.ListFormMembers(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFormMembersNamespaced200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) AddFormMemberNamespaced(ctx context.Context, req openapi.AddFormMemberNamespacedRequestObject) (openapi.AddFormMemberNamespacedResponseObject, error) {
	err := h.ops.AddFormMember(ctx, models.AddNamespaceFormMemberInput{
		UserID:      req.Body.UserId,
		NamespaceID: req.NamespaceId,
		FormID:      req.FormId,
		Role:        req.Body.Role,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddFormMemberNamespaced201Response{}, nil
}

func (h *Handlers) RemoveFormMemberNamespaced(ctx context.Context, req openapi.RemoveFormMemberNamespacedRequestObject) (openapi.RemoveFormMemberNamespacedResponseObject, error) {
	err := h.ops.RemoveFormMember(ctx, models.RemoveNamespaceFormMemberInput{
		UserID:      req.Body.UserId,
		NamespaceID: req.NamespaceId,
		FormID:      req.FormId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveFormMemberNamespaced200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) OpenFormNamespaced(ctx context.Context, req openapi.OpenFormNamespacedRequestObject) (openapi.OpenFormNamespacedResponseObject, error) {
	form, err := h.ops.OpenForm(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.OpenFormNamespaced200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CloseFormNamespaced(ctx context.Context, req openapi.CloseFormNamespacedRequestObject) (openapi.CloseFormNamespacedResponseObject, error) {
	form, err := h.ops.CloseForm(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.CloseFormNamespaced200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ArchiveFormNamespaced(ctx context.Context, req openapi.ArchiveFormNamespacedRequestObject) (openapi.ArchiveFormNamespacedResponseObject, error) {
	form, err := h.ops.ArchiveForm(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ArchiveFormNamespaced200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) RedraftFormNamespaced(ctx context.Context, req openapi.RedraftFormNamespacedRequestObject) (openapi.RedraftFormNamespacedResponseObject, error) {
	form, err := h.ops.ReDraftForm(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.RedraftFormNamespaced200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetFormResponseCountNamespaced(ctx context.Context, req openapi.GetFormResponseCountNamespacedRequestObject) (openapi.GetFormResponseCountNamespacedResponseObject, error) {
	count, err := h.ops.GetFormResponseCount(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFormResponseCountNamespaced200JSONResponse{
		Code: 200,
		Data: &struct {
			Count int `json:"count"`
		}{Count: count},
		Timestamp: time.Now(), Module: module,
	}, nil
}
