package forms

import (
	"context"

	"Informd/internal/features/forms/service"
	"Informd/internal/openapi"
	"Informd/models"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) AddFormMember(ctx context.Context, req openapi.AddFormMemberRequestObject) (openapi.AddFormMemberResponseObject, error) {
	err := h.ops.AddMember(ctx, models.AddFormMemberInput{
		UserID: req.Body.UserId,
		FormID: req.FormId,
		Role:   req.Body.Role,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddFormMember201Response{}, nil
}

func (h *Handlers) ArchiveForm(ctx context.Context, req openapi.ArchiveFormRequestObject) (openapi.ArchiveFormResponseObject, error) {
	form, err := h.ops.Archive(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ArchiveForm200JSONResponse(*form), nil
}

func (h *Handlers) CloseForm(ctx context.Context, req openapi.CloseFormRequestObject) (openapi.CloseFormResponseObject, error) {
	form, err := h.ops.Close(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.CloseForm200JSONResponse(*form), nil
}

func (h *Handlers) CreateForm(ctx context.Context, req openapi.CreateFormRequestObject) (openapi.CreateFormResponseObject, error) {
	form, err := h.ops.Create(ctx, req.Body.Title)
	if err != nil {
		return nil, err
	}
	return openapi.CreateForm201JSONResponse(*form), nil
}

func (h *Handlers) GetAnswerableForm(ctx context.Context, req openapi.GetAnswerableFormRequestObject) (openapi.GetAnswerableFormResponseObject, error) {
	form, err := h.ops.GetAnswerable(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetAnswerableForm200JSONResponse(*form), nil
}

func (h *Handlers) GetFormResponseCount(ctx context.Context, req openapi.GetFormResponseCountRequestObject) (openapi.GetFormResponseCountResponseObject, error) {
	count, err := h.ops.GetResponseCount(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFormResponseCount200JSONResponse{Count: count}, nil
}

func (h *Handlers) GetFullForm(ctx context.Context, req openapi.GetFullFormRequestObject) (openapi.GetFullFormResponseObject, error) {
	form, err := h.ops.GetFull(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFullForm200JSONResponse(*form), nil
}

func (h *Handlers) ListFormMembers(ctx context.Context, req openapi.ListFormMembersRequestObject) (openapi.ListFormMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFormMembers200JSONResponse{Items: members}, nil
}

func (h *Handlers) ListMyArchivedForms(ctx context.Context, _ openapi.ListMyArchivedFormsRequestObject) (openapi.ListMyArchivedFormsResponseObject, error) {
	forms, err := h.ops.ListArchivedForms(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListMyArchivedForms200JSONResponse{Items: forms}, nil
}

func (h *Handlers) ListMyForms(ctx context.Context, _ openapi.ListMyFormsRequestObject) (openapi.ListMyFormsResponseObject, error) {
	forms, err := h.ops.ListForms(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListMyForms200JSONResponse{Items: forms}, nil
}

func (h *Handlers) OpenForm(ctx context.Context, req openapi.OpenFormRequestObject) (openapi.OpenFormResponseObject, error) {
	form, err := h.ops.Open(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.OpenForm200JSONResponse(*form), nil
}

func (h *Handlers) RedraftForm(ctx context.Context, req openapi.RedraftFormRequestObject) (openapi.RedraftFormResponseObject, error) {
	form, err := h.ops.ReDraft(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.RedraftForm200JSONResponse(*form), nil
}

func (h *Handlers) RemoveFormMember(ctx context.Context, req openapi.RemoveFormMemberRequestObject) (openapi.RemoveFormMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, models.RemoveFormMemberInput{
		UserID: req.Body.UserId,
		FormID: req.FormId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveFormMember204Response{}, nil
}
